// Copyright (c) 2024 Sorint.lab S.p.A.
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
)

func SizePDB(cmdOutput []byte) (model.OracleDatabasePdbSize, error) {
	scanner := bufio.NewScanner(bytes.NewReader(cmdOutput))

	pdbsize := model.OracleDatabasePdbSize{}

	for scanner.Scan() {
		line := scanner.Text()
		splitted := strings.Split(line, "|||")
		iter := marshal.NewIter(splitted)

		if len(splitted) != 5 {
			return pdbsize, ercutils.NewErrorf("invalid line")
		}

		segmentsSize, err := marshal.TrimParseFloat64(iter())
		if err != nil {
			return pdbsize, err
		}

		pdbsize.SegmentsSize = segmentsSize

		datafileSize, err := marshal.TrimParseFloat64(iter())
		if err != nil {
			return pdbsize, err
		}

		pdbsize.DatafileSize = datafileSize

		allocable, err := marshal.TrimParseFloat64(iter())
		if err != nil {
			return pdbsize, err
		}

		pdbsize.Allocable = allocable

		sgatarget, err := marshal.TrimParseFloat64(iter())
		if err != nil {
			return pdbsize, err
		}

		pdbsize.SGATarget = sgatarget

		pgaAggregateTarget, err := marshal.TrimParseFloat64(iter())
		if err != nil {
			return pdbsize, err
		}

		pdbsize.PGAAggregateTarget = pgaAggregateTarget
	}

	return pdbsize, nil
}
