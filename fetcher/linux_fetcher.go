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
	"github.com/ercole-io/ercole-agent/marshal"
	"github.com/ercole-io/ercole-agent/model"
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

// Execute execute bash script by name
func (lf *LinuxFetcherImpl) Execute(fetcherName string, params ...string) []byte {
	cmdName := config.GetBaseDir() + "/fetch/linux/" + fetcherName + ".sh"
	lf.log.Infof("Fetching %s %s", cmdName, strings.Join(params, " "))

	cmd := exec.Command(cmdName, params...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if len(stdout.Bytes()) > 0 {
		lf.log.Debugf("Fetcher [%s] stdout: [%v]", fetcherName, string(stdout.Bytes()))
	}

	if len(stderr.Bytes()) > 0 {
		lf.log.Errorf("Fetcher [%s] stderr: [%v]", fetcherName, string(stderr.Bytes()))
	}

	if err != nil {
		if fetcherName == "dbstatus" {
			return []byte("UNREACHABLE")
		}

		lf.log.Fatalf("Fatal error running [%s %s]: [%v]", cmdName, strings.Join(params, " "), err)
	}

	return stdout.Bytes()
}

// executePwsh execute pwsh script by name
func (lf *LinuxFetcherImpl) executePwsh(fetcherName string, args ...string) []byte {
	scriptPath := config.GetBaseDir() + "/fetch/linux/" + fetcherName
	cmdSlice := append([]string{"/usr/bin/pwsh", scriptPath}, args...)
	cmd := strings.Join(cmdSlice, " ")

	lf.log.Infof("Fetching [%v]", cmd)

	stdout, err := exec.Command(cmd).Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			lf.log.Fatal("Fetcher [%s] exitCode: [%v] stderr: [%v]", fetcherName, exitErr.ExitCode, string(exitErr.Stderr))
		} else {
			lf.log.Fatal("Fetcher [%s] error: [%v]", fetcherName, err.Error())
		}
	}

	lf.log.Debugf("Fetcher [%s] stdout: [%v]", fetcherName, stdout)

	return stdout
}

// GetClusters return VMWare clusters from the given hyperVisor
func (lf *LinuxFetcherImpl) GetClusters(hv config.Hypervisor) []model.ClusterInfo {
	var out []byte

	switch hv.Type {
	case "vmware":
		out = lf.executePwsh("vmware.ps1", "-s", "cluster", hv.Endpoint, hv.Username, hv.Password)

	case "ovm":
		out = lf.Execute("ovm", "cluster", hv.Endpoint, hv.Username, hv.Password, hv.OvmUserKey, hv.OvmControl)

	default:
		lf.log.Errorf("Hypervisor not supported: %v (%v)", hv.Type, hv)
		return make([]model.ClusterInfo, 0)
	}

	fetchedClusters := marshal.Clusters(out)
	for i := range fetchedClusters {
		fetchedClusters[i].Type = hv.Type
	}

	return fetchedClusters
}

// GetVirtualMachines return VMWare virtual machines infos from the given hyperVisor
func (lf *LinuxFetcherImpl) GetVirtualMachines(hv config.Hypervisor) []model.VMInfo {
	var vms []model.VMInfo

	switch hv.Type {
	case "vmware":
		out := lf.executePwsh("vmware.ps1", "-s", "vms", hv.Endpoint, hv.Username, hv.Password)
		vms = marshal.VmwareVMs(out)

	case "ovm":
		out := lf.Execute("ovm", "vms", hv.Endpoint, hv.Username, hv.Password, hv.OvmUserKey, hv.OvmControl)
		vms = marshal.OvmVMs(out)

	default:
		lf.log.Errorf("Hypervisor not supported: %v (%v)", hv.Type, hv)
		return make([]model.VMInfo, 0)
	}

	lf.log.Debugf("Got %d vms from hypervisor: %s", len(vms), hv.Endpoint)

	return vms
}
