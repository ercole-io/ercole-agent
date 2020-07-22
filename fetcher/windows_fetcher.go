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
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/ercole-io/ercole-agent/agentmodel"
	"github.com/ercole-io/ercole-agent/config"
	"github.com/ercole-io/ercole-agent/logger"
	"github.com/ercole-io/ercole-agent/marshal"
	marshal_microsoft "github.com/ercole-io/ercole-agent/marshal/microsoft"
	marshal_oracle "github.com/ercole-io/ercole-agent/marshal/oracle"
	"github.com/ercole-io/ercole/model"
)

// WindowsFetcherImpl SpecializedFetcher implementation for windows
type WindowsFetcherImpl struct {
	configuration config.Configuration
	log           logger.Logger
}

const notImplementedWindows = "Not yet implemented for windows"

// NewWindowsFetcherImpl constructor
func NewWindowsFetcherImpl(conf config.Configuration, log logger.Logger) *WindowsFetcherImpl {
	return &WindowsFetcherImpl{
		conf,
		log,
	}
}

// Execute Execute specific fetcher by name
func (wf *WindowsFetcherImpl) execute(fetcherName string, params ...string) []byte {
	var (
		cmd    *exec.Cmd
		err    error
		psexe  string
		stdout bytes.Buffer
		stderr bytes.Buffer
	)

	baseDir := config.GetBaseDir()

	psexe, err = exec.LookPath("powershell.exe")
	if err != nil {
		wf.log.Fatal(psexe)
	}

	if wf.configuration.ForcePwshVersion == "0" {
		params = append([]string{"-ExecutionPolicy", "Bypass", "-File", baseDir + "\\fetch\\windows\\" + fetcherName}, params...)
	} else {
		params = append([]string{"-version", wf.configuration.ForcePwshVersion, "-ExecutionPolicy", "Bypass", "-File", baseDir + "\\fetch\\windows\\" + fetcherName}, params...)
	}

	wf.log.Info("Fetching " + psexe + " " + strings.Join(params, " "))

	cmd = exec.Command(psexe, params...)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if len(stderr.Bytes()) > 0 {
		wf.log.Error(string(stderr.Bytes()))
	}

	if err != nil {
		if fetcherName != "dbstatus" {
			return []byte("UNREACHABLE")
		}

		wf.log.Fatal(err)
	}

	return stdout.Bytes()
}

// SetUser not implemented
func (wf *WindowsFetcherImpl) SetUser(username string) error {
	wf.log.Error(notImplementedWindows)
	return fmt.Errorf(notImplementedWindows)
}

// SetUserAsCurrent set user used by fetcher to run commands as current process user
func (wf *WindowsFetcherImpl) SetUserAsCurrent() error {
	wf.log.Error(notImplementedWindows)
	return fmt.Errorf(notImplementedWindows)
}

// GetClusters not implemented
func (wf *WindowsFetcherImpl) GetClusters(hv config.Hypervisor) []model.ClusterInfo {
	wf.log.Error(notImplementedWindows)

	return make([]model.ClusterInfo, 0)
}

// GetVirtualMachines return VMWare virtual machines infos from the given hyperVisor
func (wf *WindowsFetcherImpl) GetVirtualMachines(hv config.Hypervisor) map[string][]model.VMInfo {
	wf.log.Error(notImplementedWindows)

	return make(map[string][]model.VMInfo, 0)
}

// GetOracleExadataComponents get
func (wf *WindowsFetcherImpl) GetOracleExadataComponents() []model.OracleExadataComponent {
	wf.log.Error(notImplementedWindows)

	return make([]model.OracleExadataComponent, 0)
}

// GetOracleExadataCellDisks get
func (wf *WindowsFetcherImpl) GetOracleExadataCellDisks() map[agentmodel.StorageServerName][]model.OracleExadataCellDisk {
	wf.log.Error(notImplementedWindows)

	return make(map[agentmodel.StorageServerName][]model.OracleExadataCellDisk, 0)
}

// GetClustersMembershipStatus get
func (wf *WindowsFetcherImpl) GetClustersMembershipStatus() model.ClusterMembershipStatus {
	return model.ClusterMembershipStatus{
		OracleClusterware:    false,
		SunCluster:           false,
		VeritasClusterServer: false,
		HACMP:                false,
		OtherInfo:            nil,
	}
}

// GetHost get
func (wf *WindowsFetcherImpl) GetHost() model.Host {
	out := wf.execute("win.ps1", "-s", "host")
	return marshal.Host(out)
}

// GetFilesystems get
func (wf *WindowsFetcherImpl) GetFilesystems() []model.Filesystem {
	out := wf.execute("win.ps1", "-s", "filesystem")
	return marshal.Filesystems(out)
}

// GetOracleDatabaseOratabEntries get
func (wf *WindowsFetcherImpl) GetOracleDatabaseOratabEntries() []agentmodel.OratabEntry {
	out := wf.execute("win.ps1", "-s", "oratab", wf.configuration.Features.OracleDatabase.Oratab)
	return marshal_oracle.Oratab(out)
}

// GetOracleDatabaseDbStatus get
func (wf *WindowsFetcherImpl) GetOracleDatabaseDbStatus(entry agentmodel.OratabEntry) string {
	out := wf.execute("win.ps1", "-s", "dbstatus", entry.DBName, entry.OracleHome)
	return strings.TrimSpace(string(out))
}

// GetOracleDatabaseMountedDb get
func (wf *WindowsFetcherImpl) GetOracleDatabaseMountedDb(entry agentmodel.OratabEntry) model.OracleDatabase {
	out := wf.execute("win.ps1", "-s", "dbmounted", entry.DBName, entry.OracleHome)
	return marshal_oracle.Database(out)
}

// GetOracleDatabaseDbVersion get
func (wf *WindowsFetcherImpl) GetOracleDatabaseDbVersion(entry agentmodel.OratabEntry) string {
	out := wf.execute("win.ps1", "-s", "dbversion", entry.DBName, entry.OracleHome)
	return strings.Split(string(out), ".")[0]
}

// RunOracleDatabaseStats Execute stats script
func (wf *WindowsFetcherImpl) RunOracleDatabaseStats(entry agentmodel.OratabEntry) {
	wf.execute("win.ps1", "-s", "stats", entry.DBName, entry.OracleHome)
}

// GetOracleDatabaseOpenDb get
func (wf *WindowsFetcherImpl) GetOracleDatabaseOpenDb(entry agentmodel.OratabEntry) model.OracleDatabase {
	out := wf.execute("win.ps1", "-s", "db", entry.DBName, entry.OracleHome, strconv.Itoa(wf.configuration.Features.OracleDatabase.AWR))
	return marshal_oracle.Database(out)
}

// GetOracleDatabaseTablespaces get
func (wf *WindowsFetcherImpl) GetOracleDatabaseTablespaces(entry agentmodel.OratabEntry) []model.OracleDatabaseTablespace {
	out := wf.execute("win.ps1", "-s", "tablespace", entry.DBName, entry.OracleHome)
	return marshal_oracle.Tablespaces(out)
}

// GetOracleDatabaseSchemas get
func (wf *WindowsFetcherImpl) GetOracleDatabaseSchemas(entry agentmodel.OratabEntry) []model.OracleDatabaseSchema {
	out := wf.execute("win.ps1", "-s", "schema", entry.DBName, entry.OracleHome)
	return marshal_oracle.Schemas(out)
}

// GetOracleDatabasePatches get
func (wf *WindowsFetcherImpl) GetOracleDatabasePatches(entry agentmodel.OratabEntry, dbVersion string) []model.OracleDatabasePatch {
	out := wf.execute("win.ps1", "-s", "patch", entry.DBName, dbVersion, entry.OracleHome)
	return marshal_oracle.Patches(out)
}

// GetOracleDatabaseFeatureUsageStat get
func (wf *WindowsFetcherImpl) GetOracleDatabaseFeatureUsageStat(entry agentmodel.OratabEntry, dbVersion string) []model.OracleDatabaseFeatureUsageStat {
	out := wf.execute("win.ps1", "-s", "opt", entry.DBName, dbVersion, entry.OracleHome)
	return marshal_oracle.DatabaseFeatureUsageStat(out)
}

// GetOracleDatabaseLicenses get
func (wf *WindowsFetcherImpl) GetOracleDatabaseLicenses(entry agentmodel.OratabEntry, dbVersion, hardwareAbstractionTechnology string) []model.OracleDatabaseLicense {
	out := wf.execute("win.ps1", "-s", "license", entry.DBName, dbVersion, hardwareAbstractionTechnology, entry.OracleHome)
	return marshal_oracle.Licenses(out)
}

// GetOracleDatabaseADDMs get
func (wf *WindowsFetcherImpl) GetOracleDatabaseADDMs(entry agentmodel.OratabEntry) []model.OracleDatabaseAddm {
	out := wf.execute("win.ps1", "-s", "addm", entry.DBName, entry.OracleHome)
	return marshal_oracle.Addms(out)
}

// GetOracleDatabaseSegmentAdvisors get
func (wf *WindowsFetcherImpl) GetOracleDatabaseSegmentAdvisors(entry agentmodel.OratabEntry) []model.OracleDatabaseSegmentAdvisor {
	out := wf.execute("win.ps1", "-s", "segmentadvisor", entry.DBName, entry.OracleHome)
	return marshal_oracle.SegmentAdvisor(out)
}

// GetOracleDatabasePSUs get
func (wf *WindowsFetcherImpl) GetOracleDatabasePSUs(entry agentmodel.OratabEntry, dbVersion string) []model.OracleDatabasePSU {
	out := wf.execute("win.ps1", "-s", "psu", entry.DBName, dbVersion, entry.OracleHome)
	return marshal_oracle.PSU(out)
}

// GetOracleDatabaseBackups get
func (wf *WindowsFetcherImpl) GetOracleDatabaseBackups(entry agentmodel.OratabEntry) []model.OracleDatabaseBackup {
	out := wf.execute("win.ps1", "-s", "backup", entry.DBName, entry.OracleHome)
	return marshal_oracle.Backups(out)
}

// GetMicrosoftSQLServerInstances get
func (wf *WindowsFetcherImpl) GetMicrosoftSQLServerInstances() []agentmodel.ListInstanceOutputModel {
	out := wf.execute("ercoleAgentMsSQLServer-Gather.ps1", "-action", "listInstances")
	return marshal_microsoft.ListInstances(out)
}

// GetMicrosoftSQLServerInstanceInfo get
func (wf *WindowsFetcherImpl) GetMicrosoftSQLServerInstanceInfo(conn string, inst *model.MicrosoftSQLServerInstance) {
	out := wf.execute("ercoleAgentMsSQLServer-Gather.ps1", "-action", "dbmounted", "-instance", conn)
	marshal_microsoft.DbMounted(out, inst)
}

// GetMicrosoftSQLServerInstanceEdition get
func (wf *WindowsFetcherImpl) GetMicrosoftSQLServerInstanceEdition(conn string, inst *model.MicrosoftSQLServerInstance) {
	out := wf.execute("ercoleAgentMsSQLServer-Gather.ps1", "-action", "edition", "-instance", conn)
	marshal_microsoft.Edition(out, inst)
}

// GetMicrosoftSQLServerInstanceLicensingInfo get
func (wf *WindowsFetcherImpl) GetMicrosoftSQLServerInstanceLicensingInfo(conn string, inst *model.MicrosoftSQLServerInstance) {
	out := wf.execute("ercoleAgentMsSQLServer-Gather.ps1", "-action", "licensingInfo", "-instance", conn)
	marshal_microsoft.LicensingInfo(out, inst)
}

// GetMicrosoftSQLServerInstanceDatabase get
func (wf *WindowsFetcherImpl) GetMicrosoftSQLServerInstanceDatabase(conn string) []model.MicrosoftSQLServerDatabase {
	out := wf.execute("ercoleAgentMsSQLServer-Gather.ps1", "-action", "db", "-instance", conn)
	return marshal_microsoft.ListDatabases(out)
}

// GetMicrosoftSQLServerInstanceDatabaseBackups get
func (wf *WindowsFetcherImpl) GetMicrosoftSQLServerInstanceDatabaseBackups(conn string) []agentmodel.DbBackupsModel {
	out := wf.execute("ercoleAgentMsSQLServer-Gather.ps1", "-action", "backup_schedule", "-instance", conn)
	return marshal_microsoft.BackupSchedule(out)
}

// GetMicrosoftSQLServerInstanceDatabaseSchemas get
func (wf *WindowsFetcherImpl) GetMicrosoftSQLServerInstanceDatabaseSchemas(conn string) []agentmodel.DbSchemasModel {
	out := wf.execute("ercoleAgentMsSQLServer-Gather.ps1", "-action", "schema", "-instance", conn)
	return marshal_microsoft.Schemas(out)
}

// GetMicrosoftSQLServerInstanceDatabaseTablespaces get
func (wf *WindowsFetcherImpl) GetMicrosoftSQLServerInstanceDatabaseTablespaces(conn string) []agentmodel.DbTablespacesModel {
	out := wf.execute("ercoleAgentMsSQLServer-Gather.ps1", "-action", "ts", "-instance", conn)
	return marshal_microsoft.Tablespaces(out)
}

// GetMicrosoftSQLServerInstancePatches get
func (wf *WindowsFetcherImpl) GetMicrosoftSQLServerInstancePatches(conn string) []model.MicrosoftSQLServerPatch {
	out := wf.execute("ercoleAgentMsSQLServer-Gather.ps1", "-action", "patches", "-instance", conn)
	return marshal_microsoft.Patches(out)
}

// GetMicrosoftSQLServerProductFeatures get
func (wf *WindowsFetcherImpl) GetMicrosoftSQLServerProductFeatures(conn string) []model.MicrosoftSQLServerProductFeature {
	out := wf.execute("ercoleAgentMsSQLServer-Gather.ps1", "-action", "sqlFeatures", "-instance", conn)
	return marshal_microsoft.Features(out)
}
