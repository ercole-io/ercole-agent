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
	"testing"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/stretchr/testify/assert"
)

const testExadataComponentData string = `HOST_TYPE|||RACK_ID|||HOSTNAME|||HOST_ID|||CPU_ENABLED|||CPU_TOT|||MEMORY_GB|||IMAGEVERSION|||KERNEL|||MODEL|||FAN_USED|||FAN_TOTAL|||PSU_USED|||PSU_TOTAL|||MS_STATUS|||RS_STATUS
KVM_HOST|||AK00467954|||exafakedatav01|||2043XCB05P|||96|||96|||1510|||21.2.12.0.0.220513|||4.14.35-2047.511.5.5.1.el7uek.x86_64|||Exadata Database Machine X8M-2|||16|||16|||2|||2|||running|||running`

func TestExadataComponents(t *testing.T) {
	expected := []model.OracleExadataComponent{
		{
			HostType:     "KVM_HOST",
			RackID:       "AK00467954",
			Hostname:     "exafakedatav01",
			HostID: "2043XCB05P",
			CPUEnabled:   96,
			TotalCPU:     96,
			Memory:       1510,
			ImageVersion: "21.2.12.0.0.220513",
			Kernel:       "4.14.35-2047.511.5.5.1.el7uek.x86_64",
			Model:        "Exadata Database Machine X8M-2",
			FanUsed:      16,
			FanTotal:     16,
			PsuUsed:      2,
			PsuTotal:     2,
			MsStatus:     "running",
			RsStatus:     "running",
		},
	}

	cmdOutput := []byte(testExadataComponentData)
	actual, errs := ExadataComponents(cmdOutput)
	assert.Nil(t, errs)
	assert.Equal(t, expected, actual)
}
