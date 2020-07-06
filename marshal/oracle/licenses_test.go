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

const testLicensesData = `Oracle EXE;;
Oracle ENT; 2.00;
WebLogic Server Management Pack Enterprise Edition;;
Partitioning;;
not spooling currently`

func TestLicenses(t *testing.T) {
	expected := []model.OracleDatabaseLicense{
		{Name: "Oracle EXE", Count: 0},
		{Name: "Oracle ENT", Count: 2.00},
		{Name: "WebLogic Server Management Pack Enterprise Edition", Count: 0},
		{Name: "Partitioning", Count: 0},
	}

	actual := Licenses([]byte(testLicensesData))

	assert.Equal(t, expected, actual)
}
