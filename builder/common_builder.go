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

package builder

import (
	"context"
	"runtime"
	"strings"
	"sync"

	"github.com/ercole-io/ercole-agent/config"
	"github.com/ercole-io/ercole-agent/fetcher"
	"github.com/ercole-io/ercole-agent/logger"
	"github.com/ercole-io/ercole-agent/model"
	"github.com/ercole-io/ercole-agent/utils"
)

// CommonBuilder for Linux and Windows hosts
type CommonBuilder struct {
	fetcher       fetcher.Fetcher
	configuration config.Configuration
	log           logger.Logger
}

// NewCommonBuilder initialize an appropriate builder for Linux or Windows
func NewCommonBuilder(configuration config.Configuration, log logger.Logger) CommonBuilder {
	var specializedFetcher fetcher.SpecializedFetcher

	log.Debugf("runtime.GOOS: [%v]", runtime.GOOS)

	if runtime.GOOS == "windows" {
		wf := fetcher.NewWindowsFetcherImpl(configuration, log)
		specializedFetcher = &wf
	} else {
		if runtime.GOOS != "linux" {
			log.Errorf("Unknow runtime.GOOS: [%v], I'll try with linux\n", runtime.GOOS)
		}

		wf := fetcher.NewLinuxFetcherImpl(configuration, log)
		specializedFetcher = &wf
	}

	fetcherImpl := &fetcher.CommonFetcherImpl{
		SpecializedFetcher: specializedFetcher,
		Configuration:      configuration,
		Log:                log,
	}

	builder := CommonBuilder{
		fetcher:       fetcherImpl,
		configuration: configuration,
		log:           log,
	}

	return builder
}

// Run fill hostData
func (b *CommonBuilder) Run(hostData *model.HostData) {
	hostData.Info = *b.getHost()

	hostData.Hostname = hostData.Info.Hostname
	if b.configuration.Hostname != "default" {
		hostData.Hostname = b.configuration.Hostname
	}

	if b.configuration.Features.Databases.Enabled {
		b.log.Debug("Databases mode enabled")
		hostData.Extra.Filesystems = b.fetcher.GetFilesystems()

		hostData.Extra.Databases = b.getOracleDBs(hostData.Info.Type)
		hostData.Databases, hostData.Schemas = b.getDatabasesAndSchemaNames(hostData.Extra.Databases)
	}

	if b.configuration.Features.Virtualization.Enabled {
		b.log.Debug("Virtualization mode enabled")
		hostData.Extra.Clusters = b.getClustersInfos()
	}
}

func (b *CommonBuilder) getHost() *model.Host {
	host := b.fetcher.GetHost()

	host.Environment = b.configuration.Envtype
	host.Location = b.configuration.Location

	return &host
}

func (b *CommonBuilder) getOracleDBs(hostType string) []model.Database {
	oratabEntries := b.fetcher.GetOratabEntries()

	databaseChannel := make(chan *model.Database, len(oratabEntries))

	for i := range oratabEntries {
		entry := oratabEntries[i]

		utils.RunRoutine(b.configuration, func() {
			b.log.Debugf("oratab entry: [%v]", entry)

			databaseChannel <- b.getOracleDB(entry, hostType)
		})
	}

	var databases = []model.Database{}
	for i := 0; i < len(oratabEntries); i++ {
		db := (<-databaseChannel)
		if db != nil {
			databases = append(databases, *db)
		}
	}

	return databases
}

func (b *CommonBuilder) getOracleDB(entry model.OratabEntry, hostType string) *model.Database {
	dbStatus := b.fetcher.GetDbStatus(entry)
	var database *model.Database

	switch dbStatus {
	case "OPEN":
		database = b.getOpenDatabase(entry, hostType)
	case "MOUNTED":
		{
			db := b.fetcher.GetMountedDb(entry)
			database = &db

			database.Tablespaces = []model.Tablespace{}
			database.Schemas = []model.Schema{}
			database.Patches = []model.Patch{}
			database.Licenses = []model.License{}
			database.ADDMs = []model.Addm{}
			database.SegmentAdvisors = []model.SegmentAdvisor{}
			database.LastPSUs = []model.PSU{}
			database.Backups = []model.Backup{}
		}
	default:
		b.log.Warnf("Unknown dbStatus: [%s] DBName: [%s] OracleHome: [%s]",
			dbStatus, entry.DBName, entry.OracleHome)
		return nil
	}

	return database
}

func (b *CommonBuilder) getOpenDatabase(entry model.OratabEntry, hostType string) *model.Database {
	dbVersion := b.fetcher.GetDbVersion(entry)

	statsCtx, cancelStatsCtx := context.WithCancel(context.Background())
	if b.configuration.Forcestats {
		utils.RunRoutine(b.configuration, func() {
			b.fetcher.RunStats(entry)

			cancelStatsCtx()
		})
	} else {
		cancelStatsCtx()
	}

	database := b.fetcher.GetOpenDb(entry)

	var wg sync.WaitGroup

	utils.RunRoutineInGroup(b.configuration, func() {
		database.Tablespaces = b.fetcher.GetTablespaces(entry)
	}, &wg)

	utils.RunRoutineInGroup(b.configuration, func() {
		database.Schemas = b.fetcher.GetSchemas(entry)
	}, &wg)

	utils.RunRoutineInGroup(b.configuration, func() {
		database.Patches = b.fetcher.GetPatches(entry, dbVersion)
	}, &wg)

	utils.RunRoutineInGroup(b.configuration, func() {
		<-statsCtx.Done()

		database.Features2 = b.fetcher.GetFeatures2(entry, dbVersion)
	}, &wg)

	utils.RunRoutineInGroup(b.configuration, func() {
		<-statsCtx.Done()

		database.Licenses = b.fetcher.GetLicenses(entry, dbVersion, hostType)
	}, &wg)

	utils.RunRoutineInGroup(b.configuration, func() {
		database.ADDMs = b.fetcher.GetADDMs(entry)
	}, &wg)

	utils.RunRoutineInGroup(b.configuration, func() {
		database.SegmentAdvisors = b.fetcher.GetSegmentAdvisors(entry)
	}, &wg)

	utils.RunRoutineInGroup(b.configuration, func() {
		database.LastPSUs = b.fetcher.GetLastPSUs(entry, dbVersion)
	}, &wg)

	utils.RunRoutineInGroup(b.configuration, func() {
		database.Backups = b.fetcher.GetBackups(entry)
	}, &wg)

	wg.Wait()

	return &database
}

func (b *CommonBuilder) getDatabasesAndSchemaNames(databases []model.Database) (databasesNames, schemasNames string) {
	for _, db := range databases {
		databasesNames += db.Name + " "

		for _, sc := range db.Schemas {
			schemasNames += sc.User + " "
		}
	}

	databasesNames = strings.TrimSpace(databasesNames)
	schemasNames = strings.TrimSpace(schemasNames)

	return
}

func (b *CommonBuilder) getClustersInfos() []model.ClusterInfo {
	countHypervisors := len(b.configuration.Features.Virtualization.Hypervisors)

	clustersChan := make(chan []model.ClusterInfo, countHypervisors)
	vmsChan := make(chan []model.VMInfo, countHypervisors)

	for _, hv := range b.configuration.Features.Virtualization.Hypervisors {
		utils.RunRoutine(b.configuration, func() {
			clustersChan <- b.fetcher.GetClusters(hv)
		})

		utils.RunRoutine(b.configuration, func() {
			vmsChan <- b.fetcher.GetVirtualMachines(hv)
		})
	}

	clusters := make([]model.ClusterInfo, 0)
	for i := 0; i < countHypervisors; i++ {
		clusters = append(clusters, (<-clustersChan)...)
	}

	vms := make([]model.VMInfo, 0)
	for i := 0; i < countHypervisors; i++ {
		vms = append(vms, (<-vmsChan)...)
	}

	clusters = setVMsInClusterInfo(clusters, vms)

	return clusters
}

func setVMsInClusterInfo(clusters []model.ClusterInfo, vms []model.VMInfo) []model.ClusterInfo {
	clusters = append(clusters, model.ClusterInfo{
		Name:    "not_in_cluster",
		Type:    "unknown",
		CPU:     0,
		Sockets: 0,
		VMs:     []model.VMInfo{},
	})

	clusterMap := make(map[string][]model.VMInfo)

	for _, vm := range vms {
		if vm.ClusterName == "" {
			vm.ClusterName = "not_in_cluster"
		}
		clusterMap[vm.ClusterName] = append(clusterMap[vm.ClusterName], vm)
	}

	for i := range clusters {
		if clusterMap[clusters[i].Name] != nil {
			clusters[i].VMs = clusterMap[clusters[i].Name]
		} else {
			clusters[i].VMs = []model.VMInfo{}
		}
	}

	return clusters
}
