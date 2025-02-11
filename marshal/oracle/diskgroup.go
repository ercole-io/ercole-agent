// Copyright (c) 2025 Sorint.lab S.p.A.
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

// Diskgroup marshaller
func DiskGroups(cmdOutput []byte) ([]model.OracleDatabaseDiskGroup, error) {
	diskgroups := []model.OracleDatabaseDiskGroup{}
	scanner := bufio.NewScanner(bytes.NewReader(cmdOutput))

	var merr, err error

	for scanner.Scan() {
		line := scanner.Text()

		splitted := strings.Split(line, "|||")
		if len(splitted) == 4 {
			diskgroup := model.OracleDatabaseDiskGroup{}

			diskgroup.DiskGroupName = strings.TrimSpace(splitted[0])

			if diskgroup.TotalSpace, err = marshal.TrimParseUnsafeFloat64(splitted[1], marshal.TrimParseFloat64SafeComma); err != nil {
				merr = multierror.Append(merr, ercutils.NewError(err))
			}

			if diskgroup.UsedSpace, err = marshal.TrimParseUnsafeFloat64(splitted[2], marshal.TrimParseFloat64SafeComma); err != nil {
				merr = multierror.Append(merr, ercutils.NewError(err))
			}

			if diskgroup.FreeSpace, err = marshal.TrimParseUnsafeFloat64(splitted[3], marshal.TrimParseFloat64SafeComma); err != nil {
				merr = multierror.Append(merr, ercutils.NewError(err))
			}

			diskgroups = append(diskgroups, diskgroup)
		}
	}

	if merr != nil {
		return nil, merr
	}

	return diskgroups, nil
}
