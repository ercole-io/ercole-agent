// Copyright (c) 2022 Sorint.lab S.p.A.
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
	"time"

	"github.com/ercole-io/ercole-agent/v2/agentmodel"
	"github.com/ercole-io/ercole-agent/v2/config"
	"github.com/ercole-io/ercole/v2/model"
)

const (
	FetcherStandardTimeOut = 12 * time.Hour
)

// Fetcher interface for Linux and Windows
type Fetcher interface {
	SetUser(username string) error
	SetUserAsCurrent() error

	// Operating system fetchers

	GetHost() (*model.Host, error)
	GetCwVersion() (string, error)
	GetFilesystems() ([]model.Filesystem, error)
	GetClustersMembershipStatus() (*model.ClusterMembershipStatus, error)
	GetCpuConsumption() ([]model.CpuConsumption, error)
	GetDiskConsumption() ([]model.DiskConsumption, error)

	// Virtualization fetcher

	GetClusters(hv config.Hypervisor) ([]model.ClusterInfo, error)
	GetVirtualMachines(hv config.Hypervisor) (map[string][]model.VMInfo, error)

	// Oracle/Database fetchers

	GetOracleDatabaseOratabEntries() ([]agentmodel.OratabEntry, error)
	GetOracleDatabaseRunningDatabases() ([]string, error)
	GetOraclePmonInstances() (map[string]string, error)
	GetOracleEntry(proc, instanceName string) (*agentmodel.OratabEntry, error)
	GetOracleDatabaseDbStatus(entry agentmodel.OratabEntry) (string, error)
	GetOracleDatabaseRac(entry agentmodel.OratabEntry) (string, error)
	GetOracleDatabaseMountedDb(entry agentmodel.OratabEntry) (*model.OracleDatabase, error)
	GetOracleDatabaseDbVersion(entry agentmodel.OratabEntry) (string, error)
	RunOracleDatabaseStats(entry agentmodel.OratabEntry) error
	GetOracleDatabaseOpenDb(entry agentmodel.OratabEntry) (*model.OracleDatabase, error)
	GetOracleDatabaseTablespaces(entry agentmodel.OratabEntry) ([]model.OracleDatabaseTablespace, error)
	GetOracleDatabaseSchemas(entry agentmodel.OratabEntry) ([]model.OracleDatabaseSchema, error)
	GetOracleDatabasePatches(entry agentmodel.OratabEntry, dbVersion string) ([]model.OracleDatabasePatch, error)
	GetOracleDatabaseFeatureUsageStat(entry agentmodel.OratabEntry, dbVersion string) ([]model.OracleDatabaseFeatureUsageStat, error)
	GetOracleDatabaseLicenses(entry agentmodel.OratabEntry, dbVersion, hardwareAbstractionTechnology string, hostCoreFactor float64) ([]model.OracleDatabaseLicense, error)
	GetOracleDatabaseADDMs(entry agentmodel.OratabEntry) ([]model.OracleDatabaseAddm, error)
	GetOracleDatabaseSegmentAdvisors(entry agentmodel.OratabEntry) ([]model.OracleDatabaseSegmentAdvisor, error)
	GetOracleDatabasePSUs(entry agentmodel.OratabEntry, dbVersion string) ([]model.OracleDatabasePSU, error)
	GetOracleDatabaseBackups(entry agentmodel.OratabEntry) ([]model.OracleDatabaseBackup, error)
	GetOracleDatabaseCheckPDB(entry agentmodel.OratabEntry) (bool, error)
	GetOracleDatabasePDBs(entry agentmodel.OratabEntry) ([]model.OracleDatabasePluggableDatabase, error)
	GetOracleDatabasePDBTablespaces(entry agentmodel.OratabEntry, pdb string) ([]model.OracleDatabaseTablespace, error)
	GetOracleDatabasePDBSchemas(entry agentmodel.OratabEntry, pdb string) ([]model.OracleDatabaseSchema, error)
	GetOracleDatabasePDBSegmentAdvisors(entry agentmodel.OratabEntry, pdb string) ([]model.OracleDatabaseSegmentAdvisor, error)
	GetOracleDatabaseServices(entry agentmodel.OratabEntry) ([]model.OracleDatabaseService, error)
	GetOracleDatabaseGrantsDba(entry agentmodel.OratabEntry) ([]model.OracleGrantDba, error)
	GetOracleDatabasePDBSize(entry agentmodel.OratabEntry, pdb string) (model.OracleDatabasePdbSize, error)
	GetOracleDatabasePDBCharset(entry agentmodel.OratabEntry, pdb string) (string, error)
	GetOracleDatabasePartitionings(entry agentmodel.OratabEntry) ([]model.OracleDatabasePartitioning, error)
	GetOracleDatabasePDBPartitionings(entry agentmodel.OratabEntry, pdb string) ([]model.OracleDatabasePartitioning, error)
	GetOracleDatabaseCpuDiskConsumptions(entry agentmodel.OratabEntry) ([]model.CpuDiskConsumption, error)
	GetOracleDatabaseCpuDiskConsumptionPdbs(entry agentmodel.OratabEntry, pdb string) ([]model.CpuDiskConsumptionPdb, error)
	GetOracleDatabasePgsqlMigrability(entry agentmodel.OratabEntry) ([]model.PgsqlMigrability, error)
	GetOracleDatabaseMemoryAdvisor(entry agentmodel.OratabEntry) (*model.OracleDatabaseMemoryAdvisor, error)
	GetOracleDatabasePgsqlMigrabilityPdbs(entry agentmodel.OratabEntry, pdb string) ([]model.PgsqlMigrability, error)
	GetOracleDatabasePdbServices(entry agentmodel.OratabEntry, pdb string) ([]model.OracleDatabasePdbService, error)
	GetOracleDatabasePoliciesAudit(entry agentmodel.OratabEntry) ([]string, error)
	GetOracleDatabasePoliciesAuditPdbs(entry agentmodel.OratabEntry, pdb string) ([]string, error)
	GetOracleDatabaseDiskGroups(entry agentmodel.OratabEntry) ([]model.OracleDatabaseDiskGroup, error)
	// Oracle/Exadata fetchers
	GetOracleExadataComponents() ([]model.OracleExadataComponent, error)

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

	GetMySQLVersion(connection config.MySQLInstanceConnection) (string, error)
	GetMySQLInstance(connection config.MySQLInstanceConnection) (*model.MySQLInstance, error)
	GetMySQLOldInstance(connection config.MySQLInstanceConnection) (*model.MySQLInstance, error)
	GetMySQLDatabases(connection config.MySQLInstanceConnection) ([]model.MySQLDatabase, error)
	GetMySQLTableSchemas(connection config.MySQLInstanceConnection) ([]model.MySQLTableSchema, error)
	GetMySQLSegmentAdvisors(connection config.MySQLInstanceConnection) ([]model.MySQLSegmentAdvisor, error)
	GetMySQLHighAvailability(connection config.MySQLInstanceConnection) (bool, error)
	GetMySQLUUID(dataDirectory string) (string, error)
	GetMySQLSlaveHosts(connection config.MySQLInstanceConnection) (bool, []string, error)
	GetMySQLSlaveStatus(connection config.MySQLInstanceConnection) (bool, *string, error)

	// Cloud

	GetCloudMembership() (string, error)

	// PostgreSQL
	GetPostgreSQLSetting(instanceConnection config.PostgreSQLInstanceConnection) (*model.PostgreSQLSetting, error)
	GetPostgreSQLInstance(instanceConnection config.PostgreSQLInstanceConnection, v10 bool) (*model.PostgreSQLInstance, error)
	GetPostgreSQLDbNameList(instanceConnection config.PostgreSQLInstanceConnection) ([]string, error)
	GetPostgreSQLDbSchemaNameList(instanceConnection config.PostgreSQLInstanceConnection, dbname string) ([]string, error)
	GetPostgreSQLDatabase(instanceConnection config.PostgreSQLInstanceConnection, dbname string, v10 bool) (*model.PostgreSQLDatabase, error)
	GetPostgreSQLSchema(instanceConnection config.PostgreSQLInstanceConnection, dbname string, schemaName string, v10 bool) (*model.PostgreSQLSchema, error)
}

// User struct
type User struct {
	Name     string
	UID, GID uint32
}
