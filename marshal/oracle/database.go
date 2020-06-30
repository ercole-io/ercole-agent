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
	"strings"

	"github.com/ercole-io/ercole-agent/marshal"
	"github.com/ercole-io/ercole/model"
)

// Database returns information about database extracted
// from the db fetcher command output.
func Database(cmdOutput []byte) model.OracleDatabase {
	var db model.OracleDatabase
	scanner := bufio.NewScanner(strings.NewReader(string(cmdOutput)))

	for scanner.Scan() {
		line := scanner.Text()
		splitted := strings.Split(line, "|||")
		if len(splitted) == 24 {
			db.Name = strings.TrimSpace(splitted[0])
			db.UniqueName = strings.TrimSpace(splitted[1])
			db.InstanceNumber = marshal.TrimParseInt(splitted[2])
			db.Status = strings.TrimSpace(splitted[3])
			db.Version = strings.TrimSpace(splitted[4])
			db.Platform = strings.TrimSpace(splitted[5])
			db.Archivelog = marshal.TrimParseBool(splitted[6])
			db.Charset = strings.TrimSpace(splitted[7])
			db.NCharset = strings.TrimSpace(splitted[8])
			db.BlockSize = marshal.TrimParseInt(splitted[9])
			db.CPUCount = marshal.TrimParseInt(splitted[10])
			db.SGATarget = marshal.TrimParseFloat64(splitted[11])
			db.PGATarget = marshal.TrimParseFloat64(splitted[12])
			db.MemoryTarget = marshal.TrimParseFloat64(splitted[13])
			db.SGAMaxSize = marshal.TrimParseFloat64(splitted[14])
			db.SegmentsSize = marshal.TrimParseFloat64(splitted[15])
			db.DatafileSize = marshal.TrimParseFloat64(splitted[16])
			db.Allocated = marshal.TrimParseFloat64(splitted[17])

			db.Elapsed = marshal.TrimParseFloat64Pointer(splitted[18], "N/A")
			db.DBTime = marshal.TrimParseFloat64Pointer(splitted[19], "N/A")
			db.DailyCPUUsage = marshal.TrimParseFloat64Pointer(splitted[20], "N/A")
			db.Work = marshal.TrimParseFloat64Pointer(splitted[21], "N/A")

			db.ASM = marshal.TrimParseBool(splitted[22])
			db.Dataguard = marshal.TrimParseBool(splitted[23])

			if *db.DailyCPUUsage == 0 {
				*db.DailyCPUUsage = *db.Work
			}
		}
	}
	return db
}
