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

package marshal

import (
	"testing"

	"github.com/ercole-io/ercole/model"
	"github.com/stretchr/testify/assert"
)

const testHostData string = `Hostname: littletony
CPUModel: AMD Opteron(tm) Processor 6380
CPUFrequency: 2500.000Mhz
CPUSockets: 2
CPUCores: 16
CPUThreads: 32
ThreadsPerCore: 2
CoresPerSocket: 8
HardwareAbstraction: PH
HardwareAbstractionTechnology: PH
Kernel: Linux
KernelVersion: 2.6.32-754.25.1.el6.x86_64
OS: Oracle Linux Server
OSVersion: 7.8
MemoryTotal: 189
SwapTotal: 31`

func TestHost(t *testing.T) {
	cmdOutput := []byte(testHostData)

	actual := Host(cmdOutput)

	expected := model.Host{
		Hostname:                      "littletony",
		CPUModel:                      "AMD Opteron(tm) Processor 6380",
		CPUFrequency:                  "2500.000Mhz",
		CPUSockets:                    2,
		CPUCores:                      16,
		CPUThreads:                    32,
		ThreadsPerCore:                2,
		CoresPerSocket:                8,
		HardwareAbstraction:           "PH",
		HardwareAbstractionTechnology: "PH",
		Kernel:                        "Linux",
		KernelVersion:                 "2.6.32-754.25.1.el6.x86_64",
		OS:                            "Oracle Linux Server",
		OSVersion:                     "7.8",
		MemoryTotal:                   189,
		SwapTotal:                     31,
		OtherInfo:                     nil,
	}

	assert.Equal(t, expected, actual)
}
