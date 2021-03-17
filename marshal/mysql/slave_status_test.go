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

var testSlaveStatusData1 string = `mysql: [Warning] Using a password on the command line interface can be insecure.
"Waiting for master to send event";"10.100.132.119";"replication";"3306";"60";"binlog.000005";"772";"mysql1-relay-bin.000003";"937";"binlog.000005";"Yes";"Yes";"";"";"";"";"";"";"0";"";"0";"772";"1147";"None";"";"0";"Yes";"";"";"";"";"";"0";"No";"0";"";"0";"";"";"1";"a74c753e-7cfa-11eb-991a-566f86f30033";"mysql.slave_master_info";"0";"NULL";"Slave has read all relay log; waiting for more updates";"86400";"";"";"";"";"";"";"";"0";"";"";"";"";"0";""
`

var testSlaveStatusData2 string = `mysql: [Warning] Using a password on the command line interface can be insecure.
`

func TestSlaveStatus(t *testing.T) {
	masterUUID1 := "a74c753e-7cfa-11eb-991a-566f86f30033"
	testCases := []struct {
		data       string
		isSlave    bool
		masterUUID *string
	}{
		{testSlaveStatusData1, true, &masterUUID1},
		{testSlaveStatusData2, false, nil},
	}

	for _, tc := range testCases {
		cmdOutput := []byte(tc.data)

		isSlave, masterUUID := SlaveStatus(cmdOutput)

		assert.Equal(t, tc.isSlave, isSlave)
		assert.Equal(t, tc.masterUUID, masterUUID)
	}
}
