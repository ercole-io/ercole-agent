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
	"fmt"
	"strings"

	"github.com/ercole-io/ercole-agent/v2/marshal"
	"github.com/ercole-io/ercole/v2/model"
)

// Database returns information about database extracted
// from the db fetcher command output.
func Database(cmdOutput []byte) model.OracleDatabase {
	var db model.OracleDatabase
	scanner := bufio.NewScanner(bytes.NewReader(cmdOutput))

	for scanner.Scan() {
		line := scanner.Text()
		splitted := strings.Split(line, "|||")
		if len(splitted) == 27 {
			iter := marshal.NewIter(splitted)

			db.Name = strings.TrimSpace(iter())
			db.DbID = marshal.TrimParseUint(iter())
			db.Role = strings.TrimSpace(iter())
			db.UniqueName = strings.TrimSpace(iter())
			db.InstanceNumber = marshal.TrimParseInt(iter())
			db.InstanceName = strings.TrimSpace(iter())
			db.Status = strings.TrimSpace(iter())
			db.Version = strings.TrimSpace(iter())
			db.Platform = strings.TrimSpace(iter())

			archivelog := strings.TrimSpace(iter())
			if archivelog == "ARCHIVELOG" {
				db.Archivelog = true
			} else if archivelog == "NOARCHIVELOG" {
				db.Archivelog = false
			} else {
				panic(fmt.Sprintf("Invalid archivelog value: %s", archivelog))
			}

			db.Charset = strings.TrimSpace(iter())
			db.NCharset = strings.TrimSpace(iter())
			db.BlockSize = marshal.TrimParseInt(iter())
			db.CPUCount = marshal.TrimParseInt(iter())
			db.SGATarget = marshal.TrimParseFloat64(iter())
			db.PGATarget = marshal.TrimParseFloat64(iter())
			db.MemoryTarget = marshal.TrimParseFloat64(iter())
			db.SGAMaxSize = marshal.TrimParseFloat64(iter())
			db.SegmentsSize = marshal.TrimParseFloat64(iter())
			db.DatafileSize = marshal.TrimParseFloat64(iter())
			db.Allocable = marshal.TrimParseFloat64(iter())

			db.Elapsed = marshal.TrimParseFloat64PointerSafeComma(iter(), "N/A")
			db.DBTime = marshal.TrimParseFloat64PointerSafeComma(iter(), "N/A")
			db.DailyCPUUsage = marshal.TrimParseFloat64PointerSafeComma(iter(), "N/A")
			db.Work = marshal.TrimParseFloat64PointerSafeComma(iter(), "N/A")

			db.ASM = marshal.TrimParseBool(iter())
			db.Dataguard = marshal.TrimParseBool(iter())

			if db.DailyCPUUsage == nil {
				db.DailyCPUUsage = db.Work
			}
		}
	}
	return db
}
