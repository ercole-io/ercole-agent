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

func TableSchemas(cmdOutput []byte) ([]model.MySQLTableSchema, error) {
	tableSchemas := make([]model.MySQLTableSchema, 0)

	scanner := marshal.NewCsvScanner(cmdOutput, 3)

	var merr, err error

	for scanner.SafeScan() {
		var tableSchema model.MySQLTableSchema
		tableSchema.Name = strings.TrimSpace(scanner.Iter())
		tableSchema.Engine = strings.TrimSpace(scanner.Iter())

		if tableSchema.Allocation, err = marshal.TrimParseFloat64(scanner.Iter()); err != nil {
			merr = multierror.Append(merr, ercutils.NewError(err))
		}

		tableSchemas = append(tableSchemas, tableSchema)
	}

	if merr != nil {
		return nil, merr
	}

	return tableSchemas, nil
}
