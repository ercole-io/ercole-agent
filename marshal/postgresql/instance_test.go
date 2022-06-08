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

func TestInstance(t *testing.T) {
	testInstanceData := `100|256746178|2|7|0|UTF8|f|f|f|f|0|0`
	cmdOutput := []byte(testInstanceData)
	actual, err := Instance(cmdOutput)

	expected := model.PostgreSQLInstance{
		MaxConnections:  100,
		InstanceSize:    256746178,
		UsersNum:        2,
		DbNum:           7,
		TblspNum:        0,
		Charset:         "UTF8",
		Isinreplica:     false,
		Ismaster:        false,
		Isslave:         false,
		ArchiverWorking: false,
		SlavesNum:       0,
		TrustHbaEntries: 0,
	}

	assert.Equal(t, &expected, actual)
	assert.Nil(t, err)
}
