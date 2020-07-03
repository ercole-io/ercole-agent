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

	"github.com/ercole-io/ercole-agent/agentmodel"
)

// Oratab marshals a list of dbs (one per line) from the oratab command
func Oratab(cmdOutput []byte) []agentmodel.OratabEntry {

	var oratab []agentmodel.OratabEntry

	scanner := bufio.NewScanner(strings.NewReader(string(cmdOutput)))
	for scanner.Scan() {
		oratabEntry := agentmodel.OratabEntry{}
		line := scanner.Text()
		splitted := strings.Split(line, ":")

		oratabEntry.DBName = strings.TrimSpace(splitted[0])
		oratabEntry.OracleHome = strings.TrimSpace(splitted[1])

		oratab = append(oratab, oratabEntry)
	}

	return oratab
}
