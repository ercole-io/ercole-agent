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

	"github.com/stretchr/testify/assert"
)

const (
	testPoliciesAuditData0 = `

`
	testPoliciesAuditData1 = `SYS_LOGON_LOGOFF
ORA_SECURECONFIG
ORA_LOGON_FAILURES`
)

func TestPoliciesAudit_Empty(t *testing.T) {
	actual, err := PoliciesAudit([]byte(testPoliciesAuditData0))

	assert.Nil(t, err)
	assert.Nil(t, actual)
}

func TestPoliciesAudit_Success(t *testing.T) {
	actual, err := PoliciesAudit([]byte(testPoliciesAuditData1))

	expected := []string{"SYS_LOGON_LOGOFF", "ORA_SECURECONFIG", "ORA_LOGON_FAILURES"}

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}
