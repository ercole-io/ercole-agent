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

	"github.com/ercole-io/ercole-agent/v2/marshal"
	"github.com/ercole-io/ercole/v2/model"
)

// Patches returns information about database tablespaces extracted
// from the tablespaces fetcher command output.
func Patches(cmdOutput []byte) []model.OracleDatabasePatch {
	patches := []model.OracleDatabasePatch{}
	scanner := bufio.NewScanner(strings.NewReader(string(cmdOutput)))

	for scanner.Scan() {
		patch := new(model.OracleDatabasePatch)
		line := scanner.Text()
		splitted := strings.Split(line, "|||")
		if len(splitted) == 9 {
			patch.Version = strings.TrimSpace(splitted[4])

			patchID := strings.TrimSpace(splitted[5])
			if patchID != "" {
				patch.PatchID = marshal.TrimParseInt(patchID)
			}

			patch.Action = strings.TrimSpace(splitted[6])
			patch.Description = strings.TrimSpace(splitted[7])
			patch.Date = strings.TrimSpace(splitted[8])

			patches = append(patches, *patch)
		}
	}
	return patches
}
