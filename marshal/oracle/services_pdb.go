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

func ServicesPdb(cmdOutput []byte) ([]model.OracleDatabasePdbService, error) {
	services := make([]model.OracleDatabasePdbService, 0)
	scanner := bufio.NewScanner(bytes.NewReader(cmdOutput))

	var merr, err error

	for scanner.Scan() {
		service := new(model.OracleDatabasePdbService)
		line := scanner.Text()

		splitted := strings.Split(line, "|||")
		if len(splitted) == 6 {
			if strings.TrimSpace(splitted[0]) != "" {
				service.Name = marshal.TrimParseStringPointer(strings.TrimSpace(splitted[0]))
			} else {
				service.Name = nil
			}

			if strings.TrimSpace(splitted[1]) != "" {
				service.FailoverMethod = marshal.TrimParseStringPointer(strings.TrimSpace(splitted[1]))
			} else {
				service.FailoverMethod = nil
			}

			if strings.TrimSpace(splitted[2]) != "" {
				service.FailoverType = marshal.TrimParseStringPointer(strings.TrimSpace(splitted[2]))
			} else {
				service.FailoverType = nil
			}

			if strings.TrimSpace(splitted[3]) != "" {
				if service.FailoverRetries, err = marshal.TrimParseIntPointer(splitted[3]); err != nil {
					merr = multierror.Append(merr, ercutils.NewError(err))
				}
			} else {
				service.FailoverRetries = nil
			}

			if strings.TrimSpace(splitted[4]) != "" {
				service.FailoverDelay = splitted[4]
			}

			if strings.TrimSpace(splitted[5]) != "" {
				service.Enabled = marshal.TrimParseBoolPointer(splitted[5])
			} else {
				service.Enabled = nil
			}

			services = append(services, *service)
		}
	}

	if merr != nil {
		return nil, merr
	}

	return services, nil
}
