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

	"github.com/ercole-io/ercole/model"
)

// ListPDB returns information about pdbs extracted
// from the listdb fetcher command output.
func ListPDB(cmdOutput []byte) []model.OracleDatabasePluggableDatabase {
	pdbs := []model.OracleDatabasePluggableDatabase{}
	scanner := bufio.NewScanner(strings.NewReader(string(cmdOutput)))

	for scanner.Scan() {
		pdb := new(model.OracleDatabasePluggableDatabase)
		line := scanner.Text()
		splitted := strings.Split(line, "|||")

		pdb.Name = strings.TrimSpace(splitted[0])
		pdb.Status = strings.TrimSpace(splitted[1])
		pdb.Services = []model.OracleDatabaseService{}

		pdbs = append(pdbs, *pdb)
	}
	return pdbs
}
