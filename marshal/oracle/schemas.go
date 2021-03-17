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
	"github.com/ercole-io/ercole/v2/model"
)

// Schemas returns information about database tablespaces extracted
// from the tablespaces fetcher command output.
func Schemas(cmdOutput []byte) []model.OracleDatabaseSchema {
	schemas := []model.OracleDatabaseSchema{}
	scanner := bufio.NewScanner(bytes.NewReader(cmdOutput))

	for scanner.Scan() {
		schema := new(model.OracleDatabaseSchema)
		line := scanner.Text()
		splitted := strings.Split(line, "|||")
		if len(splitted) == 8 {
			schema.User = strings.TrimSpace(splitted[3])
			schema.Total = marshal.TrimParseInt(splitted[4])
			schema.Tables = marshal.TrimParseInt(splitted[5])
			schema.Indexes = marshal.TrimParseInt(splitted[6])
			schema.LOB = marshal.TrimParseInt(splitted[7])

			schemas = append(schemas, *schema)
		}
	}
	return schemas
}
