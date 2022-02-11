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

// Addms marshaller
func Addms(cmdOutput []byte) ([]model.OracleDatabaseAddm, error) {
	addms := []model.OracleDatabaseAddm{}
	scanner := bufio.NewScanner(bytes.NewReader(cmdOutput))

	var merr, err error

	for scanner.Scan() {
		addm := new(model.OracleDatabaseAddm)
		line := scanner.Text()

		splitted := strings.Split(line, "|||")
		if len(splitted) == 6 {
			addm.Finding = strings.TrimSpace(splitted[2])
			addm.Recommendation = strings.TrimSpace(splitted[3])
			addm.Action = strings.TrimSpace(splitted[4])

			if addm.Benefit, err = marshal.TrimParseFloat64(splitted[5]); err != nil {
				merr = multierror.Append(merr, ercutils.NewError(err))
			}

			addms = append(addms, *addm)
		}
	}

	if merr != nil {
		return nil, merr
	}

	return addms, nil
}
