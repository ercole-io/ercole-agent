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

// ExadataCellDisks returns information about the cell disks extracted from exadata-storage-status command.
func ExadataCellDisks(cmdOutput []byte) []model.ExadataCellDisk {
	cellDisks := []model.ExadataCellDisk{}
	scanner := bufio.NewScanner(strings.NewReader(string(cmdOutput)))

	for scanner.Scan() {
		cellDisk := new(model.ExadataCellDisk)
		line := scanner.Text()
		splitted := strings.Split(line, "|||")
		if len(splitted) == 5 {
			cellDisk.StorageServerName = strings.TrimSpace(splitted[0])
			cellDisk.Name = strings.TrimSpace(splitted[1])
			cellDisk.Status = strings.TrimSpace(splitted[2])
			cellDisk.ErrCount = strings.TrimSpace(splitted[3])
			cellDisk.UsedPerc = strings.TrimSpace(splitted[4])

			cellDisks = append(cellDisks, *cellDisk)
		}
	}
	return cellDisks
}
