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

package marshal

import (
	"testing"

	"github.com/ercole-io/ercole/model"
	"github.com/stretchr/testify/assert"
)

var testClusterMembershipStatusOutput string = `OracleClusterware: Y
VeritasClusterServer: N
SunCluster: Y`

func TestClusterMembershipStatusOutput(t *testing.T) {
	cmdOutput := []byte(testClusterMembershipStatusOutput)

	actual := ClusterMembershipStatus(cmdOutput)

	expected := model.ClusterMembershipStatus{
		OracleClusterware:    true,
		VeritasClusterServer: false,
		SunCluster:           true,
		HACMP:                false,
	}

	assert.Equal(t, expected, actual)
}

func TestClusterMembershipStatusOutputShouldCrash(t *testing.T) {
	cmdOutput := []byte("pippo")

	assert.Panics(t, func() { ClusterMembershipStatus(cmdOutput) })
}
