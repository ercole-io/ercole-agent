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
	"bufio"
	"bytes"
	"strings"

	"github.com/ercole-io/ercole/v2/model"
)

// Clusters returns a list of Clusters entries extracted
// from the clusters fetcher command output.
func Clusters(cmdOutput []byte) []model.ClusterInfo {
	scanner := bufio.NewScanner(bytes.NewReader(cmdOutput))
	clusters := []model.ClusterInfo{}
	for scanner.Scan() {
		line := scanner.Text()
		splitted := strings.Split(line, ",")

		//Check if the line is not the header line
		if len(splitted) == 3 && splitted[0] == "Name" && splitted[1] == "NumCPU" && splitted[2] == "NumSockets" {
			continue
		}

		clusterInfo := model.ClusterInfo{
			Name: strings.TrimSpace(splitted[0]),
			CPU:  parseInt(splitted[1]),
			VMs:  []model.VMInfo{},
		}

		if len(splitted) >= 3 {
			clusterInfo.Sockets = parseInt(splitted[2])
		} else {
			clusterInfo.Sockets = 1
		}

		clusters = append(clusters, clusterInfo)
	}

	return clusters
}
