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

func Instance(cmdOutput []byte) (*model.MySQLInstance, error) {
	scanner := marshal.NewCsvScanner(cmdOutput, 16)

	var instance model.MySQLInstance

	var merr, err error

	for scanner.SafeScan() {
		instance.Name = strings.TrimSpace(scanner.Iter())
		instance.Version = strings.TrimSpace(scanner.Iter())
		instance.Edition = strings.TrimSpace(scanner.Iter())
		instance.Platform = strings.TrimSpace(scanner.Iter())
		instance.Architecture = strings.TrimSpace(scanner.Iter())
		instance.Engine = strings.TrimSpace(scanner.Iter())
		instance.RedoLogEnabled = strings.TrimSpace(scanner.Iter())
		instance.CharsetServer = strings.TrimSpace(scanner.Iter())
		instance.CharsetSystem = strings.TrimSpace(scanner.Iter())

		if instance.PageSize, err = marshal.TrimParseFloat64(scanner.Iter()); err != nil {
			merr = multierror.Append(merr, ercutils.NewError(err))
		}

		if instance.ThreadsConcurrency, err = marshal.TrimParseInt(scanner.Iter()); err != nil {
			merr = multierror.Append(merr, ercutils.NewError(err))
		}

		if instance.BufferPoolSize, err = marshal.TrimParseFloat64(scanner.Iter()); err != nil {
			merr = multierror.Append(merr, ercutils.NewError(err))
		}

		if instance.LogBufferSize, err = marshal.TrimParseFloat64(scanner.Iter()); err != nil {
			merr = multierror.Append(merr, ercutils.NewError(err))
		}

		if instance.SortBufferSize, err = marshal.TrimParseFloat64(scanner.Iter()); err != nil {
			merr = multierror.Append(merr, ercutils.NewError(err))
		}

		instance.ReadOnly = marshal.TrimParseBool(scanner.Iter())
		instance.LogBin = marshal.TrimParseBool(scanner.Iter())
	}

	if merr != nil {
		return nil, merr
	}

	return &instance, nil
}
