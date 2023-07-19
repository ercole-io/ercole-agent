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

var testConsumptionsData string = `7.27
7.20
6.88
7.29
7.39
10.10
8.27
4.92
4.83
7.68
7.83
8.33
190611:00|||10.93
190611:10|||13.96`

func TestConsumption(t *testing.T) {
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

	t15, err := time.Parse("020115:04", "190611:00")
	if err != nil {
		t.Fatal(err)
	}
	t15 = t15.AddDate(time.Now().Year(), 0, 0)

	t16, err := time.Parse("020115:04", "190611:10")
	if err != nil {
		t.Fatal(err)
	}
	t16 = t16.AddDate(time.Now().Year(), 0, 0)

	cmdOutput := []byte(testConsumptionsData)

	expected := []model.Consumption{
		{
			TimeStart: &t1,
			TimeEnd:   &currentTime,
			CpuAvg:    getPointerToFloat(7.27),
		},
		{

			TimeStart: &t2,
			TimeEnd:   &currentTime,
			CpuAvg:    getPointerToFloat(7.20),
		},
		{
			TimeStart: &t3,
			TimeEnd:   &t4,
			CpuAvg:    getPointerToFloat(6.88),
		},
		{
			TimeStart: &t5,
			TimeEnd:   &t6,
			CpuAvg:    getPointerToFloat(7.29),
		},
		{
			TimeStart: &t7,
			TimeEnd:   &t8,
			CpuAvg:    getPointerToFloat(7.39),
		},
		{
			TimeStart: &currentTime,
			TimeEnd:   nil,
			CpuAvg:    getPointerToFloat(10.10),
		},
		{
			TimeStart: &t9,
			TimeEnd:   &currentTime,
			CpuAvg:    getPointerToFloat(8.27),
		},
		{
			TimeStart: &t10,
			TimeEnd:   &t9,
			CpuAvg:    getPointerToFloat(4.92),
		},
		{
			TimeStart: &t11,
			TimeEnd:   &t10,
			CpuAvg:    getPointerToFloat(4.83),
		},
		{
			TimeStart: &t12,
			TimeEnd:   &t11,
			CpuAvg:    getPointerToFloat(7.68),
		},
		{
			TimeStart: &t13,
			TimeEnd:   &t12,
			CpuAvg:    getPointerToFloat(7.83),
		},
		{
			TimeStart: &t14,
			TimeEnd:   &t13,
			CpuAvg:    getPointerToFloat(8.33),
		},
		{
			TimeStart: &t15,
			TimeEnd:   nil,
			CpuAvg:    getPointerToFloat(10.93),
		},
		{
			TimeStart: &t16,
			TimeEnd:   nil,
			CpuAvg:    getPointerToFloat(13.96),
		},
	}

	expectedJSON, err := json.Marshal(expected)
	if err != nil {
		t.Fatal(err)
	}

	actual, err := Consumption(cmdOutput)
	assert.Nil(t, err)

	actualJSON, err := json.Marshal(actual)
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEq(t, string(expectedJSON), string(actualJSON))
}
