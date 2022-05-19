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

const testGrantDbaData = `test#01|||yes|||no
test#02|||yes|||no
test#03|||yes|||no
test#04|||yes|||no`

func TestGrantDba(t *testing.T) {
	expected := []model.OracleGrantDba{
		{
			Grantee: "test#01", AdminOption: "yes", DefaultRole: "no",
		},
		{
			Grantee: "test#02", AdminOption: "yes", DefaultRole: "no",
		},
		{
			Grantee: "test#03", AdminOption: "yes", DefaultRole: "no",
		},
		{
			Grantee: "test#04", AdminOption: "yes", DefaultRole: "no",
		},
	}

	actual := GrantDba([]byte(testGrantDbaData))

	assert.Equal(t, expected, actual)
}
