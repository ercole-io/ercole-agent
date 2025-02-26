// Copyright (c) 2025 Sorint.lab S.p.A.
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
package oracle

import (
	"testing"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/stretchr/testify/assert"
)

const testOracleDiskGroupData string = `
DATA                          |||         349452|||         124002|||          35.48
DATA                          |||      2048|||       346|||          16.91
RECO                          |||       512|||       442|||          86.24
`

func TestDiskGroups(t *testing.T) {
	testCases := []struct {
		name         string
		data         string
		expected     []model.OracleDatabaseDiskGroup
		checkRsponse func(t *testing.T,
			expected, actual []model.OracleDatabaseDiskGroup,
			err error)
	}{
		{
			name: "OK",
			data: testOracleDiskGroupData,
			expected: []model.OracleDatabaseDiskGroup{
				{
					DiskGroupName: "DATA",
					TotalSpace:    349452,
					FreeSpace:     124002,
				},
				{
					DiskGroupName: "DATA",
					TotalSpace:    2048,
					FreeSpace:     346,
				},
				{
					DiskGroupName: "RECO",
					TotalSpace:    512,
					FreeSpace:     442,
				},
			},
			checkRsponse: func(t *testing.T, expected, actual []model.OracleDatabaseDiskGroup, err error) {
				assert.Nil(t, err)
				assert.Equal(t, expected, actual)
			},
		},
		{
			name:     "Empty",
			data:     "",
			expected: []model.OracleDatabaseDiskGroup{},
			checkRsponse: func(t *testing.T, expected, actual []model.OracleDatabaseDiskGroup, err error) {
				assert.Nil(t, err)
				assert.Empty(t, actual)
			},
		},
		{
			name:     "Short output",
			data:     "DATA                          |||         349452|||         124002",
			expected: []model.OracleDatabaseDiskGroup{},
			checkRsponse: func(t *testing.T, expected, actual []model.OracleDatabaseDiskGroup, err error) {
				assert.Nil(t, err)
				assert.Empty(t, actual)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmdOutput := []byte(tc.data)
			actual, errs := DiskGroups(cmdOutput)
			tc.checkRsponse(t, tc.expected, actual, errs)
		})
	}
}
