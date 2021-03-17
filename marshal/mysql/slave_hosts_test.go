// Copyright (c) 2021 Sorint.lab S.p.A.
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

	"github.com/stretchr/testify/assert"
)

var testSlaveHostsData1 string = `mysql: [Warning] Using a password on the command line interface can be insecure.
"2";"";"3306";"1";"4a4e892d-80f3-11eb-9b32-566f86f3003d"
`

var testSlaveHostsData2 string = `mysql: [Warning] Using a password on the command line interface can be insecure.
`

func TestSlaveHosts(t *testing.T) {
	testCases := []struct {
		data       string
		isMaster   bool
		slaveUUIDs []string
	}{
		{testSlaveHostsData1, true, []string{"4a4e892d-80f3-11eb-9b32-566f86f3003d"}},
		{testSlaveHostsData2, false, nil},
	}

	for _, tc := range testCases {
		cmdOutput := []byte(tc.data)

		isMaster, slaveUUIDs := SlaveHosts(cmdOutput)

		assert.Equal(t, tc.isMaster, isMaster)
		assert.ElementsMatch(t, tc.slaveUUIDs, slaveUUIDs)
	}
}
