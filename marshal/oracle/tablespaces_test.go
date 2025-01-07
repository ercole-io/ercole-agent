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

const testTableSpacesBData string = `ercloradb9|||erc919|||erc919|||SYSTEM|||32767.98|||1140.00|||1129.63|||3.45|||ONLINE
ercloradb9|||erc919|||erc919|||USERS|||32767.98|||5.00|||	,63|||.45|||ONLINE
ercloradb9|||erc919|||erc919|||UNDOTBS1||||||  |||		|||.45|||ONLINE
`

func TestTablespaces(t *testing.T) {
	expected := []model.OracleDatabaseTablespace{
		{
			Name:     "SYSTEM",
			MaxSize:  32767.98,
			Total:    1140,
			Used:     1129.63,
			UsedPerc: 3.45,
			Status:   "ONLINE",
		},
		{
			Name:     "USERS",
			MaxSize:  32767.98,
			Total:    5,
			Used:     0.63,
			UsedPerc: 0.45,
			Status:   "ONLINE",
		},
		{
			Name:     "UNDOTBS1",
			MaxSize:  0,
			Total:    0,
			Used:     0,
			UsedPerc: 0.45,
			Status:   "ONLINE",
		},
	}

	actual, err := Tablespaces([]byte(testTableSpacesBData))

	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}
