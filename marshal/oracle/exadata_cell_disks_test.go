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

	"github.com/ercole-io/ercole-agent/v2/agentmodel"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/stretchr/testify/assert"
)

const testOracleExadataCellDisksData string = `fcax1sf1|||CD_00_fcax1sf1|||normal|||0|||54
fcax1sf1|||CD_01_fcax1sf1|||normal|||2|||42
fcax1sf2|||CD_02_fcax1sf3|||normal|||103|||54`

func TestOracleExadataCellDisks(t *testing.T) {
	cmdOutput := []byte(testOracleExadataCellDisksData)

	actual, errs := ExadataCellDisks(cmdOutput)
	assert.Nil(t, errs)

	expected := map[agentmodel.StorageServerName][]model.OracleExadataCellDisk{
		agentmodel.StorageServerName("fcax1sf1"): {
			{
				ErrCount: 0,
				Name:     "CD_00_fcax1sf1",
				Status:   "normal",
				UsedPerc: 54,
			},
			{
				ErrCount: 2,
				Name:     "CD_01_fcax1sf1",
				Status:   "normal",
				UsedPerc: 42,
			},
		},
		agentmodel.StorageServerName("fcax1sf2"): {
			{
				ErrCount: 103,
				Name:     "CD_02_fcax1sf3",
				Status:   "normal",
				UsedPerc: 54,
			},
		},
	}

	assert.Equal(t, expected, actual)
}

func TestEmptyOracleExadataCellDisks(t *testing.T) {
	cmdOutput := []byte("")

	actual, errs := ExadataCellDisks(cmdOutput)
	assert.Nil(t, errs)

	expected := make(map[agentmodel.StorageServerName][]model.OracleExadataCellDisk)

	assert.Equal(t, expected, actual)
}
