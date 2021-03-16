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

package mysql

import (
	"testing"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/stretchr/testify/assert"
)

var testInstanceData = `mysql: [Warning] Using a password on the command line interface can be insecure.
"erclinmysql:3306";"8.0.23";"COMMUNITY";"Linux";"x86_64";"InnoDB";"ON";"utf8mb4";"utf8";"16.0000";"0";"128.00000000";"16.00000000";"1.00000000";"0";"1"
`

func TestInstance(t *testing.T) {
	cmdOutput := []byte(testInstanceData)

	actual := Instance(cmdOutput)

	expected := &model.MySQLInstance{
		Name:               "erclinmysql:3306",
		Version:            "8.0.23",
		Edition:            "COMMUNITY",
		Platform:           "Linux",
		Architecture:       "x86_64",
		Engine:             "InnoDB",
		RedoLogEnabled:     "ON",
		CharsetServer:      "utf8mb4",
		CharsetSystem:      "utf8",
		PageSize:           16,
		ThreadsConcurrency: 0,
		BufferPoolSize:     128,
		LogBufferSize:      16,
		SortBufferSize:     1,
		ReadOnly:           false,
		LogBin:             true,
		Databases:          nil,
	}

	assert.Equal(t, expected, actual)
}
