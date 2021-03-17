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

func Instance(cmdOutput []byte) (instance *model.MySQLInstance) {
	scanner := marshal.NewCsvScanner(cmdOutput, 16)

	for scanner.SafeScan() {
		instance := model.MySQLInstance{
			Name:               strings.TrimSpace(scanner.Iter()),
			Version:            strings.TrimSpace(scanner.Iter()),
			Edition:            strings.TrimSpace(scanner.Iter()),
			Platform:           strings.TrimSpace(scanner.Iter()),
			Architecture:       strings.TrimSpace(scanner.Iter()),
			Engine:             strings.TrimSpace(scanner.Iter()),
			RedoLogEnabled:     strings.TrimSpace(scanner.Iter()),
			CharsetServer:      strings.TrimSpace(scanner.Iter()),
			CharsetSystem:      strings.TrimSpace(scanner.Iter()),
			PageSize:           marshal.TrimParseFloat64(scanner.Iter()),
			ThreadsConcurrency: marshal.TrimParseInt(scanner.Iter()),
			BufferPoolSize:     marshal.TrimParseFloat64(scanner.Iter()),
			LogBufferSize:      marshal.TrimParseFloat64(scanner.Iter()),
			SortBufferSize:     marshal.TrimParseFloat64(scanner.Iter()),
			ReadOnly:           marshal.TrimParseBool(scanner.Iter()),
			LogBin:             marshal.TrimParseBool(scanner.Iter()),
		}

		return &instance
	}

	return nil
}
