// Copyright (c) 2022 Sorint.lab S.p.A.
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

package postgresql

import (
	"testing"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/stretchr/testify/assert"
)

func TestSetting(t *testing.T) {
	testSettingData := `PostgreSQL 10.20|4194304|off|(disabled)|83886080|1073741824|100|50%|100|4|67108864|16777216|4294967296|1|8|8`
	cmdOutput := []byte(testSettingData)
	actual, err := Setting(cmdOutput)

	expected := model.PostgreSQLSetting{
		DbVersion:                  "PostgreSQL 10.20",
		WorkMem:                    4194304,
		ArchiveMode:                false,
		ArchiveCommand:             "(disabled)",
		MinWalSize:                 83886080,
		MaxWalSize:                 1073741824,
		MaxConnections:             100,
		CheckpointCompletionTarget: "50%",
		DefaultStatisticsTarget:    100,
		RandomPageCost:             4,
		MaintenanceWorkMem:         67108864,
		SharedBuffers:              16777216,
		EffectiveCacheSize:         4294967296,
		EffectiveIoConcurrency:     1,
		MaxWorkerProcesses:         8,
		MaxParallelWorkers:         8,
	}

	assert.Equal(t, &expected, actual)
	assert.Nil(t, err)
}
