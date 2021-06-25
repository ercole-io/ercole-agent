// Copyright (c) 2020 Sorint.lab S.p.A.
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
	"testing"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/stretchr/testify/assert"
)

const testExadataComponentData string = `fcax1hf1|||DBServer|||X8-2|||19.2.7.0.0.191012|||4/96|||766|||online|||2/2|||normal|||16/16|||normal|||22.0|||normal|||-|||running|||running|||-
fcax1sf1|||StorageServer|||X8-2L_High_Capacity|||19.2.7.0.0.191012|||64/64|||188|||online|||2/2|||normal|||8/8|||normal|||20.0|||normal|||running|||running|||running|||WriteBack
fcax1bb0|||IBSwitch|||SUN_DCS_36p|||2.2.13-2.190326|||-|||-|||-|||-|||-|||-|||-|||-|||-|||-|||-|||-|||-`

func TestExadataComponent(t *testing.T) {
	v2 := 2
	v4 := 4
	v8 := 8
	v16 := 16
	v20_0 := 20.0
	v22_0 := 22.0
	v64 := 64
	v96 := 96
	v188 := 188
	v766 := 766

	online := "online"
	normal := "normal"
	running := "running"
	writeBack := "WriteBack"

	expected := []model.OracleExadataComponent{
		{
			Hostname:             "fcax1hf1",
			ServerType:           "DBServer",
			Model:                "X8-2",
			SwVersion:            "19.2.7.0.0.191012",
			SwReleaseDate:        "191012",
			RunningCPUCount:      &v4,
			TotalCPUCount:        &v96,
			Memory:               &v766,
			Status:               &online,
			RunningPowerSupply:   &v2,
			TotalPowerSupply:     &v2,
			PowerStatus:          &normal,
			RunningFanCount:      &v16,
			TotalFanCount:        &v16,
			FanStatus:            &normal,
			TempActual:           &v22_0,
			TempStatus:           &normal,
			CellsrvServiceStatus: nil,
			MsServiceStatus:      &running,
			RsServiceStatus:      &running,
			FlashcacheMode:       nil,
		},
		{
			Hostname:             "fcax1sf1",
			ServerType:           "StorageServer",
			Model:                "X8-2L_High_Capacity",
			SwVersion:            "19.2.7.0.0.191012",
			SwReleaseDate:        "191012",
			RunningCPUCount:      &v64,
			TotalCPUCount:        &v64,
			Memory:               &v188,
			Status:               &online,
			RunningPowerSupply:   &v2,
			TotalPowerSupply:     &v2,
			PowerStatus:          &normal,
			RunningFanCount:      &v8,
			TotalFanCount:        &v8,
			FanStatus:            &normal,
			TempActual:           &v20_0,
			TempStatus:           &normal,
			CellsrvServiceStatus: &running,
			MsServiceStatus:      &running,
			RsServiceStatus:      &running,
			FlashcacheMode:       &writeBack,
		},
		{
			Hostname:             "fcax1bb0",
			ServerType:           "IBSwitch",
			Model:                "SUN_DCS_36p",
			SwVersion:            "2.2.13-2.190326",
			SwReleaseDate:        "190326",
			RunningCPUCount:      nil,
			TotalCPUCount:        nil,
			Memory:               nil,
			Status:               nil,
			RunningPowerSupply:   nil,
			TotalPowerSupply:     nil,
			PowerStatus:          nil,
			RunningFanCount:      nil,
			TotalFanCount:        nil,
			FanStatus:            nil,
			TempActual:           nil,
			TempStatus:           nil,
			CellsrvServiceStatus: nil,
			MsServiceStatus:      nil,
			RsServiceStatus:      nil,
			FlashcacheMode:       nil,
		},
	}

	cmdOutput := []byte(testExadataComponentData)
	actual, errs := ExadataComponent(cmdOutput)
	assert.Nil(t, errs)
	assert.Equal(t, expected, actual)
}
