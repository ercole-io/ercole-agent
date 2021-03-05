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

var testTableSchemasData = `mysql: [Warning] Using a password on the command line interface can be insecure.
"classicmodels";"InnoDB";"1.094"
"innodb_temporary";"InnoDB";"12.000"
"innodb_undo_001";"InnoDB";"16.000"
"innodb_undo_002";"InnoDB";"16.000"
"mysql";"InnoDB";"24.004"
"sys";"InnoDB";"0.078"
`

func TestTableSchemas(t *testing.T) {
	cmdOutput := []byte(testTableSchemasData)

	actual := TableSchemas(cmdOutput)

	expected := []model.MySQLTableSchema{
		{
			Name:       "classicmodels",
			Engine:     "InnoDB",
			Allocation: 1.094,
		},
		{
			Name:       "innodb_temporary",
			Engine:     "InnoDB",
			Allocation: 12.0,
		},
		{
			Name:       "innodb_undo_001",
			Engine:     "InnoDB",
			Allocation: 16.0,
		},
		{
			Name:       "innodb_undo_002",
			Engine:     "InnoDB",
			Allocation: 16.0,
		},
		{
			Name:       "mysql",
			Engine:     "InnoDB",
			Allocation: 24.004,
		},
		{
			Name:       "sys",
			Engine:     "InnoDB",
			Allocation: 0.078,
		},
	}

	assert.Equal(t, expected, actual)
}
