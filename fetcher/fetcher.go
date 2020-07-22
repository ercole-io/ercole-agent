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

	// Operating system fetchers
	GetHost() model.Host
	GetFilesystems() []model.Filesystem
	GetClustersMembershipStatus() model.ClusterMembershipStatus

	// Virtualization fetcher
	GetClusters(hv config.Hypervisor) []model.ClusterInfo
	GetVirtualMachines(hv config.Hypervisor) map[string][]model.VMInfo

	// Oracle/Database fetchers
	GetOracleDatabaseOratabEntries() []agentmodel.OratabEntry
	GetOracleDatabaseDbStatus(entry agentmodel.OratabEntry) string
	GetOracleDatabaseMountedDb(entry agentmodel.OratabEntry) model.OracleDatabase
	GetOracleDatabaseDbVersion(entry agentmodel.OratabEntry) string
	RunOracleDatabaseStats(entry agentmodel.OratabEntry)
	GetOracleDatabaseOpenDb(entry agentmodel.OratabEntry) model.OracleDatabase
	GetOracleDatabaseTablespaces(entry agentmodel.OratabEntry) []model.OracleDatabaseTablespace
	GetOracleDatabaseSchemas(entry agentmodel.OratabEntry) []model.OracleDatabaseSchema
	GetOracleDatabasePatches(entry agentmodel.OratabEntry, dbVersion string) []model.OracleDatabasePatch
	GetOracleDatabaseFeatureUsageStat(entry agentmodel.OratabEntry, dbVersion string) []model.OracleDatabaseFeatureUsageStat
	GetOracleDatabaseLicenses(entry agentmodel.OratabEntry, dbVersion, hardwareAbstractionTechnology string) []model.OracleDatabaseLicense
	GetOracleDatabaseADDMs(entry agentmodel.OratabEntry) []model.OracleDatabaseAddm
	GetOracleDatabaseSegmentAdvisors(entry agentmodel.OratabEntry) []model.OracleDatabaseSegmentAdvisor
	GetOracleDatabasePSUs(entry agentmodel.OratabEntry, dbVersion string) []model.OracleDatabasePSU
	GetOracleDatabaseBackups(entry agentmodel.OratabEntry) []model.OracleDatabaseBackup
	// Oracle/Exadata fetchers
	GetOracleExadataComponents() []model.OracleExadataComponent
	GetOracleExadataCellDisks() map[agentmodel.StorageServerName][]model.OracleExadataCellDisk

	// Microsoft/SQLServer fetchers
}

// User struct
type User struct {
	Name     string
	UID, GID uint32
}
