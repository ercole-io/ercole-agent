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
	"github.com/ercole-io/ercole-agent/v2/agentmodel"
	"github.com/ercole-io/ercole-agent/v2/config"
	"github.com/ercole-io/ercole/v2/model"
)

// Fetcher interface for Linux and Windows
type Fetcher interface {
	SetUser(username string) error
	SetUserAsCurrent() error

	// Operating system fetchers

	GetHost() (*model.Host, error)
	GetFilesystems() ([]model.Filesystem, error)
	GetClustersMembershipStatus() model.ClusterMembershipStatus

	// Virtualization fetcher

	GetClusters(hv config.Hypervisor) []model.ClusterInfo
	GetVirtualMachines(hv config.Hypervisor) map[string][]model.VMInfo

	// Oracle/Database fetchers

	GetOracleDatabaseOratabEntries() []agentmodel.OratabEntry
	GetOracleDatabaseRunningDatabases() []string
	GetOracleDatabaseDbStatus(entry agentmodel.OratabEntry) string
	GetOracleDatabaseMountedDb(entry agentmodel.OratabEntry) (*model.OracleDatabase, error)
	GetOracleDatabaseDbVersion(entry agentmodel.OratabEntry) string
	RunOracleDatabaseStats(entry agentmodel.OratabEntry)
	GetOracleDatabaseOpenDb(entry agentmodel.OratabEntry) (*model.OracleDatabase, error)
	GetOracleDatabaseTablespaces(entry agentmodel.OratabEntry) []model.OracleDatabaseTablespace
	GetOracleDatabaseSchemas(entry agentmodel.OratabEntry) ([]model.OracleDatabaseSchema, error)
	GetOracleDatabasePatches(entry agentmodel.OratabEntry, dbVersion string) ([]model.OracleDatabasePatch, error)
	GetOracleDatabaseFeatureUsageStat(entry agentmodel.OratabEntry, dbVersion string) ([]model.OracleDatabaseFeatureUsageStat, error)
	GetOracleDatabaseLicenses(entry agentmodel.OratabEntry, dbVersion, hardwareAbstractionTechnology string) []model.OracleDatabaseLicense
	GetOracleDatabaseADDMs(entry agentmodel.OratabEntry) []model.OracleDatabaseAddm
	GetOracleDatabaseSegmentAdvisors(entry agentmodel.OratabEntry) []model.OracleDatabaseSegmentAdvisor
	GetOracleDatabasePSUs(entry agentmodel.OratabEntry, dbVersion string) []model.OracleDatabasePSU
	GetOracleDatabaseBackups(entry agentmodel.OratabEntry) []model.OracleDatabaseBackup
	GetOracleDatabaseCheckPDB(entry agentmodel.OratabEntry) bool
	GetOracleDatabasePDBs(entry agentmodel.OratabEntry) []model.OracleDatabasePluggableDatabase
	GetOracleDatabasePDBTablespaces(entry agentmodel.OratabEntry, pdb string) []model.OracleDatabaseTablespace
	GetOracleDatabasePDBSchemas(entry agentmodel.OratabEntry, pdb string) ([]model.OracleDatabaseSchema, error)

	// Oracle/Exadata fetchers

	GetOracleExadataComponents() ([]model.OracleExadataComponent, error)
	GetOracleExadataCellDisks() (map[agentmodel.StorageServerName][]model.OracleExadataCellDisk, error)

	// Microsoft/SQLServer fetchers

	GetMicrosoftSQLServerInstances() []agentmodel.ListInstanceOutputModel
	GetMicrosoftSQLServerInstanceInfo(conn string, inst *model.MicrosoftSQLServerInstance)
	GetMicrosoftSQLServerInstanceEdition(conn string, inst *model.MicrosoftSQLServerInstance)
	GetMicrosoftSQLServerInstanceLicensingInfo(conn string, inst *model.MicrosoftSQLServerInstance)
	GetMicrosoftSQLServerInstanceDatabase(conn string) []model.MicrosoftSQLServerDatabase
	GetMicrosoftSQLServerInstanceDatabaseBackups(conn string) []agentmodel.DbBackupsModel
	GetMicrosoftSQLServerInstanceDatabaseSchemas(conn string) []agentmodel.DbSchemasModel
	GetMicrosoftSQLServerInstanceDatabaseTablespaces(conn string) []agentmodel.DbTablespacesModel
	GetMicrosoftSQLServerInstancePatches(conn string) []model.MicrosoftSQLServerPatch
	GetMicrosoftSQLServerProductFeatures(conn string) []model.MicrosoftSQLServerProductFeature

	// MySQL fetchers

	GetMySQLInstance(connection config.MySQLInstanceConnection) (*model.MySQLInstance, error)
	GetMySQLDatabases(connection config.MySQLInstanceConnection) []model.MySQLDatabase
	GetMySQLTableSchemas(connection config.MySQLInstanceConnection) []model.MySQLTableSchema
	GetMySQLSegmentAdvisors(connection config.MySQLInstanceConnection) []model.MySQLSegmentAdvisor
	GetMySQLHighAvailability(connection config.MySQLInstanceConnection) bool
	GetMySQLUUID() string
	GetMySQLSlaveHosts(connection config.MySQLInstanceConnection) (bool, []string)
	GetMySQLSlaveStatus(connection config.MySQLInstanceConnection) (bool, *string)
}

// User struct
type User struct {
	Name     string
	UID, GID uint32
}
