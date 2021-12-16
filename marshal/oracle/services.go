// Copyright (c) 2021 Sorint.lab S.p.A.
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

// Services returns information about database services extracted
// from the services fetcher command output.

func Services(cmdOutput []byte) ([]model.OracleDatabaseService, error) {
	services := []model.OracleDatabaseService{}
	scanner := bufio.NewScanner(bytes.NewReader(cmdOutput))
	var merr, err error

	for scanner.Scan() {
		service := new(model.OracleDatabaseService)
		line := scanner.Text()
		splitted := strings.Split(line, "|||")
		if len(splitted) == 7 {
			if strings.TrimSpace(splitted[0]) != "" {
				service.Name = marshal.TrimParseStringPointer(strings.TrimSpace(splitted[0]))
			} else {
				service.Name = nil
			}
			if strings.TrimSpace(splitted[1]) != "" {
				service.CreationDate, err = marshal.TrimParseDatePointer(splitted[1])
				if err != nil {
					merr = multierror.Append(merr, ercutils.NewError(err))
				}
			} else {
				service.CreationDate = nil
			}
			if strings.TrimSpace(splitted[2]) != "" {
				service.FailoverMethod = marshal.TrimParseStringPointer(strings.TrimSpace(splitted[2]))
			} else {
				service.FailoverMethod = nil
			}
			if strings.TrimSpace(splitted[3]) != "" {
				service.FailoverType = marshal.TrimParseStringPointer(strings.TrimSpace(splitted[3]))
			} else {
				service.FailoverType = nil
			}
			if strings.TrimSpace(splitted[4]) != "" {
				if service.FailoverRetries, err = marshal.TrimParseIntPointer(splitted[4]); err != nil {
					merr = multierror.Append(merr, ercutils.NewError(err))
				}
			} else {
				service.FailoverRetries = nil
			}
			if strings.TrimSpace(splitted[5]) != "" {
				if service.FailoverDelay, err = marshal.TrimParseIntPointer(splitted[5]); err != nil {
					merr = multierror.Append(merr, ercutils.NewError(err))
				}
			} else {
				service.FailoverDelay = nil
			}
			if strings.TrimSpace(splitted[6]) != "" {
				service.Enabled = marshal.TrimParseBoolPointer(splitted[6])
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
