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
	"sort"
	"strings"
	"sync"

	"github.com/ercole-io/ercole-agent/agentmodel"
	"github.com/ercole-io/ercole-agent/utils"
	"github.com/ercole-io/ercole/model"
	"github.com/hashicorp/go-version"
)

func (b *CommonBuilder) getOracleDatabaseFeature(hardwareAbstractionTechnology string, cpuCores int, cpuSockets int) *model.OracleDatabaseFeature {
	oracleDatabaseFeature := new(model.OracleDatabaseFeature)
	oracleDatabaseFeature.Databases = b.getOracleDBs(
		hardwareAbstractionTechnology,
		cpuCores,
		cpuSockets,
	)
	oracleDatabaseFeature.UnlistedRunningDatabases = b.getUnlistedRunningOracleDBs(oracleDatabaseFeature.Databases)

	return oracleDatabaseFeature
}

func (b *CommonBuilder) getUnlistedRunningOracleDBs(listedDBs []model.OracleDatabase) []string {
	runningDBs := b.fetcher.GetOracleDatabaseRunningDatabases()

	// copy listedDBs names to listedDBNames
	listedDBNames := make([]string, len(listedDBs))
	for i, s := range listedDBs {
		listedDBNames[i] = s.Name
	}
	sort.Strings(listedDBNames)

	// make the subtraction
	unlistedRunningDBs := make([]string, 0)
	for _, r := range runningDBs {
		if len(listedDBNames) == sort.SearchStrings(listedDBNames, r) {
			unlistedRunningDBs = append(unlistedRunningDBs, r)
		}
	}

	return unlistedRunningDBs
}

func (b *CommonBuilder) getOracleDBs(hardwareAbstractionTechnology string, cpuCores int, cpuSockets int) []model.OracleDatabase {
	oratabEntries := b.fetcher.GetOracleDatabaseOratabEntries()

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
	dbStatus := b.fetcher.GetOracleDatabaseDbStatus(entry)
	var database *model.OracleDatabase

	switch dbStatus {
	case "OPEN":
		database = b.getOpenDatabase(entry, hardwareAbstractionTechnology)
	case "MOUNTED":
		{
			db := b.fetcher.GetOracleDatabaseMountedDb(entry)
			database = &db

			database.Tablespaces = []model.OracleDatabaseTablespace{}
			database.Schemas = []model.OracleDatabaseSchema{}
			database.Patches = []model.OracleDatabasePatch{}
			database.ADDMs = []model.OracleDatabaseAddm{}
			database.SegmentAdvisors = []model.OracleDatabaseSegmentAdvisor{}
			database.PSUs = []model.OracleDatabasePSU{}
			database.Backups = []model.OracleDatabaseBackup{}
			database.PDBs = []model.OracleDatabasePluggableDatabase{}
			database.Services = []model.OracleDatabaseService{}
			database.FeatureUsageStats = []model.OracleDatabaseFeatureUsageStat{}

			dbEdition := computeDBEdition(database.Version)
			coreFactor := computeCoreFactor(cpuCores, cpuSockets, hardwareAbstractionTechnology, dbEdition)
			database.Licenses = computeLicenses(dbEdition, coreFactor)
		}
	default:
		b.log.Warnf("Unknown dbStatus: [%s] DBName: [%s] OracleHome: [%s]",
			dbStatus, entry.DBName, entry.OracleHome)
		return nil
	}

	return database
}

func (b *CommonBuilder) getOpenDatabase(entry agentmodel.OratabEntry, hardwareAbstractionTechnology string) *model.OracleDatabase {
	stringDbVersion := b.fetcher.GetOracleDatabaseDbVersion(entry)

	dbVersion, err := version.NewVersion(stringDbVersion)
	if err != nil {
		panic(err)
	}

	statsCtx, cancelStatsCtx := context.WithCancel(context.Background())
	if b.configuration.Features.OracleDatabase.Forcestats {
		utils.RunRoutine(b.configuration, func() {
			b.fetcher.RunOracleDatabaseStats(entry)

			cancelStatsCtx()
		})
	} else {
		cancelStatsCtx()
	}

	database := b.fetcher.GetOracleDatabaseOpenDb(entry)
	var wg sync.WaitGroup

	utils.RunRoutineInGroup(b.configuration, func() {
		b.setPDBs(&database, *dbVersion, entry)
	}, &wg)

	utils.RunRoutineInGroup(b.configuration, func() {
		database.Tablespaces = b.fetcher.GetOracleDatabaseTablespaces(entry)
	}, &wg)

	utils.RunRoutineInGroup(b.configuration, func() {
		database.Schemas = b.fetcher.GetOracleDatabaseSchemas(entry)
	}, &wg)

	utils.RunRoutineInGroup(b.configuration, func() {
		database.Patches = b.fetcher.GetOracleDatabasePatches(entry, stringDbVersion)
	}, &wg)

	utils.RunRoutineInGroup(b.configuration, func() {
		<-statsCtx.Done()

		database.FeatureUsageStats = b.fetcher.GetOracleDatabaseFeatureUsageStat(entry, stringDbVersion)
	}, &wg)

	utils.RunRoutineInGroup(b.configuration, func() {
		<-statsCtx.Done()

		database.Licenses = b.fetcher.GetOracleDatabaseLicenses(entry, stringDbVersion, hardwareAbstractionTechnology)
	}, &wg)

	utils.RunRoutineInGroup(b.configuration, func() {
		database.ADDMs = b.fetcher.GetOracleDatabaseADDMs(entry)
	}, &wg)

	utils.RunRoutineInGroup(b.configuration, func() {
		database.SegmentAdvisors = b.fetcher.GetOracleDatabaseSegmentAdvisors(entry)
	}, &wg)

	utils.RunRoutineInGroup(b.configuration, func() {
		database.PSUs = b.fetcher.GetOracleDatabasePSUs(entry, stringDbVersion)
	}, &wg)

	utils.RunRoutineInGroup(b.configuration, func() {
		database.Backups = b.fetcher.GetOracleDatabaseBackups(entry)
	}, &wg)

	wg.Wait()

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

func (b *CommonBuilder) setPDBs(database *model.OracleDatabase, dbVersion version.Version, entry agentmodel.OratabEntry) {
	database.PDBs = []model.OracleDatabasePluggableDatabase{}

	v2, _ := version.NewVersion("11.2.0.4.0")
	if dbVersion.LessThan(v2) {
		database.IsCDB = false
		return
	}

	database.IsCDB = b.fetcher.GetOracleDatabaseCheckPDB(entry)

	if database.IsCDB {
		database.PDBs = b.fetcher.GetOracleDatabasePDBs(entry)

		var wg sync.WaitGroup

		for i := range database.PDBs {
			var pdb *model.OracleDatabasePluggableDatabase = &database.PDBs[i]

			utils.RunRoutineInGroup(b.configuration, func() {
				pdb.Tablespaces = b.fetcher.GetOracleDatabasePDBTablespaces(entry, pdb.Name)
			}, &wg)

			utils.RunRoutineInGroup(b.configuration, func() {
				pdb.Schemas = b.fetcher.GetOracleDatabasePDBSchemas(entry, pdb.Name)
			}, &wg)
		}

		wg.Wait()
	}
}

func computeDBEdition(version string) (dbEdition string) {
	if strings.Contains(strings.ToUpper(version), "ENTERPRISE") {
		dbEdition = "ENT"
	} else if strings.Contains(strings.ToUpper(version), "EXTREME") {
		dbEdition = "EXE"
	} else {
		dbEdition = "STD"
	}

	return
}

func computeCoreFactor(cpuCores, cpuSockets int, hardwareAbstractionTechnology, dbEdition string) float64 {
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

	return coreFactor
}

func computeLicenses(dbEdition string, coreFactor float64) []model.OracleDatabaseLicense {
	licenses := make([]model.OracleDatabaseLicense, 0)

	if dbEdition == "EXE" {
		licenses = append(licenses, model.OracleDatabaseLicense{
			Name:  "Oracle EXE",
			Count: coreFactor,
		})
	} else {
		licenses = append(licenses, model.OracleDatabaseLicense{
			Name:  "Oracle EXE",
			Count: 0,
		})
	}

	if dbEdition == "ENT" {
		licenses = append(licenses, model.OracleDatabaseLicense{
			Name:  "Oracle ENT",
			Count: coreFactor,
		})
	} else {
		licenses = append(licenses, model.OracleDatabaseLicense{
			Name:  "Oracle ENT",
			Count: 0,
		})
	}

	if dbEdition == "STD" {
		licenses = append(licenses, model.OracleDatabaseLicense{
			Name:  "Oracle STD",
			Count: coreFactor,
		})
	} else {
		licenses = append(licenses, model.OracleDatabaseLicense{
			Name:  "Oracle STD",
			Count: 0,
		})
	}

	return licenses
}
