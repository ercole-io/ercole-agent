// Copyright (c) 2020 Sorint.lab S.p.A.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package fetcher

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/ercole-io/ercole-agent/config"
	"github.com/sirupsen/logrus"
)

// LinuxFetcherImpl SpecializedFetcher implementation for linux
type LinuxFetcherImpl struct {
	configuration config.Configuration
	log           *logrus.Logger
}

// NewLinuxFetcherImpl constructor
func NewLinuxFetcherImpl(conf config.Configuration, log *logrus.Logger) LinuxFetcherImpl {
	return LinuxFetcherImpl{
		conf,
		log,
	}
}

// Execute Execute specific fetcher by name
func (lf *LinuxFetcherImpl) Execute(fetcherName string, params ...string) []byte {
	cmdName := config.GetBaseDir() + "/fetch/linux/" + fetcherName + ".sh"
	lf.log.Infof("Fetching %s %s", cmdName, strings.Join(params, " "))

	cmd := exec.Command(cmdName, params...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if len(stdout.Bytes()) > 0 {
		lf.log.Debug(string(stdout.Bytes()))
	}

	if len(stderr.Bytes()) > 0 {
		lf.log.Error(string(stderr.Bytes()))
	}

	if err != nil {
		if fetcherName == "dbstatus" {
			return []byte("UNREACHABLE")
		}

		lf.log.Fatalf("Fatal error running [%s %s]: [%v]", cmdName, strings.Join(params, " "), err)
	}

	return stdout.Bytes()
}
