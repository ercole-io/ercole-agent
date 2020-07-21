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

	"github.com/ercole-io/ercole-agent/agentmodel"
	"github.com/ercole-io/ercole-agent/utils"
	"github.com/ercole-io/ercole/model"
)

func (b *CommonBuilder) getOracleDBs(hardwareAbstractionTechnology string, cpuCores int, cpuSockets int) []model.OracleDatabase {
	oratabEntries := b.fetcher.GetOracleOratabEntries()

	databaseChannel := make(chan *model.OracleDatabase, len(oratabEntries))

	for i := range oratabEntries {
		entry := oratabEntries[i]

		utils.RunRoutine(b.configuration, func() {
			b.log.Debugf("oratab entry: [%v]", entry)

			databaseChannel <- b.getOracleDB(entry, hardwareAbstractionTechnology, cpuCores, cpuSockets)
		})
	}

	var databases = []model.OracleDatabase{}
	for i := 0; i < len(oratabEntries); i++ {
		db := (<-databaseChannel)
		if db != nil {
			databases = append(databases, *db)
		}
	}

	return databases
}

func (b *CommonBuilder) getOracleDB(entry agentmodel.OratabEntry, hardwareAbstractionTechnology string, cpuCores, cpuSockets int) *model.OracleDatabase {
	dbStatus := b.fetcher.GetOracleDbStatus(entry)
	var database *model.OracleDatabase

	switch dbStatus {
	case "OPEN":
		database = b.getOpenDatabase(entry, hardwareAbstractionTechnology)
	case "MOUNTED":
		{
			db := b.fetcher.GetOracleMountedDb(entry)
			database = &db

			database.Tablespaces = []model.OracleDatabaseTablespace{}
			database.Schemas = []model.OracleDatabaseSchema{}
			database.Patches = []model.OracleDatabasePatch{}
			database.Licenses = []model.OracleDatabaseLicense{}
			database.ADDMs = []model.OracleDatabaseAddm{}
			database.SegmentAdvisors = []model.OracleDatabaseSegmentAdvisor{}
			database.PSUs = []model.OracleDatabasePSU{}
			database.Backups = []model.OracleDatabaseBackup{}
			database.PDBs = []model.OracleDatabasePluggableDatabase{}
			database.Services = []model.OracleDatabaseService{}
			database.FeatureUsageStats = []model.OracleDatabaseFeatureUsageStat{}

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
			coreFactor := float64(-1)
			if hardwareAbstractionTechnology == "OVM" || hardwareAbstractionTechnology == "VMWARE" || hardwareAbstractionTechnology == "VMOTHER" {
				if dbEdition == "EXE" || dbEdition == "ENT" {
					coreFactor = float64(cpuCores) * 0.25
				} else if dbEdition == "STD" {
					coreFactor = 0
				}
			} else if hardwareAbstractionTechnology == "PH" {
				if dbEdition == "EXE" || dbEdition == "ENT" {
					coreFactor = float64(cpuCores) * 0.25
				} else if dbEdition == "STD" {
					coreFactor = float64(cpuSockets)
				}
			}

			if dbEdition == "EXE" {
				database.Licenses = append(database.Licenses, model.OracleDatabaseLicense{
					Name:  "Oracle EXE",
					Count: coreFactor,
				})
			} else {
				database.Licenses = append(database.Licenses, model.OracleDatabaseLicense{
					Name:  "Oracle EXE",
					Count: 0,
				})
			}

			if dbEdition == "ENT" {
				database.Licenses = append(database.Licenses, model.OracleDatabaseLicense{
					Name:  "Oracle ENT",
					Count: coreFactor,
				})
			} else {
				database.Licenses = append(database.Licenses, model.OracleDatabaseLicense{
					Name:  "Oracle ENT",
					Count: 0,
				})
			}

			if dbEdition == "STD" {
				database.Licenses = append(database.Licenses, model.OracleDatabaseLicense{
					Name:  "Oracle STD",
					Count: coreFactor,
				})
			} else {
				database.Licenses = append(database.Licenses, model.OracleDatabaseLicense{
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

func (b *CommonBuilder) getOpenDatabase(entry agentmodel.OratabEntry, hardwareAbstractionTechnology string) *model.OracleDatabase {
	dbVersion := b.fetcher.GetOracleDbVersion(entry)

	statsCtx, cancelStatsCtx := context.WithCancel(context.Background())
	if b.configuration.Features.OracleDatabase.Forcestats {
		utils.RunRoutine(b.configuration, func() {
			b.fetcher.RunOracleStats(entry)

			cancelStatsCtx()
		})
	} else {
		cancelStatsCtx()
	}

	database := b.fetcher.GetOracleOpenDb(entry)

	var wg sync.WaitGroup

	utils.RunRoutineInGroup(b.configuration, func() {
		database.Tablespaces = b.fetcher.GetOracleTablespaces(entry)
	}, &wg)

	utils.RunRoutineInGroup(b.configuration, func() {
		database.Schemas = b.fetcher.GetOracleSchemas(entry)
	}, &wg)

	utils.RunRoutineInGroup(b.configuration, func() {
		database.Patches = b.fetcher.GetOraclePatches(entry, dbVersion)
	}, &wg)

	utils.RunRoutineInGroup(b.configuration, func() {
		<-statsCtx.Done()

		database.FeatureUsageStats = b.fetcher.GetOracleDatabaseFeatureUsageStat(entry, dbVersion)
	}, &wg)

	utils.RunRoutineInGroup(b.configuration, func() {
		<-statsCtx.Done()

		database.Licenses = b.fetcher.GetOracleLicenses(entry, dbVersion, hardwareAbstractionTechnology)
	}, &wg)

	utils.RunRoutineInGroup(b.configuration, func() {
		database.ADDMs = b.fetcher.GetOracleADDMs(entry)
	}, &wg)

	utils.RunRoutineInGroup(b.configuration, func() {
		database.SegmentAdvisors = b.fetcher.GetOracleSegmentAdvisors(entry)
	}, &wg)

	utils.RunRoutineInGroup(b.configuration, func() {
		database.PSUs = b.fetcher.GetOraclePSUs(entry, dbVersion)
	}, &wg)

	utils.RunRoutineInGroup(b.configuration, func() {
		database.Backups = b.fetcher.GetOracleBackups(entry)
	}, &wg)

	wg.Wait()

	database.PDBs = []model.OracleDatabasePluggableDatabase{}
	database.Services = []model.OracleDatabaseService{}

	return &database
}

func (b *CommonBuilder) getDatabasesAndSchemaNames(databases []model.OracleDatabase) (databasesNames, schemasNames string) {
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
