// Copyright (c) 2021 Sorint.lab S.p.A.
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

package mysql

import (
	"bufio"
	"bytes"
	"strings"
)

func HighAvailability(cmdOutput []byte) (isMaster bool) {
	scanner := bufio.NewScanner(bytes.NewReader(cmdOutput))

	for scanner.Scan() {
		if strings.TrimSpace(scanner.Text()) == `"mysql_innodb_cluster_metadata"` {
			return true
		}
	}

	return false
}
