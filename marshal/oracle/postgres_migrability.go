// Copyright (c) 2023 Sorint.lab S.p.A.
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

// PostgresMigrability returns information about database migrability to pgsql
func PostgresMigrability(cmdOutput []byte) ([]model.PgsqlMigrability, error) {
	res := make([]model.PgsqlMigrability, 0)
	scanner := bufio.NewScanner(bytes.NewReader(cmdOutput))

	var merr error

	i := 0

	for scanner.Scan() {
		line := scanner.Text()

		splitted := strings.Split(line, "|||")

		if i < 10 {
			count, err := marshal.TrimParseInt(splitted[1])

			if err != nil {
				merr = multierror.Append(merr, ercutils.NewError(err))
			}

			m := strings.TrimSpace(splitted[0])

			res = append(res, model.PgsqlMigrability{
				Metric: &m,
				Count:  count,
			})
		}

		if i > 9 {
			count, err := marshal.TrimParseInt(splitted[2])

			if err != nil {
				merr = multierror.Append(merr, ercutils.NewError(err))
			}

			s := strings.TrimSpace(splitted[0])
			o := strings.TrimSpace(splitted[1])

			res = append(res, model.PgsqlMigrability{
				Schema:     &s,
				ObjectType: &o,
				Count:      count,
			})
		}

		i++
	}

	if merr != nil {
		return nil, merr
	}

	return res, nil
}
