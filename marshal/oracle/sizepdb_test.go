// Copyright (c) 2024 Sorint.lab S.p.A.
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

const testSizePDBData string = `139||| 175||| 339||| 0||| 7.625`

func TestSizePDB(t *testing.T) {
	expected := model.OracleDatabasePdbSize{
		SegmentsSize:       139,
		DatafileSize:       175,
		Allocable:          339,
		SGATarget:          0,
		PGAAggregateTarget: 7.625,
	}

	actual, err := SizePDB([]byte(testSizePDBData))

	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}
