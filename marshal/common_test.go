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
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package marshal

import (
	//"fmt"

	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrimParseInt(t *testing.T) {
	testCases := []struct {
		s        string
		expected int
		expErr   error
	}{
		{
			s:        "",
			expected: 0,
			expErr:   strconv.ErrSyntax,
		},
		{
			s:        "N/A",
			expected: 0,
			expErr:   strconv.ErrSyntax,
		},
		{
			s:        "42.42",
			expected: 0,
			expErr:   strconv.ErrSyntax,
		},
		{
			s:        "42",
			expected: 42,
			expErr:   nil,
		},
	}

	for _, tc := range testCases {
		actual, actErr := TrimParseInt(tc.s)
		assert.Equal(t, tc.expected, actual)

		if tc.expErr == nil {
			assert.Nil(t, actErr)
			continue
		}

		assert.ErrorIs(t, actErr, tc.expErr)
	}
}

func TestTrimParseFloat64PointerSafeComma(t *testing.T) {
	testCases := []struct {
		s           string
		nils        []string
		expected    *float64
		expectedErr error
	}{
		{
			s:           "N/A",
			nils:        []string{"N/A"},
			expected:    nil,
			expectedErr: nil,
		},
		{
			s:           "42",
			nils:        []string{"N/A"},
			expected:    getPointerToFloat(42),
			expectedErr: nil,
		},
		{
			s:           "42.42",
			nils:        []string{"N/A"},
			expected:    getPointerToFloat(42.42),
			expectedErr: nil,
		},
		{
			s:           "42,42",
			nils:        []string{"N/A"},
			expected:    getPointerToFloat(42.42),
			expectedErr: nil,
		},
		{
			s:           "42!42",
			nils:        []string{"N/A"},
			expected:    nil,
			expectedErr: &strconv.NumError{},
		},
	}

	for _, tc := range testCases {
		v, err := TrimParseFloat64PointerSafeComma(tc.s, tc.nils...)

		assert.Equal(t, tc.expected, v)

		if tc.expectedErr == nil {
			assert.Nil(t, err)
			continue
		}

		assert.ErrorAs(t, err, &tc.expectedErr)
	}
}

func getPointerToFloat(f float64) *float64 {
	return &f
}
