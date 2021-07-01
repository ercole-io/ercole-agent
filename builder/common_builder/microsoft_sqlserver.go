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

	"github.com/ercole-io/ercole-agent/v2/agentmodel"
	"github.com/ercole-io/ercole/v2/model"
	ercutils "github.com/ercole-io/ercole/v2/utils"
	"github.com/hashicorp/go-multierror"
)

func (b *CommonBuilder) getMicrosoftSQLServerFeature() (*model.MicrosoftSQLServerFeature, error) {
	instances := b.fetcher.GetMicrosoftSQLServerInstances()

	sqlServer := model.MicrosoftSQLServerFeature{
		Patches: b.fetcher.GetMicrosoftSQLServerInstancePatches(instances[0].ConnString),
	}

	var merr, err error
	sqlServer.Instances, err = b.getMicrosoftSQLServerInstances(instances)
	if err != nil {
		merr = multierror.Append(merr, ercutils.NewError(err))
	}

	sqlServer.Features, err = b.fetcher.GetMicrosoftSQLServerProductFeatures(instances[0].ConnString)
	if err != nil {
		merr = multierror.Append(merr, ercutils.NewError(err))
	}

	if merr != nil {
		return nil, merr
	}
	return &sqlServer, nil
}

func (b *CommonBuilder) getMicrosoftSQLServerInstances(instanceList []agentmodel.ListInstanceOutputModel) ([]model.MicrosoftSQLServerInstance, error) {
	instances := make([]model.MicrosoftSQLServerInstance, len(instanceList))
	var merr, err error

	for i, v := range instanceList {
		instances[i].Name = v.Name
		instances[i].Status = v.Status
		instances[i].DisplayName = v.DisplayName

		if v.Status == "Running" {
			instances[i].Platform = "Windows"
			b.fetcher.GetMicrosoftSQLServerInstanceInfo(v.ConnString, &instances[i])

			b.fetcher.GetMicrosoftSQLServerInstanceEdition(v.ConnString, &instances[i])
			b.fetcher.GetMicrosoftSQLServerInstanceLicensingInfo(v.ConnString, &instances[i])

			if instances[i].Databases, err = b.fetcher.GetMicrosoftSQLServerInstanceDatabase(v.ConnString); err != nil {
				merr = multierror.Append(merr, ercutils.NewError(err))
			}
			instances[i].CollationName = instances[i].Databases[0].CollationName
			dbsMap := make(map[string]*model.MicrosoftSQLServerDatabase)

			for j, db := range instances[i].Databases {
				instances[i].Databases[j].Backups = make([]model.MicrosoftSQLServerDatabaseBackup, 0)
				instances[i].Databases[j].Schemas = make([]model.MicrosoftSQLServerDatabaseSchema, 0)
				instances[i].Databases[j].Tablespaces = make([]model.MicrosoftSQLServerDatabaseTablespace, 0)
				dbsMap[db.Name] = &instances[i].Databases[j]
			}

			backupSchedules, err := b.fetcher.GetMicrosoftSQLServerInstanceDatabaseBackups(v.ConnString)
			if err != nil {
				merr = multierror.Append(merr, ercutils.NewError(err))
			}
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

	if merr != nil {
		return nil, merr
	}
	return instances, nil
}
