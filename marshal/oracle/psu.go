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

	"github.com/ercole-io/ercole/v2/model"
)

// PSU returns informations about PSU parsed from fetcher command output.
func PSU(cmdOutput []byte) []model.OracleDatabasePSU {
	psuS := []model.OracleDatabasePSU{}
	scanner := bufio.NewScanner(bytes.NewReader(cmdOutput))

	for scanner.Scan() {
		psu := new(model.OracleDatabasePSU)
		line := scanner.Text()

		splitted := strings.Split(line, "|||")
		if len(splitted) == 2 {
			psu.Description = strings.TrimSpace(splitted[0])
			psu.Date = strings.TrimSpace(splitted[1])
			psuS = append(psuS, *psu)
		}
	}

	return psuS
}
