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

// Database returns information about database extracted
// from the db fetcher command output.
func Database(cmdOutput []byte) (*model.OracleDatabase, error) {
	var db model.OracleDatabase
	var merr, err error
	scanner := bufio.NewScanner(bytes.NewReader(cmdOutput))

	for scanner.Scan() {
		line := scanner.Text()
		splitted := strings.Split(line, "|||")
		if len(splitted) == 27 {
			iter := marshal.NewIter(splitted)

			db.Name = strings.TrimSpace(iter())
			if db.DbID, err = marshal.TrimParseUint(iter()); err != nil {
				merr = multierror.Append(merr, ercutils.NewError(err))
			}
			db.Role = strings.TrimSpace(iter())
			db.UniqueName = strings.TrimSpace(iter())
			if db.InstanceNumber, err = marshal.TrimParseInt(iter()); err != nil {
				merr = multierror.Append(merr, ercutils.NewError(err))
			}
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
				merr = multierror.Append(merr, ercutils.NewErrorf("Invalid archivelog value: %s", archivelog))
			}

			db.Charset = strings.TrimSpace(iter())
			db.NCharset = strings.TrimSpace(iter())
			if db.BlockSize, err = marshal.TrimParseInt(iter()); err != nil {
				merr = multierror.Append(merr, ercutils.NewError(err))
			}
			if db.CPUCount, err = marshal.TrimParseInt(iter()); err != nil {
				merr = multierror.Append(merr, ercutils.NewError(err))
			}
			if db.SGATarget, err = marshal.TrimParseFloat64(iter()); err != nil {
				merr = multierror.Append(merr, ercutils.NewError(err))
			}
			if db.PGATarget, err = marshal.TrimParseFloat64(iter()); err != nil {
				merr = multierror.Append(merr, ercutils.NewError(err))
			}
			if db.MemoryTarget, err = marshal.TrimParseFloat64(iter()); err != nil {
				merr = multierror.Append(merr, ercutils.NewError(err))
			}
			if db.SGAMaxSize, err = marshal.TrimParseFloat64(iter()); err != nil {
				merr = multierror.Append(merr, ercutils.NewError(err))
			}
			if db.SegmentsSize, err = marshal.TrimParseFloat64(iter()); err != nil {
				merr = multierror.Append(merr, ercutils.NewError(err))
			}
			if db.DatafileSize, err = marshal.TrimParseFloat64(iter()); err != nil {
				merr = multierror.Append(merr, ercutils.NewError(err))
			}
			if db.Allocable, err = marshal.TrimParseFloat64(iter()); err != nil {
				merr = multierror.Append(merr, ercutils.NewError(err))
			}
			if db.Elapsed, err = marshal.TrimParseFloat64PointerSafeComma(iter(), "N/A"); err != nil {
				merr = multierror.Append(merr, ercutils.NewError(err))
			}
			if db.DBTime, err = marshal.TrimParseFloat64PointerSafeComma(iter(), "N/A"); err != nil {
				merr = multierror.Append(merr, ercutils.NewError(err))
			}
			if db.DailyCPUUsage, err = marshal.TrimParseFloat64PointerSafeComma(iter(), "N/A"); err != nil {
				merr = multierror.Append(merr, ercutils.NewError(err))
			}
			if db.Work, err = marshal.TrimParseFloat64PointerSafeComma(iter(), "N/A"); err != nil {
				merr = multierror.Append(merr, ercutils.NewError(err))
			}

			db.ASM = marshal.TrimParseBool(iter())
			db.Dataguard = marshal.TrimParseBool(iter())

			if db.DailyCPUUsage == nil {
				db.DailyCPUUsage = db.Work
			}
		}
	}

	if merr != nil {
		return nil, merr
	}
	return &db, nil
}
