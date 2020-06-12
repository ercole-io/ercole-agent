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

	"github.com/ercole-io/ercole-agent/model"
)

// ExadataDevices returns information about devices extracted from exadata-info command.
func ExadataDevices(cmdOutput []byte) []model.ExadataDevice {
	devices := []model.ExadataDevice{}
	scanner := bufio.NewScanner(strings.NewReader(string(cmdOutput)))

	for scanner.Scan() {
		device := new(model.ExadataDevice)
		line := scanner.Text()
		splitted := strings.Split(line, "|||")
		if len(splitted) == 17 {
			device.Hostname = strings.TrimSpace(splitted[0])
			device.ServerType = strings.TrimSpace(splitted[1])
			device.Model = strings.TrimSpace(splitted[2])
			device.ExaSwVersion = strings.TrimSpace(splitted[3])
			device.CPUEnabled = strings.TrimSpace(splitted[4])
			device.Memory = strings.TrimSpace(splitted[5])
			device.Status = strings.TrimSpace(splitted[6])
			device.PowerCount = strings.TrimSpace(splitted[7])
			device.PowerStatus = strings.TrimSpace(splitted[8])
			device.FanCount = strings.TrimSpace(splitted[9])
			device.FanStatus = strings.TrimSpace(splitted[10])
			device.TempActual = strings.TrimSpace(splitted[11])
			device.TempStatus = strings.TrimSpace(splitted[12])
			device.CellsrvService = strings.TrimSpace(splitted[13])
			device.MsService = strings.TrimSpace(splitted[14])
			device.RsService = strings.TrimSpace(splitted[15])
			device.FlashcacheMode = strings.TrimSpace(splitted[16])

			devices = append(devices, *device)
		}
	}
	return devices
}
