// Copyright (c) 2023 Sorint.lab S.p.A.
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
	"encoding/json"
	"testing"
	"time"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/stretchr/testify/assert"
)

var testDiskConsumptionsData string = `88.81|||.40
80.10|||.40
75.28|||.38
75.26|||.38
127.18|||.42
81.89|||.42
85.04|||.42
82.48|||.41
77.95|||.39
76.94|||.39
77.15|||.39
79.26|||.39
220611:00|||78.54|||.38
220611:10|||94.36|||N/A`

func TestDiskConsumption(t *testing.T) {
	currentTime = time.Now()
	t1 := currentTime.AddDate(0, 0, -30)
	t2 := currentTime.AddDate(0, 0, -7)
	t3 := currentTime.AddDate(0, 0, -14)
	t4 := currentTime.AddDate(0, 0, -8)
	t5 := currentTime.AddDate(0, 0, -21)
	t6 := currentTime.AddDate(0, 0, -15)
	t7 := currentTime.AddDate(0, 0, -28)
	t8 := currentTime.AddDate(0, 0, -22)
	t9 := currentTime.AddDate(0, 0, -1)
	t10 := currentTime.AddDate(0, 0, -2)
	t11 := currentTime.AddDate(0, 0, -3)
	t12 := currentTime.AddDate(0, 0, -4)
	t13 := currentTime.AddDate(0, 0, -5)
	t14 := currentTime.AddDate(0, 0, -6)

	t15, err := time.Parse("020115:04", "220611:00")
	if err != nil {
		t.Fatal(err)
	}
	t15 = t15.AddDate(time.Now().Year(), 0, 0)

	t16, err := time.Parse("020115:04", "220611:10")
	if err != nil {
		t.Fatal(err)
	}
	t16 = t16.AddDate(time.Now().Year(), 0, 0)

	cmdOutput := []byte(testDiskConsumptionsData)

	expected := []model.DiskConsumption{
		{
			TimeStart:      &t1,
			TimeEnd:        &currentTime,
			IopsHostDayAvg: getPointerToFloat(88.81),
			IombHostDayAvg: getPointerToFloat(0.40),
		},
		{

			TimeStart:      &t2,
			TimeEnd:        &currentTime,
			IopsHostDayAvg: getPointerToFloat(80.1),
			IombHostDayAvg: getPointerToFloat(0.40),
		},
		{
			TimeStart:      &t3,
			TimeEnd:        &t4,
			IopsHostDayAvg: getPointerToFloat(75.28),
			IombHostDayAvg: getPointerToFloat(0.38),
		},
		{
			TimeStart:      &t5,
			TimeEnd:        &t6,
			IopsHostDayAvg: getPointerToFloat(75.26),
			IombHostDayAvg: getPointerToFloat(0.38),
		},
		{
			TimeStart:      &t7,
			TimeEnd:        &t8,
			IopsHostDayAvg: getPointerToFloat(127.18),
			IombHostDayAvg: getPointerToFloat(0.42),
		},
		{
			TimeStart:      &currentTime,
			TimeEnd:        nil,
			IopsHostDayAvg: getPointerToFloat(81.89),
			IombHostDayAvg: getPointerToFloat(0.42),
		},
		{
			TimeStart:      &t9,
			TimeEnd:        &currentTime,
			IopsHostDayAvg: getPointerToFloat(85.04),
			IombHostDayAvg: getPointerToFloat(0.42),
		},
		{
			TimeStart:      &t10,
			TimeEnd:        &t9,
			IopsHostDayAvg: getPointerToFloat(82.48),
			IombHostDayAvg: getPointerToFloat(0.41),
		},
		{
			TimeStart:      &t11,
			TimeEnd:        &t10,
			IopsHostDayAvg: getPointerToFloat(77.95),
			IombHostDayAvg: getPointerToFloat(0.39),
		},
		{
			TimeStart:      &t12,
			TimeEnd:        &t11,
			IopsHostDayAvg: getPointerToFloat(76.94),
			IombHostDayAvg: getPointerToFloat(0.39),
		},
		{
			TimeStart:      &t13,
			TimeEnd:        &t12,
			IopsHostDayAvg: getPointerToFloat(77.15),
			IombHostDayAvg: getPointerToFloat(0.39),
		},
		{
			TimeStart:      &t14,
			TimeEnd:        &t13,
			IopsHostDayAvg: getPointerToFloat(79.26),
			IombHostDayAvg: getPointerToFloat(0.39),
		},
		{
			TimeStart:      &t15,
			TimeEnd:        nil,
			IopsHostDayAvg: getPointerToFloat(78.54),
			IombHostDayAvg: getPointerToFloat(0.38),
		},

		{
			TimeStart:      &t16,
			TimeEnd:        nil,
			IopsHostDayAvg: getPointerToFloat(94.36),
			IombHostDayAvg: nil,
		},
	}

	expectedJSON, err := json.Marshal(expected)
	if err != nil {
		t.Fatal(err)
	}

	actual, err := DiskConsumption(cmdOutput)
	assert.Nil(t, err)

	actualJSON, err := json.Marshal(actual)
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEq(t, string(expectedJSON), string(actualJSON))
}
