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

	"github.com/hashicorp/go-multierror"

	"github.com/ercole-io/ercole-agent/v2/marshal"
	"github.com/ercole-io/ercole/v2/model"
	ercutils "github.com/ercole-io/ercole/v2/utils"
)

// Backups marshals a backup output list into a struct.
func Backups(cmdOutput []byte) ([]model.OracleDatabaseBackup, error) {
	backups := []model.OracleDatabaseBackup{}

	scanner := bufio.NewScanner(bytes.NewReader(cmdOutput))

	var merr error

	for scanner.Scan() {
		backup := new(model.OracleDatabaseBackup)
		line := scanner.Text()

		splitted := strings.Split(line, "|||")
		if len(splitted) == 5 {
			backup.BackupType = strings.TrimSpace(splitted[0])
			backup.Hour = strings.TrimSpace(splitted[1])

			weekDays := strings.TrimSpace(splitted[2])
			backup.WeekDays = strings.Split(weekDays, ",")

			avgBckSize, err := marshal.TrimParseFloat64(splitted[3])
			if err != nil {
				merr = multierror.Append(merr, ercutils.NewError(err))
			}

			backup.AvgBckSize = avgBckSize
			backup.Retention = strings.TrimSpace(splitted[4])
			backups = append(backups, *backup)
		}
	}

	if merr != nil {
		return nil, merr
	}

	return backups, nil
}
