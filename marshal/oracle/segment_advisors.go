// Copyright (c) 2019 Sorint.lab S.p.A.
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
	"bufio"
	"strings"

	"github.com/ercole-io/ercole-agent/v2/marshal"
	"github.com/ercole-io/ercole/v2/model"
)

// SegmentAdvisor returns informations about SegmentAdvisor parsed from fetcher command output.
func SegmentAdvisor(cmdOutput []byte) []model.OracleDatabaseSegmentAdvisor {
	segmentadvisors := []model.OracleDatabaseSegmentAdvisor{}
	scanner := bufio.NewScanner(strings.NewReader(string(cmdOutput)))

	for scanner.Scan() {
		line := scanner.Text()

		splitted := strings.Split(line, "|||")

		if len(splitted) == 8 {
			segmentadvisor := new(model.OracleDatabaseSegmentAdvisor)

			segmentadvisor.SegmentOwner = strings.TrimSpace(splitted[2])
			segmentadvisor.SegmentName = strings.TrimSpace(splitted[3])
			segmentadvisor.SegmentType = strings.TrimSpace(splitted[4])
			segmentadvisor.PartitionName = strings.TrimSpace(splitted[5])
			segmentadvisor.Reclaimable = marshal.TrimParseFloat64(splitted[6])
			segmentadvisor.Recommendation = strings.TrimSpace(splitted[7])
			segmentadvisors = append(segmentadvisors, *segmentadvisor)
		}
	}

	return segmentadvisors
}
