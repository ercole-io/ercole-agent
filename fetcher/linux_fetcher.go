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

// LinuxFetcherImpl implementation
type LinuxFetcherImpl struct {
	Configuration config.Configuration
}

// Execute Execute specific fetcher by name
func (lf *LinuxFetcherImpl) Execute(fetcherName string, params ...string) []byte {
	cmdName := config.GetBaseDir() + "/fetch/linux" + fetcherName + ".sh"
	log.Println("Fetching", cmdName, strings.Join(params, " "))

	cmd := exec.Command(cmdName, params...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

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

func (lf *LinuxFetcherImpl) GetSpecializedFetcherName() string {
	return "linux"
}

//// GetHost get
//func (lf *LinuxFetcherImpl) GetHost() model.Host {
//	out := lf.Execute("host")
//	return marshal.Host(out)
//}
//
//// GetFilesystems get
//func (lf *LinuxFetcherImpl) GetFilesystems() []model.Filesystem {
//	out := lf.Execute("filesystem")
//	return marshal.Filesystems(out)
//}
//
//// GetOratabEntries get
//func (lf *LinuxFetcherImpl) GetOratabEntries() []model.OratabEntry {
//	out := lf.Execute("oratab", lf.Configuration.Oratab)
//	return marshal.Oratab(out)
//}
//
//// GetDbStatus get
//func (lf *LinuxFetcherImpl) GetDbStatus(entry model.OratabEntry) string {
//	out := lf.Execute("dbstatus", entry.DBName, entry.OracleHome)
//	return strings.TrimSpace(string(out))
//}
//
//// GetMountedDb get
//func (lf *LinuxFetcherImpl) GetMountedDb(entry model.OratabEntry) model.Database {
//	out := lf.Execute("dbmounted", entry.DBName, entry.OracleHome)
//	return marshal.Database(out)
//}
//
//// GetDbVersion get
//func (lf *LinuxFetcherImpl) GetDbVersion(entry model.OratabEntry) string {
//	out := lf.Execute("dbversion", entry.DBName, entry.OracleHome)
//	return strings.Split(string(out), ".")[0]
//}
//
//// RunStats execute stats script
//func (lf *LinuxFetcherImpl) RunStats(entry model.OratabEntry) {
//	lf.Execute("stats", entry.DBName, entry.OracleHome)
//}
//
//// GetOpenDb get
//func (lf *LinuxFetcherImpl) GetOpenDb(entry model.OratabEntry) model.Database {
//	out := lf.Execute("db", entry.DBName, entry.OracleHome, strconv.Itoa(lf.Configuration.AWR))
//	return marshal.Database(out)
//}
//
//// GetTablespaces get
//func (lf *LinuxFetcherImpl) GetTablespaces(entry model.OratabEntry) []model.Tablespace {
//	out := lf.Execute("tablespace", entry.DBName, entry.OracleHome)
//	return marshal.Tablespaces(out)
//}
//
//// GetSchemas get
//func (lf *LinuxFetcherImpl) GetSchemas(entry model.OratabEntry) []model.Schema {
//	out := lf.Execute("schema", entry.DBName, entry.OracleHome)
//	return marshal.Schemas(out)
//}
//
//// GetPatches get
//func (lf *LinuxFetcherImpl) GetPatches(entry model.OratabEntry, dbVersion string) []model.Patch {
//	out := lf.Execute("patch", entry.DBName, dbVersion, entry.OracleHome)
//	return marshal.Patches(out)
//}
//
//// GetFeatures get
//func (lf *LinuxFetcherImpl) GetFeatures(entry model.OratabEntry, dbVersion string) (features []model.Feature) {
//	out := lf.Execute("feature", entry.DBName, dbVersion, entry.OracleHome)
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
//func (lf *LinuxFetcherImpl) GetFeatures2(entry model.OratabEntry, dbVersion string) []model.Feature2 {
//	out := lf.Execute("opt", entry.DBName, dbVersion, entry.OracleHome)
//	return marshal.Features2(out)
//}
//
//// GetLicenses get
//func (lf *LinuxFetcherImpl) GetLicenses(entry model.OratabEntry, dbVersion, hostType string) []model.License {
//	out := lf.Execute("license", entry.DBName, dbVersion, hostType, entry.OracleHome)
//	return marshal.Licenses(out)
//}
//
//// GetADDMs get
//func (lf *LinuxFetcherImpl) GetADDMs(entry model.OratabEntry) []model.Addm {
//	out := lf.Execute("addm", entry.DBName, entry.OracleHome)
//	return marshal.Addms(out)
//}
//
//// GetSegmentAdvisors get
//func (lf *LinuxFetcherImpl) GetSegmentAdvisors(entry model.OratabEntry) []model.SegmentAdvisor {
//	out := lf.Execute("segmentadvisor", entry.DBName, entry.OracleHome)
//	return marshal.SegmentAdvisor(out)
//}
//
//// GetLastPSUs get
//func (lf *LinuxFetcherImpl) GetLastPSUs(entry model.OratabEntry, dbVersion string) []model.PSU {
//	out := lf.Execute("psu", entry.DBName, dbVersion, entry.OracleHome)
//	return marshal.PSU(out)
//}
//
//// GetBackups get
//func (lf *LinuxFetcherImpl) GetBackups(entry model.OratabEntry) []model.Backup {
//	out := lf.Execute("backup", entry.DBName, entry.OracleHome)
//	return marshal.Backups(out)
//}
