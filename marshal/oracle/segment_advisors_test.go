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

	"github.com/ercole-io/ercole/model"
	"github.com/stretchr/testify/assert"
)

const testSegmentAdvisorData string = `erclin7dbx |||ERC718 |||ERCOLE |||ERC_TABLE_2 |||TABLE ||| ||| .36|||Enable row movement of the table ERCOLE.ERC_TABLE_2 and perform shrink, estimated savings is 389430466 bytes.`

func TestSegmentAdvisor(t *testing.T) {
	expected := []model.OracleDatabaseSegmentAdvisor{
		{
			SegmentOwner:   "ERCOLE",
			SegmentName:    "ERC_TABLE_2",
			SegmentType:    "TABLE",
			PartitionName:  "",
			Reclaimable:    0.36,
			Recommendation: "Enable row movement of the table ERCOLE.ERC_TABLE_2 and perform shrink, estimated savings is 389430466 bytes.",
			OtherInfo:      map[string]interface{}(nil),
		},
	}

	cmdOutput := []byte(testSegmentAdvisorData)
	actual := SegmentAdvisor(cmdOutput)

	assert.Equal(t, expected, actual)
}
