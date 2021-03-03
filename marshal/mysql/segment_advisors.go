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

package mysql

import (
	"strings"

	"github.com/ercole-io/ercole-agent/v2/marshal"
	"github.com/ercole-io/ercole/v2/model"
)

func SegmentAdvisors(cmdOutput []byte) []model.MySQLSegmentAdvisor {
	segmentAdvs := make([]model.MySQLSegmentAdvisor, 0)

	scanner := marshal.NewCsvScanner(cmdOutput, 8)

	for scanner.SafeScan() {
		_ = scanner.Iter() // throw away instance name

		segmentAdv := model.MySQLSegmentAdvisor{
			TableSchema: strings.TrimSpace(scanner.Iter()),
			TableName:   strings.TrimSpace(scanner.Iter()),
			Engine:      strings.TrimSpace(scanner.Iter()),
			Allocation:  marshal.TrimParseFloat64(scanner.Iter()),
			Data:        marshal.TrimParseFloat64(scanner.Iter()),
			Index:       marshal.TrimParseFloat64(scanner.Iter()),
			Free:        marshal.TrimParseFloat64(scanner.Iter()),
		}

		segmentAdvs = append(segmentAdvs, segmentAdv)
	}

	return segmentAdvs
}
