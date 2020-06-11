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
)

func (b *CommonBuilder) getExadataDevices() []model.ExadataDevice {
	exadataDevices := b.fetcher.GetExadataDevices()
	exadataCellDisks := b.fetcher.GetExadataCellDisks()

	//Join exadataDevices with exadataCellDisks
	for _, cd := range exadataCellDisks {
		for i := range exadataDevices {
			if cd.StorageServerName == exadataDevices[i].Hostname {
				if exadataDevices[i].CellDisks == nil {
					exadataDevices[i].CellDisks = []model.ExadataCellDisk{}
				}
				exadataDevices[i].CellDisks = append(exadataDevices[i].CellDisks, cd)
			}
		}
	}

	return exadataDevices
}
