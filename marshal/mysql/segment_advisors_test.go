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

package mysql

import (
	"testing"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/stretchr/testify/assert"
)

var testSegmentAdvisorsData = `mysql: [Warning] Using a password on the command line interface can be insecure.
"erclinmysql:3306";"classicmodels";"customers";"InnoDB";"0.109";"0.078";"0.031";"0.000"
"erclinmysql:3306";"classicmodels";"employees";"InnoDB";"0.141";"0.094";"0.047";"0.000"
"erclinmysql:3306";"classicmodels";"offices";"InnoDB";"0.078";"0.062";"0.016";"0.000"
`

func TestSegmentAdvisors(t *testing.T) {
	cmdOutput := []byte(testSegmentAdvisorsData)

	actual := SegmentAdvisors(cmdOutput)

	expected := []model.MySQLSegmentAdvisor{
		{
			TableSchema: "classicmodels",
			TableName:   "customers",
			Engine:      "InnoDB",
			Allocation:  0.109,
			Data:        0.078,
			Index:       0.031,
			Free:        0,
		},
		{
			TableSchema: "classicmodels",
			TableName:   "employees",
			Engine:      "InnoDB",
			Allocation:  0.141,
			Data:        0.094,
			Index:       0.047,
			Free:        0,
		},
		{
			TableSchema: "classicmodels",
			TableName:   "offices",
			Engine:      "InnoDB",
			Allocation:  0.078,
			Data:        0.062,
			Index:       0.016,
			Free:        0,
		},
	}

	assert.Equal(t, expected, actual)
}
