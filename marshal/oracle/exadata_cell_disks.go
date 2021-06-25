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

	"github.com/ercole-io/ercole-agent/v2/agentmodel"
	"github.com/ercole-io/ercole/v2/model"
	ercutils "github.com/ercole-io/ercole/v2/utils"
)

// ExadataCellDisks returns information about the cell disks extracted from exadata-storage-status command.
func ExadataCellDisks(cmdOutput []byte) (map[agentmodel.StorageServerName][]model.OracleExadataCellDisk, []error) {
	cellDisks := make(map[agentmodel.StorageServerName][]model.OracleExadataCellDisk)
	scanner := bufio.NewScanner(bytes.NewReader(cmdOutput))
	var errs []error
	var err error

	for scanner.Scan() {
		cellDisk := new(model.OracleExadataCellDisk)
		line := scanner.Text()
		splitted := strings.Split(line, "|||")
		if len(splitted) == 5 {
			storageServerName := strings.TrimSpace(splitted[0])

			cellDisk.Name = strings.TrimSpace(splitted[1])
			cellDisk.Status = strings.TrimSpace(splitted[2])
			if cellDisk.ErrCount, err = marshal.TrimParseInt(splitted[3]); err != nil {
				errs = append(errs, ercutils.NewError(err))
			}
			if cellDisk.UsedPerc, err = marshal.TrimParseInt(splitted[4]); err != nil {
				errs = append(errs, ercutils.NewError(err))
			}
			addCellDisk(cellDisks, storageServerName, cellDisk)
		}
	}

	if len(errs) > 0 {
		return nil, errs
	}

	return cellDisks, nil
}

func addCellDisk(cellDisks map[agentmodel.StorageServerName][]model.OracleExadataCellDisk,
	storageServerName string, cellDisk *model.OracleExadataCellDisk) {
	ssn := agentmodel.StorageServerName(storageServerName)

	storageServerCellDisks := cellDisks[ssn]
	storageServerCellDisks = append(storageServerCellDisks, *cellDisk)
	cellDisks[ssn] = storageServerCellDisks
}
