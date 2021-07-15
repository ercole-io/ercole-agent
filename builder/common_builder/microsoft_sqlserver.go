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

	"github.com/ercole-io/ercole/v2/model"
	"github.com/hashicorp/go-multierror"
)

func (b *CommonBuilder) getMicrosoftSQLServerFeature() (*model.MicrosoftSQLServerFeature, error) {
	var sqlServer model.MicrosoftSQLServerFeature
	var merr, err error
	var connString string

	sqlServer.Instances, connString, err = b.getMicrosoftSQLServerInstances()
	if err != nil {
		b.log.Error(err)
		merr = multierror.Append(merr, err)
	}

	if len(sqlServer.Instances) > 0 {
		sqlServer.Features, err = b.fetcher.GetMicrosoftSQLServerProductFeatures(connString)
		if err != nil {
			b.log.Error(err)
			merr = multierror.Append(merr, err)
		}

		if sqlServer.Patches, err = b.fetcher.GetMicrosoftSQLServerInstancePatches(connString); err != nil {
			b.log.Error(err)
			merr = multierror.Append(merr, err)
		}
	}

	return &sqlServer, merr
}

func (b *CommonBuilder) getMicrosoftSQLServerInstances() ([]model.MicrosoftSQLServerInstance, string, error) {
	fetchedInstances, err := b.fetcher.GetMicrosoftSQLServerInstances()
	if err != nil {
		return nil, "", err
	}

	instances := make([]model.MicrosoftSQLServerInstance, 0)
	var merr error

	for _, out := range fetchedInstances {
		instance := model.MicrosoftSQLServerInstance{}
		var instanceErr error

		instance.Name = out.Name
		instance.Status = out.Status
		instance.DisplayName = out.DisplayName

		if out.Status != "Running" {
			continue
		}

		instance.Platform = "Windows"

		if err := b.fetcher.GetMicrosoftSQLServerInstanceInfo(out.ConnString, &instance); err != nil {
			instanceErr = multierror.Append(instanceErr, err)
		}

		if err := b.fetcher.GetMicrosoftSQLServerInstanceEdition(out.ConnString, &instance); err != nil {
			instanceErr = multierror.Append(instanceErr, err)
		}

		if err := b.fetcher.GetMicrosoftSQLServerInstanceLicensingInfo(out.ConnString, &instance); err != nil {
			instanceErr = multierror.Append(instanceErr, err)
		}

		if instance.Databases, err = b.fetcher.GetMicrosoftSQLServerInstanceDatabase(out.ConnString); err != nil {
			instanceErr = multierror.Append(instanceErr, err)
		}

		if len(instance.Databases) > 0 {
			instance.CollationName = instance.Databases[0].CollationName
		}

		dbsMap := make(map[string]*model.MicrosoftSQLServerDatabase)

		for j, db := range instance.Databases {
			instance.Databases[j].Backups = make([]model.MicrosoftSQLServerDatabaseBackup, 0)
			instance.Databases[j].Schemas = make([]model.MicrosoftSQLServerDatabaseSchema, 0)
			instance.Databases[j].Tablespaces = make([]model.MicrosoftSQLServerDatabaseTablespace, 0)
			dbsMap[db.Name] = &instance.Databases[j]
		}

		backupSchedules, err := b.fetcher.GetMicrosoftSQLServerInstanceDatabaseBackups(out.ConnString)
		if err != nil {
			instanceErr = multierror.Append(instanceErr, err)
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

		schemas, err := b.fetcher.GetMicrosoftSQLServerInstanceDatabaseSchemas(out.ConnString)
		if err != nil {
			instanceErr = multierror.Append(instanceErr, err)
		}

		for _, v := range schemas {
			dbsMap[v.DatabaseName].Schemas = make([]model.MicrosoftSQLServerDatabaseSchema, len(v.Data))
			for i, b := range v.Data {
				dbsMap[v.DatabaseName].Schemas[i].AllocatedSpace = int(b.AllocatedMB * 1024 * 1024)
				dbsMap[v.DatabaseName].Schemas[i].UsedSpace = int(b.UsedMB * 1024 * 1024)
				dbsMap[v.DatabaseName].Schemas[i].AllocationType = b.AllocationType
			}
		}

		tablespaces, err := b.fetcher.GetMicrosoftSQLServerInstanceDatabaseTablespaces(out.ConnString)
		if err != nil {
			instanceErr = multierror.Append(instanceErr, err)
		}

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

		if instanceErr != nil {
			merr = multierror.Append(merr, instanceErr)
			continue
		}

		instances = append(instances, instance)
	}

	var connString string
	if len(fetchedInstances) > 0 {
		connString = fetchedInstances[0].ConnString
	}

	return instances, connString, merr
}
