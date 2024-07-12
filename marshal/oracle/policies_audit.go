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
)

func PoliciesAudit(cmdOutput []byte) ([]string, error) {
	scanner := bufio.NewScanner(bytes.NewReader(cmdOutput))

	res := make([]string, 0)

	for scanner.Scan() {
		line := scanner.Text()

		line = strings.TrimSpace(line)

		if line != "" {
			res = append(res, line)
		}
	}

	if len(res) == 0 {
		return nil, nil
	}

	return res, nil
}
