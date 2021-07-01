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
	ercutils "github.com/ercole-io/ercole/v2/utils"
	"github.com/hashicorp/go-multierror"
)

func SegmentAdvisors(cmdOutput []byte) ([]model.MySQLSegmentAdvisor, error) {
	segmentAdvs := make([]model.MySQLSegmentAdvisor, 0)

	scanner := marshal.NewCsvScanner(cmdOutput, 7)

	var merr, err error

	for scanner.SafeScan() {
		var segmentAdv model.MySQLSegmentAdvisor
		segmentAdv.TableSchema = strings.TrimSpace(scanner.Iter())
		segmentAdv.TableName = strings.TrimSpace(scanner.Iter())
		segmentAdv.Engine = strings.TrimSpace(scanner.Iter())
		if segmentAdv.Allocation, err = marshal.TrimParseFloat64(scanner.Iter()); err != nil {
			merr = multierror.Append(merr, ercutils.NewError(err))
		}
		if segmentAdv.Data, err = marshal.TrimParseFloat64(scanner.Iter()); err != nil {
			merr = multierror.Append(merr, ercutils.NewError(err))
		}
		if segmentAdv.Index, err = marshal.TrimParseFloat64(scanner.Iter()); err != nil {
			merr = multierror.Append(merr, ercutils.NewError(err))
		}
		if segmentAdv.Free, err = marshal.TrimParseFloat64(scanner.Iter()); err != nil {
			merr = multierror.Append(merr, ercutils.NewError(err))
		}

		segmentAdvs = append(segmentAdvs, segmentAdv)
	}

	if merr != nil {
		return nil, merr
	}
	return segmentAdvs, nil
}
