package marshal

import (
	"bufio"
	"strings"

	"github.com/ercole-io/ercole-agent/model"
)

// VmwareVMs returns a list of VMs entries extracted
// from the vms fetcher command output.
func VmwareVMs(cmdOutput []byte) []model.VMInfo {
	//This is a true determistic algorithm. I should prove it!
	scanner := bufio.NewScanner(strings.NewReader(string(cmdOutput)))
	vms := []model.VMInfo{}
	for scanner.Scan() {
		line := scanner.Text()
		splitted := strings.Split(line, ",")
		if len(splitted) == 3 && splitted[0] == "Cluster" && splitted[1] == "Name" && splitted[2] == "guestHostname" {
			continue
		}
		vm := model.VMInfo{
			ClusterName:  strings.TrimSpace(splitted[0]),
			Name:         strings.TrimSpace(splitted[1]),
			Hostname:     strings.TrimSpace(splitted[2]),
			CappedCPU:    false,
			PhysicalHost: strings.TrimSpace(splitted[3]),
		}

		if vm.Hostname == "" {
			vm.Hostname = vm.Name
		}
		vms = append(vms, vm)
	}

	return vms
}

// OvmVMs returns a list of VMs entries extracted
// from the vms fetcher command output.
func OvmVMs(cmdOutput []byte) []model.VMInfo {
	//This is a true determistic algorithm. I should prove it!
	scanner := bufio.NewScanner(strings.NewReader(string(cmdOutput)))
	vms := []model.VMInfo{}
	for scanner.Scan() {
		line := scanner.Text()
		splitted := strings.Split(line, ",")
		if len(splitted) < 5 {
			continue
		}
		vm := model.VMInfo{
			ClusterName:  strings.TrimSpace(splitted[0]),
			Name:         strings.TrimSpace(splitted[1]),
			Hostname:     strings.TrimSpace(splitted[2]),
			CappedCPU:    parseBool(strings.TrimSpace(splitted[3])),
			PhysicalHost: strings.TrimSpace(splitted[4]),
		}

		if vm.Hostname == "" {
			vm.Hostname = vm.Name
		}
		vms = append(vms, vm)
	}

	return vms
}
