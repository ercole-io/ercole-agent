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
	"strings"

	"github.com/ercole-io/ercole-agent/agentmodel"
	"github.com/ercole-io/ercole/model"
)

func (b *CommonBuilder) getMicrosoftSQLServerFeature() *model.MicrosoftSQLServerFeature {
	instances := b.fetcher.GetMicrosoftSQLServerInstances()

	return &model.MicrosoftSQLServerFeature{
		Instances: b.getMicrosoftSQLServerInstances(instances),
		Features:  b.getMicrosoftSQLServerProductFeatures(instances[0].ConnString),
		Patches:   b.fetcher.GetMicrosoftSQLServerInstancePatches(instances[0].ConnString),
	}
}

func (b *CommonBuilder) getMicrosoftSQLServerInstances(instanceList []agentmodel.ListInstanceOutputModel) []model.MicrosoftSQLServerInstance {
	instances := make([]model.MicrosoftSQLServerInstance, len(instanceList))

	for i, v := range instanceList {
		instances[i].Name = v.Name
		instances[i].Status = v.Status
		instances[i].DisplayName = v.DisplayName

		if v.Status == "Running" {
			instances[i].Platform = "Windows"
			b.fetcher.GetMicrosoftSQLServerInstanceInfo(v.ConnString, &instances[i])

			b.fetcher.GetMicrosoftSQLServerInstanceEdition(v.ConnString, &instances[i])
			b.fetcher.GetMicrosoftSQLServerInstanceLicensingInfo(v.ConnString, &instances[i])

			instances[i].Databases = b.fetcher.GetMicrosoftSQLServerInstanceDatabase(v.ConnString)
			instances[i].CollationName = instances[i].Databases[0].CollationName
			dbsMap := make(map[string]*model.MicrosoftSQLServerDatabase)

			for j, db := range instances[i].Databases {
				instances[i].Databases[j].Backups = make([]model.MicrosoftSQLServerDatabaseBackup, 0)
				instances[i].Databases[j].Schemas = make([]model.MicrosoftSQLServerDatabaseSchema, 0)
				instances[i].Databases[j].Tablespaces = make([]model.MicrosoftSQLServerDatabaseTablespace, 0)
				dbsMap[db.Name] = &instances[i].Databases[j]
			}

			backupSchedules := b.fetcher.GetMicrosoftSQLServerInstanceDatabaseBackups(v.ConnString)
			for _, v := range backupSchedules {
				dbsMap[v.DatabaseName].Backups = make([]model.MicrosoftSQLServerDatabaseBackup, len(v.Data))
				for i, b := range v.Data {
					dbsMap[v.DatabaseName].Backups[i].AvgBckSize = b.AvgBckSizeGB
					dbsMap[v.DatabaseName].Backups[i].BackupType = b.BackupType
					dbsMap[v.DatabaseName].Backups[i].Hour = b.Hour
					dbsMap[v.DatabaseName].Backups[i].WeekDays = strings.Split(b.WeekDays, ",")
				}
			}

			schemas := b.fetcher.GetMicrosoftSQLServerInstanceDatabaseSchemas(v.ConnString)
			for _, v := range schemas {
				dbsMap[v.DatabaseName].Schemas = make([]model.MicrosoftSQLServerDatabaseSchema, len(v.Data))
				for i, b := range v.Data {
					dbsMap[v.DatabaseName].Schemas[i].AllocatedSpace = int(b.AllocatedMB * 1024 * 1024)
					dbsMap[v.DatabaseName].Schemas[i].UsedSpace = int(b.UsedMB * 1024 * 1024)
					dbsMap[v.DatabaseName].Schemas[i].AllocationType = b.AllocationType
				}
			}

			tablespaces := b.fetcher.GetMicrosoftSQLServerInstanceDatabaseTablespaces(v.ConnString)
			for _, v := range tablespaces {
				dbsMap[v.DatabaseName].Tablespaces = make([]model.MicrosoftSQLServerDatabaseTablespace, len(v.Data))
				for i, b := range v.Data {
					dbsMap[v.DatabaseName].Tablespaces[i].Alloc = int(b.AllocMB * 1024 * 1024)
					dbsMap[v.DatabaseName].Tablespaces[i].Used = int(b.UsedMB * 1024 * 1024)
					dbsMap[v.DatabaseName].Tablespaces[i].Growth = b.Growth
					dbsMap[v.DatabaseName].Tablespaces[i].GrowthUnit = b.GrowthUnit
					dbsMap[v.DatabaseName].Tablespaces[i].FileType = b.FileType
					dbsMap[v.DatabaseName].Tablespaces[i].Status = b.Status
					dbsMap[v.DatabaseName].Tablespaces[i].Filename = b.Filename
				}
			}
		}
	}

	return instances
}

func (b *CommonBuilder) getMicrosoftSQLServerProductFeatures(connString string) []model.MicrosoftSQLServerProductFeature {
	return b.fetcher.GetMicrosoftSQLServerProductFeatures(connString)
}
