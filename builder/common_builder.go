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
	"strconv"
	"strings"
	"sync"

	"github.com/ercole-io/ercole-agent/config"
	"github.com/ercole-io/ercole-agent/fetcher"
	"github.com/ercole-io/ercole-agent/marshal"
	"github.com/ercole-io/ercole-agent/model"
	"github.com/ercole-io/ercole-agent/utils"
)

// CommonBuilder for Linux and Windows hosts
type CommonBuilder struct {
	fetcher       fetcher.CommonFetcher
	configuration config.Configuration
}

// NewCommonBuilder initialize an appropriate builder for Linux or Windows
func NewCommonBuilder(configuration config.Configuration) CommonBuilder {
	var fetcherImpl fetcher.CommonFetcher

	if runtime.GOOS == "windows" {
		fetcherImpl = &fetcher.WindowsFetcherImpl{
			Configuration: configuration,
		}

	} else {
		if runtime.GOOS != "linux" {
			log.Printf("Unknow runtime.GOOS: [%v], I'll try with linux\n", runtime.GOOS)
		}

		fetcherImpl = &fetcher.LinuxFetcherImpl{
			Configuration: configuration,
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

	hostData.Extra.Filesystems = b.getFilesystems()
	hostData.Extra.Databases = b.getDatabases(hostData.Info.Type)

	hostData.Databases, hostData.Schemas = b.getDatabasesAndSchemaNames(hostData.Extra.Databases)
}

func (b *CommonBuilder) getHost() *model.Host {
	out := b.fetcher.Execute("host")
	host := marshal.Host(out)

	host.Environment = b.configuration.Envtype
	host.Location = b.configuration.Location

	return &host
}

func (b *CommonBuilder) getFilesystems() []model.Filesystem {
	out := b.fetcher.Execute("filesystem")
	return marshal.Filesystems(out)
}

//TODO Rename in OracleDatabases
func (b *CommonBuilder) getDatabases(hostType string) []model.Database {
	out := b.fetcher.Execute("oratab", b.configuration.Oratab)
	oratabEntries := marshal.Oratab(out)

	databaseChannel := make(chan *model.Database, len(oratabEntries))

	for i := range oratabEntries {
		entry := oratabEntries[i]

		utils.RunRoutine(b.configuration, func() {
			databaseChannel <- b.getDatabase(entry, hostType)
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

func (b *CommonBuilder) getDatabase(entry model.OratabEntry, hostType string) *model.Database {
	dbStatusOut := b.fetcher.Execute("dbstatus", entry.DBName, entry.OracleHome)
	dbStatus := strings.TrimSpace(string(dbStatusOut))

	var database *model.Database

	switch dbStatus {
	case "OPEN":
		database = b.getOpenDatabase(entry, hostType)
	case "MOUNTED":
		{
			out := b.fetcher.Execute("dbmounted", entry.DBName, entry.OracleHome)
			tmp := marshal.Database(out)
			database = &tmp

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
	dbVersionOut := b.fetcher.Execute("dbversion", entry.DBName, entry.OracleHome)
	dbVersion := strings.Split(string(dbVersionOut), ".")[0]

	statsCtx, cancelStatsCtx := context.WithCancel(context.Background())
	if b.configuration.Forcestats {
		utils.RunRoutine(b.configuration, func() {
			b.fetcher.Execute("stats", entry.DBName, entry.OracleHome)

			cancelStatsCtx()
		})
	} else {
		cancelStatsCtx()
	}

	out := b.fetcher.Execute("db", entry.DBName, entry.OracleHome, strconv.Itoa(b.configuration.AWR))
	tmp := marshal.Database(out)
	database := &tmp

	var wg sync.WaitGroup

	utils.RunRoutineInGroup(b.configuration, func() {
		out := b.fetcher.Execute("tablespace", entry.DBName, entry.OracleHome)
		database.Tablespaces = marshal.Tablespaces(out)
	}, &wg)

	utils.RunRoutineInGroup(b.configuration, func() {
		out := b.fetcher.Execute("schema", entry.DBName, entry.OracleHome)
		database.Schemas = marshal.Schemas(out)
	}, &wg)

	utils.RunRoutineInGroup(b.configuration, func() {
		out := b.fetcher.Execute("patch", entry.DBName, dbVersion, entry.OracleHome)
		database.Patches = marshal.Patches(out)
	}, &wg)

	utils.RunRoutineInGroup(b.configuration, func() {
		<-statsCtx.Done()

		out := b.fetcher.Execute("feature", entry.DBName, dbVersion, entry.OracleHome)

		if strings.Contains(string(out), "deadlocked on readable physical standby") {
			log.Println("Detected bug active dataguard 2311894.1!")
			database.Features = []model.Feature{}

		} else if strings.Contains(string(out), "ORA-01555: snapshot too old: rollback segment number") {
			log.Println("Detected error on active dataguard ORA-01555!")
			database.Features = []model.Feature{}

		} else {
			database.Features = marshal.Features(out)
		}
	}, &wg)

	utils.RunRoutineInGroup(b.configuration, func() {
		<-statsCtx.Done()

		out := b.fetcher.Execute("opt", entry.DBName, dbVersion, entry.OracleHome)
		database.Features2 = marshal.Features2(out)
	}, &wg)

	utils.RunRoutineInGroup(b.configuration, func() {
		<-statsCtx.Done()

		out := b.fetcher.Execute("license", entry.DBName, dbVersion, hostType, entry.OracleHome)
		database.Licenses = marshal.Licenses(out)
	}, &wg)

	utils.RunRoutineInGroup(b.configuration, func() {
		out := b.fetcher.Execute("addm", entry.DBName, entry.OracleHome)
		database.ADDMs = marshal.Addms(out)
	}, &wg)

	utils.RunRoutineInGroup(b.configuration, func() {
		out := b.fetcher.Execute("segmentadvisor", entry.DBName, entry.OracleHome)
		database.SegmentAdvisors = marshal.SegmentAdvisor(out)
	}, &wg)

	utils.RunRoutineInGroup(b.configuration, func() {
		out := b.fetcher.Execute("psu", entry.DBName, dbVersion, entry.OracleHome)
		database.LastPSUs = marshal.PSU(out)
	}, &wg)

	utils.RunRoutineInGroup(b.configuration, func() {
		out := b.fetcher.Execute("backup", entry.DBName, entry.OracleHome)
		database.Backups = marshal.Backups(out)
	}, &wg)

	wg.Wait()

	return database
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
