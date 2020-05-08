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
	"log"
	"os/exec"
	"strings"

	"github.com/ercole-io/ercole-agent/config"
)

// LinuxFetcherImpl implemenentation
type LinuxFetcherImpl struct {
	Configuration config.Configuration
}

// Execute execute specific fetcher by name
func (cf *LinuxFetcherImpl) Execute(fetcherName string, params ...string) []byte {
	var (
		cmd    *exec.Cmd
		err    error
		stdout bytes.Buffer
		stderr bytes.Buffer
	)

	baseDir := config.GetBaseDir()

	log.Println("Fetching " + baseDir + "/fetch/" + fetcherName + " " + strings.Join(params, " "))
	cmd = exec.Command(baseDir+"/fetch/"+fetcherName, params...)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()

	if len(stderr.Bytes()) > 0 {
		log.Print(string(stderr.Bytes()))
	}

	if err != nil {
		if fetcherName != "dbstatus" {
			log.Fatal(err)
		} else {
			return []byte("UNREACHABLE") // fallback
		}
	}

	return stdout.Bytes()
}
