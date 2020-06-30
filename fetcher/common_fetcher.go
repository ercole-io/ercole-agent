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
	GetExadataComponents() []model.OracleExadataComponent
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

// GetOratabEntries get
func (cf *CommonFetcherImpl) GetOratabEntries() []agentmodel.OratabEntry {
	out := cf.Execute("oratab", cf.Configuration.Oratab)
	return marshal.Oratab(out)
}

// GetDbStatus get
func (cf *CommonFetcherImpl) GetDbStatus(entry agentmodel.OratabEntry) string {
	out := cf.Execute("dbstatus", entry.DBName, entry.OracleHome)
	return strings.TrimSpace(string(out))
}

// GetMountedDb get
func (cf *CommonFetcherImpl) GetMountedDb(entry agentmodel.OratabEntry) model.OracleDatabase {
	out := cf.Execute("dbmounted", entry.DBName, entry.OracleHome)
	return marshal.Database(out)
}

// GetDbVersion get
func (cf *CommonFetcherImpl) GetDbVersion(entry agentmodel.OratabEntry) string {
	out := cf.Execute("dbversion", entry.DBName, entry.OracleHome)
	return strings.Split(string(out), ".")[0]
}

// RunStats Execute stats script
func (cf *CommonFetcherImpl) RunStats(entry agentmodel.OratabEntry) {
	cf.Execute("stats", entry.DBName, entry.OracleHome)
}

// GetOpenDb get
func (cf *CommonFetcherImpl) GetOpenDb(entry agentmodel.OratabEntry) model.OracleDatabase {
	out := cf.Execute("db", entry.DBName, entry.OracleHome, strconv.Itoa(cf.Configuration.AWR))
	return marshal.Database(out)
}

// GetTablespaces get
func (cf *CommonFetcherImpl) GetTablespaces(entry agentmodel.OratabEntry) []model.OracleDatabaseTablespace {
	out := cf.Execute("tablespace", entry.DBName, entry.OracleHome)
	return marshal.Tablespaces(out)
}

// GetSchemas get
func (cf *CommonFetcherImpl) GetSchemas(entry agentmodel.OratabEntry) []model.OracleDatabaseSchema {
	out := cf.Execute("schema", entry.DBName, entry.OracleHome)
	return marshal.Schemas(out)
}

// GetPatches get
func (cf *CommonFetcherImpl) GetPatches(entry agentmodel.OratabEntry, dbVersion string) []model.OracleDatabasePatch {
	out := cf.Execute("patch", entry.DBName, dbVersion, entry.OracleHome)
	return marshal.Patches(out)
}

// GetDatabaseFeatureUsageStat get
func (cf *CommonFetcherImpl) GetDatabaseFeatureUsageStat(entry agentmodel.OratabEntry, dbVersion string) []model.OracleDatabaseFeatureUsageStat {
	out := cf.Execute("opt", entry.DBName, dbVersion, entry.OracleHome)
	return marshal.DatabaseFeatureUsageStat(out)
}

// GetLicenses get
func (cf *CommonFetcherImpl) GetLicenses(entry agentmodel.OratabEntry, dbVersion, hardwareAbstractionTechnology string) []model.OracleDatabaseLicense {
	out := cf.Execute("license", entry.DBName, dbVersion, hardwareAbstractionTechnology, entry.OracleHome)
	return marshal.Licenses(out)
}

// GetADDMs get
func (cf *CommonFetcherImpl) GetADDMs(entry agentmodel.OratabEntry) []model.OracleDatabaseAddm {
	out := cf.Execute("addm", entry.DBName, entry.OracleHome)
	return marshal.Addms(out)
}

// GetSegmentAdvisors get
func (cf *CommonFetcherImpl) GetSegmentAdvisors(entry agentmodel.OratabEntry) []model.OracleDatabaseSegmentAdvisor {
	out := cf.Execute("segmentadvisor", entry.DBName, entry.OracleHome)
	return marshal.SegmentAdvisor(out)
}

// GetPSUs get
func (cf *CommonFetcherImpl) GetPSUs(entry agentmodel.OratabEntry, dbVersion string) []model.OracleDatabasePSU {
	out := cf.Execute("psu", entry.DBName, dbVersion, entry.OracleHome)
	return marshal.PSU(out)
}

// GetBackups get
func (cf *CommonFetcherImpl) GetBackups(entry agentmodel.OratabEntry) []model.OracleDatabaseBackup {
	out := cf.Execute("backup", entry.DBName, entry.OracleHome)
	return marshal.Backups(out)
}
