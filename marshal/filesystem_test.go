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

	"github.com/ercole-io/ercole/v2/model"
	"github.com/stretchr/testify/assert"
)

var testFilesystemsData string = ` /dev/nvme0n1p7 ext4       256981444 172784996  71072864      71% /
/dev/nvme0n1p1 vfat          661504     68640    592864      11% /boot/efi
`

func TestFilesystems(t *testing.T) {
	cmdOutput := []byte(testFilesystemsData)

	actual, err := Filesystems(cmdOutput)

	expected := []model.Filesystem{
		{
			Filesystem:     "/dev/nvme0n1p7",
			Type:           "ext4",
			Size:           256981444,
			UsedSpace:      172784996,
			AvailableSpace: 71072864,
			MountedOn:      "/",
		},
		{
			Filesystem:     "/dev/nvme0n1p1",
			Type:           "vfat",
			Size:           661504,
			UsedSpace:      68640,
			AvailableSpace: 592864,
			MountedOn:      "/boot/efi",
		},
	}

	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}
