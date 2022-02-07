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
	"strings"

	"github.com/ercole-io/ercole/v2/model"
)

// ClusterMembershipStatus returns this struct filled from the output of the script
func ClusterMembershipStatus(cmdOutput []byte) *model.ClusterMembershipStatus {
	data := parseKeyValueColonSeparated(cmdOutput)

	var clusterMembershipStatus model.ClusterMembershipStatus
	clusterMembershipStatus.OracleClusterware = TrimParseBool(data["OracleClusterware"])
	clusterMembershipStatus.SunCluster = TrimParseBool(data["SunCluster"])
	clusterMembershipStatus.HACMP = false

	clusterMembershipStatus.VeritasClusterServer = TrimParseBool(data["VeritasClusterServer"])

	hostnames := make([]string, 0)

	for _, s := range strings.Split(data["VeritasClusterHostnames"], ";") {
		fields := strings.Fields(s)
		if len(fields) != 2 {
			continue
		}

		hostnames = append(hostnames, fields[1])
	}

	if len(hostnames) > 0 {
		clusterMembershipStatus.VeritasClusterHostnames = hostnames
	}

	return &clusterMembershipStatus
}
