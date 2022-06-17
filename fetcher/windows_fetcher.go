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

package fetcher

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/ercole-io/ercole-agent/v2/agentmodel"
	"github.com/ercole-io/ercole-agent/v2/config"
	"github.com/ercole-io/ercole-agent/v2/logger"
	"github.com/ercole-io/ercole-agent/v2/marshal"
	marshal_microsoft "github.com/ercole-io/ercole-agent/v2/marshal/microsoft"
	marshal_oracle "github.com/ercole-io/ercole-agent/v2/marshal/oracle"
	"github.com/ercole-io/ercole/v2/model"
	ercutils "github.com/ercole-io/ercole/v2/utils"
)

// WindowsFetcherImpl SpecializedFetcher implementation for windows
type WindowsFetcherImpl struct {
	configuration config.Configuration
	log           logger.Logger
}

var notImplementedWindows = errors.New("Not yet implemented for windows")

// NewWindowsFetcherImpl constructor
func NewWindowsFetcherImpl(conf config.Configuration, log logger.Logger) *WindowsFetcherImpl {
	return &WindowsFetcherImpl{
		conf,
		log,
	}
}

func (wf *WindowsFetcherImpl) executeWithContext(ctx context.Context, fetcherName string, args ...string) ([]byte, error) {
	var (
		cmd    *exec.Cmd
		err    error
		psexe  string
		stdout bytes.Buffer
		stderr bytes.Buffer
	)

	baseDir, err := config.GetBaseDir(wf.log)
	if err != nil {
		return nil, err
	}

	psexe, err = exec.LookPath("powershell.exe")
	if err != nil {
		wf.log.Fatal(psexe)
	}

	if wf.configuration.ForcePwshVersion == "0" {
		args = append([]string{"-ExecutionPolicy", "Bypass", "-File", baseDir + "\\fetch\\windows\\" + fetcherName}, args...)
	} else {
		args = append([]string{"-version", wf.configuration.ForcePwshVersion, "-ExecutionPolicy", "Bypass", "-File", baseDir + "\\fetch\\windows\\" + fetcherName}, args...)
	}

	wf.log.Info("Fetching " + psexe + " " + strings.Join(args, " "))

	cmd = exec.CommandContext(ctx, psexe, args...)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()

	msg := fmt.Sprintf("Fetcher [%s] stdout: [%v] stderr: [%v] exitCode: [%v] err: [%v]",
		fetcherName,
		strings.TrimSpace(stdout.String()), strings.TrimSpace(stderr.String()),
		cmd.ProcessState.ExitCode(), err)

	if stderr.Len() > 0 || cmd.ProcessState.ExitCode() > 0 {
		wf.log.Errorf(msg)
	} else {
		wf.log.Debugf(msg)
	}

	if err != nil {
		if fetcherName != "dbstatus" {
			return []byte("UNREACHABLE"), nil
		}

		err = fmt.Errorf("Error running [%s %s]: [%v]", psexe, strings.Join(args, " "), err)

		return nil, err
	}

	return stdout.Bytes(), nil
}

func (wf *WindowsFetcherImpl) executeWithDeadline(duration time.Duration, fetcherName string, args ...string) ([]byte, error) {
	type execResult struct {
		bytes []byte
		err   error
	}

	c := make(chan execResult, 1)

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(duration))
	defer cancel()

	go func() {
		bytes, err := wf.executeWithContext(ctx, fetcherName, args...)
		c <- execResult{
			bytes: bytes,
			err:   err,
		}
	}()

	select {
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, fmt.Errorf("fetcher last more than %s, timer has exceeded", duration)
		}

		return nil, ctx.Err()

	case result := <-c:
		return result.bytes, result.err
	}
}

func (wf *WindowsFetcherImpl) CreateOracleArgs(args ...string) []string {
	if wf.configuration.Features.OracleDatabase.IsOracleUser() {
		args = append(args, wf.configuration.Features.OracleDatabase.OracleUser.Username, wf.configuration.Features.OracleDatabase.OracleUser.Password)
	}

	return args
}

// SetUser not implemented
func (wf *WindowsFetcherImpl) SetUser(username string) error {
	wf.log.Error(notImplementedWindows)
	return ercutils.NewError(notImplementedWindows)
}

// SetUserAsCurrent set user used by fetcher to run commands as current process user
func (wf *WindowsFetcherImpl) SetUserAsCurrent() error {
	wf.log.Error(notImplementedWindows)
	return ercutils.NewError(notImplementedWindows)
}

// GetClusters not implemented
func (wf *WindowsFetcherImpl) GetClusters(hv config.Hypervisor) ([]model.ClusterInfo, error) {
	wf.log.Error(notImplementedWindows)
	return nil, ercutils.NewError(notImplementedWindows)
}

// GetVirtualMachines return VMWare virtual machines infos from the given hyperVisor
func (wf *WindowsFetcherImpl) GetVirtualMachines(hv config.Hypervisor) (map[string][]model.VMInfo, error) {
	wf.log.Error(notImplementedWindows)
	return nil, ercutils.NewError(notImplementedWindows)
}

// GetOracleExadataComponents get
func (wf *WindowsFetcherImpl) GetOracleExadataComponents() ([]model.OracleExadataComponent, error) {
	wf.log.Error(notImplementedWindows)
	return nil, ercutils.NewError(notImplementedWindows)
}

// GetOracleExadataCellDisks get
func (wf *WindowsFetcherImpl) GetOracleExadataCellDisks() (map[agentmodel.StorageServerName][]model.OracleExadataCellDisk, error) {
	wf.log.Error(notImplementedWindows)
	return nil, ercutils.NewError(notImplementedWindows)
}

// GetClustersMembershipStatus get
func (wf *WindowsFetcherImpl) GetClustersMembershipStatus() (*model.ClusterMembershipStatus, error) {
	cms := model.ClusterMembershipStatus{
		OracleClusterware:    false,
		SunCluster:           false,
		VeritasClusterServer: false,
		HACMP:                false,
	}

	return &cms, nil
}

// GetHost get
func (wf *WindowsFetcherImpl) GetHost() (*model.Host, error) {
	out, err := wf.executeWithDeadline(FetcherStandardTimeOut, "win.ps1", "-s", "host")
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	return marshal.Host(out)
}

// GetFilesystems get
func (wf *WindowsFetcherImpl) GetFilesystems() ([]model.Filesystem, error) {
	out, err := wf.executeWithDeadline(20*time.Second, "win.ps1", "-s", "filesystem")
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	return marshal.Filesystems(out)
}

// GetOracleDatabaseOratabEntries get
func (wf *WindowsFetcherImpl) GetOracleDatabaseOratabEntries() ([]agentmodel.OratabEntry, error) {
	out, err := wf.executeWithDeadline(FetcherStandardTimeOut, "win.ps1", "-s", "oratab", wf.configuration.Features.OracleDatabase.Oratab)
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	return marshal_oracle.Oratab(out), nil
}

// GetOracleDatabaseRunningDatabases get
func (wf *WindowsFetcherImpl) GetOracleDatabaseRunningDatabases() ([]string, error) {
	return []string{}, nil
}

// GetOracleDatabaseDbStatus get
func (wf *WindowsFetcherImpl) GetOracleDatabaseDbStatus(entry agentmodel.OratabEntry) (string, error) {
	out, err := wf.executeWithDeadline(FetcherStandardTimeOut, "win.ps1", wf.CreateOracleArgs("-s", "dbstatus", entry.DBName, entry.OracleHome)...)
	if err != nil {
		return "", ercutils.NewError(err)
	}

	return strings.TrimSpace(string(out)), nil
}

// GetOracleDatabaseRac get
func (wf *WindowsFetcherImpl) GetOracleDatabaseRac(entry agentmodel.OratabEntry) (string, error) {
	wf.log.Error(notImplementedWindows)
	return "", ercutils.NewError(notImplementedWindows)
}

// GetOracleDatabaseMountedDb get
func (wf *WindowsFetcherImpl) GetOracleDatabaseMountedDb(entry agentmodel.OratabEntry) (*model.OracleDatabase, error) {
	out, err := wf.executeWithDeadline(FetcherStandardTimeOut, "win.ps1", wf.CreateOracleArgs("-s", "dbmounted", entry.DBName, entry.OracleHome)...)
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	return marshal_oracle.Database(out)
}

func (wf *WindowsFetcherImpl) GetOracleDatabaseGrantsDba(entry agentmodel.OratabEntry) ([]model.OracleGrantDba, error) {
	out, err := wf.executeWithDeadline(FetcherStandardTimeOut, "win.ps1", wf.CreateOracleArgs("grant_dba", entry.DBName, entry.OracleHome)...)
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	return marshal_oracle.GrantDba(out), nil
}

// GetOracleDatabaseDbVersion get
func (wf *WindowsFetcherImpl) GetOracleDatabaseDbVersion(entry agentmodel.OratabEntry) (string, error) {
	out, err := wf.executeWithDeadline(FetcherStandardTimeOut, "win.ps1", "-s", "dbversion", entry.DBName, entry.OracleHome)
	if err != nil {
		return "", ercutils.NewError(err)
	}

	return strings.Split(string(out), ".")[0], nil
}

// RunOracleDatabaseStats Execute stats script
func (wf *WindowsFetcherImpl) RunOracleDatabaseStats(entry agentmodel.OratabEntry) error {
	if _, err := wf.executeWithDeadline(FetcherStandardTimeOut, "win.ps1", wf.CreateOracleArgs("-s", "stats", entry.DBName, entry.OracleHome)...); err != nil {
		return ercutils.NewError(err)
	}

	return nil
}

// GetOracleDatabaseOpenDb get
func (wf *WindowsFetcherImpl) GetOracleDatabaseOpenDb(entry agentmodel.OratabEntry) (*model.OracleDatabase, error) {
	out, err := wf.executeWithDeadline(FetcherStandardTimeOut, "win.ps1", wf.CreateOracleArgs("-s", "db", entry.DBName, strconv.Itoa(wf.configuration.Features.OracleDatabase.AWR))...)
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	return marshal_oracle.Database(out)
}

// GetOracleDatabaseTablespaces get
func (wf *WindowsFetcherImpl) GetOracleDatabaseTablespaces(entry agentmodel.OratabEntry) ([]model.OracleDatabaseTablespace, error) {
	out, err := wf.executeWithDeadline(FetcherStandardTimeOut, "win.ps1", wf.CreateOracleArgs("-s", "tablespace", entry.DBName, entry.OracleHome)...)
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	return marshal_oracle.Tablespaces(out)
}

// GetOracleDatabaseSchemas get
func (wf *WindowsFetcherImpl) GetOracleDatabaseSchemas(entry agentmodel.OratabEntry) ([]model.OracleDatabaseSchema, error) {
	out, err := wf.executeWithDeadline(FetcherStandardTimeOut, "win.ps1", wf.CreateOracleArgs("-s", "schema", entry.DBName, entry.OracleHome)...)
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	return marshal_oracle.Schemas(out)
}

// GetOracleDatabasePatches get
func (wf *WindowsFetcherImpl) GetOracleDatabasePatches(entry agentmodel.OratabEntry, dbVersion string) ([]model.OracleDatabasePatch, error) {
	out, err := wf.executeWithDeadline(FetcherStandardTimeOut, "win.ps1", wf.CreateOracleArgs("-s", "patch", entry.DBName, dbVersion, entry.OracleHome)...)
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	return marshal_oracle.Patches(out)
}

// GetOracleDatabaseFeatureUsageStat get
func (wf *WindowsFetcherImpl) GetOracleDatabaseFeatureUsageStat(entry agentmodel.OratabEntry, dbVersion string) ([]model.OracleDatabaseFeatureUsageStat, error) {
	out, err := wf.executeWithDeadline(FetcherStandardTimeOut, "win.ps1", wf.CreateOracleArgs("-s", "opt", entry.DBName, dbVersion, entry.OracleHome)...)
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	return marshal_oracle.DatabaseFeatureUsageStat(out)
}

// GetOracleDatabaseLicenses get
func (wf *WindowsFetcherImpl) GetOracleDatabaseLicenses(entry agentmodel.OratabEntry,
	dbVersion, hardwareAbstractionTechnology string, hostCoreFactor float64,
) ([]model.OracleDatabaseLicense, error) {
	out, err := wf.executeWithDeadline(FetcherStandardTimeOut, "win.ps1", wf.CreateOracleArgs("-s", "license", entry.DBName, dbVersion, hardwareAbstractionTechnology, entry.OracleHome)...)
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	return marshal_oracle.Licenses(out)
}

// GetOracleDatabaseADDMs get
func (wf *WindowsFetcherImpl) GetOracleDatabaseADDMs(entry agentmodel.OratabEntry) ([]model.OracleDatabaseAddm, error) {
	out, err := wf.executeWithDeadline(FetcherStandardTimeOut, "win.ps1", wf.CreateOracleArgs("-s", "addm", entry.DBName, entry.OracleHome)...)
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	return marshal_oracle.Addms(out)
}

// GetOracleDatabaseSegmentAdvisors get
func (wf *WindowsFetcherImpl) GetOracleDatabaseSegmentAdvisors(entry agentmodel.OratabEntry) ([]model.OracleDatabaseSegmentAdvisor, error) {
	out, err := wf.executeWithDeadline(FetcherStandardTimeOut, "win.ps1", wf.CreateOracleArgs("-s", "segmentadvisor", entry.DBName, entry.OracleHome)...)
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	return marshal_oracle.SegmentAdvisor(out)
}

// GetOracleDatabasePSUs get
func (wf *WindowsFetcherImpl) GetOracleDatabasePSUs(entry agentmodel.OratabEntry, dbVersion string) ([]model.OracleDatabasePSU, error) {
	out, err := wf.executeWithDeadline(FetcherStandardTimeOut, "win.ps1", wf.CreateOracleArgs("-s", "psu", entry.DBName, dbVersion, entry.OracleHome)...)
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	return marshal_oracle.PSU(out), nil
}

// GetOracleDatabaseBackups get
func (wf *WindowsFetcherImpl) GetOracleDatabaseBackups(entry agentmodel.OratabEntry) ([]model.OracleDatabaseBackup, error) {
	out, err := wf.executeWithDeadline(FetcherStandardTimeOut, "win.ps1", wf.CreateOracleArgs("-s", "backup", entry.DBName, entry.OracleHome)...)
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	return marshal_oracle.Backups(out)
}

// GetOracleDatabaseCheckPDB get
func (wf *WindowsFetcherImpl) GetOracleDatabaseCheckPDB(entry agentmodel.OratabEntry) (bool, error) {
	wf.log.Error(notImplementedWindows)
	return false, ercutils.NewError(notImplementedWindows)
}

// GetOracleDatabasePDBs get
func (wf *WindowsFetcherImpl) GetOracleDatabasePDBs(entry agentmodel.OratabEntry) ([]model.OracleDatabasePluggableDatabase, error) {
	wf.log.Error(notImplementedWindows)
	return nil, ercutils.NewError(notImplementedWindows)
}

// GetOracleDatabasePDBTablespaces get
func (wf *WindowsFetcherImpl) GetOracleDatabasePDBTablespaces(entry agentmodel.OratabEntry, pdb string) ([]model.OracleDatabaseTablespace, error) {
	wf.log.Error(notImplementedWindows)
	return nil, ercutils.NewError(notImplementedWindows)
}

// GetOracleDatabasePDBSchemas get
func (wf *WindowsFetcherImpl) GetOracleDatabasePDBSchemas(entry agentmodel.OratabEntry, pdb string) ([]model.OracleDatabaseSchema, error) {
	wf.log.Error(notImplementedWindows)
	return nil, ercutils.NewError(notImplementedWindows)
}

// GetOracleDatabaseServices get
func (wf *WindowsFetcherImpl) GetOracleDatabaseServices(entry agentmodel.OratabEntry) ([]model.OracleDatabaseService, error) {
	out, err := wf.executeWithDeadline(FetcherStandardTimeOut, "win.ps1", wf.CreateOracleArgs("-s", "services", entry.DBName, entry.OracleHome)...)
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	return marshal_oracle.Services(out)
}

// GetMicrosoftSQLServerInstances get
func (wf *WindowsFetcherImpl) GetMicrosoftSQLServerInstances() ([]agentmodel.ListInstanceOutputModel, error) {
	out, err := wf.executeWithDeadline(FetcherStandardTimeOut, "ercoleAgentMsSQLServer-Gather.ps1", "-action", "listInstances")
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	return marshal_microsoft.ListInstances(out)
}

// GetMicrosoftSQLServerInstanceInfo get
func (wf *WindowsFetcherImpl) GetMicrosoftSQLServerInstanceInfo(conn string, inst *model.MicrosoftSQLServerInstance) error {
	out, err := wf.executeWithDeadline(FetcherStandardTimeOut, "ercoleAgentMsSQLServer-Gather.ps1", "-action", "dbmounted", "-instance", conn)
	if err != nil {
		return ercutils.NewError(err)
	}

	return marshal_microsoft.DbMounted(out, inst)
}

// GetMicrosoftSQLServerInstanceEdition get
func (wf *WindowsFetcherImpl) GetMicrosoftSQLServerInstanceEdition(conn string, inst *model.MicrosoftSQLServerInstance) error {
	out, err := wf.executeWithDeadline(FetcherStandardTimeOut, "ercoleAgentMsSQLServer-Gather.ps1", "-action", "edition", "-instance", conn)
	if err != nil {
		return ercutils.NewError(err)
	}

	if err := marshal_microsoft.Edition(out, inst); err != nil {
		return ercutils.NewError(err)
	}

	return nil
}

// GetMicrosoftSQLServerInstanceLicensingInfo get
func (wf *WindowsFetcherImpl) GetMicrosoftSQLServerInstanceLicensingInfo(conn string, inst *model.MicrosoftSQLServerInstance) error {
	out, err := wf.executeWithDeadline(FetcherStandardTimeOut, "ercoleAgentMsSQLServer-Gather.ps1", "-action", "licensingInfo", "-instance", conn)
	if err != nil {
		return ercutils.NewError(err)
	}

	if err := marshal_microsoft.LicensingInfo(out, inst); err != nil {
		return ercutils.NewError(err)
	}

	return nil
}

// GetMicrosoftSQLServerInstanceDatabase get
func (wf *WindowsFetcherImpl) GetMicrosoftSQLServerInstanceDatabase(conn string) ([]model.MicrosoftSQLServerDatabase, error) {
	out, err := wf.executeWithDeadline(FetcherStandardTimeOut, "ercoleAgentMsSQLServer-Gather.ps1", "-action", "db", "-instance", conn)
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	return marshal_microsoft.ListDatabases(out)
}

// GetMicrosoftSQLServerInstanceDatabaseBackups get
func (wf *WindowsFetcherImpl) GetMicrosoftSQLServerInstanceDatabaseBackups(conn string) ([]agentmodel.DbBackupsModel, error) {
	out, err := wf.executeWithDeadline(FetcherStandardTimeOut, "ercoleAgentMsSQLServer-Gather.ps1", "-action", "backup_schedule", "-instance", conn)
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	return marshal_microsoft.BackupSchedule(out)
}

// GetMicrosoftSQLServerInstanceDatabaseSchemas get
func (wf *WindowsFetcherImpl) GetMicrosoftSQLServerInstanceDatabaseSchemas(conn string) ([]agentmodel.DbSchemasModel, error) {
	out, err := wf.executeWithDeadline(FetcherStandardTimeOut, "ercoleAgentMsSQLServer-Gather.ps1", "-action", "schema", "-instance", conn)
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	return marshal_microsoft.Schemas(out)
}

// GetMicrosoftSQLServerInstanceDatabaseTablespaces get
func (wf *WindowsFetcherImpl) GetMicrosoftSQLServerInstanceDatabaseTablespaces(conn string) ([]agentmodel.DbTablespacesModel, error) {
	out, err := wf.executeWithDeadline(FetcherStandardTimeOut, "ercoleAgentMsSQLServer-Gather.ps1", "-action", "ts", "-instance", conn)
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	return marshal_microsoft.Tablespaces(out)
}

// GetMicrosoftSQLServerInstancePatches get
func (wf *WindowsFetcherImpl) GetMicrosoftSQLServerInstancePatches(conn string) ([]model.MicrosoftSQLServerPatch, error) {
	out, err := wf.executeWithDeadline(FetcherStandardTimeOut, "ercoleAgentMsSQLServer-Gather.ps1", "-action", "patches", "-instance", conn)
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	return marshal_microsoft.Patches(out)
}

// GetMicrosoftSQLServerProductFeatures get
func (wf *WindowsFetcherImpl) GetMicrosoftSQLServerProductFeatures(conn string) ([]model.MicrosoftSQLServerProductFeature, error) {
	out, err := wf.executeWithDeadline(FetcherStandardTimeOut, "ercoleAgentMsSQLServer-Gather.ps1", "-action", "sqlFeatures", "-instance", conn)
	if err != nil {
		return nil, ercutils.NewError(err)
	}

	return marshal_microsoft.Features(out)
}

func (wf *WindowsFetcherImpl) GetMySQLVersion(connection config.MySQLInstanceConnection) (string, error) {
	wf.log.Error(notImplementedWindows)
	return "", ercutils.NewError(notImplementedWindows)
}

func (wf *WindowsFetcherImpl) GetMySQLInstance(connection config.MySQLInstanceConnection) (*model.MySQLInstance, error) {
	wf.log.Error(notImplementedWindows)
	return nil, ercutils.NewError(notImplementedWindows)
}

func (wf *WindowsFetcherImpl) GetMySQLOldInstance(connection config.MySQLInstanceConnection) (*model.MySQLInstance, error) {
	wf.log.Error(notImplementedWindows)
	return nil, ercutils.NewError(notImplementedWindows)
}

func (wf *WindowsFetcherImpl) GetMySQLDatabases(connection config.MySQLInstanceConnection) ([]model.MySQLDatabase, error) {
	wf.log.Error(notImplementedWindows)
	return nil, ercutils.NewError(notImplementedWindows)
}

func (wf *WindowsFetcherImpl) GetMySQLTableSchemas(connection config.MySQLInstanceConnection) ([]model.MySQLTableSchema, error) {
	wf.log.Error(notImplementedWindows)
	return nil, ercutils.NewError(notImplementedWindows)
}

func (wf *WindowsFetcherImpl) GetMySQLSegmentAdvisors(connection config.MySQLInstanceConnection) ([]model.MySQLSegmentAdvisor, error) {
	wf.log.Error(notImplementedWindows)
	return nil, ercutils.NewError(notImplementedWindows)
}

func (wf *WindowsFetcherImpl) GetMySQLHighAvailability(connection config.MySQLInstanceConnection) (bool, error) {
	wf.log.Error(notImplementedWindows)
	return false, ercutils.NewError(notImplementedWindows)
}

func (wf *WindowsFetcherImpl) GetMySQLUUID() (string, error) {
	wf.log.Error(notImplementedWindows)
	return "", ercutils.NewError(notImplementedWindows)
}

func (wf *WindowsFetcherImpl) GetMySQLSlaveHosts(connection config.MySQLInstanceConnection) (bool, []string, error) {
	wf.log.Error(notImplementedWindows)
	return false, nil, ercutils.NewError(notImplementedWindows)
}

func (wf *WindowsFetcherImpl) GetMySQLSlaveStatus(connection config.MySQLInstanceConnection) (bool, *string, error) {
	wf.log.Error(notImplementedWindows)
	return false, nil, ercutils.NewError(notImplementedWindows)
}

func (wf *WindowsFetcherImpl) GetCloudMembership() (string, error) {
	return model.CloudMembershipUnknown, nil
}

func (wf *WindowsFetcherImpl) GetPostgreSQLSetting(instanceConnection config.PostgreSQLInstanceConnection) (*model.PostgreSQLSetting, error) {
	wf.log.Error(notImplementedWindows)
	return nil, ercutils.NewError(notImplementedWindows)
}

func (wf *WindowsFetcherImpl) GetPostgreSQLInstance(instanceConnection config.PostgreSQLInstanceConnection, v10 bool) (*model.PostgreSQLInstance, error) {
	wf.log.Error(notImplementedWindows)
	return nil, ercutils.NewError(notImplementedWindows)
}

func (wf *WindowsFetcherImpl) GetPostgreSQLDbNameList(instanceConnection config.PostgreSQLInstanceConnection) ([]string, error) {
	wf.log.Error(notImplementedWindows)
	return nil, ercutils.NewError(notImplementedWindows)
}
func (wf *WindowsFetcherImpl) GetPostgreSQLDbSchemaNameList(instanceConnection config.PostgreSQLInstanceConnection, dbname string) ([]string, error) {
	wf.log.Error(notImplementedWindows)
	return nil, ercutils.NewError(notImplementedWindows)
}
func (wf *WindowsFetcherImpl) GetPostgreSQLDatabase(instanceConnection config.PostgreSQLInstanceConnection, dbname string, v10 bool) (*model.PostgreSQLDatabase, error) {
	wf.log.Error(notImplementedWindows)
	return nil, ercutils.NewError(notImplementedWindows)
}
func (wf *WindowsFetcherImpl) GetPostgreSQLSchema(instanceConnection config.PostgreSQLInstanceConnection, dbname string, schemaName string, v10 bool) (*model.PostgreSQLSchema, error) {
	wf.log.Error(notImplementedWindows)
	return nil, ercutils.NewError(notImplementedWindows)
}
