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

package fetcher

import (
	"os"
	"os/user"
	"strconv"
	"strings"

	"github.com/ercole-io/ercole-agent/v2/agentmodel"
	"github.com/ercole-io/ercole-agent/v2/config"
	"github.com/ercole-io/ercole-agent/v2/logger"
	"github.com/ercole-io/ercole-agent/v2/marshal"
	marshal_mysql "github.com/ercole-io/ercole-agent/v2/marshal/mysql"
	marshal_oracle "github.com/ercole-io/ercole-agent/v2/marshal/oracle"
	"github.com/ercole-io/ercole/v2/model"
)

// LinuxFetcherImpl fetcher implementation for linux
type LinuxFetcherImpl struct {
	configuration config.Configuration
	log           logger.Logger
	fetcherUser   *User
}

const notImplementedLinux = "Not yet implemented for GNU/Linux"

// NewLinuxFetcherImpl constructor
func NewLinuxFetcherImpl(conf config.Configuration, log logger.Logger) *LinuxFetcherImpl {
	return &LinuxFetcherImpl{
		conf,
		log,
		nil,
	}
}

// SetUser set user used by fetcher to run commands
func (lf *LinuxFetcherImpl) SetUser(username string) error {
	u, err := lf.getUserInfo(username)
	if err != nil {
		return err
	}

	lf.fetcherUser = u
	return nil
}

func (lf *LinuxFetcherImpl) getUserInfo(username string) (*User, error) {
	u, err := user.Lookup(username)
	if err != nil {
		lf.log.Errorf("Can't lookup username [%s], error: [%v]", username, err)
		return nil, err
	}

	intUID, err := strconv.Atoi(u.Uid)
	if err != nil {
		lf.log.Errorf("Can't convert uid [%s], error: [%v]", u.Uid, err)
		return nil, err
	}

	intGID, err := strconv.Atoi(u.Gid)
	if err != nil {
		lf.log.Errorf("Can't convert gid [%s], error: [%v]", u.Gid, err)
		return nil, err
	}

	return &User{u.Name, uint32(intUID), uint32(intGID)}, nil
}

// SetUserAsCurrent set user used by fetcher to run commands as current process user
func (lf *LinuxFetcherImpl) SetUserAsCurrent() error {
	lf.fetcherUser = nil
	return nil
}

// Execute execute bash script by name
func (lf *LinuxFetcherImpl) execute(fetcherName string, args ...string) []byte {
	commandName := config.GetBaseDir() + "/fetch/linux/" + fetcherName + ".sh"
	lf.log.Infof("Fetching %s %s", commandName, strings.Join(args, " "))

	stdout, stderr, exitCode, err := runCommandAs(lf.log, lf.fetcherUser, commandName, args...)

	lf.log.Debugf("Fetcher [%s] stdout: [%v]", fetcherName, strings.TrimSpace(string(stdout)))

	if len(stderr) > 0 {
		format := "Fetcher [%s] exitCode: [%v] stderr: [%v]"
		args := []interface{}{fetcherName, exitCode, strings.TrimSpace(string(stderr))}

		if exitCode == 0 {
			lf.log.Debugf(format, args...)
		} else {
			lf.log.Errorf(format, args...)
		}
	}

	if err != nil {
		if fetcherName == "dbstatus" {
			return []byte("UNREACHABLE")
		}

		lf.log.Fatalf("Fatal error running [%s %s]: [%v]", commandName, strings.Join(args, " "), err)
	}

	return stdout
}

// executePwsh execute pwsh script by name
func (lf *LinuxFetcherImpl) executePwsh(fetcherName string, args ...string) []byte {
	scriptPath := config.GetBaseDir() + "/fetch/linux/" + fetcherName
	args = append([]string{scriptPath}, args...)

	lf.log.Infof("Fetching %v", scriptPath, strings.Join(args, " "))

	stdout, stderr, exitCode, err := runCommandAs(lf.log, lf.fetcherUser, "/usr/bin/pwsh", args...)

	if len(stdout) > 0 {
		lf.log.Debugf("Fetcher [%s] stdout: [%v]", fetcherName, strings.TrimSpace(string(stdout)))
	}

	if len(stderr) > 0 {
		lf.log.Errorf("Fetcher [%s] exitCode: [%v] stderr: [%v]", fetcherName, exitCode, strings.TrimSpace(string(stderr)))
	}

	if err != nil {
		lf.log.Fatalf("Fatal error running [%s %s]: [%v]", scriptPath, strings.Join(args, " "), err)
	}

	return stdout
}

// GetHost get
func (lf *LinuxFetcherImpl) GetHost() model.Host {
	out := lf.execute("host")
	return marshal.Host(out)
}

// GetFilesystems get
func (lf *LinuxFetcherImpl) GetFilesystems() []model.Filesystem {
	out := lf.execute("filesystem")
	return marshal.Filesystems(out)
}

// GetOracleDatabaseOratabEntries get
func (lf *LinuxFetcherImpl) GetOracleDatabaseOratabEntries() []agentmodel.OratabEntry {
	out := lf.execute("oratab", lf.configuration.Features.OracleDatabase.Oratab)
	return marshal_oracle.Oratab(out)
}

// GetOracleDatabaseRunningDatabases get
func (lf *LinuxFetcherImpl) GetOracleDatabaseRunningDatabases() []string {
	out := lf.execute("oracle_running_databases")

	dbs := strings.Split(string(out), "\n")

	ret := make([]string, 0)
	for _, db := range dbs {
		tmp := strings.TrimSpace(db)
		if len(tmp) > 0 {
			ret = append(ret, db)
		}
	}

	return ret
}

// GetOracleDatabaseDbStatus get
func (lf *LinuxFetcherImpl) GetOracleDatabaseDbStatus(entry agentmodel.OratabEntry) string {
	out := lf.execute("dbstatus", entry.DBName, entry.OracleHome)
	return strings.TrimSpace(string(out))
}

// GetOracleDatabaseMountedDb get
func (lf *LinuxFetcherImpl) GetOracleDatabaseMountedDb(entry agentmodel.OratabEntry) model.OracleDatabase {
	out := lf.execute("dbmounted", entry.DBName, entry.OracleHome)
	return marshal_oracle.Database(out)
}

// GetOracleDatabaseDbVersion get
func (lf *LinuxFetcherImpl) GetOracleDatabaseDbVersion(entry agentmodel.OratabEntry) string {
	out := lf.execute("dbversion", entry.DBName, entry.OracleHome)
	return strings.Split(string(out), ".")[0]
}

// RunOracleDatabaseStats Execute stats script
func (lf *LinuxFetcherImpl) RunOracleDatabaseStats(entry agentmodel.OratabEntry) {
	lf.execute("stats", entry.DBName, entry.OracleHome)
}

// GetOracleDatabaseOpenDb get
func (lf *LinuxFetcherImpl) GetOracleDatabaseOpenDb(entry agentmodel.OratabEntry) model.OracleDatabase {
	out := lf.execute("db", entry.DBName, entry.OracleHome, strconv.Itoa(lf.configuration.Features.OracleDatabase.AWR))
	return marshal_oracle.Database(out)
}

// GetOracleDatabaseTablespaces get
func (lf *LinuxFetcherImpl) GetOracleDatabaseTablespaces(entry agentmodel.OratabEntry) []model.OracleDatabaseTablespace {
	out := lf.execute("tablespace", entry.DBName, entry.OracleHome)
	return marshal_oracle.Tablespaces(out)
}

// GetOracleDatabaseSchemas get
func (lf *LinuxFetcherImpl) GetOracleDatabaseSchemas(entry agentmodel.OratabEntry) []model.OracleDatabaseSchema {
	out := lf.execute("schema", entry.DBName, entry.OracleHome)
	return marshal_oracle.Schemas(out)
}

// GetOracleDatabasePatches get
func (lf *LinuxFetcherImpl) GetOracleDatabasePatches(entry agentmodel.OratabEntry, dbVersion string) []model.OracleDatabasePatch {
	out := lf.execute("patch", entry.DBName, dbVersion, entry.OracleHome)
	return marshal_oracle.Patches(out)
}

// GetOracleDatabaseFeatureUsageStat get
func (lf *LinuxFetcherImpl) GetOracleDatabaseFeatureUsageStat(entry agentmodel.OratabEntry, dbVersion string) []model.OracleDatabaseFeatureUsageStat {
	out := lf.execute("opt", entry.DBName, dbVersion, entry.OracleHome)
	return marshal_oracle.DatabaseFeatureUsageStat(out)
}

// GetOracleDatabaseLicenses get
func (lf *LinuxFetcherImpl) GetOracleDatabaseLicenses(entry agentmodel.OratabEntry, dbVersion, hardwareAbstractionTechnology string) []model.OracleDatabaseLicense {
	out := lf.execute("license", entry.DBName, dbVersion, hardwareAbstractionTechnology, entry.OracleHome)
	return marshal_oracle.Licenses(out)
}

// GetOracleDatabaseADDMs get
func (lf *LinuxFetcherImpl) GetOracleDatabaseADDMs(entry agentmodel.OratabEntry) []model.OracleDatabaseAddm {
	out := lf.execute("addm", entry.DBName, entry.OracleHome)
	return marshal_oracle.Addms(out)
}

// GetOracleDatabaseSegmentAdvisors get
func (lf *LinuxFetcherImpl) GetOracleDatabaseSegmentAdvisors(entry agentmodel.OratabEntry) []model.OracleDatabaseSegmentAdvisor {
	out := lf.execute("segmentadvisor", entry.DBName, entry.OracleHome)
	return marshal_oracle.SegmentAdvisor(out)
}

// GetOracleDatabasePSUs get
func (lf *LinuxFetcherImpl) GetOracleDatabasePSUs(entry agentmodel.OratabEntry, dbVersion string) []model.OracleDatabasePSU {
	out := lf.execute("psu", entry.DBName, dbVersion, entry.OracleHome)
	return marshal_oracle.PSU(out)
}

// GetOracleDatabaseBackups get
func (lf *LinuxFetcherImpl) GetOracleDatabaseBackups(entry agentmodel.OratabEntry) []model.OracleDatabaseBackup {
	out := lf.execute("backup", entry.DBName, entry.OracleHome)
	return marshal_oracle.Backups(out)
}

// GetOracleDatabaseCheckPDB get
func (lf *LinuxFetcherImpl) GetOracleDatabaseCheckPDB(entry agentmodel.OratabEntry) bool {
	out := lf.execute("checkpdb", entry.DBName, entry.OracleHome)
	return strings.TrimSpace(string(out)) == "TRUE"
}

// GetOracleDatabasePDBs get
func (lf *LinuxFetcherImpl) GetOracleDatabasePDBs(entry agentmodel.OratabEntry) []model.OracleDatabasePluggableDatabase {
	out := lf.execute("listpdb", entry.DBName, entry.OracleHome)
	return marshal_oracle.ListPDB(out)
}

// GetOracleDatabasePDBTablespaces get
func (lf *LinuxFetcherImpl) GetOracleDatabasePDBTablespaces(entry agentmodel.OratabEntry, pdb string) []model.OracleDatabaseTablespace {
	out := lf.execute("tablespace_pdb", entry.DBName, entry.OracleHome, pdb)
	return marshal_oracle.Tablespaces(out)
}

// GetOracleDatabasePDBSchemas get
func (lf *LinuxFetcherImpl) GetOracleDatabasePDBSchemas(entry agentmodel.OratabEntry, pdb string) []model.OracleDatabaseSchema {
	out := lf.execute("schema_pdb", entry.DBName, entry.OracleHome, pdb)
	return marshal_oracle.Schemas(out)
}

// GetClusters return VMWare clusters from the given hyperVisor
func (lf *LinuxFetcherImpl) GetClusters(hv config.Hypervisor) []model.ClusterInfo {
	var out []byte

	switch hv.Type {
	case model.TechnologyVMWare:
		out = lf.executePwsh("vmware.ps1", "-s", "cluster", hv.Endpoint, hv.Username, hv.Password)

	case model.TechnologyOracleVM:
		out = lf.execute("ovm", "cluster", hv.Endpoint, hv.Username, hv.Password, hv.OvmUserKey, hv.OvmControl)

	default:
		lf.log.Errorf("Hypervisor not supported: %v (%v)", hv.Type, hv)
		return make([]model.ClusterInfo, 0)
	}

	fetchedClusters := marshal.Clusters(out)
	for i := range fetchedClusters {
		fetchedClusters[i].Type = hv.Type
		fetchedClusters[i].FetchEndpoint = hv.Endpoint
	}

	return fetchedClusters
}

// GetVirtualMachines return VMWare virtual machines infos from the given hyperVisor
func (lf *LinuxFetcherImpl) GetVirtualMachines(hv config.Hypervisor) map[string][]model.VMInfo {
	var vms map[string][]model.VMInfo

	switch hv.Type {
	case model.TechnologyVMWare:
		out := lf.executePwsh("vmware.ps1", "-s", "vms", hv.Endpoint, hv.Username, hv.Password)
		vms = marshal.VmwareVMs(out)

	case model.TechnologyOracleVM:
		out := lf.execute("ovm", "vms", hv.Endpoint, hv.Username, hv.Password, hv.OvmUserKey, hv.OvmControl)
		vms = marshal.OvmVMs(out)

	default:
		lf.log.Errorf("Hypervisor not supported: %v (%v)", hv.Type, hv)
		return make(map[string][]model.VMInfo)
	}

	lf.log.Debugf("Got %d vms from hypervisor: %s", len(vms), hv.Endpoint)

	return vms
}

// GetOracleExadataComponents get
func (lf *LinuxFetcherImpl) GetOracleExadataComponents() []model.OracleExadataComponent {
	out := lf.execute("exadata/info")
	return marshal_oracle.ExadataComponent(out)
}

// GetOracleExadataCellDisks get
func (lf *LinuxFetcherImpl) GetOracleExadataCellDisks() map[agentmodel.StorageServerName][]model.OracleExadataCellDisk {
	out := lf.execute("exadata/storage-status")
	return marshal_oracle.ExadataCellDisks(out)
}

// GetClustersMembershipStatus get
func (lf *LinuxFetcherImpl) GetClustersMembershipStatus() model.ClusterMembershipStatus {
	out := lf.execute("cluster_membership_status")
	return marshal.ClusterMembershipStatus(out)
}

// GetMicrosoftSQLServerInstances get
func (lf *LinuxFetcherImpl) GetMicrosoftSQLServerInstances() []agentmodel.ListInstanceOutputModel {
	lf.log.Error(notImplementedLinux)
	return nil
}

// GetMicrosoftSQLServerInstanceInfo get
func (lf *LinuxFetcherImpl) GetMicrosoftSQLServerInstanceInfo(conn string, inst *model.MicrosoftSQLServerInstance) {
	lf.log.Error(notImplementedLinux)
}

// GetMicrosoftSQLServerInstanceEdition get
func (lf *LinuxFetcherImpl) GetMicrosoftSQLServerInstanceEdition(conn string, inst *model.MicrosoftSQLServerInstance) {
	lf.log.Error(notImplementedLinux)
}

// GetMicrosoftSQLServerInstanceLicensingInfo get
func (lf *LinuxFetcherImpl) GetMicrosoftSQLServerInstanceLicensingInfo(conn string, inst *model.MicrosoftSQLServerInstance) {
	lf.log.Error(notImplementedLinux)
}

// GetMicrosoftSQLServerInstanceDatabase get
func (lf *LinuxFetcherImpl) GetMicrosoftSQLServerInstanceDatabase(conn string) []model.MicrosoftSQLServerDatabase {
	lf.log.Error(notImplementedLinux)
	return nil
}

// GetMicrosoftSQLServerInstanceDatabaseBackups get
func (lf *LinuxFetcherImpl) GetMicrosoftSQLServerInstanceDatabaseBackups(conn string) []agentmodel.DbBackupsModel {
	lf.log.Error(notImplementedLinux)
	return nil
}

// GetMicrosoftSQLServerInstanceDatabaseSchemas get
func (lf *LinuxFetcherImpl) GetMicrosoftSQLServerInstanceDatabaseSchemas(conn string) []agentmodel.DbSchemasModel {
	lf.log.Error(notImplementedLinux)
	return nil
}

// GetMicrosoftSQLServerInstanceDatabaseTablespaces get
func (lf *LinuxFetcherImpl) GetMicrosoftSQLServerInstanceDatabaseTablespaces(conn string) []agentmodel.DbTablespacesModel {
	lf.log.Error(notImplementedLinux)
	return nil
}

// GetMicrosoftSQLServerInstancePatches get
func (lf *LinuxFetcherImpl) GetMicrosoftSQLServerInstancePatches(conn string) []model.MicrosoftSQLServerPatch {
	lf.log.Error(notImplementedLinux)
	return nil
}

// GetMicrosoftSQLServerProductFeatures get
func (lf *LinuxFetcherImpl) GetMicrosoftSQLServerProductFeatures(conn string) []model.MicrosoftSQLServerProductFeature {
	lf.log.Error(notImplementedLinux)
	return nil
}

func (lf *LinuxFetcherImpl) GetMySQLInstance(connection config.MySQLInstanceConnection) *model.MySQLInstance {
	out := lf.execute("mysql/mysql_gather", "-h", connection.Host, "-u", connection.User, "-p", connection.Password, "-a", "instance")

	return marshal_mysql.Instance(out)
}

func (lf *LinuxFetcherImpl) GetMySQLDatabases(connection config.MySQLInstanceConnection) []model.MySQLDatabase {
	out := lf.execute("mysql/mysql_gather", "-h", connection.Host, "-u", connection.User, "-p", connection.Password, "-a", "databases")

	return marshal_mysql.Databases(out)
}

func (lf *LinuxFetcherImpl) GetMySQLTableSchemas(connection config.MySQLInstanceConnection) []model.MySQLTableSchema {
	out := lf.execute("mysql/mysql_gather", "-h", connection.Host, "-u", connection.User, "-p", connection.Password, "-a", "table_schemas")

	return marshal_mysql.TableSchemas(out)
}

func (lf *LinuxFetcherImpl) GetMySQLSegmentAdvisors(connection config.MySQLInstanceConnection) []model.MySQLSegmentAdvisor {
	out := lf.execute("mysql/mysql_gather", "-h", connection.Host, "-u", connection.User, "-p", connection.Password, "-a", "segment_advisors")

	return marshal_mysql.SegmentAdvisors(out)
}

func (lf *LinuxFetcherImpl) GetMySQLHighAvailability(connection config.MySQLInstanceConnection) bool {
	out := lf.execute("mysql/mysql_gather", "-h", connection.Host, "-u", connection.User, "-p", connection.Password, "-a", "high_availability")

	return marshal_mysql.HighAvailability(out)
}

func (lf *LinuxFetcherImpl) GetMySQLUUID() string {
	file := "/var/lib/mysql/auto.cnf"
	out, err := os.ReadFile(file)
	if err != nil {
		lf.log.Fatal("Can't get MySQL UUID from ", file, ": ", err)
	}

	uuid, err := marshal_mysql.UUID(out)
	if err != nil {
		lf.log.Fatal("Can't get MySQL UUID: ", err)
	}

	return uuid
}

func (lf *LinuxFetcherImpl) GetMySQLSlaveHosts(connection config.MySQLInstanceConnection) (bool, []string) {
	out := lf.execute("mysql/mysql_gather", "-h", connection.Host, "-u", connection.User, "-p", connection.Password, "-a", "slave_hosts")

	return marshal_mysql.SlaveHosts(out)
}

func (lf *LinuxFetcherImpl) GetMySQLSlaveStatus(connection config.MySQLInstanceConnection) (bool, *string) {
	out := lf.execute("mysql/mysql_gather", "-h", connection.Host, "-u", connection.User, "-p", connection.Password, "-a", "slave_status")

	return marshal_mysql.SlaveStatus(out)
}
