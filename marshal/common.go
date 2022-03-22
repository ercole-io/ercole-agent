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

package marshal

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

func parseBool(s string) bool {
	return strings.EqualFold(s, "y") ||
		strings.EqualFold(s, "yes") ||
		strings.EqualFold(s, "true") ||
		strings.EqualFold(s, "1")
}

func parseInt(s string) int {
	i, err := strconv.Atoi(s)

	if err != nil {
		return 0
	}

	return i
}

func TrimParseInt(s string) (int, error) {
	s = strings.TrimSpace(s)

	val, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("Can't parse value \"%s\" as int; err: %w", s, err)
	}

	return val, nil
}

func TrimParseIntPointer(s string, nils ...string) (*int, error) {
	for _, aNil := range nils {
		if s == aNil {
			return nil, nil
		}
	}

	i, err := TrimParseInt(s)

	return &i, err
}

func TrimParseInt64(s string) (int64, error) {
	s = strings.TrimSpace(s)

	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("Can't parse value \"%s\" as int64; err: %w", s, err)
	}

	return val, nil
}

func TrimParseUint(s string) (uint, error) {
	s = strings.TrimSpace(s)

	val, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}

	return uint(val), nil
}

func TrimParseFloat64(s string) (float64, error) {
	s = strings.TrimSpace(s)

	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}

	return val, nil
}

func TrimParseFloat64SafeComma(s string) (float64, error) {
	s = strings.TrimSpace(s)

	if val, err := strconv.ParseFloat(s, 64); err == nil {
		return val, nil
	}

	s = strings.Replace(s, ",", ".", 1)
	val, err := strconv.ParseFloat(s, 64)

	if err == nil {
		return val, nil
	}

	return 0, err
}

func TrimParseFloat64Pointer(s string, nils ...string) (*float64, error) {
	s = strings.TrimSpace(s)

	for _, aNil := range nils {
		if s == aNil {
			return nil, nil
		}
	}

	f, err := TrimParseFloat64(s)
	if err != nil {
		return nil, err
	}

	return &f, nil
}

func TrimParseFloat64PointerSafeComma(s string, nils ...string) (*float64, error) {
	s = strings.TrimSpace(s)

	for _, aNil := range nils {
		if s == aNil {
			return nil, nil
		}
	}

	s = strings.Replace(s, ",", ".", 1)

	f, err := TrimParseFloat64(s)
	if err != nil {
		return nil, err
	}

	return &f, nil
}

func TrimParseBool(s string) bool {
	s = strings.TrimSpace(s)
	return parseBool(s)
}

func TrimParseBoolPointer(s string, nils ...string) *bool {
	for _, aNil := range nils {
		if s == aNil {
			return nil
		}
	}

	b := TrimParseBool(s)

	return &b
}

func TrimParseStringPointer(s string, nils ...string) *string {
	for _, aNil := range nils {
		if s == aNil {
			return nil
		}
	}

	return &s
}

func TrimParseDatePointer(s string, nils ...string) (*time.Time, error) {
	for _, aNil := range nils {
		if s == aNil {
			return nil, nil
		}
	}

	d, err := time.Parse("02-Jan-06", strings.TrimSpace(s))
	if err != nil {
		return nil, err
	}

	return &d, nil
}

func parseKeyValueColonSeparated(b []byte) map[string]string {
	return ParseKeyValue(b, ":")
}

// ParseKeyValue scan lines from b and put key values in map
func ParseKeyValue(b []byte, sep string) map[string]string {
	scanner := bufio.NewScanner(bytes.NewBuffer(b))

	data := make(map[string]string, 20)

	for scanner.Scan() {
		line := scanner.Text()
		splitted := strings.SplitN(line, sep, 2)

		if len(splitted) == 2 {
			data[strings.TrimSpace(splitted[0])] = strings.TrimSpace(splitted[1])
		}
	}

	return data
}

type Iterator func() string

type CsvScanner struct {
	reader  *csv.Reader
	records []string
	iter    Iterator
}

// SafeScan advances the CsvScanner to the next line with correct number of fields,
// which will then be available through the Iter method.
// It returns false when the scan stops by reaching the end of the input.
func (s *CsvScanner) SafeScan() bool {
	var err error

	for err != io.EOF {
		s.records, err = s.reader.Read()

		if err == nil {
			s.iter = NewIter(s.records)
			return true
		}
	}

	s.iter = nil

	return false
}

func (s *CsvScanner) Iter() string {
	return s.iter()
}

func (s *CsvScanner) Get(i int) string {
	return s.records[i]
}

func NewCsvScanner(cmdOutput []byte, fieldsPerRecord int) CsvScanner {
	reader := csv.NewReader(bytes.NewReader(cmdOutput))
	reader.FieldsPerRecord = fieldsPerRecord
	reader.Comma = ';'

	scanner := CsvScanner{
		reader: reader,
	}

	return scanner
}

// NewIter return a an iterator on each string of a slice
func NewIter(splitted []string) Iterator {
	i := -1

	return func() string {
		i++
		return splitted[i]
	}
}

// NewIter return a an iterator on each string of a slice
func NewSplitIter(s, sep string) func() string {
	splitted := strings.Split(s, sep)

	i := -1

	return func() string {
		i++
		return splitted[i]
	}
}
