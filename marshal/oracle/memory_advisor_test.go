// Copyright (c) 2024 Sorint.lab S.p.A.
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
	"github.com/stretchr/testify/require"
)

const (
	testMemoryAdvisorData01 = `
BEGINOUTPUT
MEMORY_SIZE_LOWER_GB|||1.078
ENDOUTPUT
`
	testMemoryAdvisorData02 = `
BEGINOUTPUT
PGA_TARGET_AGGREGATE_LOWER_GB|||7.5
SGA_SIZE_LOWER_GB|||N/A
ENDOUTPUT
	`
)

func TestMemoryAdvisor_Output01(t *testing.T) {
	expected := &model.OracleDatabaseMemoryAdvisor{
		AutomaticMemoryManagement: true,
		MemorySizeLowerGb:         "1.078",
	}

	actual, err := MemoryAdvisor([]byte(testMemoryAdvisorData01))
	require.NoError(t, err)

	assert.Equal(t, expected, actual)
}

func TestMemoryAdvisor_Output02(t *testing.T) {
	expected := &model.OracleDatabaseMemoryAdvisor{
		AutomaticMemoryManagement: false,
		PgaTargetAggregateLowerGb: "7.5",
		SgaSizeLowerGb:            "N/A",
	}

	actual, err := MemoryAdvisor([]byte(testMemoryAdvisorData02))
	require.NoError(t, err)

	assert.Equal(t, expected, actual)
}
