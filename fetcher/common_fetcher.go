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
	"strconv"
	"strings"

	"github.com/ercole-io/ercole-agent/agentmodel"
	"github.com/ercole-io/ercole-agent/config"
	"github.com/ercole-io/ercole-agent/logger"
	"github.com/ercole-io/ercole-agent/marshal"
	marshal_oracle "github.com/ercole-io/ercole-agent/marshal/oracle"
	"github.com/ercole-io/ercole/model"
)

// CommonFetcherImpl implement common behaviour between Linux and Windows fetchers
type CommonFetcherImpl struct {
	SpecializedFetcher
	Configuration config.Configuration
	Log           logger.Logger
}

// SpecializedFetcher define specific behaviour of Linux and Windows fetchers
type SpecializedFetcher interface {
	SetUser(username string) error
	SetUserAsCurrent() error

	Execute(fetcherName string, params ...string) []byte
	GetClusters(hv config.Hypervisor) []model.ClusterInfo
	GetVirtualMachines(hv config.Hypervisor) map[string][]model.VMInfo
	GetOracleExadataComponents() []model.OracleExadataComponent
	GetOracleExadataCellDisks() map[agentmodel.StorageServerName][]model.OracleExadataCellDisk
	GetClustersMembershipStatus() model.ClusterMembershipStatus
}

// GetHost get
func (cf *CommonFetcherImpl) GetHost() model.Host {
	out := cf.Execute("host")
	return marshal.Host(out)
}

// GetFilesystems get
func (cf *CommonFetcherImpl) GetFilesystems() []model.Filesystem {
	out := cf.Execute("filesystem")
	return marshal.Filesystems(out)
}

// GetOracleOratabEntries get
func (cf *CommonFetcherImpl) GetOracleOratabEntries() []agentmodel.OratabEntry {
	out := cf.Execute("oratab", cf.Configuration.Features.OracleDatabase.Oratab)
	return marshal_oracle.Oratab(out)
}

// GetOracleDbStatus get
func (cf *CommonFetcherImpl) GetOracleDbStatus(entry agentmodel.OratabEntry) string {
	out := cf.Execute("dbstatus", entry.DBName, entry.OracleHome)
	return strings.TrimSpace(string(out))
}

// GetOracleMountedDb get
func (cf *CommonFetcherImpl) GetOracleMountedDb(entry agentmodel.OratabEntry) model.OracleDatabase {
	out := cf.Execute("dbmounted", entry.DBName, entry.OracleHome)
	return marshal_oracle.Database(out)
}

// GetOracleDbVersion get
func (cf *CommonFetcherImpl) GetOracleDbVersion(entry agentmodel.OratabEntry) string {
	out := cf.Execute("dbversion", entry.DBName, entry.OracleHome)
	return strings.Split(string(out), ".")[0]
}

// RunOracleStats Execute stats script
func (cf *CommonFetcherImpl) RunOracleStats(entry agentmodel.OratabEntry) {
	cf.Execute("stats", entry.DBName, entry.OracleHome)
}

// GetOracleOpenDb get
func (cf *CommonFetcherImpl) GetOracleOpenDb(entry agentmodel.OratabEntry) model.OracleDatabase {
	out := cf.Execute("db", entry.DBName, entry.OracleHome, strconv.Itoa(cf.Configuration.Features.OracleDatabase.AWR))
	return marshal_oracle.Database(out)
}

// GetOracleTablespaces get
func (cf *CommonFetcherImpl) GetOracleTablespaces(entry agentmodel.OratabEntry) []model.OracleDatabaseTablespace {
	out := cf.Execute("tablespace", entry.DBName, entry.OracleHome)
	return marshal_oracle.Tablespaces(out)
}

// GetOracleSchemas get
func (cf *CommonFetcherImpl) GetOracleSchemas(entry agentmodel.OratabEntry) []model.OracleDatabaseSchema {
	out := cf.Execute("schema", entry.DBName, entry.OracleHome)
	return marshal_oracle.Schemas(out)
}

// GetOraclePatches get
func (cf *CommonFetcherImpl) GetOraclePatches(entry agentmodel.OratabEntry, dbVersion string) []model.OracleDatabasePatch {
	out := cf.Execute("patch", entry.DBName, dbVersion, entry.OracleHome)
	return marshal_oracle.Patches(out)
}

// GetOracleDatabaseFeatureUsageStat get
func (cf *CommonFetcherImpl) GetOracleDatabaseFeatureUsageStat(entry agentmodel.OratabEntry, dbVersion string) []model.OracleDatabaseFeatureUsageStat {
	out := cf.Execute("opt", entry.DBName, dbVersion, entry.OracleHome)
	return marshal_oracle.DatabaseFeatureUsageStat(out)
}

// GetOracleLicenses get
func (cf *CommonFetcherImpl) GetOracleLicenses(entry agentmodel.OratabEntry, dbVersion, hardwareAbstractionTechnology string) []model.OracleDatabaseLicense {
	out := cf.Execute("license", entry.DBName, dbVersion, hardwareAbstractionTechnology, entry.OracleHome)
	return marshal_oracle.Licenses(out)
}

// GetOracleADDMs get
func (cf *CommonFetcherImpl) GetOracleADDMs(entry agentmodel.OratabEntry) []model.OracleDatabaseAddm {
	out := cf.Execute("addm", entry.DBName, entry.OracleHome)
	return marshal_oracle.Addms(out)
}

// GetOracleSegmentAdvisors get
func (cf *CommonFetcherImpl) GetOracleSegmentAdvisors(entry agentmodel.OratabEntry) []model.OracleDatabaseSegmentAdvisor {
	out := cf.Execute("segmentadvisor", entry.DBName, entry.OracleHome)
	return marshal_oracle.SegmentAdvisor(out)
}

// GetOraclePSUs get
func (cf *CommonFetcherImpl) GetOraclePSUs(entry agentmodel.OratabEntry, dbVersion string) []model.OracleDatabasePSU {
	out := cf.Execute("psu", entry.DBName, dbVersion, entry.OracleHome)
	return marshal_oracle.PSU(out)
}

// GetOracleBackups get
func (cf *CommonFetcherImpl) GetOracleBackups(entry agentmodel.OratabEntry) []model.OracleDatabaseBackup {
	out := cf.Execute("backup", entry.DBName, entry.OracleHome)
	return marshal_oracle.Backups(out)
}
