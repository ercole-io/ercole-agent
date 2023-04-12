// Copyright (c) 2022 Sorint.lab S.p.A.
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
	"bytes"
	"strings"

	"github.com/ercole-io/ercole-agent/v2/marshal"
	"github.com/ercole-io/ercole/v2/model"
	ercutils "github.com/ercole-io/ercole/v2/utils"
	"github.com/hashicorp/go-multierror"
)

// Partitionings returns information about database partitionings extracted
// from the partitionings fetcher command output.
func Partitionings(cmdOutput []byte) ([]model.OracleDatabasePartitioning, error) {
	partitionings := []model.OracleDatabasePartitioning{}
	scanner := bufio.NewScanner(bytes.NewReader(cmdOutput))

	var merr, err error

	for scanner.Scan() {
		line := scanner.Text()

		splitted := strings.Split(line, "|||")
		if len(splitted) == 4 {
			partitioning := model.OracleDatabasePartitioning{}

			partitioning.Owner = strings.TrimSpace(splitted[0])

			partitioning.SegmentName = strings.TrimSpace(splitted[1])

			if partitioning.Count, err = marshal.TrimParseInt(splitted[2]); err != nil {
				merr = multierror.Append(merr, ercutils.NewError(err))
			}

			if partitioning.Mb, err = marshal.TrimParseUnsafeFloat64(splitted[3], marshal.TrimParseFloat64); err != nil {
				merr = multierror.Append(merr, ercutils.NewError(err))
			}

			partitionings = append(partitionings, partitioning)
		}
	}

	if merr != nil {
		return nil, merr
	}

	return partitionings, nil
}
