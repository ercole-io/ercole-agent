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

func TrimParseInt(s string) int {
	s = strings.TrimSpace(s)

	val, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}

	return val
}

func TrimParseIntPointer(s string, nils ...string) *int {
	for _, aNil := range nils {
		if s == aNil {
			return nil
		}
	}

	i := TrimParseInt(s)
	return &i
}

func TrimParseInt64(s string) int64 {
	s = strings.TrimSpace(s)

	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}

	return val
}

func TrimParseUint(s string) uint {
	s = strings.TrimSpace(s)

	val, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		panic(err)
	}

	return uint(val)
}

func TrimParseFloat64(s string) float64 {
	s = strings.TrimSpace(s)

	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}

	return val
}

func TrimParseFloat64Pointer(s string, nils ...string) *float64 {
	s = strings.TrimSpace(s)

	for _, aNil := range nils {
		if s == aNil {
			return nil
		}
	}

	f := TrimParseFloat64(s)
	return &f
}

func TrimParseBool(s string) bool {
	s = strings.TrimSpace(s)

	return parseBool(s)
}

func TrimParseStringPointer(s string, nils ...string) *string {
	for _, aNil := range nils {
		if s == aNil {
			return nil
		}
	}

	return &s
}

func parseKeyValueColonSeparated(b []byte) map[string]string {
	scanner := bufio.NewScanner(strings.NewReader(string(b)))

	data := make(map[string]string, 20)

	for scanner.Scan() {
		line := scanner.Text()
		splitted := strings.Split(line, ":")
		key := strings.TrimSpace(splitted[0])
		value := strings.TrimSpace(splitted[1])

		data[key] = value
	}

	return data
}

// NewIter return a an iterator on each string of a slice
func NewIter(splitted []string) func() string {
	i := -1
	return func() string {
		i++

		return splitted[i]
	}
}
