package builder

import (
	"log"
	"strconv"
	"strings"

	"github.com/ercole-io/ercole-agent/config"
	"github.com/ercole-io/ercole-agent/marshal"
	"github.com/ercole-io/ercole-agent/model"
	"github.com/ercole-io/ercole-agent/utils"
)

func BuildData(configuration config.Configuration, version string, hostDataSchemaVersion int) *model.HostData {
	hostData := new(model.HostData)

	if configuration.Hostname != "default" {
		hostData.Hostname = configuration.Hostname
	}

	hostData.Environment = configuration.Envtype
	hostData.Location = configuration.Location
	hostData.HostType = configuration.HostType
	hostData.Version = version
	hostData.HostDataSchemaVersion = hostDataSchemaVersion
	hostData.Info = *getHost(configuration)
	hostData.Hostname = hostData.Info.Hostname

	hostData.Extra.Filesystems = getFilesystems(configuration)
	hostData.Extra.Databases = getDatabases(configuration, hostData.Info.Type)

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

func getDatabases(configuration config.Configuration, hostType string) []model.Database {
	out := fetcher(configuration, "oratab", configuration.Oratab)
	oratabEntries := marshal.Oratab(out)

	databaseChannel := make(chan model.Database, len(oratabEntries))

	for _, entry := range oratabEntries {
		utils.RunInRoutine(configuration, func() {
			databaseChannel <- getDatabase(configuration, entry, hostType)
		})
	}

	var databases = []model.Database{}
	for i := 0; i < len(oratabEntries); i++ {
		databases = append(databases, <-databaseChannel)
	}

	return databases
}

func getDatabase(configuration config.Configuration, entry model.OratabEntry, hostType string) model.Database {
	out := fetcher(configuration, "dbstatus", entry.DBName, entry.OracleHome)
	dbStatus := strings.TrimSpace(string(out))
	var database model.Database

	if dbStatus == "OPEN" {
		out = fetcher(configuration, "dbversion", entry.DBName, entry.OracleHome)
		outVersion := string(out)

		dbVersion := strings.Split(outVersion, ".")[0]

		if configuration.Forcestats {
			fetcher(configuration, "stats", entry.DBName, entry.OracleHome)
		}

		out = fetcher(configuration, "db", entry.DBName, entry.OracleHome, strconv.Itoa(configuration.AWR))
		database = marshal.Database(out)

		out = fetcher(configuration, "tablespace", entry.DBName, entry.OracleHome)
		database.Tablespaces = marshal.Tablespaces(out)

		out = fetcher(configuration, "schema", entry.DBName, entry.OracleHome)
		database.Schemas = marshal.Schemas(out)

		out = fetcher(configuration, "patch", entry.DBName, dbVersion, entry.OracleHome)
		database.Patches = marshal.Patches(out)

		out = fetcher(configuration, "feature", entry.DBName, dbVersion, entry.OracleHome)
		if strings.Contains(string(out), "deadlocked on readable physical standby") {
			log.Println("Detected bug active dataguard 2311894.1!")
			database.Features = []model.Feature{}
		} else if strings.Contains(string(out), "ORA-01555: snapshot too old: rollback segment number") {
			log.Println("Detected error on active dataguard ORA-01555!")
			database.Features = []model.Feature{}
		} else {
			database.Features = marshal.Features(out)
		}

		out = fetcher(configuration, "opt", entry.DBName, dbVersion, entry.OracleHome)
		database.Features2 = marshal.Features2(out)

		out = fetcher(configuration, "license", entry.DBName, dbVersion, hostType, entry.OracleHome)
		database.Licenses = marshal.Licenses(out)

		out = fetcher(configuration, "addm", entry.DBName, entry.OracleHome)
		database.ADDMs = marshal.Addms(out)

		out = fetcher(configuration, "segmentadvisor", entry.DBName, entry.OracleHome)
		database.SegmentAdvisors = marshal.SegmentAdvisor(out)

		out = fetcher(configuration, "psu", entry.DBName, dbVersion, entry.OracleHome)
		database.LastPSUs = marshal.PSU(out)

		out = fetcher(configuration, "backup", entry.DBName, entry.OracleHome)
		database.Backups = marshal.Backups(out)

	} else if dbStatus == "MOUNTED" {
		out = fetcher(configuration, "dbmounted", entry.DBName, entry.OracleHome)
		database = marshal.Database(out)

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
