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
	"sync"

	"github.com/ercole-io/ercole-agent/v2/agentmodel"
	"github.com/ercole-io/ercole-agent/v2/utils"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/hashicorp/go-version"
)

func (b *CommonBuilder) getOracleDatabaseFeature(host model.Host) *model.OracleDatabaseFeature {
	oracleDatabaseFeature := new(model.OracleDatabaseFeature)

	oratabEntries := b.fetcher.GetOracleDatabaseOratabEntries()
	oracleDatabaseFeature.UnlistedRunningDatabases = b.getUnlistedRunningOracleDBs(oratabEntries)

	oracleDatabaseFeature.Databases = b.getOracleDBs(oratabEntries, host)

	return oracleDatabaseFeature
}

func (b *CommonBuilder) getUnlistedRunningOracleDBs(listedDBs []agentmodel.OratabEntry) []string {
	runningDBs := b.fetcher.GetOracleDatabaseRunningDatabases()

	listedDBNames := make([]string, len(listedDBs))
	for i, db := range listedDBs {
		listedDBNames[i] = db.DBName
	}
	sort.Strings(listedDBNames)

	unlistedRunningDBs := make([]string, 0)
	for _, r := range runningDBs {
		if len(listedDBNames) == sort.SearchStrings(listedDBNames, r) {
			unlistedRunningDBs = append(unlistedRunningDBs, r)
		}
	}

	return unlistedRunningDBs
}

func (b *CommonBuilder) getOracleDBs(oratabEntries []agentmodel.OratabEntry, host model.Host) []model.OracleDatabase {

	databaseChannel := make(chan *model.OracleDatabase, len(oratabEntries))

	for i := range oratabEntries {
		entry := oratabEntries[i]

		utils.RunRoutine(b.configuration, func() {
			b.log.Debugf("oratab entry: [%v]", entry)

			databaseChannel <- b.getOracleDB(entry, host)
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

func (b *CommonBuilder) getOracleDB(entry agentmodel.OratabEntry, host model.Host) *model.OracleDatabase {
	dbStatus := b.fetcher.GetOracleDatabaseDbStatus(entry)
	var database *model.OracleDatabase

	switch dbStatus {
	case "OPEN":
		database = b.getOpenDatabase(entry, host.HardwareAbstractionTechnology)
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

			database.Licenses = computeLicenses(database.Edition(), database.CoreFactor(host), host.CPUCores)
		}
	default:
		b.log.Errorf("Unknown dbStatus: [%s] DBName: [%s] OracleHome: [%s]",
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

func computeLicenses(dbEdition string, coreFactor float64, cpuCores int) []model.OracleDatabaseLicense {
	licenses := make([]model.OracleDatabaseLicense, 0)
	numLicenses := coreFactor * float64(cpuCores)

	if dbEdition == "EXE" {
		licenses = append(licenses, model.OracleDatabaseLicense{
			Name:  "Oracle EXE",
			Count: numLicenses,
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
			Count: numLicenses,
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
			Count: numLicenses,
		})
	} else {
		licenses = append(licenses, model.OracleDatabaseLicense{
			Name:  "Oracle STD",
			Count: 0,
		})
	}

	return licenses
}
