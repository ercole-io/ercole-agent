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

package oracle

import (
	"testing"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/stretchr/testify/assert"
)

func TestPatches(t *testing.T) {
	testPatchesData := `erclin8dbx							|||ND|||ERC819							|||ERC819			 |||19.0.0.0.0	     |||  30869156|||APPLY	    |||Database Release Update : 19.7.0.0.200414 (30869156)						   |||2020-05-28`

	cmdOutput := []byte(testPatchesData)
	actual, err := Patches(cmdOutput)

	expected := []model.OracleDatabasePatch{
		{
			Version:     "19.0.0.0.0",
			PatchID:     30869156,
			Action:      "APPLY",
			Description: "Database Release Update : 19.7.0.0.200414 (30869156)",
			Date:        "2020-05-28",
			OtherInfo:   nil,
		},
	}

	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestPatches_WrongDates(t *testing.T) {
	testPatchesData := `erclin8dbx				         	|||ND|||ERC819							|||ERC819			 |||19.0.0.0.0	     |||  30869156|||APPLY	    |||Database Release Update : 19.7.0.0.200414 (30869156)						   |||2020-05-28
	erclin8dbx							|||ND|||ERC819							|||ERC819			 |||19.0.0.0.0	     |||  30869156|||APPLY	    |||Database Release Update : 19.7.0.0.200414 (30869156)						   |||2020-28
	erclin8dbx							|||ND|||ERC819							|||ERC819			 |||19.0.0.0.0	     |||  30869156|||APPLY	    |||Database Release Update : 19.7.0.0.200414 (30869156)						   |||2020-05-28	`

	cmdOutput := []byte(testPatchesData)
	actual, err := Patches(cmdOutput)

	expected := []model.OracleDatabasePatch{
		{
			Version:     "19.0.0.0.0",
			PatchID:     30869156,
			Action:      "APPLY",
			Description: "Database Release Update : 19.7.0.0.200414 (30869156)",
			Date:        "2020-05-28",
			OtherInfo:   nil,
		},
		{
			Version:     "19.0.0.0.0",
			PatchID:     30869156,
			Action:      "APPLY",
			Description: "Database Release Update : 19.7.0.0.200414 (30869156)",
			Date:        "2020-05-28",
			OtherInfo:   nil,
		},
	}

	assert.Equal(t, expected, actual)
	assert.NotNil(t, err)
}
