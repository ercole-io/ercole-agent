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

	"github.com/ercole-io/ercole/v2/model"
)

func MemoryAdvisor(cmdOutput []byte) (*model.OracleDatabaseMemoryAdvisor, error) {
	res := &model.OracleDatabaseMemoryAdvisor{}

	scanner := bufio.NewScanner(bytes.NewReader(cmdOutput))

	// check if the current line is in the designed marker or not
	isBegin := false

	for scanner.Scan() {
		line := scanner.Text()

		if line == "BEGINOUTPUT" {
			isBegin = true
			continue
		}

		if line == "ENDOUTPUT" {
			isBegin = false
			continue
		}

		if !isBegin {
			continue
		}

		splitted := strings.Split(line, "|||")

		if len(splitted) == 2 {
			switch splitted[0] {
			case "MEMORY_SIZE_LOWER_GB":
				res.AutomaticMemoryManagement = true
				res.MemorySizeLowerGb = splitted[1]

			case "PGA_TARGET_AGGREGATE_LOWER_GB":
				res.AutomaticMemoryManagement = false
				res.PgaTargetAggregateLowerGb = splitted[1]

			case "SGA_SIZE_LOWER_GB":
				res.AutomaticMemoryManagement = false
				res.SgaSizeLowerGb = splitted[1]
			}
		}
	}

	return res, nil
}
