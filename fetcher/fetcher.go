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

	// Oracle related getters
	GetOratabEntries() []agentmodel.OratabEntry
	GetDbStatus(entry agentmodel.OratabEntry) string
	GetMountedDb(entry agentmodel.OratabEntry) model.OracleDatabase
	GetDbVersion(entry agentmodel.OratabEntry) string
	RunStats(entry agentmodel.OratabEntry)
	GetOpenDb(entry agentmodel.OratabEntry) model.OracleDatabase
	GetTablespaces(entry agentmodel.OratabEntry) []model.OracleDatabaseTablespace
	GetSchemas(entry agentmodel.OratabEntry) []model.OracleDatabaseSchema
	GetPatches(entry agentmodel.OratabEntry, dbVersion string) []model.OracleDatabasePatch
	GetDatabaseFeatureUsageStat(entry agentmodel.OratabEntry, dbVersion string) []model.OracleDatabaseFeatureUsageStat
	GetLicenses(entry agentmodel.OratabEntry, dbVersion, hardwareAbstractionTechnology string) []model.OracleDatabaseLicense
	GetADDMs(entry agentmodel.OratabEntry) []model.OracleDatabaseAddm
	GetSegmentAdvisors(entry agentmodel.OratabEntry) []model.OracleDatabaseSegmentAdvisor
	GetPSUs(entry agentmodel.OratabEntry, dbVersion string) []model.OracleDatabasePSU
	GetBackups(entry agentmodel.OratabEntry) []model.OracleDatabaseBackup
	GetExadataComponents() []model.OracleExadataComponent
	GetOracleExadataCellDisks() map[agentmodel.StorageServerName][]model.OracleExadataCellDisk
}

// User struct
type User struct {
	Name     string
	UID, GID uint32
}
