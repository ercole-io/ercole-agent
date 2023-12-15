// Copyright (c) 2023 Sorint.lab S.p.A.
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

const testServicePDBData string = `ORCLPDB|||||||||||||||NO
PLUTOSERVICE|||||||||||||||NO`

func TestServicesPdb(t *testing.T) {
	n1 := "ORCLPDB"
	e1 := false

	n2 := "PLUTOSERVICE"
	e2 := false

	expected := []model.OracleDatabasePdbService{
		{
			Name:    &n1,
			Enabled: &e1,
		},
		{
			Name:    &n2,
			Enabled: &e2,
		},
	}

	cmdOutput := []byte(testServicePDBData)
	actual, err := ServicesPdb(cmdOutput)

	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}
