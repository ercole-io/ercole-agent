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
	"log"
	"strconv"
	"strings"

	"github.com/ercole-io/ercole/v2/model"
)

// Backups marshals a backup output list into a struct.
func Backups(cmdOutput []byte) []model.OracleDatabaseBackup {
	backups := []model.OracleDatabaseBackup{}

	scanner := bufio.NewScanner(bytes.NewReader(cmdOutput))

	for scanner.Scan() {
		backup := new(model.OracleDatabaseBackup)
		line := scanner.Text()
		splitted := strings.Split(line, "|||")
		if len(splitted) == 5 {
			backup.BackupType = strings.TrimSpace(splitted[0])
			backup.Hour = strings.TrimSpace(splitted[1])

			weekDays := strings.TrimSpace(splitted[2])
			backup.WeekDays = strings.Split(weekDays, ",")

			avgBckSize, err := strconv.ParseFloat(splitted[3], 64)
			if err != nil {
				log.Printf("%v\n", err)
			}

			backup.AvgBckSize = avgBckSize
			backup.Retention = strings.TrimSpace(splitted[4])
			backups = append(backups, *backup)
		}
	}
	return backups
}
