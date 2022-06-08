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

package postgresql

import (
	"bufio"
	"bytes"
	"strconv"
	"strings"

	"github.com/ercole-io/ercole-agent/v2/marshal"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/hashicorp/go-multierror"
)

func Schema(cmdOutput []byte) (*model.PostgreSQLSchema, error) {
	scanner := bufio.NewScanner(bytes.NewReader(cmdOutput))
	result := model.PostgreSQLSchema{}

	var merr, err error

	for scanner.Scan() {
		line := scanner.Text()

		splitted := strings.Split(line, "|")
		if len(splitted) < 7 {
			continue
		}

		iter := marshal.NewIter(splitted)

		result.SchemaName = iter()
		result.SchemaOwner = iter()

		if result.TablesCount, err = strconv.Atoi(iter()); err != nil {
			merr = multierror.Append(merr, err)
		}

		if result.TablesSize, err = strconv.Atoi(iter()); err != nil {
			merr = multierror.Append(merr, err)
		}

		if result.IndexesCount, err = strconv.Atoi(iter()); err != nil {
			merr = multierror.Append(merr, err)
		}

		if result.IndexesSize, err = strconv.Atoi(iter()); err != nil {
			merr = multierror.Append(merr, err)
		}

		if result.ViewsCount, err = strconv.Atoi(iter()); err != nil {
			merr = multierror.Append(merr, err)
		}

		if len(splitted) > 7 {
			if result.SchemaSize, err = strconv.Atoi(iter()); err != nil {
				merr = multierror.Append(merr, err)
			}

			if result.MatviewsCount, err = strconv.Atoi(iter()); err != nil {
				merr = multierror.Append(merr, err)
			}

			if result.MatviewsSize, err = strconv.Atoi(iter()); err != nil {
				merr = multierror.Append(merr, err)
			}
		}
	}

	if merr != nil {
		return nil, merr
	}

	return &result, nil
}
