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
	"errors"
	"fmt"
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
	ercutils "github.com/ercole-io/ercole/v2/utils"
)

// LinuxFetcherImpl fetcher implementation for linux
type LinuxFetcherImpl struct {
	configuration config.Configuration
	log           logger.Logger
	fetcherUser   *User
}

var notImplementedLinux = errors.New("Not yet implemented for GNU/Linux")

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
func (lf *LinuxFetcherImpl) execute(fetcherName string, args ...string) ([]byte, error) {
	commandName := config.GetBaseDir() + "/fetch/linux/" + fetcherName + ".sh"
	lf.log.Infof("Fetching %s %s", commandName, strings.Join(args, " "))

	stdout, stderr, exitCode, err := runCommandAs(lf.log, lf.fetcherUser, commandName, args...)

	msg := fmt.Sprintf("Fetcher [%s] stdout: [%v] stderr: [%v] exitCode: [%v] err: [%v]",
		fetcherName,
		strings.TrimSpace(string(stdout)), strings.TrimSpace(string(stderr)),
		exitCode, err)

	if len(stderr) > 0 || exitCode > 0 {
		lf.log.Errorf(msg)
	} else {
		lf.log.Debugf(msg)
	}

	if err != nil {
		if fetcherName == "dbstatus" {
			return []byte("UNREACHABLE"), nil
		}

		err = fmt.Errorf("Error running [%s %s]: [%v]", commandName, strings.Join(args, " "), err)
		return nil, err
	}

	return stdout, nil
}

// executePwsh execute pwsh script by name
func (lf *LinuxFetcherImpl) executePwsh(fetcherName string, args ...string) []byte {
	scriptPath := config.GetBaseDir() + "/fetch/linux/" + fetcherName
	args = append([]string{scriptPath}, args...)

	lf.log.Infof("Fetching %s %s", scriptPath, strings.Join(args, " "))

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
func (lf *LinuxFetcherImpl) GetHost() (*model.Host, error) {
	out, err := lf.execute("host")
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	return marshal.Host(out)
}

// GetFilesystems get
func (lf *LinuxFetcherImpl) GetFilesystems() ([]model.Filesystem, error) {
	out, err := lf.execute("filesystem")
	err = errors.New("fake error")
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	return marshal.Filesystems(out)
}

// GetOracleDatabaseOratabEntries get
func (lf *LinuxFetcherImpl) GetOracleDatabaseOratabEntries() ([]agentmodel.OratabEntry, error) {
	out, err := lf.execute("oratab", lf.configuration.Features.OracleDatabase.Oratab)
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	return marshal_oracle.Oratab(out), nil
}

// GetOracleDatabaseRunningDatabases get
func (lf *LinuxFetcherImpl) GetOracleDatabaseRunningDatabases() ([]string, error) {
	out, err := lf.execute("oracle_running_databases")
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	dbs := strings.Split(string(out), "\n")

	ret := make([]string, 0)
	for _, db := range dbs {
		tmp := strings.TrimSpace(db)
		if len(tmp) > 0 {
			ret = append(ret, db)
		}
	}

	return ret, nil
}

// GetOracleDatabaseDbStatus get
func (lf *LinuxFetcherImpl) GetOracleDatabaseDbStatus(entry agentmodel.OratabEntry) (string, error) {
	out, err := lf.execute("dbstatus", entry.DBName, entry.OracleHome)
	if err != nil {
		return "", ercutils.NewError(err)
	}

	return strings.TrimSpace(string(out)), nil
}

// GetOracleDatabaseMountedDb get
func (lf *LinuxFetcherImpl) GetOracleDatabaseMountedDb(entry agentmodel.OratabEntry) (*model.OracleDatabase, error) {
	out, err := lf.execute("dbmounted", entry.DBName, entry.OracleHome)
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	return marshal_oracle.Database(out)
}

// GetOracleDatabaseDbVersion get
func (lf *LinuxFetcherImpl) GetOracleDatabaseDbVersion(entry agentmodel.OratabEntry) (string, error) {
	out, err := lf.execute("dbversion", entry.DBName, entry.OracleHome)
	if err != nil {
		return "", ercutils.NewError(err)
	}

	return strings.Split(string(out), ".")[0], nil
}

// RunOracleDatabaseStats Execute stats script
func (lf *LinuxFetcherImpl) RunOracleDatabaseStats(entry agentmodel.OratabEntry) error {
	_, err := lf.execute("stats", entry.DBName, entry.OracleHome)
	if err != nil {
		return ercutils.NewError(err)
	}

	return nil
}

// GetOracleDatabaseOpenDb get
func (lf *LinuxFetcherImpl) GetOracleDatabaseOpenDb(entry agentmodel.OratabEntry) (*model.OracleDatabase, error) {
	out, err := lf.execute("db", entry.DBName, entry.OracleHome, strconv.Itoa(lf.configuration.Features.OracleDatabase.AWR))
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	return marshal_oracle.Database(out)
}

// GetOracleDatabaseTablespaces get
func (lf *LinuxFetcherImpl) GetOracleDatabaseTablespaces(entry agentmodel.OratabEntry) ([]model.OracleDatabaseTablespace, error) {
	out, err := lf.execute("tablespace", entry.DBName, entry.OracleHome)
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	return marshal_oracle.Tablespaces(out)
}

// GetOracleDatabaseSchemas get
func (lf *LinuxFetcherImpl) GetOracleDatabaseSchemas(entry agentmodel.OratabEntry) ([]model.OracleDatabaseSchema, error) {
	out, err := lf.execute("schema", entry.DBName, entry.OracleHome)
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	return marshal_oracle.Schemas(out)
}

// GetOracleDatabasePatches get
func (lf *LinuxFetcherImpl) GetOracleDatabasePatches(entry agentmodel.OratabEntry, dbVersion string) ([]model.OracleDatabasePatch, error) {
	out, err := lf.execute("patch", entry.DBName, dbVersion, entry.OracleHome)
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	return marshal_oracle.Patches(out)
}

// GetOracleDatabaseFeatureUsageStat get
func (lf *LinuxFetcherImpl) GetOracleDatabaseFeatureUsageStat(entry agentmodel.OratabEntry, dbVersion string) ([]model.OracleDatabaseFeatureUsageStat, error) {
	out, err := lf.execute("opt", entry.DBName, dbVersion, entry.OracleHome)
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	return marshal_oracle.DatabaseFeatureUsageStat(out)
}

// GetOracleDatabaseLicenses get
func (lf *LinuxFetcherImpl) GetOracleDatabaseLicenses(entry agentmodel.OratabEntry, dbVersion, hardwareAbstractionTechnology string) ([]model.OracleDatabaseLicense, error) {
	out, err := lf.execute("license", entry.DBName, dbVersion, hardwareAbstractionTechnology, entry.OracleHome)
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	return marshal_oracle.Licenses(out)
}

// GetOracleDatabaseADDMs get
func (lf *LinuxFetcherImpl) GetOracleDatabaseADDMs(entry agentmodel.OratabEntry) ([]model.OracleDatabaseAddm, error) {
	out, err := lf.execute("addm", entry.DBName, entry.OracleHome)
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	return marshal_oracle.Addms(out)
}

// GetOracleDatabaseSegmentAdvisors get
func (lf *LinuxFetcherImpl) GetOracleDatabaseSegmentAdvisors(entry agentmodel.OratabEntry) ([]model.OracleDatabaseSegmentAdvisor, error) {
	out, err := lf.execute("segmentadvisor", entry.DBName, entry.OracleHome)
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	return marshal_oracle.SegmentAdvisor(out)
}

// GetOracleDatabasePSUs get
func (lf *LinuxFetcherImpl) GetOracleDatabasePSUs(entry agentmodel.OratabEntry, dbVersion string) ([]model.OracleDatabasePSU, error) {
	out, err := lf.execute("psu", entry.DBName, dbVersion, entry.OracleHome)
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	return marshal_oracle.PSU(out), nil
}

// GetOracleDatabaseBackups get
func (lf *LinuxFetcherImpl) GetOracleDatabaseBackups(entry agentmodel.OratabEntry) ([]model.OracleDatabaseBackup, error) {
	out, err := lf.execute("backup", entry.DBName, entry.OracleHome)
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	return marshal_oracle.Backups(out), nil
}

// GetOracleDatabaseCheckPDB get
func (lf *LinuxFetcherImpl) GetOracleDatabaseCheckPDB(entry agentmodel.OratabEntry) (bool, error) {
	out, err := lf.execute("checkpdb", entry.DBName, entry.OracleHome)
	if err != nil {
		return false, ercutils.NewError(err)
	}

	return strings.TrimSpace(string(out)) == "TRUE", nil
}

// GetOracleDatabasePDBs get
func (lf *LinuxFetcherImpl) GetOracleDatabasePDBs(entry agentmodel.OratabEntry) ([]model.OracleDatabasePluggableDatabase, error) {
	out, err := lf.execute("listpdb", entry.DBName, entry.OracleHome)
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	return marshal_oracle.ListPDB(out), nil
}

// GetOracleDatabasePDBTablespaces get
func (lf *LinuxFetcherImpl) GetOracleDatabasePDBTablespaces(entry agentmodel.OratabEntry, pdb string) ([]model.OracleDatabaseTablespace, error) {
	out, err := lf.execute("tablespace_pdb", entry.DBName, entry.OracleHome, pdb)
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	return marshal_oracle.Tablespaces(out)
}

// GetOracleDatabasePDBSchemas get
func (lf *LinuxFetcherImpl) GetOracleDatabasePDBSchemas(entry agentmodel.OratabEntry, pdb string) ([]model.OracleDatabaseSchema, error) {
	out, err := lf.execute("schema_pdb", entry.DBName, entry.OracleHome, pdb)
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	return marshal_oracle.Schemas(out)
}

// GetClusters return VMWare clusters from the given hyperVisor
func (lf *LinuxFetcherImpl) GetClusters(hv config.Hypervisor) ([]model.ClusterInfo, error) {
	var out []byte
	var err error

	switch hv.Type {
	case model.TechnologyVMWare:
		out = lf.executePwsh("vmware.ps1", "-s", "cluster", hv.Endpoint, hv.Username, hv.Password)

	case model.TechnologyOracleVM:
		out, err = lf.execute("ovm", "cluster", hv.Endpoint, hv.Username, hv.Password, hv.OvmUserKey, hv.OvmControl)
		if err != nil {
			return nil, ercutils.NewError(err)
		}

	default:
		return nil, ercutils.NewErrorf("Hypervisor not supported: %v (%v)", hv.Type, hv)
	}

	fetchedClusters := marshal.Clusters(out)
	for i := range fetchedClusters {
		fetchedClusters[i].Type = hv.Type
		fetchedClusters[i].FetchEndpoint = hv.Endpoint
	}

	return fetchedClusters, nil
}

// GetVirtualMachines return VMWare virtual machines infos from the given hyperVisor
func (lf *LinuxFetcherImpl) GetVirtualMachines(hv config.Hypervisor) (map[string][]model.VMInfo, error) {
	var vms map[string][]model.VMInfo

	switch hv.Type {
	case model.TechnologyVMWare:
		out := lf.executePwsh("vmware.ps1", "-s", "vms", hv.Endpoint, hv.Username, hv.Password)
		vms = marshal.VmwareVMs(out)

	case model.TechnologyOracleVM:
		out, err := lf.execute("ovm", "vms", hv.Endpoint, hv.Username, hv.Password, hv.OvmUserKey, hv.OvmControl)
		if err != nil {
			return nil, ercutils.NewError(err)
		}

		vms = marshal.OvmVMs(out)

	default:
		return nil, ercutils.NewErrorf("Hypervisor not supported: %v (%v)", hv.Type, hv)
	}

	lf.log.Debugf("Got %d vms from hypervisor: %s", len(vms), hv.Endpoint)

	return vms, nil
}

// GetOracleExadataComponents get
func (lf *LinuxFetcherImpl) GetOracleExadataComponents() ([]model.OracleExadataComponent, error) {
	out, err := lf.execute("exadata/info")
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	return marshal_oracle.ExadataComponent(out)
}

// GetOracleExadataCellDisks get
func (lf *LinuxFetcherImpl) GetOracleExadataCellDisks() (map[agentmodel.StorageServerName][]model.OracleExadataCellDisk, error) {
	out, err := lf.execute("exadata/storage-status")
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	return marshal_oracle.ExadataCellDisks(out)
}

// GetClustersMembershipStatus get
func (lf *LinuxFetcherImpl) GetClustersMembershipStatus() (*model.ClusterMembershipStatus, error) {
	out, err := lf.execute("cluster_membership_status")
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	return marshal.ClusterMembershipStatus(out), nil
}

// GetMicrosoftSQLServerInstances get
func (lf *LinuxFetcherImpl) GetMicrosoftSQLServerInstances() ([]agentmodel.ListInstanceOutputModel, error) {
	lf.log.Error(notImplementedLinux)
	return nil, ercutils.NewError(notImplementedLinux)
}

// GetMicrosoftSQLServerInstanceInfo get
func (lf *LinuxFetcherImpl) GetMicrosoftSQLServerInstanceInfo(conn string, inst *model.MicrosoftSQLServerInstance) error {
	lf.log.Error(notImplementedLinux)
	return ercutils.NewError(notImplementedLinux)
}

// GetMicrosoftSQLServerInstanceEdition get
func (lf *LinuxFetcherImpl) GetMicrosoftSQLServerInstanceEdition(conn string, inst *model.MicrosoftSQLServerInstance) error {
	lf.log.Error(notImplementedLinux)
	return ercutils.NewError(notImplementedLinux)
}

// GetMicrosoftSQLServerInstanceLicensingInfo get
func (lf *LinuxFetcherImpl) GetMicrosoftSQLServerInstanceLicensingInfo(conn string, inst *model.MicrosoftSQLServerInstance) error {
	lf.log.Error(notImplementedLinux)
	return ercutils.NewError(notImplementedLinux)
}

// GetMicrosoftSQLServerInstanceDatabase get
func (lf *LinuxFetcherImpl) GetMicrosoftSQLServerInstanceDatabase(conn string) ([]model.MicrosoftSQLServerDatabase, error) {
	lf.log.Error(notImplementedLinux)
	return nil, ercutils.NewError(notImplementedLinux)
}

// GetMicrosoftSQLServerInstanceDatabaseBackups get
func (lf *LinuxFetcherImpl) GetMicrosoftSQLServerInstanceDatabaseBackups(conn string) ([]agentmodel.DbBackupsModel, error) {
	lf.log.Error(notImplementedLinux)
	return nil, ercutils.NewError(notImplementedLinux)
}

// GetMicrosoftSQLServerInstanceDatabaseSchemas get
func (lf *LinuxFetcherImpl) GetMicrosoftSQLServerInstanceDatabaseSchemas(conn string) ([]agentmodel.DbSchemasModel, error) {
	lf.log.Error(notImplementedLinux)
	return nil, ercutils.NewError(notImplementedLinux)
}

// GetMicrosoftSQLServerInstanceDatabaseTablespaces get
func (lf *LinuxFetcherImpl) GetMicrosoftSQLServerInstanceDatabaseTablespaces(conn string) ([]agentmodel.DbTablespacesModel, error) {
	lf.log.Error(notImplementedLinux)
	return nil, ercutils.NewError(notImplementedLinux)
}

// GetMicrosoftSQLServerInstancePatches get
func (lf *LinuxFetcherImpl) GetMicrosoftSQLServerInstancePatches(conn string) ([]model.MicrosoftSQLServerPatch, error) {
	lf.log.Error(notImplementedLinux)
	return nil, ercutils.NewError(notImplementedLinux)
}

// GetMicrosoftSQLServerProductFeatures get
func (lf *LinuxFetcherImpl) GetMicrosoftSQLServerProductFeatures(conn string) ([]model.MicrosoftSQLServerProductFeature, error) {
	lf.log.Error(notImplementedLinux)
	return nil, ercutils.NewError(notImplementedLinux)
}

func (lf *LinuxFetcherImpl) GetMySQLInstance(connection config.MySQLInstanceConnection) (*model.MySQLInstance, error) {
	out, err := lf.execute("mysql/mysql_gather", "-h", connection.Host, "-u", connection.User, "-p", connection.Password, "-a", "instance")
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	return marshal_mysql.Instance(out)
}

func (lf *LinuxFetcherImpl) GetMySQLDatabases(connection config.MySQLInstanceConnection) ([]model.MySQLDatabase, error) {
	out, err := lf.execute("mysql/mysql_gather", "-h", connection.Host, "-u", connection.User, "-p", connection.Password, "-a", "databases")
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	return marshal_mysql.Databases(out), nil
}

func (lf *LinuxFetcherImpl) GetMySQLTableSchemas(connection config.MySQLInstanceConnection) ([]model.MySQLTableSchema, error) {
	out, err := lf.execute("mysql/mysql_gather", "-h", connection.Host, "-u", connection.User, "-p", connection.Password, "-a", "table_schemas")
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	return marshal_mysql.TableSchemas(out)
}

func (lf *LinuxFetcherImpl) GetMySQLSegmentAdvisors(connection config.MySQLInstanceConnection) ([]model.MySQLSegmentAdvisor, error) {
	out, err := lf.execute("mysql/mysql_gather", "-h", connection.Host, "-u", connection.User, "-p", connection.Password, "-a", "segment_advisors")
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	return marshal_mysql.SegmentAdvisors(out)
}

func (lf *LinuxFetcherImpl) GetMySQLHighAvailability(connection config.MySQLInstanceConnection) (bool, error) {
	out, err := lf.execute("mysql/mysql_gather", "-h", connection.Host, "-u", connection.User, "-p", connection.Password, "-a", "high_availability")
	if err != nil {
		return false, ercutils.NewError(err)
	}

	return marshal_mysql.HighAvailability(out), nil
}

func (lf *LinuxFetcherImpl) GetMySQLUUID() (string, error) {
	file := "/var/lib/mysql/auto.cnf"
	out, err := os.ReadFile(file)
	if err != nil {
		err = fmt.Errorf("Can't get MySQL UUID from %s: %w", file, err)
		lf.log.Error(err)
		return "", ercutils.NewError(err)
	}

	uuid, err := marshal_mysql.UUID(out)
	if err != nil {
		err = fmt.Errorf("Can't get MySQL UUID from %s: %w", file, err)
		lf.log.Error(err)
		return "", ercutils.NewError(err)
	}

	return uuid, nil
}

func (lf *LinuxFetcherImpl) GetMySQLSlaveHosts(connection config.MySQLInstanceConnection) (bool, []string, error) {
	out, err := lf.execute("mysql/mysql_gather", "-h", connection.Host, "-u", connection.User, "-p", connection.Password, "-a", "slave_hosts")
	if err != nil {
		return false, nil, ercutils.NewError(err)
	}

	isMaster, slaveUUIDs := marshal_mysql.SlaveHosts(out)
	return isMaster, slaveUUIDs, nil
}

func (lf *LinuxFetcherImpl) GetMySQLSlaveStatus(connection config.MySQLInstanceConnection) (bool, *string, error) {
	out, err := lf.execute("mysql/mysql_gather", "-h", connection.Host, "-u", connection.User, "-p", connection.Password, "-a", "slave_status")
	if err != nil {
		return false, nil, ercutils.NewError(err)
	}

	isSlave, masterUUID := marshal_mysql.SlaveStatus(out)
	return isSlave, masterUUID, nil
}
