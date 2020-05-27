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

package builder

//
//import (
//	"testing"
//
//	"github.com/ercole-io/ercole-agent/config"
//	"github.com/ercole-io/ercole-agent/model"
//	"github.com/stretchr/testify/assert"
//)
//
//func TestLinuxHostGetFetcherName(t *testing.T) {
//	var fetcher DataFetcher = &LinuxHostFetcher{}
//
//	assert.Equal(t, "linux_host", fetcher.GetFetcherName())
//}
//
//func TestLinuxHostMarshal(t *testing.T) {
//	var linuxFetcher *LinuxHostFetcher = &LinuxHostFetcher{}
//	linuxFetcher.fetcher = func() []byte {
//		return []byte(`hostname: amreo-ubuntu-pc
//			cpumodel: Intel(R) Core(TM) i5-8265U CPU @ 1.60GHz
//			cputhreads: 8
//			cpucores: 4
//			socket: 1
//			type: PH
//			virtual: N
//			kernel: 5.3.0-45-generic
//			os: Ubuntu 18.04.4 LTS
//			memorytotal: 15
//			swaptotal: 1
//			oraclecluster: N
//			veritascluster: N
//			suncluster: N
//aixcluster: N`)
//	}
//
//	conf := config.Configuration{
//		Envtype:  "PRD",
//		Location: "Italy",
//	}
//
//	hostdata := model.HostData{}
//	linuxFetcher.getLinuxHost(conf, &hostdata)
//
//	assert.Equal(t, model.Host{
//		Hostname:       "amreo-ubuntu-pc",
//		Environment:    "PRD",
//		Location:       "Italy",
//		CPUModel:       "Intel(R) Core(TM) i5-8265U CPU @ 1.60GHz",
//		CPUCores:       4,
//		CPUThreads:     8,
//		Socket:         1,
//		Type:           "PH",
//		Virtual:        false,
//		Kernel:         "5.3.0-45-generic",
//		OS:             "Ubuntu 18.04.4 LTS",
//		MemoryTotal:    15,
//		SwapTotal:      1,
//		OracleCluster:  false,
//		VeritasCluster: false,
//		SunCluster:     false,
//		AixCluster:     false,
//	}, hostdata.Hostname)
//}
//
