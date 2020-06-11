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
	"github.com/ercole-io/ercole-agent/config"
	"github.com/ercole-io/ercole-agent/model"
)

// Fetcher interface for Linux and Windows
type Fetcher interface {
	SetUser(username string) error
	SetUserAsCurrent() error

	GetHost() model.Host
	GetFilesystems() []model.Filesystem
	GetOratabEntries() []model.OratabEntry
	GetDbStatus(entry model.OratabEntry) string
	GetMountedDb(entry model.OratabEntry) model.Database
	GetDbVersion(entry model.OratabEntry) string
	RunStats(entry model.OratabEntry)
	GetOpenDb(entry model.OratabEntry) model.Database
	GetTablespaces(entry model.OratabEntry) []model.Tablespace
	GetSchemas(entry model.OratabEntry) []model.Schema
	GetPatches(entry model.OratabEntry, dbVersion string) []model.Patch
	GetFeatures2(entry model.OratabEntry, dbVersion string) []model.Feature2
	GetLicenses(entry model.OratabEntry, dbVersion, hostType string) []model.License
	GetADDMs(entry model.OratabEntry) []model.Addm
	GetSegmentAdvisors(entry model.OratabEntry) []model.SegmentAdvisor
	GetLastPSUs(entry model.OratabEntry, dbVersion string) []model.PSU
	GetBackups(entry model.OratabEntry) []model.Backup
	GetClusters(hv config.Hypervisor) []model.ClusterInfo
	GetVirtualMachines(hv config.Hypervisor) []model.VMInfo
	GetExadataDevices() []model.ExadataDevice
	GetExadataCellDisks() []model.ExadataCellDisk
}
