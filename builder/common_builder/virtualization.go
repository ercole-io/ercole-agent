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

package common

import (
	"github.com/ercole-io/ercole-agent/model"
	"github.com/ercole-io/ercole-agent/utils"
)

func (b *CommonBuilder) getClustersInfos() []model.ClusterInfo {
	countHypervisors := len(b.configuration.Features.Virtualization.Hypervisors)

	clustersChan := make(chan []model.ClusterInfo, countHypervisors)
	vmsChan := make(chan []model.VMInfo, countHypervisors)

	for _, hv := range b.configuration.Features.Virtualization.Hypervisors {
		utils.RunRoutine(b.configuration, func() {
			clustersChan <- b.fetcher.GetClusters(hv)
		})

		utils.RunRoutine(b.configuration, func() {
			vmsChan <- b.fetcher.GetVirtualMachines(hv)
		})
	}

	clusters := make([]model.ClusterInfo, 0)
	for i := 0; i < countHypervisors; i++ {
		clusters = append(clusters, (<-clustersChan)...)
	}

	vms := make([]model.VMInfo, 0)
	for i := 0; i < countHypervisors; i++ {
		vms = append(vms, (<-vmsChan)...)
	}

	clusters = setVMsInClusterInfo(clusters, vms)

	return clusters
}

func setVMsInClusterInfo(clusters []model.ClusterInfo, vms []model.VMInfo) []model.ClusterInfo {
	clusters = append(clusters, model.ClusterInfo{
		Name:    "not_in_cluster",
		Type:    "unknown",
		CPU:     0,
		Sockets: 0,
		VMs:     []model.VMInfo{},
	})

	clusterMap := make(map[string][]model.VMInfo)

	for _, vm := range vms {
		if vm.ClusterName == "" {
			vm.ClusterName = "not_in_cluster"
		}
		clusterMap[vm.ClusterName] = append(clusterMap[vm.ClusterName], vm)
	}

	for i := range clusters {
		if clusterMap[clusters[i].Name] != nil {
			clusters[i].VMs = clusterMap[clusters[i].Name]
		} else {
			clusters[i].VMs = []model.VMInfo{}
		}
	}

	return clusters
}
