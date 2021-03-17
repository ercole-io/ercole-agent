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

var testHighAvailabilityData1 string = `mysql: [Warning] Using a password on the command line interface can be insecure.
"mysql_innodb_cluster_metadata"
`

var testHighAvailabilityData2 string = `mysql: [Warning] Using a password on the command line interface can be insecure.
`

func TestHighAvailability(t *testing.T) {
	testCases := []struct {
		data     string
		expected bool
	}{
		{testHighAvailabilityData1, true},
		{testHighAvailabilityData2, false},
	}

	for _, tc := range testCases {
		cmdOutput := []byte(tc.data)

		actual := HighAvailability(cmdOutput)

		assert.Equal(t, tc.expected, actual)
	}
}
