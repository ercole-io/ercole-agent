// Copyright (c) 2019 Sorint.lab S.p.A.
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

// Tablespaces returns information about database tablespaces extracted
// from the tablespaces fetcher command output.
func Tablespaces(cmdOutput []byte) ([]model.OracleDatabaseTablespace, error) {
	tablespaces := []model.OracleDatabaseTablespace{}
	scanner := bufio.NewScanner(bytes.NewReader(cmdOutput))
	var merr, err error

	for scanner.Scan() {
		tablespace := new(model.OracleDatabaseTablespace)
		line := scanner.Text()
		splitted := strings.Split(line, "|||")
		if len(splitted) == 9 {
			tablespace.Name = strings.TrimSpace(splitted[3])
			if tablespace.MaxSize, err = marshal.TrimParseFloat64(splitted[4]); err != nil {
				merr = multierror.Append(merr, ercutils.NewError(err))
			}
			if tablespace.Total, err = marshal.TrimParseFloat64(splitted[5]); err != nil {
				merr = multierror.Append(merr, ercutils.NewError(err))
			}
			if tablespace.Used, err = marshal.TrimParseFloat64(splitted[6]); err != nil {
				merr = multierror.Append(merr, ercutils.NewError(err))
			}
			if tablespace.UsedPerc, err = marshal.TrimParseFloat64(splitted[7]); err != nil {
				merr = multierror.Append(merr, ercutils.NewError(err))
			}
			tablespace.Status = strings.TrimSpace(splitted[8])

			tablespaces = append(tablespaces, *tablespace)
		}
	}

	if merr != nil {
		return nil, merr
	}
	return tablespaces, nil
}
