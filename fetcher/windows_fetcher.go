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
	"fmt"
	"os/exec"
	"strings"

	"github.com/ercole-io/ercole-agent/config"
	"github.com/ercole-io/ercole-agent/logger"
	"github.com/ercole-io/ercole-agent/model"
)

// WindowsFetcherImpl SpecializedFetcher implementation for windows
type WindowsFetcherImpl struct {
	Configuration config.Configuration
	log           logger.Logger
}

const not_implemented = "Not yet implemented for Windows"

// NewWindowsFetcherImpl constructor
func NewWindowsFetcherImpl(conf config.Configuration, log logger.Logger) WindowsFetcherImpl {
	return WindowsFetcherImpl{
		conf,
		log,
	}
}

// Execute Execute specific fetcher by name
func (wf *WindowsFetcherImpl) Execute(fetcherName string, params ...string) []byte {
	var (
		cmd    *exec.Cmd
		err    error
		psexe  string
		stdout bytes.Buffer
		stderr bytes.Buffer
	)

	baseDir := config.GetBaseDir()

	psexe, err = exec.LookPath("powershell.exe")
	if err != nil {
		wf.log.Fatal(psexe)
	}

	if wf.Configuration.ForcePwshVersion == "0" {
		params = append([]string{"-ExecutionPolicy", "Bypass", "-File", baseDir + "\\fetch\\win.ps1", "-s", fetcherName}, params...)
	} else {
		params = append([]string{"-version", wf.Configuration.ForcePwshVersion, "-ExecutionPolicy", "Bypass", "-File", baseDir + "\\fetch\\win.ps1", "-s", fetcherName}, params...)
	}

	wf.log.Info("Fetching " + psexe + " " + strings.Join(params, " "))

	cmd = exec.Command(psexe, params...)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if len(stderr.Bytes()) > 0 {
		wf.log.Error(string(stderr.Bytes()))
	}

	if err != nil {
		if fetcherName != "dbstatus" {
			return []byte("UNREACHABLE")
		}

		wf.log.Fatal(err)
	}

	return stdout.Bytes()
}

// SetUser not implemented
func (wf *WindowsFetcherImpl) SetUser(username string) error {
	wf.log.Error(not_implemented)
	return fmt.Errorf(not_implemented)
}

// SetUserAsCurrent set user used by fetcher to run commands as current process user
func (wf *WindowsFetcherImpl) SetUserAsCurrent() error {
	wf.log.Error(not_implemented)
	return fmt.Errorf(not_implemented)
}

// GetClusters not implemented
func (wf *WindowsFetcherImpl) GetClusters(hv config.Hypervisor) []model.ClusterInfo {
	wf.log.Error(not_implemented)

	return make([]model.ClusterInfo, 0)
}

// GetVirtualMachines return VMWare virtual machines infos from the given hyperVisor
func (wf *WindowsFetcherImpl) GetVirtualMachines(hv config.Hypervisor) []model.VMInfo {
	wf.log.Error(not_implemented)

	return make([]model.VMInfo, 0)
}

// GetExadataDevices get
func (wf *WindowsFetcherImpl) GetExadataDevices() []model.ExadataDevice {
	wf.log.Error(not_implemented)

	return make([]model.ExadataDevice, 0)
}

// GetExadataCellDisks get
func (wf *WindowsFetcherImpl) GetExadataCellDisks() []model.ExadataCellDisk {
	wf.log.Error(not_implemented)

	return make([]model.ExadataCellDisk, 0)
}
