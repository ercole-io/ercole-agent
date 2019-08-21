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

package marshal

import (
	"bufio"
	"encoding/json"
	"log"
	"strings"

	"github.com/ercole-io/ercole-agent/model"
)

// Host returns a Host struct from the output of the host
// fetcher command. Host fields output is in key: value format separated by a newline
func Host(cmdOutput []byte) model.Host {

	lines := "{"
	scanner := bufio.NewScanner(strings.NewReader(string(cmdOutput)))
	for scanner.Scan() {
		line := scanner.Text()
		splitted := strings.Split(line, ":")
		key := strings.TrimSpace(splitted[0])
		value := strings.TrimSpace(splitted[1])
		lines += marshalKey(key) + marshalValue(value) + ", "
	}

	lines += "}"
	lines = strings.Replace(lines, ", }", "}", -1)

	b := []byte(lines)
	var m model.Host
	err := json.Unmarshal(b, &m)

	if err != nil {
		log.Fatal(err)
	}

	return m
}
