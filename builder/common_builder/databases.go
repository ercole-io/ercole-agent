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

package common

import (
	"context"
	"strings"
	"sync"

	"github.com/ercole-io/ercole-agent/model"
	"github.com/ercole-io/ercole-agent/utils"
)

func (b *CommonBuilder) getOracleDBs(hostType string, cpuCores int, socket int) []model.Database {
	oratabEntries := b.fetcher.GetOratabEntries()

	databaseChannel := make(chan *model.Database, len(oratabEntries))

	for i := range oratabEntries {
		entry := oratabEntries[i]

		utils.RunRoutine(b.configuration, func() {
			b.log.Debugf("oratab entry: [%v]", entry)

			databaseChannel <- b.getOracleDB(entry, hostType, cpuCores, socket)
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

func (b *CommonBuilder) getOracleDB(entry model.OratabEntry, hostType string, cpuCores int, socket int) *model.Database {
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

			// compute db edition
			var dbEdition string
			if strings.Contains(strings.ToUpper(database.Version), "ENTERPRISE") {
				dbEdition = "ENT"
			} else if strings.Contains(strings.ToUpper(database.Version), "EXTREME") {
				dbEdition = "EXE"
			} else {
				dbEdition = "STD"
			}

			// compute coreFactor/factor
			coreFactor := float32(-1)
			if hostType == "OVM" || hostType == "VMWARE" || hostType == "VMOTHER" {
				if dbEdition == "EXE" || dbEdition == "ENT" {
					coreFactor = float32(cpuCores) * 0.25
				} else if dbEdition == "STD" {
					coreFactor = 0
				}
			} else if hostType == "PH" {
				if dbEdition == "EXE" || dbEdition == "ENT" {
					coreFactor = float32(cpuCores) * 0.25
				} else if dbEdition == "STD" {
					coreFactor = float32(socket)
				}
			}

			if dbEdition == "EXE" {
				database.Licenses = append(database.Licenses, model.License{
					Name:  "Oracle EXE",
					Count: coreFactor,
				})
			} else {
				database.Licenses = append(database.Licenses, model.License{
					Name:  "Oracle EXE",
					Count: 0,
				})
			}

			if dbEdition == "ENT" {
				database.Licenses = append(database.Licenses, model.License{
					Name:  "Oracle ENT",
					Count: coreFactor,
				})
			} else {
				database.Licenses = append(database.Licenses, model.License{
					Name:  "Oracle ENT",
					Count: 0,
				})
			}

			if dbEdition == "STD" {
				database.Licenses = append(database.Licenses, model.License{
					Name:  "Oracle STD",
					Count: coreFactor,
				})
			} else {
				database.Licenses = append(database.Licenses, model.License{
					Name:  "Oracle STD",
					Count: 0,
				})
			}
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
