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

package fetcher

import (
	"github.com/ercole-io/ercole-agent/agentmodel"
	"github.com/ercole-io/ercole-agent/config"
	"github.com/ercole-io/ercole/model"
)

// Fetcher interface for Linux and Windows
type Fetcher interface {
	SetUser(username string) error
	SetUserAsCurrent() error

	GetHost() model.Host
	GetFilesystems() []model.Filesystem
	GetClusters(hv config.Hypervisor) []model.ClusterInfo
	GetVirtualMachines(hv config.Hypervisor) map[string][]model.VMInfo
	GetClustersMembershipStatus() model.ClusterMembershipStatus

	// Oracle related functions
	GetOracleOratabEntries() []agentmodel.OratabEntry
	GetOracleDbStatus(entry agentmodel.OratabEntry) string
	GetOracleMountedDb(entry agentmodel.OratabEntry) model.OracleDatabase
	GetOracleDbVersion(entry agentmodel.OratabEntry) string
	RunOracleStats(entry agentmodel.OratabEntry)
	GetOracleOpenDb(entry agentmodel.OratabEntry) model.OracleDatabase
	GetOracleTablespaces(entry agentmodel.OratabEntry) []model.OracleDatabaseTablespace
	GetOracleSchemas(entry agentmodel.OratabEntry) []model.OracleDatabaseSchema
	GetOraclePatches(entry agentmodel.OratabEntry, dbVersion string) []model.OracleDatabasePatch
	GetOracleDatabaseFeatureUsageStat(entry agentmodel.OratabEntry, dbVersion string) []model.OracleDatabaseFeatureUsageStat
	GetOracleLicenses(entry agentmodel.OratabEntry, dbVersion, hardwareAbstractionTechnology string) []model.OracleDatabaseLicense
	GetOracleADDMs(entry agentmodel.OratabEntry) []model.OracleDatabaseAddm
	GetOracleSegmentAdvisors(entry agentmodel.OratabEntry) []model.OracleDatabaseSegmentAdvisor
	GetOraclePSUs(entry agentmodel.OratabEntry, dbVersion string) []model.OracleDatabasePSU
	GetOracleBackups(entry agentmodel.OratabEntry) []model.OracleDatabaseBackup
	GetOracleExadataComponents() []model.OracleExadataComponent
	GetOracleExadataCellDisks() map[agentmodel.StorageServerName][]model.OracleExadataCellDisk
}

// User struct
type User struct {
	Name     string
	UID, GID uint32
}
