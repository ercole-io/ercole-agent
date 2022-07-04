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

func Instance(cmdOutput []byte) (*model.PostgreSQLInstance, error) {
	scanner := bufio.NewScanner(bytes.NewReader(cmdOutput))
	result := model.PostgreSQLInstance{}

	var merr, err error

	for scanner.Scan() {
		line := scanner.Text()

		splitted := strings.Split(line, "|")
		if len(splitted) < 5 {
			continue
		}

		iter := marshal.NewIter(splitted)

		if result.MaxConnections, err = strconv.Atoi(iter()); err != nil {
			merr = multierror.Append(merr, err)
		}

		if result.InstanceSize, err = strconv.Atoi(iter()); err != nil {
			merr = multierror.Append(merr, err)
		}

		if result.UsersNum, err = strconv.Atoi(iter()); err != nil {
			merr = multierror.Append(merr, err)
		}

		if result.DbNum, err = strconv.Atoi(iter()); err != nil {
			merr = multierror.Append(merr, err)
		}

		if result.TblspNum, err = strconv.Atoi(iter()); err != nil {
			merr = multierror.Append(merr, err)
		}

		if len(splitted) > 5 {
			if result.Isinreplica, err = strconv.ParseBool(iter()); err != nil {
				merr = multierror.Append(merr, err)
			}

			if result.Ismaster, err = strconv.ParseBool(iter()); err != nil {
				merr = multierror.Append(merr, err)
			}

			if result.Isslave, err = strconv.ParseBool(iter()); err != nil {
				merr = multierror.Append(merr, err)
			}

			if result.ArchiverWorking, err = strconv.ParseBool(iter()); err != nil {
				merr = multierror.Append(merr, err)
			}

			if result.SlavesNum, err = strconv.Atoi(iter()); err != nil {
				merr = multierror.Append(merr, err)
			}

			result.Charset = iter()
		}
	}

	if merr != nil {
		return nil, merr
	}

	return &result, nil
}
