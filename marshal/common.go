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
	"regexp"
	"strconv"
	"strings"
)

func marshalValue(s string) string {

	if s == "Y" {
		return "true"
	}

	if s == "N" {
		return "false"
	}

	re := regexp.MustCompile(`^[0-9]+(\.[0-9]+)?$`)
	isNum := re.Match([]byte(s))

	if isNum {
		return s
	}

	return "\"" + s + "\""

}

func marshalString(s string) string {
	s = strings.Replace(s, "\\", "\\\\", -1)
	return "\"" + s + "\""
}

func marshalKey(s string) string {

	return "\"" + s + "\" : "

}

func cleanTr(s string) string {
	value := strings.Trim(s, " ")
	value = strings.Replace(value, "\n", "", -1)
	value = strings.Replace(value, "\t", "", -1)
	value = strings.Trim(value, " ")

	return value
}

func parseBool(s string) bool {
	return s == "Y" || s == "TRUE"
}

func parseInt(s string) int {

	i, err := strconv.Atoi(s)

	if err != nil {
		return 0
	}

	return i

}

// s is supposed to be non null and already trimmed
func parseCount(s string) float32 {

	if s == "" {
		return 0
	}

	count, err := strconv.ParseFloat(s, 32)

	if err != nil {
		return 0
	}

	return float32(count)

}
