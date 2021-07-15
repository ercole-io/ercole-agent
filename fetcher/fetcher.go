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
	GetClustersMembershipStatus() (*model.ClusterMembershipStatus, error)

	// Virtualization fetcher

	GetClusters(hv config.Hypervisor) ([]model.ClusterInfo, error)
	GetVirtualMachines(hv config.Hypervisor) (map[string][]model.VMInfo, error)

	// Oracle/Database fetchers

	GetOracleDatabaseOratabEntries() ([]agentmodel.OratabEntry, error)
	GetOracleDatabaseRunningDatabases() ([]string, error)
	GetOracleDatabaseDbStatus(entry agentmodel.OratabEntry) (string, error)
	GetOracleDatabaseMountedDb(entry agentmodel.OratabEntry) (*model.OracleDatabase, error)
	GetOracleDatabaseDbVersion(entry agentmodel.OratabEntry) (string, error)
	RunOracleDatabaseStats(entry agentmodel.OratabEntry) error
	GetOracleDatabaseOpenDb(entry agentmodel.OratabEntry) (*model.OracleDatabase, error)
	GetOracleDatabaseTablespaces(entry agentmodel.OratabEntry) ([]model.OracleDatabaseTablespace, error)
	GetOracleDatabaseSchemas(entry agentmodel.OratabEntry) ([]model.OracleDatabaseSchema, error)
	GetOracleDatabasePatches(entry agentmodel.OratabEntry, dbVersion string) ([]model.OracleDatabasePatch, error)
	GetOracleDatabaseFeatureUsageStat(entry agentmodel.OratabEntry, dbVersion string) ([]model.OracleDatabaseFeatureUsageStat, error)
	GetOracleDatabaseLicenses(entry agentmodel.OratabEntry, dbVersion, hardwareAbstractionTechnology string) ([]model.OracleDatabaseLicense, error)
	GetOracleDatabaseADDMs(entry agentmodel.OratabEntry) ([]model.OracleDatabaseAddm, error)
	GetOracleDatabaseSegmentAdvisors(entry agentmodel.OratabEntry) ([]model.OracleDatabaseSegmentAdvisor, error)
	GetOracleDatabasePSUs(entry agentmodel.OratabEntry, dbVersion string) ([]model.OracleDatabasePSU, error)
	GetOracleDatabaseBackups(entry agentmodel.OratabEntry) ([]model.OracleDatabaseBackup, error)
	GetOracleDatabaseCheckPDB(entry agentmodel.OratabEntry) (bool, error)
	GetOracleDatabasePDBs(entry agentmodel.OratabEntry) ([]model.OracleDatabasePluggableDatabase, error)
	GetOracleDatabasePDBTablespaces(entry agentmodel.OratabEntry, pdb string) ([]model.OracleDatabaseTablespace, error)
	GetOracleDatabasePDBSchemas(entry agentmodel.OratabEntry, pdb string) ([]model.OracleDatabaseSchema, error)

	// Oracle/Exadata fetchers

	GetOracleExadataComponents() ([]model.OracleExadataComponent, error)
	GetOracleExadataCellDisks() (map[agentmodel.StorageServerName][]model.OracleExadataCellDisk, error)

	// Microsoft/SQLServer fetchers

	GetMicrosoftSQLServerInstances() ([]agentmodel.ListInstanceOutputModel, error)
	GetMicrosoftSQLServerInstanceInfo(conn string, inst *model.MicrosoftSQLServerInstance) error
	GetMicrosoftSQLServerInstanceEdition(conn string, inst *model.MicrosoftSQLServerInstance) error
	GetMicrosoftSQLServerInstanceLicensingInfo(conn string, inst *model.MicrosoftSQLServerInstance) error
	GetMicrosoftSQLServerInstanceDatabase(conn string) ([]model.MicrosoftSQLServerDatabase, error)
	GetMicrosoftSQLServerInstanceDatabaseBackups(conn string) ([]agentmodel.DbBackupsModel, error)
	GetMicrosoftSQLServerInstanceDatabaseSchemas(conn string) ([]agentmodel.DbSchemasModel, error)
	GetMicrosoftSQLServerInstanceDatabaseTablespaces(conn string) ([]agentmodel.DbTablespacesModel, error)
	GetMicrosoftSQLServerInstancePatches(conn string) ([]model.MicrosoftSQLServerPatch, error)
	GetMicrosoftSQLServerProductFeatures(conn string) ([]model.MicrosoftSQLServerProductFeature, error)

	// MySQL fetchers

	GetMySQLInstance(connection config.MySQLInstanceConnection) (*model.MySQLInstance, error)
	GetMySQLDatabases(connection config.MySQLInstanceConnection) ([]model.MySQLDatabase, error)
	GetMySQLTableSchemas(connection config.MySQLInstanceConnection) ([]model.MySQLTableSchema, error)
	GetMySQLSegmentAdvisors(connection config.MySQLInstanceConnection) ([]model.MySQLSegmentAdvisor, error)
	GetMySQLHighAvailability(connection config.MySQLInstanceConnection) (bool, error)
	GetMySQLUUID() (string, error)
	GetMySQLSlaveHosts(connection config.MySQLInstanceConnection) (bool, []string, error)
	GetMySQLSlaveStatus(connection config.MySQLInstanceConnection) (bool, *string, error)
}

// User struct
type User struct {
	Name     string
	UID, GID uint32
}
