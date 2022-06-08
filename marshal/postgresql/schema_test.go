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

package postgresql

import (
	"testing"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/stretchr/testify/assert"
)

func TestSchema(t *testing.T) {
	testSchemaData := `public|postgres|35|9519104|0|0|0|9519104|0|0`
	cmdOutput := []byte(testSchemaData)
	actual, err := Schema(cmdOutput)

	expected := model.PostgreSQLSchema{
		SchemaName:    "public",
		SchemaOwner:   "postgres",
		TablesCount:   35,
		TablesSize:    9519104,
		IndexesCount:  0,
		IndexesSize:   0,
		ViewsCount:    0,
		SchemaSize:    9519104,
		MatviewsCount: 0,
		MatviewsSize:  0,
	}

	assert.Equal(t, &expected, actual)
	assert.Nil(t, err)
}
