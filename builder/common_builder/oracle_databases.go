// Copyright (c) 2022 Sorint.lab S.p.A.
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

	"github.com/ercole-io/ercole-agent/v2/agentmodel"
	"github.com/ercole-io/ercole-agent/v2/utils"
	"github.com/ercole-io/ercole/v2/model"
	ercutils "github.com/ercole-io/ercole/v2/utils"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-version"
)

func (b *CommonBuilder) getOracleDatabaseFeature(host model.Host, hostCoreFactor float64,
) (*model.OracleDatabaseFeature, error) {
	oracleDatabaseFeature := new(model.OracleDatabaseFeature)

	oratabEntries, err := b.fetcher.GetOracleDatabaseOratabEntries()
	if err != nil {
		b.log.Errorf("Can't get oratab entries")
		return nil, err
	}

	uniqueOratabEntries := b.RemoveDuplicatedOratabEntries(oratabEntries)

	oracleDatabaseFeature.UnlistedRunningDatabases = b.getUnlistedRunningOracleDBs(uniqueOratabEntries)

	oracleDatabaseFeature.Databases, err = b.getOracleDBs(uniqueOratabEntries, host, hostCoreFactor)

	return oracleDatabaseFeature, err
}

func (b *CommonBuilder) getUnlistedRunningOracleDBs(oratabEntries []agentmodel.OratabEntry) []string {
	runningDBs, err := b.fetcher.GetOracleDatabaseRunningDatabases()
	if err != nil {
		b.log.Errorf("Can't get running Oracle databases")
		return []string{}
	}

	oratabEntriesNames := make(map[string]bool, len(oratabEntries))
	for _, db := range oratabEntries {
		oratabEntriesNames[db.DBName] = true
	}

	unlistedRunningDBs := make([]string, 0)

	for _, runningDB := range runningDBs {
		if !oratabEntriesNames[runningDB] {
			unlistedRunningDBs = append(unlistedRunningDBs, runningDB)
		}
	}

	return unlistedRunningDBs
}

func (b *CommonBuilder) getOracleDBs(oratabEntries []agentmodel.OratabEntry, host model.Host, hostCoreFactor float64,
) ([]model.OracleDatabase, error) {
	databaseChan := make(chan *model.OracleDatabase, len(oratabEntries))
	errChan := make(chan error, len(oratabEntries))

	for i := range oratabEntries {
		entry := oratabEntries[i]

		utils.RunRoutine(b.configuration, func() {
			b.log.Debugf("oratab entry: [%v]", entry)
			database, err := b.getOracleDB(entry, host, hostCoreFactor)
			if err != nil && database == nil {
				b.log.Errorf("Oracle db, blocking error (db discarded): %s\n Errors: %s\n", entry, err)
				errChan <- err
				databaseChan <- nil
				return
			} else if err != nil {
				b.log.Warnf("Oracle db, non-blocking error: %s\n Errors: %s\n", entry, err)
				errChan <- err
			}

			databaseChan <- database
		})
	}

	var databases = []model.OracleDatabase{}

	for i := 0; i < len(oratabEntries); i++ {
		db := <-databaseChan
		if db != nil {
			databases = append(databases, *db)
		}
	}

	var merr error
	for len(errChan) > 0 {
		merr = multierror.Append(merr, <-errChan)
	}

	return databases, merr
}

func (b *CommonBuilder) getOracleDB(entry agentmodel.OratabEntry, host model.Host, hostCoreFactor float64) (*model.OracleDatabase, error) {
	dbStatus, err := b.fetcher.GetOracleDatabaseDbStatus(entry)
	if err != nil {
		b.log.Errorf("Oracle db [%s]: can't get db status, failed", entry.DBName)
		return nil, err
	}

	var database *model.OracleDatabase

	switch {
	case dbStatus == "READ WRITE" || dbStatus == "READ ONLY":
		database, err = b.getOpenDatabase(entry, host.HardwareAbstractionTechnology, hostCoreFactor)

	case dbStatus == "MOUNTED" || dbStatus == "READ ONLY WITH APPLY" || strings.Contains(dbStatus, "ORA-03170") || strings.Contains(dbStatus, "ORA-01555"):
		database, err = b.getMountedDatabase(entry, host, hostCoreFactor)

	case dbStatus == "unreachable!":
		b.log.Infof("dbStatus: [%s] DBName: [%s] OracleHome: [%s]",
			dbStatus, entry.DBName, entry.OracleHome)
		return nil, nil

	default:
		if strings.Contains(dbStatus, "ORA-01034") {
			b.log.Debugf("Connection Error: DBName: [%s] OracleHome: [%s]", entry.DBName, entry.OracleHome)
			return nil, nil
		}

		err := ercutils.NewErrorf("Unknown dbStatus: [%s] DBName: [%s] OracleHome: [%s]",
			dbStatus, entry.DBName, entry.OracleHome)

		return nil, err
	}

	if database != nil && err == nil {
		grantsDba, errGrant := b.fetcher.GetOracleDatabaseGrantsDba(entry)
		if errGrant != nil {
			b.log.Errorf("Oracle db [%s]: can't get dba grants, failed", entry.DBName)
			return nil, errGrant
		}

		database.GrantDba = grantsDba
	}

	return database, err
}

func (b *CommonBuilder) getOpenDatabase(entry agentmodel.OratabEntry, hardwareAbstractionTechnology string,
	hostCoreFactor float64) (*model.OracleDatabase, error) {
	stringDbVersion, err := b.fetcher.GetOracleDatabaseDbVersion(entry)
	if err != nil {
		b.log.Errorf("Oracle db [%s]: can't get db version, failed", entry.DBName)
		return nil, err
	}

	dbVersion, err := version.NewVersion(stringDbVersion)
	if err != nil {
		err = ercutils.NewErrorf("Can't parse db version [%s]: %w", stringDbVersion, err)

		b.log.Error(err)

		return nil, err
	}

	blockingErrs := make(chan error, 4)    // database errs are serious, must not be returned
	nonBlockingErrs := make(chan error, 8) // database errs, but not blocking ones

	statsCtx, cancelStatsCtx := context.WithCancel(context.Background())

	if b.configuration.Features.OracleDatabase.Forcestats {
		utils.RunRoutine(b.configuration, func() {
			if err := b.fetcher.RunOracleDatabaseStats(entry); err != nil {
				blockingErrs <- err
			}

			cancelStatsCtx()
		})
	} else {
		cancelStatsCtx()
	}

	database, err := b.fetcher.GetOracleDatabaseOpenDb(entry)
	if err != nil {
		b.log.Errorf("Oracle db [%s]: can't get open db, failed", entry.DBName)
		return nil, err
	}

	var wg sync.WaitGroup

	utils.RunRoutineInGroup(b.configuration, func() {
		if err := b.setPDBs(database, *dbVersion, entry); err != nil {
			nonBlockingErrs <- err
		}
	}, &wg)

	utils.RunRoutineInGroup(b.configuration, func() {
		if database.Tablespaces, err = b.fetcher.GetOracleDatabaseTablespaces(entry); err != nil {
			database.Tablespaces = []model.OracleDatabaseTablespace{}
			b.log.Warnf("Oracle db [%s]: can't get tablespaces", entry.DBName)
			nonBlockingErrs <- err
		}
	}, &wg)

	utils.RunRoutineInGroup(b.configuration, func() {
		if database.Schemas, err = b.fetcher.GetOracleDatabaseSchemas(entry); err != nil {
			database.Schemas = []model.OracleDatabaseSchema{}
			b.log.Warnf("Oracle db [%s]: can't get schemas", entry.DBName)
			nonBlockingErrs <- err
		}
	}, &wg)

	utils.RunRoutineInGroup(b.configuration, func() {
		database.Patches, err = b.fetcher.GetOracleDatabasePatches(entry, stringDbVersion)
		if err != nil && database.Patches != nil {
			b.log.Warnf("Oracle db [%s]: some patches have not passed", entry.DBName)
			nonBlockingErrs <- err
			return
		}

		if err != nil {
			database.Patches = []model.OracleDatabasePatch{}
			b.log.Warnf("Oracle db [%s]: can't get patches", entry.DBName)
			nonBlockingErrs <- err
			return
		}
	}, &wg)

	utils.RunRoutineInGroup(b.configuration, func() {
		<-statsCtx.Done()

		if database.FeatureUsageStats, err = b.fetcher.GetOracleDatabaseFeatureUsageStat(entry, stringDbVersion); err != nil {
			b.log.Errorf("Oracle db [%s]: can't get feature usage stat", entry.DBName)
			blockingErrs <- err
		}
	}, &wg)

	utils.RunRoutineInGroup(b.configuration, func() {
		<-statsCtx.Done()

		database.Licenses, err = b.fetcher.GetOracleDatabaseLicenses(entry, stringDbVersion, hardwareAbstractionTechnology, hostCoreFactor)
		if err != nil {
			b.log.Errorf("Oracle db [%s]: can't get licenses, failed", entry.DBName)
			blockingErrs <- err
		}
	}, &wg)

	utils.RunRoutineInGroup(b.configuration, func() {
		if database.ADDMs, err = b.fetcher.GetOracleDatabaseADDMs(entry); err != nil {
			database.ADDMs = []model.OracleDatabaseAddm{}
			b.log.Errorf("Oracle db [%s]: can't get ADDMs, failed", entry.DBName)
			nonBlockingErrs <- err
		}
	}, &wg)

	utils.RunRoutineInGroup(b.configuration, func() {
		if database.SegmentAdvisors, err = b.fetcher.GetOracleDatabaseSegmentAdvisors(entry); err != nil {
			database.SegmentAdvisors = []model.OracleDatabaseSegmentAdvisor{}
			b.log.Warnf("Oracle db [%s]: can't get segment advisors", entry.DBName)
			nonBlockingErrs <- err
		}
	}, &wg)

	utils.RunRoutineInGroup(b.configuration, func() {
		if database.PSUs, err = b.fetcher.GetOracleDatabasePSUs(entry, stringDbVersion); err != nil {
			b.log.Errorf("Oracle db [%s]: can't get PSUs, failed", entry.DBName)
			blockingErrs <- err
		}
	}, &wg)

	utils.RunRoutineInGroup(b.configuration, func() {
		if database.Backups, err = b.fetcher.GetOracleDatabaseBackups(entry); err != nil {
			database.Backups = []model.OracleDatabaseBackup{}
			b.log.Warnf("Oracle db [%s]: can't get backups", entry.DBName)
			nonBlockingErrs <- err
		}
	}, &wg)

	utils.RunRoutineInGroup(b.configuration, func() {
		if database.Services, err = b.fetcher.GetOracleDatabaseServices(entry); err != nil {
			database.Services = []model.OracleDatabaseService{}
			b.log.Warnf("Oracle db [%s]: can't get services", entry.DBName)
			nonBlockingErrs <- err
		}
	}, &wg)

	wg.Wait()
	close(blockingErrs)
	close(nonBlockingErrs)

	var merr error

	for err := range nonBlockingErrs {
		merr = multierror.Append(merr, err)
	}

	if len(blockingErrs) > 0 {
		for err := range blockingErrs {
			merr = multierror.Append(merr, err)
		}

		return nil, merr
	}

	database.Version = checkVersion(database.Name, database.Version)

	return database, merr
}

func (b *CommonBuilder) setPDBs(database *model.OracleDatabase, dbVersion version.Version, entry agentmodel.OratabEntry) error {
	database.PDBs = []model.OracleDatabasePluggableDatabase{}

	v2, errVersion := version.NewVersion("11.2.0.4.0")
	if errVersion != nil {
		return errVersion
	}

	if dbVersion.LessThan(v2) {
		database.IsCDB = false
		return nil
	}

	var err error

	var segmentsSize, datafileSize, allocable, totalSegmentsSize, totalDatafileSize, totalAllocable float64

	if database.IsCDB, err = b.fetcher.GetOracleDatabaseCheckPDB(entry); err != nil {
		database.IsCDB = false

		b.log.Warnf("Oracle db [%s]: can't check PDB", entry.DBName)

		return err
	}

	if !database.IsCDB {
		return nil
	}

	if database.PDBs, err = b.fetcher.GetOracleDatabasePDBs(entry); err != nil {
		b.log.Warnf("Oracle db [%s]: can't get PDBs", entry.DBName)
		return err
	}

	var wg sync.WaitGroup

	errChan := make(chan error, 2*len(database.PDBs))

	for i := range database.PDBs {
		pdb := &database.PDBs[i]

		utils.RunRoutineInGroup(b.configuration, func() {
			if segmentsSize, datafileSize, allocable, err = b.fetcher.GetOracleDatabasePDBSize(entry); err != nil {
				b.log.Warnf("Oracle db [%s]: can't get PDB [%s] size", entry.DBName, pdb.Name)
				errChan <- err
			}
		}, &wg)

		pdb.SegmentsSize = segmentsSize
		pdb.DatafileSize = datafileSize
		pdb.Allocable = allocable

		totalSegmentsSize += segmentsSize
		totalDatafileSize += datafileSize
		totalAllocable += allocable

		utils.RunRoutineInGroup(b.configuration, func() {
			if pdb.Tablespaces, err = b.fetcher.GetOracleDatabasePDBTablespaces(entry, pdb.Name); err != nil {
				b.log.Warnf("Oracle db [%s]: can't get PDB [%s] tablespaces", entry.DBName, pdb.Name)
				errChan <- err
			}
		}, &wg)

		utils.RunRoutineInGroup(b.configuration, func() {
			if pdb.Schemas, err = b.fetcher.GetOracleDatabasePDBSchemas(entry, pdb.Name); err != nil {
				b.log.Warnf("Oracle db [%s]: can't get PDB [%s] schemas", entry.DBName, pdb.Name)
				errChan <- err
			}
		}, &wg)

		utils.RunRoutineInGroup(b.configuration, func() {
			if pdb.GrantDba, err = b.fetcher.GetOracleDatabaseGrantsDba(entry); err != nil {
				b.log.Warnf("Oracle db [%s]: can't get PDB [%s] grants", entry.DBName, pdb.Name)
				errChan <- err
			}
		}, &wg)
	}

	database.SegmentsSize = totalSegmentsSize
	database.DatafileSize = totalDatafileSize
	database.Allocable = totalAllocable

	wg.Wait()

	if len(errChan) > 0 {
		var merr error

		for len(errChan) > 0 {
			merr = multierror.Append(merr, <-errChan)
		}

		return merr
	}

	return nil
}

func (b *CommonBuilder) getMountedDatabase(entry agentmodel.OratabEntry, host model.Host, hostCoreFactor float64,
) (*model.OracleDatabase, error) {
	database, err := b.fetcher.GetOracleDatabaseMountedDb(entry)
	if err != nil {
		b.log.Errorf("Oracle db [%s]: can't get mounted db, failed", entry.DBName)
		return nil, err
	}

	isRac, err := b.fetcher.GetOracleDatabaseRac(entry)
	if err != nil {
		b.log.Errorf("Oracle db [%s]: can't get rac information, failed", entry.DBName)
		return nil, err
	}

	if isRac == "TRUE" {
		database.IsRAC = true
	}

	database.Version = checkVersion(database.Name, database.Version)
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

	database.Licenses = make([]model.OracleDatabaseLicense, 0)
	if database.Edition() != model.OracleDatabaseEditionExpress {
		coreFactor, err := database.CoreFactor(host, hostCoreFactor)
		if err != nil {
			b.log.Errorf("Oracle db [%s]: can't calculate coreFactor, failed", entry.DBName)
			return nil, err
		}

		database.Licenses = computeLicenses(database.Edition(), coreFactor, host.CPUCores)
	}

	return database, nil
}

func (b *CommonBuilder) RemoveDuplicatedOratabEntries(oratabEntries []agentmodel.OratabEntry) []agentmodel.OratabEntry {
	m := map[agentmodel.OratabEntry]struct{}{}
	uniqueOratabEntries := []agentmodel.OratabEntry{}

	for _, d := range oratabEntries {
		if _, ok := m[d]; !ok {
			uniqueOratabEntries = append(uniqueOratabEntries, d)
			m[d] = struct{}{}
		} else {
			b.log.Warnf("Duplicated oratab entries %s", d.DBName)
		}
	}

	return uniqueOratabEntries
}

func computeLicenses(dbEdition string, coreFactor float64, cpuCores int) []model.OracleDatabaseLicense {
	licenses := make([]model.OracleDatabaseLicense, 0)
	numLicenses := coreFactor * float64(cpuCores)

	editions := []struct {
		name      string
		dbEdition string
	}{
		{
			name:      "Oracle EXE",
			dbEdition: model.OracleDatabaseEditionExtreme,
		},
		{
			name:      "Oracle ENT",
			dbEdition: model.OracleDatabaseEditionEnterprise,
		},
		{
			name:      "Oracle STD",
			dbEdition: model.OracleDatabaseEditionStandard,
		},
	}

	for _, edition := range editions {
		license := model.OracleDatabaseLicense{
			Name: edition.name,
		}

		if dbEdition == edition.dbEdition {
			license.Count = numLicenses
		}

		licenses = append(licenses, license)
	}

	return licenses
}

func checkVersion(dbName, dbVersion string) string {
	if strings.Contains(strings.ToUpper(dbVersion), "ENTERPRISE") {
		return dbVersion
	}

	if dbName != "XE" {
		return dbVersion
	}

	return "Express Edition"
}
