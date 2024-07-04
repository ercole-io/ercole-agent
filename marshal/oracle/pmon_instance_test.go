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

	"github.com/stretchr/testify/assert"
)

const testPmonInstancesData = `oracle      5118       1  0 06:02 ?        00:00:00 ora_pmon_erc919`

func TestPmonInstances(t *testing.T) {
	expected := map[string]string{
		"5118": "ora_pmon_erc919",
	}

	actual := PmonInstances([]byte(testPmonInstancesData))

	assert.Equal(t, expected, actual)
}
