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
	"os/exec"
	"os/user"
	"strconv"
	"strings"

	"github.com/ercole-io/ercole-agent/config"
	"github.com/ercole-io/ercole-agent/logger"
	"github.com/ercole-io/ercole-agent/marshal"
	"github.com/ercole-io/ercole-agent/model"
)

// LinuxFetcherImpl SpecializedFetcher implementation for linux
type LinuxFetcherImpl struct {
	configuration config.Configuration
	log           logger.Logger
	fetcherUser   *User
}

// NewLinuxFetcherImpl constructor
func NewLinuxFetcherImpl(conf config.Configuration, log logger.Logger) LinuxFetcherImpl {
	return LinuxFetcherImpl{
		conf,
		log,
		nil,
	}
}

// Execute execute bash script by name
func (lf *LinuxFetcherImpl) Execute(fetcherName string, params ...string) []byte {
	commandName := config.GetBaseDir() + "/fetch/linux/" + fetcherName + ".sh"
	lf.log.Infof("Fetching %s %s", commandName, strings.Join(params, " "))

	stdout, stderr, exitCode, err := runCommandAs(lf.log, lf.fetcherUser, commandName, params...)

	if len(stdout) > 0 {
		lf.log.Debugf("Fetcher [%s] stdout: [%v]", fetcherName, string(stdout))
	}

	if len(stderr) > 0 {
		lf.log.Errorf("Fetcher [%s] exitCode: [%v] stderr: [%v]", fetcherName, exitCode, string(stderr))
	}

	if err != nil {
		if fetcherName == "dbstatus" {
			return []byte("UNREACHABLE")
		}

		lf.log.Fatalf("Fatal error running [%s %s]: [%v]", commandName, strings.Join(params, " "), err)
	}

	return stdout
}

//func (lf *LinuxFetcherImpl) runCommand(commandName string, args ...string) (stdout, stderr []byte, exitCode int, err error) {
//	cmd := exec.Command(commandName, args...)
//
//	if lf.fetcherUser != nil {
//		lf.log.Debugf("runCommand [%v] with user [%v]", commandName, lf.fetcherUser)
//
//		cmd.SysProcAttr = &syscall.SysProcAttr{}
//		if err != nil {
//			lf.log.Errorf("Can't set process attributes at command [%v]", commandName)
//			return nil, nil, -1, err
//		}
//		cmd.SysProcAttr.Credential = &syscall.Credential{Uid: lf.fetcherUser.uid, Gid: lf.fetcherUser.gid}
//	}
//
//	stdout, err = cmd.Output()
//
//	if err != nil {
//		if exitErr, ok := err.(*exec.ExitError); ok {
//			exitCode = exitErr.ExitCode()
//			stderr = exitErr.Stderr
//		} else {
//			exitCode = -1
//		}
//	}
//
//	return
//}

// SetUser set user used by fetcher to run commands
func (lf *LinuxFetcherImpl) SetUser(username string) error {
	u, err := lf.getUserInfo(username)
	if err != nil {
		return err
	}

	lf.fetcherUser = u
	return nil
}

func (lf *LinuxFetcherImpl) getUserInfo(username string) (*User, error) {
	u, err := user.Lookup(username)
	if err != nil {
		lf.log.Errorf("Can't lookup username [%s], error: [%v]", username, err)
		return nil, err
	}

	intUID, err := strconv.Atoi(u.Uid)
	if err != nil {
		lf.log.Errorf("Can't convert uid [%s], error: [%v]", u.Uid, err)
		return nil, err
	}

	intGID, err := strconv.Atoi(u.Gid)
	if err != nil {
		lf.log.Errorf("Can't convert gid [%s], error: [%v]", u.Gid, err)
		return nil, err
	}

	return &User{u.Name, uint32(intUID), uint32(intGID)}, nil
}

// SetUserAsCurrent set user used by fetcher to run commands as current process user
func (lf *LinuxFetcherImpl) SetUserAsCurrent() error {
	lf.fetcherUser = nil
	return nil
}

//TODO Use fetcher user
// executePwsh execute pwsh script by name
func (lf *LinuxFetcherImpl) executePwsh(fetcherName string, args ...string) []byte {
	scriptPath := config.GetBaseDir() + "/fetch/linux/" + fetcherName
	cmdSlice := append([]string{"/usr/bin/pwsh", scriptPath}, args...)
	cmd := strings.Join(cmdSlice, " ")

	lf.log.Infof("Fetching [%v]", cmd)

	stdout, err := exec.Command(cmd).Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			lf.log.Fatalf("Fetcher [%s] exitCode: [%v] stderr: [%v]", fetcherName, exitErr.ExitCode(), string(exitErr.Stderr))
		} else {
			lf.log.Fatalf("Fetcher [%s] error: [%v]", fetcherName, err.Error())
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

// GetExadataDevices get
func (lf *LinuxFetcherImpl) GetExadataDevices() []model.ExadataDevice {
	out := lf.Execute("exadata/info")
	return marshal.ExadataDevices(out)
}

// GetExadataCellDisks get
func (lf *LinuxFetcherImpl) GetExadataCellDisks() []model.ExadataCellDisk {
	out := lf.Execute("exadata/storage-status")
	return marshal.ExadataCellDisks(out)
}
