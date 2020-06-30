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

	"github.com/ercole-io/ercole/model"
)

// Filesystems returns a list of Filesystem entries extracted
// from the filesystem fetcher command output.
// Filesystem output is a list of filesystem entries with positional attribute columns
// separated by one or more spaces
func Filesystems(cmdOutput []byte) []model.Filesystem {

	lines := "["
	scanner := bufio.NewScanner(strings.NewReader(string(cmdOutput)))
	for scanner.Scan() {
		lines += "{"
		line := scanner.Text()
		line = strings.Join(strings.Fields(line), " ")
		splitted := strings.Split(line, " ")
		lines += marshalKey("filesystem") + marshalString(strings.TrimSpace(splitted[0])) + ", "
		lines += marshalKey("fstype") + marshalString(strings.TrimSpace(splitted[1])) + ", "
		lines += marshalKey("size") + marshalString(strings.TrimSpace(splitted[2])) + ", "
		lines += marshalKey("used") + marshalString(strings.TrimSpace(splitted[3])) + ", "
		lines += marshalKey("available") + marshalString(strings.TrimSpace(splitted[4])) + ", "
		lines += marshalKey("usedperc") + marshalString(strings.TrimSpace(splitted[5])) + ", "
		lines += marshalKey("mountedon") + marshalString(strings.TrimSpace(splitted[6])) + ", "
		lines += "},"
	}

	lines += "]"
	lines = strings.Replace(lines, ", }", "}", -1)
	lines = strings.Replace(lines, "},]", "}]", -1)

	b := []byte(lines)
	var m []model.Filesystem
	err := json.Unmarshal(b, &m)

	if err != nil {
		log.Fatal(err)
	}

	return m
}
