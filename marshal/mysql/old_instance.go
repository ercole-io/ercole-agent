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

package mysql

import (
	"strings"

	"github.com/ercole-io/ercole-agent/v2/marshal"
	"github.com/ercole-io/ercole/v2/model"
	ercutils "github.com/ercole-io/ercole/v2/utils"
	"github.com/hashicorp/go-multierror"
)

func Old_Instance(cmdOutput []byte) (*model.MySQLInstance, error) {
	scanner := marshal.NewCsvScanner(cmdOutput, 2)

	var instance model.MySQLInstance

	var merr, err error

	var hostname, port string

	for scanner.SafeScan() {
		field := scanner.Get(0)

		switch field {
		case "hostname":
			hostname = scanner.Get(1)
		case "port":
			port = scanner.Get(1)
		case "version":
			instance.Version = scanner.Get(1)
		case "version_comment":
			edition := scanner.Get(1)
			if strings.Contains(edition, "Community") {
				instance.Edition = "COMMUNITY"
			} else {
				instance.Edition = "ENTERPRISE"
			}

		case "version_compile_os":
			instance.Platform = scanner.Get(1)
		case "version_compile_machine":
			instance.Architecture = scanner.Get(1)
		case "default_storage_engine":
			instance.Engine = scanner.Get(1)
		case "storage_engine":
			instance.Engine = scanner.Get(1)
		case "character_set_server":
			instance.CharsetServer = scanner.Get(1)
		case "character_set_system":
			instance.CharsetSystem = scanner.Get(1)
		case "innodb_page_size":
			var pagesize float64

			if pagesize, err = marshal.TrimParseFloat64(scanner.Get(1)); err != nil {
				merr = multierror.Append(merr, ercutils.NewError(err))
			}

			instance.PageSize = pagesize / 1024

		case "innodb_thread_concurrency":
			if instance.ThreadsConcurrency, err = marshal.TrimParseInt(scanner.Get(1)); err != nil {
				merr = multierror.Append(merr, ercutils.NewError(err))
			}

		case "innodb_buffer_pool_size":
			var bufferPoolSize float64

			if bufferPoolSize, err = marshal.TrimParseFloat64(scanner.Get(1)); err != nil {
				merr = multierror.Append(merr, ercutils.NewError(err))
			}

			instance.BufferPoolSize = bufferPoolSize / 1024 / 1024

		case "innodb_log_buffer_size":
			var logBufferSize float64

			if logBufferSize, err = marshal.TrimParseFloat64(scanner.Get(1)); err != nil {
				merr = multierror.Append(merr, ercutils.NewError(err))
			}

			instance.LogBufferSize = logBufferSize / 1024 / 1024

		case "innodb_sort_buffer_size":
			var sortBufferSize float64

			if sortBufferSize, err = marshal.TrimParseFloat64(scanner.Get(1)); err != nil {
				merr = multierror.Append(merr, ercutils.NewError(err))
			}

			instance.SortBufferSize = sortBufferSize / 1024 / 1024

		case "read_only":
			instance.ReadOnly = marshal.TrimParseBool(scanner.Get(1))
		case "log_bin":
			instance.LogBin = marshal.TrimParseBool(scanner.Get(1))
		}
	}

	instance.Name = hostname + ":" + port
	instance.RedoLogEnabled = "test"

	if merr != nil {
		return nil, merr
	}

	return &instance, nil
}
