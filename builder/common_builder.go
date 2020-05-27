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
	"log"
	"runtime"
	"strings"
	"sync"

	"github.com/ercole-io/ercole-agent/config"
	"github.com/ercole-io/ercole-agent/fetcher"
	"github.com/ercole-io/ercole-agent/model"
	"github.com/ercole-io/ercole-agent/utils"
)

// CommonBuilder for Linux and Windows hosts
type CommonBuilder struct {
	fetcher       fetcher.Fetcher
	configuration config.Configuration
}

// NewCommonBuilder initialize an appropriate builder for Linux or Windows
func NewCommonBuilder(configuration config.Configuration) CommonBuilder {
	var fetcherImpl fetcher.Fetcher

	if runtime.GOOS == "windows" {
		fetcherImpl = &fetcher.CommonFetcherImpl{
			Configuration: configuration,
			SpecializedFetcher: &fetcher.WindowsFetcherImpl{
				Configuration: configuration,
			},
		}

	} else {
		if runtime.GOOS != "linux" {
			log.Printf("Unknow runtime.GOOS: [%v], I'll try with linux\n", runtime.GOOS)
		}

		fetcherImpl = &fetcher.CommonFetcherImpl{
			Configuration: configuration,
			SpecializedFetcher: &fetcher.LinuxFetcherImpl{
				Configuration: configuration,
			},
		}
	}

	builder := CommonBuilder{
		fetcher:       fetcherImpl,
		configuration: configuration,
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

	hostData.Extra.Filesystems = b.fetcher.GetFilesystems()
	hostData.Extra.Databases = b.getOracleDBs(hostData.Info.Type)

	hostData.Databases, hostData.Schemas = b.getDatabasesAndSchemaNames(hostData.Extra.Databases)
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
			database.Features = []model.Feature{}
			database.Licenses = []model.License{}
			database.ADDMs = []model.Addm{}
			database.SegmentAdvisors = []model.SegmentAdvisor{}
			database.LastPSUs = []model.PSU{}
			database.Backups = []model.Backup{}
		}
	default:
		log.Println("Error! DBName: [" + entry.DBName + "] OracleHome: [" + entry.OracleHome + "]  Wrong dbStatus: [" + dbStatus + "]")
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

		database.Features = b.fetcher.GetFeatures(entry, dbVersion)
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
