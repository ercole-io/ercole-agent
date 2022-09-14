// Copyright (c) 2022 Sorint.lab S.p.A.
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
	ercutils "github.com/ercole-io/ercole/v2/utils"
)

func CharsetPDB(cmdOutput []byte) (string, error) {
	scanner := bufio.NewScanner(bytes.NewReader(cmdOutput))

	var charset string

	for scanner.Scan() {
		line := scanner.Text()
		splitted := strings.Split(line, "|||")
		iter := marshal.NewIter(splitted)

		if len(splitted) != 1 {
			return "", ercutils.NewErrorf("Invalid line")
		}

		charset = strings.TrimSpace(iter())
	}

	return charset, nil
}
