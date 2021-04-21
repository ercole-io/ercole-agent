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

package common

import (
	"testing"

	"github.com/ercole-io/ercole-agent/v2/agentmodel"
	"github.com/ercole-io/ercole-agent/v2/config"
	"github.com/ercole-io/ercole-agent/v2/logger"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:generate mockgen -source ../../fetcher/fetcher.go -destination=fake_fetcher_test.go -package=common

func TestGetUnlistedRunningOracleDBs(t *testing.T) {
	ctrl := gomock.NewController(t)
	fakeFetcher := NewMockFetcher(ctrl)
	fakeFetcher.
		EXPECT().GetOracleDatabaseRunningDatabases().
		Return([]string{"pippo", "topolino", "pluto"})

	b := CommonBuilder{
		fetcher:       fakeFetcher,
		configuration: config.Configuration{},
		log:           nil,
	}

	listedDBs := []agentmodel.OratabEntry{
		{DBName: "pippo"},
		{DBName: "pluto"},
	}

	expected := []string{"topolino"}
	actual := b.getUnlistedRunningOracleDBs(listedDBs)

	assert.Equal(t, expected, actual)
}

func TestGetUnlistedRunningOracleDBs2(t *testing.T) {
	ctrl := gomock.NewController(t)
	fakeFetcher := NewMockFetcher(ctrl)
	fakeFetcher.
		EXPECT().GetOracleDatabaseRunningDatabases().
		Return([]string{"ERC002", "ERC001"})

	log, err := logger.NewLogger("TEST")
	require.Nil(t, err)

	b := CommonBuilder{
		fetcher:       fakeFetcher,
		configuration: config.Configuration{},
		log:           log,
	}

	listedDBs := []agentmodel.OratabEntry{
		{
			DBName:     "ERC002",
			OracleHome: "",
		},
	}

	expected := []string{"ERC001"}
	actual := b.getUnlistedRunningOracleDBs(listedDBs)

	assert.Equal(t, expected, actual)
}
