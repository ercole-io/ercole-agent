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
	"log"
	"os/exec"
	"strings"

	"github.com/ercole-io/ercole-agent/config"
)

// WindowsFetcherImpl implemenentation
type WindowsFetcherImpl struct {
	Configuration config.Configuration
}

// Execute Execute specific fetcher by name
func (wf *WindowsFetcherImpl) Execute(fetcherName string, params ...string) []byte {
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
		log.Fatal(psexe)
	}
	if wf.Configuration.ForcePwshVersion == "0" {
		params = append([]string{"-ExecutionPolicy", "Bypass", "-File", baseDir + "\\fetch\\win.ps1", "-s", fetcherName}, params...)
	} else {
		params = append([]string{"-version", wf.Configuration.ForcePwshVersion, "-ExecutionPolicy", "Bypass", "-File", baseDir + "\\fetch\\win.ps1", "-s", fetcherName}, params...)
	}
	log.Println("Fetching " + psexe + " " + strings.Join(params, " "))

	cmd = exec.Command(psexe, params...)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if len(stderr.Bytes()) > 0 {
		log.Print(string(stderr.Bytes()))
	}

	if err != nil {
		if fetcherName != "dbstatus" {
			log.Fatal(err)
		} else {
			return []byte("UNREACHABLE") // fallback
		}
	}

	return stdout.Bytes()
}

func (wf *WindowsFetcherImpl) GetSpecializedFetcherName() string {
	return "windows"
}

//// GetHost get
//func (wf *WindowsFetcherImpl) GetHost() model.Host {
//	out := wf.Execute("host")
//	return marshal.Host(out)
//}
//
//// GetFilesystems get
//func (wf *WindowsFetcherImpl) GetFilesystems() []model.Filesystem {
//	out := wf.Execute("filesystem")
//	return marshal.Filesystems(out)
//}
//
//// GetOratabEntries get
//func (wf *WindowsFetcherImpl) GetOratabEntries() []model.OratabEntry {
//	out := wf.Execute("oratab", wf.Configuration.Oratab)
//	return marshal.Oratab(out)
//}
//
//// GetDbStatus get
//func (wf *WindowsFetcherImpl) GetDbStatus(entry model.OratabEntry) string {
//	out := wf.Execute("dbstatus", entry.DBName, entry.OracleHome)
//	return strings.TrimSpace(string(out))
//}
//
//// GetMountedDb get
//func (wf *WindowsFetcherImpl) GetMountedDb(entry model.OratabEntry) model.Database {
//	out := wf.Execute("dbmounted", entry.DBName, entry.OracleHome)
//	return marshal.Database(out)
//}
//
//// GetDbVersion get
//func (wf *WindowsFetcherImpl) GetDbVersion(entry model.OratabEntry) string {
//	out := wf.Execute("dbversion", entry.DBName, entry.OracleHome)
//	return strings.Split(string(out), ".")[0]
//}
//
//// RunStats execute stats script
//func (wf *WindowsFetcherImpl) RunStats(entry model.OratabEntry) {
//	wf.Execute("stats", entry.DBName, entry.OracleHome)
//}
//
//// GetOpenDb get
//func (wf *WindowsFetcherImpl) GetOpenDb(entry model.OratabEntry) model.Database {
//	out := wf.Execute("db", entry.DBName, entry.OracleHome, strconv.Itoa(wf.Configuration.AWR))
//	return marshal.Database(out)
//}
//
//// GetTablespaces get
//func (wf *WindowsFetcherImpl) GetTablespaces(entry model.OratabEntry) []model.Tablespace {
//	out := wf.Execute("tablespace", entry.DBName, entry.OracleHome)
//	return marshal.Tablespaces(out)
//}
//
//// GetSchemas get
//func (wf *WindowsFetcherImpl) GetSchemas(entry model.OratabEntry) []model.Schema {
//	out := wf.Execute("schema", entry.DBName, entry.OracleHome)
//	return marshal.Schemas(out)
//}
//
//// GetPatches get
//func (wf *WindowsFetcherImpl) GetPatches(entry model.OratabEntry, dbVersion string) []model.Patch {
//	out := wf.Execute("patch", entry.DBName, dbVersion, entry.OracleHome)
//	return marshal.Patches(out)
//}
//
//// GetFeatures get
//func (wf *WindowsFetcherImpl) GetFeatures(entry model.OratabEntry, dbVersion string) (features []model.Feature) {
//	out := wf.Execute("feature", entry.DBName, dbVersion, entry.OracleHome)
//
//	if strings.Contains(string(out), "deadlocked on readable physical standby") {
//		log.Println("Detected bug active dataguard 2311894.1!")
//		features = []model.Feature{}
//
//	} else if strings.Contains(string(out), "ORA-01555: snapshot too old: rollback segment number") {
//		log.Println("Detected error on active dataguard ORA-01555!")
//		features = []model.Feature{}
//
//	} else {
//		features = marshal.Features(out)
//	}
//
//	return
//}
//
//// GetFeatures2 get
//func (wf *WindowsFetcherImpl) GetFeatures2(entry model.OratabEntry, dbVersion string) []model.Feature2 {
//	out := wf.Execute("opt", entry.DBName, dbVersion, entry.OracleHome)
//	return marshal.Features2(out)
//}
//
//// GetLicenses get
//func (wf *WindowsFetcherImpl) GetLicenses(entry model.OratabEntry, dbVersion, hostType string) []model.License {
//	out := wf.Execute("license", entry.DBName, dbVersion, hostType, entry.OracleHome)
//	return marshal.Licenses(out)
//}
//
//// GetADDMs get
//func (wf *WindowsFetcherImpl) GetADDMs(entry model.OratabEntry) []model.Addm {
//	out := wf.Execute("addm", entry.DBName, entry.OracleHome)
//	return marshal.Addms(out)
//}
//
//// GetSegmentAdvisors get
//func (wf *WindowsFetcherImpl) GetSegmentAdvisors(entry model.OratabEntry) []model.SegmentAdvisor {
//	out := wf.Execute("segmentadvisor", entry.DBName, entry.OracleHome)
//	return marshal.SegmentAdvisor(out)
//}
//
//// GetLastPSUs get
//func (wf *WindowsFetcherImpl) GetLastPSUs(entry model.OratabEntry, dbVersion string) []model.PSU {
//	out := wf.Execute("psu", entry.DBName, dbVersion, entry.OracleHome)
//	return marshal.PSU(out)
//}
//
//// GetBackups get
//func (wf *WindowsFetcherImpl) GetBackups(entry model.OratabEntry) []model.Backup {
//	out := wf.Execute("backup", entry.DBName, entry.OracleHome)
//	return marshal.Backups(out)
//}
//
