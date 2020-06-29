// Copyright (c) 2019 Sorint.lab S.p.A.
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
	"bufio"
	"strings"

	"github.com/ercole-io/ercole/model"
)

// ExadataComponent returns information about devices extracted from exadata-info command.
func ExadataComponent(cmdOutput []byte) []model.OracleExadataComponent {
	devices := []model.OracleExadataComponent{}
	scanner := bufio.NewScanner(strings.NewReader(string(cmdOutput)))

	for scanner.Scan() {
		device := new(model.OracleExadataComponent)
		line := scanner.Text()
		splitted := strings.Split(line, "|||")
		if len(splitted) == 17 {
			device.Hostname = strings.TrimSpace(splitted[0])
			device.ServerType = strings.TrimSpace(splitted[1])
			device.Model = strings.TrimSpace(splitted[2])
			device.SwVersion = strings.TrimSpace(splitted[3])
			tmp := strings.Split(device.SwVersion, ".")
			device.SwReleaseDate = tmp[len(tmp)-1]

			cpuEnabled := strings.Split(splitted[4], "/")
			if len(cpuEnabled) == 2 {
				device.RunningCPUCount = trimParseIntPointer(cpuEnabled[0], "-")
				device.TotalCPUCount = trimParseIntPointer(cpuEnabled[1], "-")
			}

			device.Memory = trimParseIntPointer(splitted[5], "-")
			device.Status = trimParseStringPointer(splitted[6], "-")

			powerCount := strings.Split(splitted[7], "/")
			if len(powerCount) == 2 {
				device.RunningPowerSupply = trimParseIntPointer(powerCount[0], "-")
				device.TotalPowerSupply = trimParseIntPointer(powerCount[1], "-")
			}

			device.PowerStatus = trimParseStringPointer(splitted[8], "-")

			fanCount := strings.Split(splitted[9], "/")
			if len(fanCount) == 2 {
				device.RunningFanCount = trimParseIntPointer(fanCount[0], "-")
				device.TotalFanCount = trimParseIntPointer(fanCount[1], "-")
			}

			device.FanStatus = trimParseStringPointer(splitted[10], "-")

			device.TempActual = trimParseFloat64Pointer(splitted[11], "-")
			device.TempStatus = trimParseStringPointer(splitted[12], "-")
			device.CellsrvServiceStatus = trimParseStringPointer(splitted[13], "-")
			device.MsServiceStatus = trimParseStringPointer(splitted[14], "-")
			device.RsServiceStatus = trimParseStringPointer(splitted[15], "-")
			device.FlashcacheMode = trimParseStringPointer(splitted[16], "-")

			devices = append(devices, *device)
		}
	}
	return devices
}
