// Copyright (c) 2022 Sorint.lab S.p.A.
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

package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/ercole-io/ercole-agent/v2/utils"
	"github.com/ercole-io/ercole/v2/model"
	ercutils "github.com/ercole-io/ercole/v2/utils"
	"github.com/hashicorp/go-multierror"
)

type ResponseOLVMClusters struct {
	Clusters []OLVMCluster `json:"cluster"`
}

type OLVMCluster struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type ResponseOLVMHosts struct {
	Hosts []OLVMHost `json:"host"`
}

type OLVMHost struct {
	Name string `json:"name"`
	Id   string `json:"id"`
	Cpu  struct {
		Topology struct {
			Cores   string `json:"cores"`
			Sockets string `json:"sockets"`
		}
	}
	Cluster struct {
		Id string `json:"id"`
	}
}

type ResponseOLVMVMs struct {
	VMs []OLVMVM `json:"vm"`
}

type OLVMVM struct {
	Name string `json:"name"`
	Fqdn string `json:"fqdn"`
	Id   string `json:"id"`
	Cpu  struct {
		Cpu_tune struct {
			Vcpu_pins struct {
				Vcpu_pin []OLVMVcpu_pins
			}
		}
		Topology struct {
			Cores   string `json:"cores"`
			Sockets string `json:"sockets"`
			Threads string `json:"threads"`
		}
	}
	Host struct {
		Id string `json:"id"`
	}
	Capped bool
}

type OLVMVcpu_pins struct {
	Vcpu string `json:"vcpu"`
}

func (b *CommonBuilder) getClustersInfos() ([]model.ClusterInfo, error) {
	countHypervisors := len(b.configuration.Features.Virtualization.Hypervisors)

	clustersChan := make(chan []model.ClusterInfo, countHypervisors)
	vmsChan := make(chan map[string][]model.VMInfo, countHypervisors)
	errsChan := make(chan error)

	for i := range b.configuration.Features.Virtualization.Hypervisors {
		hv := b.configuration.Features.Virtualization.Hypervisors[i]

		if hv.Type == model.TechnologyOracleLVM {
			utils.RunRoutine(b.configuration, func() {
				clusters, err := getOlvmClustersData(hv.Endpoint, hv.Type, hv.Username, hv.Password)
				if err != nil {
					errsChan <- err
					clustersChan <- nil
					return
				}

				clustersChan <- clusters
			})

			utils.RunRoutine(b.configuration, func() {
				vms, err := getOlvmVMsData(hv.Endpoint, hv.Username, hv.Password)
				if err != nil {
					errsChan <- err
					vmsChan <- nil
					return
				}

				vmsChan <- vms
			})
		} else {
			utils.RunRoutine(b.configuration, func() {
				clusters, err := b.fetcher.GetClusters(hv)
				if err != nil {
					errsChan <- err
					clustersChan <- nil
					return
				}

				clustersChan <- clusters
			})

			utils.RunRoutine(b.configuration, func() {
				vms, err := b.fetcher.GetVirtualMachines(hv)
				if err != nil {
					errsChan <- err
					vmsChan <- nil
					return
				}

				vmsChan <- vms
			})
		}
	}

	clusters := make([]model.ClusterInfo, 0)

	for i := 0; i < countHypervisors; i++ {
		c := <-clustersChan
		if c == nil {
			continue
		}

		clusters = append(clusters, c...)
	}

	allVMs := make(map[string][]model.VMInfo)

	for i := 0; i < countHypervisors; i++ {
		vmsPerCluster := <-vmsChan
		if vmsPerCluster == nil {
			continue
		}

		for clusterName, vms := range vmsPerCluster {
			thisClusterVMs := allVMs[clusterName]
			thisClusterVMs = append(thisClusterVMs, vms...)

			allVMs[clusterName] = thisClusterVMs
		}
	}

	clusters = setVMsInClusterInfo(clusters, allVMs)

	var merr error
	for len(errsChan) > 0 {
		merr = multierror.Append(merr, <-errsChan)
	}

	return clusters, merr
}

func setVMsInClusterInfo(clusters []model.ClusterInfo, clusterMap map[string][]model.VMInfo) []model.ClusterInfo {
	for i := range clusters {
		if clusterMap[clusters[i].Name] != nil {
			clusters[i].VMs = clusterMap[clusters[i].Name]
		} else {
			clusters[i].VMs = []model.VMInfo{}
		}
	}

	return clusters
}

func getOlvmClustersData(endpoint string, vmType string, username string, password string) ([]model.ClusterInfo, error) {
	responseOLVMClusters, err := getResponseOLVMClusters(endpoint, username, password)
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	responseOLVMHosts, err := getResponseOLVMHosts(endpoint, username, password)
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	var cpu, socket int

	clusters := []model.ClusterInfo{}

	for _, vC := range responseOLVMClusters.Clusters {
		cpu, socket = 0, 0

		for _, vH := range responseOLVMHosts.Hosts {
			if vH.Cluster.Id == vC.Id {
				cores, err := strconv.Atoi(vH.Cpu.Topology.Cores)
				if err != nil {
					return nil, ercutils.NewError(err)
				}

				cpu += cores

				sockets, err := strconv.Atoi(vH.Cpu.Topology.Sockets)
				if err != nil {
					return nil, ercutils.NewError(err)
				}

				socket += sockets
			}
		}

		clusterInfo := model.ClusterInfo{
			Type:          vmType,
			FetchEndpoint: endpoint,
			Name:          vC.Name,
			CPU:           cpu,
			Sockets:       socket,
			VMs:           []model.VMInfo{},
		}
		clusters = append(clusters, clusterInfo)
	}

	return clusters, nil
}

func getOlvmVMsData(endpoint string, username string, password string) (map[string][]model.VMInfo, error) {
	responseClusters, err := getResponseOLVMClusters(endpoint, username, password)
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	responseHosts, err := getResponseOLVMHosts(endpoint, username, password)
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	responseVMs, err := getResponseOLVMVMs(endpoint, username, password)
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	vms := map[string][]model.VMInfo{}

	for _, vC := range responseClusters.Clusters {
		for _, vH := range responseHosts.Hosts {
			if vH.Cluster.Id == vC.Id {
				for _, vV := range responseVMs.VMs {
					if vV.Host.Id == vH.Id {
						if vV.Fqdn == "" {
							vV.Fqdn = vV.Name
						}

						vmCores, err := strconv.Atoi(vV.Cpu.Topology.Cores)
						if err != nil {
							return nil, ercutils.NewError(err)
						}

						vmSockets, err := strconv.Atoi(vV.Cpu.Topology.Sockets)
						if err != nil {
							return nil, ercutils.NewError(err)
						}

						vmThreads, err := strconv.Atoi(vV.Cpu.Topology.Threads)
						if err != nil {
							return nil, ercutils.NewError(err)
						}

						vmCappedCPU := vmCores * vmSockets * vmThreads

						vmVcpu := len(vV.Cpu.Cpu_tune.Vcpu_pins.Vcpu_pin)
						if vmVcpu != 0 && vmVcpu == vmCappedCPU {
							vV.Capped = true
						}

						vm := model.VMInfo{
							Name:               vV.Name,
							Hostname:           vV.Fqdn,
							CappedCPU:          vV.Capped,
							VirtualizationNode: vH.Name,
						}

						thisVMs := vms[vC.Name]
						thisVMs = append(thisVMs, vm)

						vms[vC.Name] = thisVMs
					}
				}
			}
		}
	}

	return vms, nil
}

func getResponseOLVMClusters(endpoint string, username string, password string) (ResponseOLVMClusters, error) {
	url := "https://" + endpoint + "/ovirt-engine/api/clusters"

	bodyClustersBytes, err := getBodyResponse(url, username, password)
	if err != nil {
		return ResponseOLVMClusters{}, ercutils.NewError(err)
	}

	var responseOLVMClusters ResponseOLVMClusters

	errClustersUnmarshal := json.Unmarshal(bodyClustersBytes, &responseOLVMClusters)
	if errClustersUnmarshal != nil {
		return ResponseOLVMClusters{},
			ercutils.NewError(errClustersUnmarshal,
				fmt.Sprintf("Can't unmarshal clusters, bodyClusterBytes:\n%s\n", string(bodyClustersBytes)))
	}

	return responseOLVMClusters, nil
}

func getResponseOLVMHosts(endpoint string, username string, password string) (ResponseOLVMHosts, error) {
	url := "https://" + endpoint + "/ovirt-engine/api/hosts"

	bodyHostsBytes, err := getBodyResponse(url, username, password)
	if err != nil {
		return ResponseOLVMHosts{}, ercutils.NewError(err)
	}

	var responseOLVMHosts ResponseOLVMHosts

	errHostUnmarshal := json.Unmarshal(bodyHostsBytes, &responseOLVMHosts)
	if errHostUnmarshal != nil {
		return ResponseOLVMHosts{},
			ercutils.NewError(errHostUnmarshal,
				fmt.Sprintf("Can't unmarshal hosts, bodyHostsBytes:\n%s\n", string(bodyHostsBytes)))
	}

	return responseOLVMHosts, nil
}

func getResponseOLVMVMs(endpoint string, username string, password string) (ResponseOLVMVMs, error) {
	url := "https://" + endpoint + "/ovirt-engine/api/vms"

	bodyVMsBytes, err := getBodyResponse(url, username, password)
	if err != nil {
		return ResponseOLVMVMs{}, ercutils.NewError(err)
	}

	var responseOLVMVMs ResponseOLVMVMs

	errVMUnmarshal := json.Unmarshal(bodyVMsBytes, &responseOLVMVMs)
	if errVMUnmarshal != nil {
		return ResponseOLVMVMs{},
			ercutils.NewError(errVMUnmarshal,
				fmt.Sprintf("Can't unmarshal vms, bodyHostsBytes:\n%s\n", string(bodyVMsBytes)))
	}

	return responseOLVMVMs, nil
}

func getBodyResponse(url string, username string, password string) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(username, password)

	resp, err := client.Do(req)
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	return bodyBytes, nil
}
