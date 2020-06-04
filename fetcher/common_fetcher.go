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

	"github.com/ercole-io/ercole-agent/config"
	"github.com/ercole-io/ercole-agent/marshal"
	"github.com/ercole-io/ercole-agent/model"
	"github.com/sirupsen/logrus"
)

// CommonFetcherImpl implement common behaviour between Linux and Windows fetchers
type CommonFetcherImpl struct {
	SpecializedFetcher
	Configuration config.Configuration
	Log           *logrus.Logger
}

// SpecializedFetcher define specific behaviour of Linux and Windows fetchers
type SpecializedFetcher interface {
	Execute(fetcherName string, params ...string) []byte
	GetClusters(hv config.Hypervisor) []model.ClusterInfo
	GetVirtualMachines(hv config.Hypervisor) []model.VMInfo
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
func (cf *CommonFetcherImpl) GetOratabEntries() []model.OratabEntry {
	out := cf.Execute("oratab", cf.Configuration.Oratab)
	return marshal.Oratab(out)
}

// GetDbStatus get
func (cf *CommonFetcherImpl) GetDbStatus(entry model.OratabEntry) string {
	out := cf.Execute("dbstatus", entry.DBName, entry.OracleHome)
	return strings.TrimSpace(string(out))
}

// GetMountedDb get
func (cf *CommonFetcherImpl) GetMountedDb(entry model.OratabEntry) model.Database {
	out := cf.Execute("dbmounted", entry.DBName, entry.OracleHome)
	return marshal.Database(out)
}

// GetDbVersion get
func (cf *CommonFetcherImpl) GetDbVersion(entry model.OratabEntry) string {
	out := cf.Execute("dbversion", entry.DBName, entry.OracleHome)
	return strings.Split(string(out), ".")[0]
}

// RunStats Execute stats script
func (cf *CommonFetcherImpl) RunStats(entry model.OratabEntry) {
	cf.Execute("stats", entry.DBName, entry.OracleHome)
}

// GetOpenDb get
func (cf *CommonFetcherImpl) GetOpenDb(entry model.OratabEntry) model.Database {
	out := cf.Execute("db", entry.DBName, entry.OracleHome, strconv.Itoa(cf.Configuration.AWR))
	return marshal.Database(out)
}

// GetTablespaces get
func (cf *CommonFetcherImpl) GetTablespaces(entry model.OratabEntry) []model.Tablespace {
	out := cf.Execute("tablespace", entry.DBName, entry.OracleHome)
	return marshal.Tablespaces(out)
}

// GetSchemas get
func (cf *CommonFetcherImpl) GetSchemas(entry model.OratabEntry) []model.Schema {
	out := cf.Execute("schema", entry.DBName, entry.OracleHome)
	return marshal.Schemas(out)
}

// GetPatches get
func (cf *CommonFetcherImpl) GetPatches(entry model.OratabEntry, dbVersion string) []model.Patch {
	out := cf.Execute("patch", entry.DBName, dbVersion, entry.OracleHome)
	return marshal.Patches(out)
}

// GetFeatures2 get
func (cf *CommonFetcherImpl) GetFeatures2(entry model.OratabEntry, dbVersion string) []model.Feature2 {
	out := cf.Execute("opt", entry.DBName, dbVersion, entry.OracleHome)
	return marshal.Features2(out)
}

// GetLicenses get
func (cf *CommonFetcherImpl) GetLicenses(entry model.OratabEntry, dbVersion, hostType string) []model.License {
	out := cf.Execute("license", entry.DBName, dbVersion, hostType, entry.OracleHome)
	return marshal.Licenses(out)
}

// GetADDMs get
func (cf *CommonFetcherImpl) GetADDMs(entry model.OratabEntry) []model.Addm {
	out := cf.Execute("addm", entry.DBName, entry.OracleHome)
	return marshal.Addms(out)
}

// GetSegmentAdvisors get
func (cf *CommonFetcherImpl) GetSegmentAdvisors(entry model.OratabEntry) []model.SegmentAdvisor {
	out := cf.Execute("segmentadvisor", entry.DBName, entry.OracleHome)
	return marshal.SegmentAdvisor(out)
}

// GetLastPSUs get
func (cf *CommonFetcherImpl) GetLastPSUs(entry model.OratabEntry, dbVersion string) []model.PSU {
	out := cf.Execute("psu", entry.DBName, dbVersion, entry.OracleHome)
	return marshal.PSU(out)
}

// GetBackups get
func (cf *CommonFetcherImpl) GetBackups(entry model.OratabEntry) []model.Backup {
	out := cf.Execute("backup", entry.DBName, entry.OracleHome)
	return marshal.Backups(out)
}
