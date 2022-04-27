// Copyright (c) 2022 Sorint.lab S.p.A.
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

const testListPDB string = `TAGGR01 ||| READ WRITE
TDPRO01     |||          READ WRITE`
const testListPDB1 string = `TAGGR01 ||| READ WRITE  ||| TDPRO01|||          READ WRITE`
const testListPDB2 string = `TAGGR01 ||| READ WRITE

TDPRO01     |||          READ WRITE`

func TestListPDB(t *testing.T) {
	expected := []model.OracleDatabasePluggableDatabase{
		{
			Name:        "TAGGR01",
			Status:      "READ WRITE",
			Tablespaces: nil,
			Schemas:     nil,
			Services:    []model.OracleDatabaseService{},
		},
		{
			Name:        "TDPRO01",
			Status:      "READ WRITE",
			Tablespaces: nil,
			Schemas:     nil,
			Services:    []model.OracleDatabaseService{},
		},
	}

	actual, err := ListPDB([]byte(testListPDB))

	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestListPDB_Error1(t *testing.T) {
	actual, err := ListPDB([]byte(testListPDB1))

	assert.Nil(t, actual)
	assert.Error(t, err)
}

func TestListPDB_Error2(t *testing.T) {
	actual, err := ListPDB([]byte(testListPDB2))

	assert.Nil(t, actual)
	assert.Error(t, err)
}
