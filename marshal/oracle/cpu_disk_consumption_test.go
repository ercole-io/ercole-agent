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

const testCpuDiskConsumptionData = `BEGINOUTPUT
N/A|||N/A|||N/A|||N/A|||N/A|||N/A|||N/A|||N/A
.21|||4.41|||6.59|||26.65|||407.68|||19431.85|||164.88|||11712.28
.21|||4.41|||6.59|||26.65|||407.68|||19431.85|||164.88|||11712.28
.08|||.75|||7.21|||17.34|||132.72|||4648.49|||36.97|||3960.1
.21|||4.41|||6.59|||26.65|||407.68|||19431.85|||164.88|||11712.28
.21|||4.41|||6.59|||26.65|||407.68|||19431.85|||164.88|||11712.28
.21|||2.87|||7.57|||25.93|||378.53|||19431.85|||229.71|||11712.28
.33|||4.41|||8.06|||26.65|||642.6|||14519.22|||153.15|||3776.37
.3|||2.2|||7.68|||18.44|||447.7|||8938.07|||203.64|||8453.28
.31|||3.21|||7.4|||17.19|||640.92|||14526.36|||199.18|||8370.3
.08|||2.14|||4.31|||16.03|||292.56|||7197.71|||204.65|||4892.83
.08|||2.39|||4.3|||15.38|||146.84|||8886.45|||46.91|||2808.41	
150609:59|||.34|||.64|||9.04|||13.35|||361.47|||4334.83|||151.05|||3732.88
150610:59|||.38|||.75|||10.89|||16.1|||746.41|||8935.51|||292.23|||5821.71
150611:59|||.43|||.7|||10.84|||16.21|||449.94|||5855.83|||162.22|||3757.54
ENDOUTPUT`

func TestCpuDiskConsumptions(t *testing.T) {
	currentTime = time.Now()
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
	// t15 := currentTime.AddDate(0, 0, -7)

	t16, err := time.Parse("020115:04", "150609:59")
	if err != nil {
		t.Fatal(err)
	}
	t16 = t16.AddDate(time.Now().Year(), 0, 0)

	t17, err := time.Parse("020115:04", "150610:59")
	if err != nil {
		t.Fatal(err)
	}
	t17 = t17.AddDate(time.Now().Year(), 0, 0)

	t18, err := time.Parse("020115:04", "150611:59")
	if err != nil {
		t.Fatal(err)
	}
	t18 = t18.AddDate(time.Now().Year(), 0, 0)

	expected := []model.CpuDiskConsumption{
		{
			TimeStart:  &t2,
			TimeEnd:    &currentTime,
			Target:     "w4",
			CpuDbAvg:   float64ToPointer(0.21),
			CpuDbMax:   float64ToPointer(4.41),
			CpuHostAvg: float64ToPointer(6.59),
			CpuHostMax: float64ToPointer(26.65),
			IopsAvg:    float64ToPointer(407.68),
			IopsMax:    float64ToPointer(19431.85),
			IombAvg:    float64ToPointer(164.88),
			IombMax:    float64ToPointer(11712.28),
		},
		{
			TimeStart:  &t3,
			TimeEnd:    &t4,
			Target:     "w3",
			CpuDbAvg:   float64ToPointer(0.21),
			CpuDbMax:   float64ToPointer(4.41),
			CpuHostAvg: float64ToPointer(6.59),
			CpuHostMax: float64ToPointer(26.65),
			IopsAvg:    float64ToPointer(407.68),
			IopsMax:    float64ToPointer(19431.85),
			IombAvg:    float64ToPointer(164.88),
			IombMax:    float64ToPointer(11712.28),
		},
		{
			TimeStart:  &t5,
			TimeEnd:    &t6,
			Target:     "w2",
			CpuDbAvg:   float64ToPointer(0.08),
			CpuDbMax:   float64ToPointer(0.75),
			CpuHostAvg: float64ToPointer(7.21),
			CpuHostMax: float64ToPointer(17.34),
			IopsAvg:    float64ToPointer(132.72),
			IopsMax:    float64ToPointer(4648.49),
			IombAvg:    float64ToPointer(36.97),
			IombMax:    float64ToPointer(3960.1),
		},
		{
			TimeStart:  &t7,
			TimeEnd:    &t8,
			Target:     "w1",
			CpuDbAvg:   float64ToPointer(0.21),
			CpuDbMax:   float64ToPointer(4.41),
			CpuHostAvg: float64ToPointer(6.59),
			CpuHostMax: float64ToPointer(26.65),
			IopsAvg:    float64ToPointer(407.68),
			IopsMax:    float64ToPointer(19431.85),
			IombAvg:    float64ToPointer(164.88),
			IombMax:    float64ToPointer(11712.28),
		},
		{
			TimeStart:  &currentTime,
			TimeEnd:    nil,
			Target:     "d7",
			CpuDbAvg:   float64ToPointer(0.21),
			CpuDbMax:   float64ToPointer(4.41),
			CpuHostAvg: float64ToPointer(6.59),
			CpuHostMax: float64ToPointer(26.65),
			IopsAvg:    float64ToPointer(407.68),
			IopsMax:    float64ToPointer(19431.85),
			IombAvg:    float64ToPointer(164.88),
			IombMax:    float64ToPointer(11712.28),
		},
		{
			TimeStart:  &t9,
			TimeEnd:    &currentTime,
			Target:     "d6",
			CpuDbAvg:   float64ToPointer(0.21),
			CpuDbMax:   float64ToPointer(2.87),
			CpuHostAvg: float64ToPointer(7.57),
			CpuHostMax: float64ToPointer(25.93),
			IopsAvg:    float64ToPointer(378.53),
			IopsMax:    float64ToPointer(19431.85),
			IombAvg:    float64ToPointer(229.71),
			IombMax:    float64ToPointer(11712.28),
		},
		{
			TimeStart:  &t10,
			TimeEnd:    &t9,
			Target:     "d5",
			CpuDbAvg:   float64ToPointer(0.33),
			CpuDbMax:   float64ToPointer(4.41),
			CpuHostAvg: float64ToPointer(8.06),
			CpuHostMax: float64ToPointer(26.65),
			IopsAvg:    float64ToPointer(642.6),
			IopsMax:    float64ToPointer(14519.22),
			IombAvg:    float64ToPointer(153.15),
			IombMax:    float64ToPointer(3776.37),
		},
		{
			TimeStart:  &t11,
			TimeEnd:    &t10,
			Target:     "d4",
			CpuDbAvg:   float64ToPointer(0.3),
			CpuDbMax:   float64ToPointer(2.2),
			CpuHostAvg: float64ToPointer(7.68),
			CpuHostMax: float64ToPointer(18.44),
			IopsAvg:    float64ToPointer(447.7),
			IopsMax:    float64ToPointer(8938.07),
			IombAvg:    float64ToPointer(203.64),
			IombMax:    float64ToPointer(8453.28),
		},
		{
			TimeStart:  &t12,
			TimeEnd:    &t11,
			Target:     "d3",
			CpuDbAvg:   float64ToPointer(0.31),
			CpuDbMax:   float64ToPointer(3.21),
			CpuHostAvg: float64ToPointer(7.4),
			CpuHostMax: float64ToPointer(17.19),
			IopsAvg:    float64ToPointer(640.92),
			IopsMax:    float64ToPointer(14526.36),
			IombAvg:    float64ToPointer(199.18),
			IombMax:    float64ToPointer(8370.3),
		},
		{
			TimeStart:  &t13,
			TimeEnd:    &t12,
			Target:     "d2",
			CpuDbAvg:   float64ToPointer(0.08),
			CpuDbMax:   float64ToPointer(2.14),
			CpuHostAvg: float64ToPointer(4.31),
			CpuHostMax: float64ToPointer(16.03),
			IopsAvg:    float64ToPointer(292.56),
			IopsMax:    float64ToPointer(7197.71),
			IombAvg:    float64ToPointer(204.65),
			IombMax:    float64ToPointer(4892.83),
		},
		{
			TimeStart:  &t14,
			TimeEnd:    &t13,
			Target:     "d1",
			CpuDbAvg:   float64ToPointer(0.08),
			CpuDbMax:   float64ToPointer(2.39),
			CpuHostAvg: float64ToPointer(4.3),
			CpuHostMax: float64ToPointer(15.38),
			IopsAvg:    float64ToPointer(146.84),
			IopsMax:    float64ToPointer(8886.45),
			IombAvg:    float64ToPointer(46.91),
			IombMax:    float64ToPointer(2808.41),
		},
		{
			TimeStart:  &t16,
			TimeEnd:    nil,
			CpuDbAvg:   float64ToPointer(0.34),
			CpuDbMax:   float64ToPointer(0.64),
			CpuHostAvg: float64ToPointer(9.04),
			CpuHostMax: float64ToPointer(13.35),
			IopsAvg:    float64ToPointer(361.47),
			IopsMax:    float64ToPointer(4334.83),
			IombAvg:    float64ToPointer(151.05),
			IombMax:    float64ToPointer(3732.88),
		},
		{
			TimeStart:  &t17,
			TimeEnd:    nil,
			CpuDbAvg:   float64ToPointer(0.38),
			CpuDbMax:   float64ToPointer(0.75),
			CpuHostAvg: float64ToPointer(10.89),
			CpuHostMax: float64ToPointer(16.1),
			IopsAvg:    float64ToPointer(746.41),
			IopsMax:    float64ToPointer(8935.51),
			IombAvg:    float64ToPointer(292.23),
			IombMax:    float64ToPointer(5821.71),
		},
		{
			TimeStart:  &t18,
			TimeEnd:    nil,
			CpuDbAvg:   float64ToPointer(0.43),
			CpuDbMax:   float64ToPointer(0.7),
			CpuHostAvg: float64ToPointer(10.84),
			CpuHostMax: float64ToPointer(16.21),
			IopsAvg:    float64ToPointer(449.94),
			IopsMax:    float64ToPointer(5855.83),
			IombAvg:    float64ToPointer(162.22),
			IombMax:    float64ToPointer(3757.54),
		},
	}

	expectedJSON, err := json.Marshal(expected)
	if err != nil {
		t.Fatal(err)
	}

	actual, err := CpuDiskConsumptions([]byte(testCpuDiskConsumptionData))
	assert.Nil(t, err)

	actualJSON, err := json.Marshal(actual)
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEq(t, string(expectedJSON), string(actualJSON))
}
