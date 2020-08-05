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

var testClustersData string = `Name,NumCPU,NumSockets
pippo,108,6
pluto,144,8
topolino,192
`

func TestClusters(t *testing.T) {
	cmdOutput := []byte(testClustersData)

	actual := Clusters(cmdOutput)

	expected := []model.ClusterInfo{
		{
			Name:    "pippo",
			CPU:     108,
			Sockets: 6,
			VMs:     []model.VMInfo{},
		},
		{
			Name:    "pluto",
			CPU:     144,
			Sockets: 8,
			VMs:     []model.VMInfo{},
		},
		{
			Name:    "topolino",
			CPU:     192,
			Sockets: 1,
			VMs:     []model.VMInfo{},
		},
	}
	assert.Equal(t, expected, actual)
}
