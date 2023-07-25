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

package oracle

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/stretchr/testify/assert"
)

const testCpuDiskConsumptionsPdbData = `BEGINOUTPUT
N/A|||N/A|||N/A|||N/A|||N/A
1.02|||5.37|||0|||30.11|||73130.31
N/A|||N/A|||N/A|||N/A|||N/A
N/A|||N/A|||N/A|||N/A|||N/A
N/A|||N/A|||N/A|||N/A|||N/A
1.03|||2.7|||0|||17.28|||243.13
.56|||2.44|||0|||71.33|||73130.31
.82|||2.61|||0|||14.15|||1694.79
.89|||2.83|||0|||22.82|||1589.23
.94|||2.39|||0|||22.91|||891.82
1.02|||2.69|||0|||25.35|||1608.81
1.84|||5.37|||0|||29.08|||553.31
160709:59|||.28|||.76|||N/A|||8.82|||160.42
ENDOUTPUT

PL/SQL procedure successfully completed.
`

func float64ToPointer(f float64) *float64 {
	return &f
}

func TestCpuDiskConsumptionsPdb(t *testing.T) {
	currentTime = time.Now()
	t2 := currentTime.AddDate(0, 0, -7)
	t9 := currentTime.AddDate(0, 0, -1)
	t10 := currentTime.AddDate(0, 0, -2)
	t11 := currentTime.AddDate(0, 0, -3)
	t12 := currentTime.AddDate(0, 0, -4)
	t13 := currentTime.AddDate(0, 0, -5)
	t14 := currentTime.AddDate(0, 0, -6)

	t16, err := time.Parse("020115:04", "160709:59")
	if err != nil {
		t.Fatal(err)
	}
	t16 = t16.AddDate(time.Now().Year(), 0, 0)

	expected := []model.CpuDiskConsumptionPdb{
		{
			TimeStart: &t2,
			TimeEnd:   &currentTime,
			Target:    "w1",
			CpuDbAvg:  float64ToPointer(1.02),
			CpuDbMax:  float64ToPointer(5.37),
			IopsAvg:   float64ToPointer(0),
			IombAvg:   float64ToPointer(30.11),
			IombMax:   float64ToPointer(73130.31),
		},
		{
			TimeStart: &currentTime,
			TimeEnd:   nil,
			Target:    "d1",
			CpuDbAvg:  float64ToPointer(1.03),
			CpuDbMax:  float64ToPointer(2.7),
			IopsAvg:   float64ToPointer(0),
			IombAvg:   float64ToPointer(17.28),
			IombMax:   float64ToPointer(243.13),
		},
		{
			TimeStart: &t9,
			TimeEnd:   &currentTime,
			Target:    "d2",
			CpuDbAvg:  float64ToPointer(0.56),
			CpuDbMax:  float64ToPointer(2.44),
			IopsAvg:   float64ToPointer(0),
			IombAvg:   float64ToPointer(71.33),
			IombMax:   float64ToPointer(73130.31),
		},
		{
			TimeStart: &t10,
			TimeEnd:   &t9,
			Target:    "d3",
			CpuDbAvg:  float64ToPointer(0.82),
			CpuDbMax:  float64ToPointer(2.61),
			IopsAvg:   float64ToPointer(0),
			IombAvg:   float64ToPointer(14.15),
			IombMax:   float64ToPointer(1694.79),
		},
		{
			TimeStart: &t11,
			TimeEnd:   &t10,
			Target:    "d4",
			CpuDbAvg:  float64ToPointer(0.89),
			CpuDbMax:  float64ToPointer(2.83),
			IopsAvg:   float64ToPointer(0),
			IombAvg:   float64ToPointer(22.82),
			IombMax:   float64ToPointer(1589.23),
		},
		{
			TimeStart: &t12,
			TimeEnd:   &t11,
			Target:    "d5",
			CpuDbAvg:  float64ToPointer(0.94),
			CpuDbMax:  float64ToPointer(2.39),
			IopsAvg:   float64ToPointer(0),
			IombAvg:   float64ToPointer(22.91),
			IombMax:   float64ToPointer(891.82),
		},
		{
			TimeStart: &t13,
			TimeEnd:   &t12,
			Target:    "d6",
			CpuDbAvg:  float64ToPointer(1.02),
			CpuDbMax:  float64ToPointer(2.69),
			IopsAvg:   float64ToPointer(0),
			IombAvg:   float64ToPointer(25.35),
			IombMax:   float64ToPointer(1608.81),
		},
		{
			TimeStart: &t14,
			TimeEnd:   &t13,
			Target:    "d7",
			CpuDbAvg:  float64ToPointer(1.84),
			CpuDbMax:  float64ToPointer(5.37),
			IopsAvg:   float64ToPointer(0),
			IombAvg:   float64ToPointer(29.08),
			IombMax:   float64ToPointer(553.31),
		},
		{
			TimeStart: &t16,
			TimeEnd:   nil,
			CpuDbAvg:  float64ToPointer(0.28),
			CpuDbMax:  float64ToPointer(0.76),
			IopsAvg:   nil,
			IombAvg:   float64ToPointer(8.82),
			IombMax:   float64ToPointer(160.42),
		},
	}

	expectedJSON, err := json.Marshal(expected)
	if err != nil {
		t.Fatal(err)
	}

	actual, err := CpuDiskConsumptionsPdb([]byte(testCpuDiskConsumptionsPdbData))
	assert.Nil(t, err)

	actualJSON, err := json.Marshal(actual)
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEq(t, string(expectedJSON), string(actualJSON))
}
