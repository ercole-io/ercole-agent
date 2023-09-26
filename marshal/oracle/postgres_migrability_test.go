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

	"github.com/ercole-io/ercole/v2/model"
	"github.com/stretchr/testify/assert"
)

const testPostgresMigrabilityData = `PLSQL LINES           |||      6212
PARTITIONED TABLES    |||         0
PARTITIONED INDEXES   |||         0
PROCEDURES            |||         1
SEQUENCES             |||         1
TRIGGERS              |||         0
FUNCTIONS             |||         0
HCC COMPRESSED TABLES |||         0
MVIEWS REWRITE ENABLED|||         0
VDP POLICIES          |||         0
PERFSTAT                      |||PACKAGE     |||       348
PERFSTAT                      |||PACKAGE BODY|||      5817
PERFSTAT                      |||PROCEDURE   |||        47`

func TestPostgresMigrability(t *testing.T) {

	s0 := "PLSQL LINES"
	s1 := "PARTITIONED TABLES"
	s2 := "PARTITIONED INDEXES"
	s3 := "PROCEDURES"
	s4 := "SEQUENCES"
	s5 := "TRIGGERS"
	s6 := "FUNCTIONS"
	s7 := "HCC COMPRESSED TABLES"
	s8 := "MVIEWS REWRITE ENABLED"
	s9 := "VDP POLICIES"
	s10 := "PERFSTAT"
	s11 := "PERFSTAT"
	s12 := "PERFSTAT"
	o0:= "PACKAGE"
	o1:= "PACKAGE BODY"
	o2:= "PROCEDURE"

	expected := []model.PgsqlMigrability{
		{
			Metric: &s0,
			Count:  6212,
		},
		{
			Metric: &s1,
			Count:  0,
		},
		{
			Metric: &s2,
			Count:  0,
		},
		{
			Metric: &s3,
			Count:  1,
		},
		{
			Metric: &s4,
			Count:  1,
		},
		{
			Metric: &s5,
			Count:  0,
		},
		{
			Metric: &s6,
			Count:  0,
		},
		{
			Metric: &s7,
			Count:  0,
		},
		{
			Metric: &s8,
			Count:  0,
		},
		{
			Metric: &s9,
			Count:  0,
		},

		{
			Schema: &s10,
			ObjectType: &o0,
			Count:  348,
		},
		{
			Schema: &s11,
			ObjectType: &o1,
			Count:  5817,
		},
		{
			Schema: &s12,
			ObjectType: &o2,
			Count:  47,
		},
	}

	expectedJSON, err := json.Marshal(expected)
	if err != nil {
		t.Fatal(err)
	}

	actual, err := PostgresMigrability([]byte(testPostgresMigrabilityData))
	assert.Nil(t, err)

	actualJSON, err := json.Marshal(actual)
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEq(t, string(expectedJSON), string(actualJSON))
}
