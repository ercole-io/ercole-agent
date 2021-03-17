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

var testUUIDData1 string = `server-uuid=a74c753e-7cfa-11eb-991a-566f86f30033`

var testUUIDData2 string = ``

func TestUUID(t *testing.T) {
	testCases := []struct {
		data string
		uuid string
	}{
		{testUUIDData1, "a74c753e-7cfa-11eb-991a-566f86f30033"},
		{testUUIDData2, ""},
	}

	for _, tc := range testCases {
		cmdOutput := []byte(tc.data)

		uuid := UUID(cmdOutput)

		assert.Equal(t, tc.uuid, uuid)
	}
}
