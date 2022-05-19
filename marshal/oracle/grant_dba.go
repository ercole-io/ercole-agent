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

func GrantDba(cmdOutput []byte) []model.OracleGrantDba {
	grants := make([]model.OracleGrantDba, 0)

	scanner := bufio.NewScanner(bytes.NewReader(cmdOutput))
	for scanner.Scan() {
		line := scanner.Text()
		grant := model.OracleGrantDba{}

		splitted := strings.Split(line, "|||")
		if len(splitted) == 3 {
			grant.Grantee = strings.TrimSpace(splitted[0])
			grant.AdminOption = strings.TrimSpace(splitted[1])
			grant.DefaultRole = strings.TrimSpace(splitted[2])
		}

		grants = append(grants, grant)
	}

	return grants
}
