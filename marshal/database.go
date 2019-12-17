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

package marshal

import (
	"bufio"
	"strings"

	"github.com/ercole-io/ercole-agent/model"
)

// Database returns information about database extracted
// from the db fetcher command output.
func Database(cmdOutput []byte) model.Database {
	var db model.Database
	scanner := bufio.NewScanner(strings.NewReader(string(cmdOutput)))

	for scanner.Scan() {
		line := scanner.Text()
		splitted := strings.Split(line, "|||")
		if len(splitted) == 23 {
			db.Name = strings.TrimSpace(splitted[0])
			db.UniqueName = strings.TrimSpace(splitted[1])
			db.InstanceNumber = strings.TrimSpace(splitted[2])
			db.Status = strings.TrimSpace(splitted[3])
			db.Version = strings.TrimSpace(splitted[4])
			db.Platform = strings.TrimSpace(splitted[5])
			db.Archivelog = strings.TrimSpace(splitted[6])
			db.Charset = strings.TrimSpace(splitted[7])
			db.NCharset = strings.TrimSpace(splitted[8])
			db.BlockSize = strings.TrimSpace(splitted[9])
			db.CPUCount = strings.TrimSpace(splitted[10])
			db.SGATarget = strings.TrimSpace(splitted[11])
			db.PGATarget = strings.TrimSpace(splitted[12])
			db.MemoryTarget = strings.TrimSpace(splitted[13])
			db.SGAMaxSize = strings.TrimSpace(splitted[14])
			db.SegmentsSize = strings.TrimSpace(splitted[15])
			db.Used = strings.TrimSpace(splitted[16])
			db.Allocated = strings.TrimSpace(splitted[17])
			db.Elapsed = strings.TrimSpace(splitted[18])
			db.DBTime = strings.TrimSpace(splitted[19])
			db.DailyCPUUsage = strings.TrimSpace(splitted[20])
			db.Work = strings.TrimSpace(splitted[21])
			db.ASM = parseBool(strings.TrimSpace(splitted[22]))
			db.Dataguard = parseBool(strings.TrimSpace(splitted[23]))

			if db.DailyCPUUsage == "" {
				db.DailyCPUUsage = db.Work
			}
		}
	}
	return db
}
