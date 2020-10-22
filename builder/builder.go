package builder

import (
	"context"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/ercole-io/ercole-agent/config"
	"github.com/ercole-io/ercole-agent/marshal"
	"github.com/ercole-io/ercole-agent/model"
	"github.com/ercole-io/ercole-agent/utils"
)

func BuildData(configuration config.Configuration, version string, hostDataSchemaVersion int) *model.HostData {
	hostData := new(model.HostData)

	hostData.Environment = configuration.Envtype
	hostData.Location = configuration.Location
	hostData.HostType = configuration.HostType
	hostData.Version = version
	hostData.HostDataSchemaVersion = hostDataSchemaVersion
	hostData.Info = *getHost(configuration)
	hostData.Hostname = hostData.Info.Hostname
	if configuration.Hostname != "default" {
		hostData.Hostname = configuration.Hostname
	}

	hostData.Extra.Filesystems = getFilesystems(configuration)
	hostData.Extra.Databases = getDatabases(configuration, hostData.Info)

	hostData.Databases, hostData.Schemas = getDatabasesAndSchemaNames(hostData.Extra.Databases)

	return hostData
}

func getHost(configuration config.Configuration) *model.Host {
	out := fetcher(configuration, "host")
	host := marshal.Host(out)

	host.Environment = configuration.Envtype
	host.Location = configuration.Location

	return &host
}

func getFilesystems(configuration config.Configuration) []model.Filesystem {
	out := fetcher(configuration, "filesystem")
	return marshal.Filesystems(out)
}

func getDatabases(configuration config.Configuration, host model.Host) []model.Database {
	out := fetcher(configuration, "oratab", configuration.Oratab)
	oratabEntries := marshal.Oratab(out)

	databaseChannel := make(chan *model.Database, len(oratabEntries))

	for i := range oratabEntries {
		entry := oratabEntries[i]

		utils.RunRoutine(configuration, func() {
			databaseChannel <- getDatabase(configuration, entry, host)
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

func getDatabase(configuration config.Configuration, entry model.OratabEntry, host model.Host) *model.Database {
	dbStatusOut := fetcher(configuration, "dbstatus", entry.DBName, entry.OracleHome)
	dbStatus := strings.TrimSpace(string(dbStatusOut))

	var database *model.Database

	switch dbStatus {
	case "OPEN":
		database = getOpenDatabase(configuration, entry, host.Type)

	case "MOUNTED":
		{
			out := fetcher(configuration, "dbmounted", entry.DBName, entry.OracleHome)
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

			dbEdition := computeDBEdition(database.Version)

			coreFactor := computeCoreFactor(host, dbEdition)

			database.Licenses = computeLicenses(dbEdition, coreFactor)
		}

	default:
		log.Println("Error! DBName: [" + entry.DBName + "] OracleHome: [" + entry.OracleHome + "]  Wrong dbStatus: [" + dbStatus + "]")
		return nil
	}

	return database
}

func getOpenDatabase(configuration config.Configuration, entry model.OratabEntry, hostType string) *model.Database {
	dbVersionOut := fetcher(configuration, "dbversion", entry.DBName, entry.OracleHome)
	dbVersion := strings.Split(string(dbVersionOut), ".")[0]

	statsCtx, cancelStatsCtx := context.WithCancel(context.Background())
	if configuration.Forcestats {
		utils.RunRoutine(configuration, func() {
			fetcher(configuration, "stats", entry.DBName, entry.OracleHome)

			cancelStatsCtx()
		})
	} else {
		cancelStatsCtx()
	}

	out := fetcher(configuration, "db", entry.DBName, entry.OracleHome, strconv.Itoa(configuration.AWR))
	tmp := marshal.Database(out)
	database := &tmp

	var wg sync.WaitGroup

	utils.RunRoutineInGroup(configuration, func() {
		out := fetcher(configuration, "tablespace", entry.DBName, entry.OracleHome)
		database.Tablespaces = marshal.Tablespaces(out)
	}, &wg)

	utils.RunRoutineInGroup(configuration, func() {
		out := fetcher(configuration, "schema", entry.DBName, entry.OracleHome)
		database.Schemas = marshal.Schemas(out)
	}, &wg)

	utils.RunRoutineInGroup(configuration, func() {
		out := fetcher(configuration, "patch", entry.DBName, dbVersion, entry.OracleHome)
		database.Patches = marshal.Patches(out)
	}, &wg)

	utils.RunRoutineInGroup(configuration, func() {
		<-statsCtx.Done()

		out := fetcher(configuration, "opt", entry.DBName, dbVersion, entry.OracleHome)
		database.Features2 = marshal.Features2(out)
	}, &wg)

	utils.RunRoutineInGroup(configuration, func() {
		<-statsCtx.Done()

		out := fetcher(configuration, "license", entry.DBName, dbVersion, hostType, entry.OracleHome)
		database.Licenses = marshal.Licenses(out)

		database.Features = make([]model.Feature, 0)
		for _, fe := range database.Licenses {
			if fe.Name == "Oracle EXE" || fe.Name == "Oracle ENT" || fe.Name == "Oracle STD" {
				continue
			}
			database.Features = append(database.Features, model.Feature{
				Name:   fe.Name,
				Status: fe.Count > 0,
			})
		}
	}, &wg)

	utils.RunRoutineInGroup(configuration, func() {
		out := fetcher(configuration, "addm", entry.DBName, entry.OracleHome)
		database.ADDMs = marshal.Addms(out)
	}, &wg)

	utils.RunRoutineInGroup(configuration, func() {
		out := fetcher(configuration, "segmentadvisor", entry.DBName, entry.OracleHome)
		database.SegmentAdvisors = marshal.SegmentAdvisor(out)
	}, &wg)

	utils.RunRoutineInGroup(configuration, func() {
		out := fetcher(configuration, "psu", entry.DBName, dbVersion, entry.OracleHome)
		database.LastPSUs = marshal.PSU(out)
	}, &wg)

	utils.RunRoutineInGroup(configuration, func() {
		out := fetcher(configuration, "backup", entry.DBName, entry.OracleHome)
		database.Backups = marshal.Backups(out)
	}, &wg)

	wg.Wait()

	return database
}

func getDatabasesAndSchemaNames(databases []model.Database) (databasesNames, schemasNames string) {
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

func computeCoreFactor(host model.Host, dbEdition string) float32 {
	coreFactor := float32(-1)
	if host.Type == "OVM" || host.Type == "VMWARE" || host.Type == "VMOTHER" {
		if dbEdition == "EXE" || dbEdition == "ENT" {
			coreFactor = float32(host.CPUCores) * 0.25
		} else if dbEdition == "STD" {
			coreFactor = 0
		}
	} else if host.Type == "PH" {
		if dbEdition == "EXE" || dbEdition == "ENT" {
			coreFactor = float32(host.CPUCores) * 0.25
		} else if dbEdition == "STD" {
			coreFactor = float32(host.Socket)
		}
	}

	return coreFactor
}

func computeLicenses(dbEdition string, coreFactor float32) []model.License {
	licenses := make([]model.License, 0)

	if dbEdition == "EXE" {
		licenses = append(licenses, model.License{
			Name:  "Oracle EXE",
			Count: coreFactor,
		})
	} else {
		licenses = append(licenses, model.License{
			Name:  "Oracle EXE",
			Count: 0,
		})
	}

	if dbEdition == "ENT" {
		licenses = append(licenses, model.License{
			Name:  "Oracle ENT",
			Count: coreFactor,
		})
	} else {
		licenses = append(licenses, model.License{
			Name:  "Oracle ENT",
			Count: 0,
		})
	}

	if dbEdition == "STD" {
		licenses = append(licenses, model.License{
			Name:  "Oracle STD",
			Count: coreFactor,
		})
	} else {
		licenses = append(licenses, model.License{
			Name:  "Oracle STD",
			Count: 0,
		})
	}

	return licenses
}
