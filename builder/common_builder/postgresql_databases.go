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

package common

import (
	"regexp"
	"strconv"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-version"
)

func (b *CommonBuilder) isv10(dbversion string) bool {
	re := regexp.MustCompile(`[0-9]+\.[0-9]+`)
	matches := re.FindStringSubmatch(dbversion)

	v10, err := version.NewVersion(matches[0])
	if err != nil {
		b.log.Error(err)
		return false
	}

	constraints, err := version.NewConstraint(">= 9.6")
	if err != nil {
		b.log.Error(err)
		return false
	}

	return constraints.Check(v10)
}

func (b *CommonBuilder) getPostgreSQLFeature() (*model.PostgreSQLFeature, error) {
	var merr error

	feature := model.PostgreSQLFeature{}

	for _, configInstance := range b.configuration.Features.PostgreSQL.Instances {
		setting, err := b.fetcher.GetPostgreSQLSetting(configInstance)
		if err != nil {
			b.log.Error(err)
			merr = multierror.Append(merr, err)
			
			continue
		}

		isv10 := b.isv10(setting.DbVersion)

		instance, err := b.fetcher.GetPostgreSQLInstance(configInstance, isv10)
		if err != nil {
			b.log.Error(err)
			merr = multierror.Append(merr, err)
		}

		port, err := strconv.Atoi(configInstance.Port)
		if err != nil {
			b.log.Error(err)
			merr = multierror.Append(merr, err)
		}

		instance.Port = port

		dbNameList, err := b.fetcher.GetPostgreSQLDbNameList(configInstance)
		if err != nil {
			b.log.Error(err)
			merr = multierror.Append(merr, err)
		}

		for _, dbName := range dbNameList {
			db, err := b.fetcher.GetPostgreSQLDatabase(configInstance, dbName, isv10)
			if err != nil {
				b.log.Error(err)
				merr = multierror.Append(merr, err)

				continue
			}

			schemaNameList, err := b.fetcher.GetPostgreSQLDbSchemaNameList(configInstance, dbName)
			if err != nil {
				b.log.Error(err)
				merr = multierror.Append(merr, err)

				continue
			}

			for _, schemaName := range schemaNameList {
				schema, err := b.fetcher.GetPostgreSQLSchema(configInstance, dbName, schemaName, isv10)
				if err != nil {
					b.log.Error(err)
					merr = multierror.Append(merr, err)

					continue
				}

				db.Schemas = append(db.Schemas, *schema)
			}

			instance.Databases = append(instance.Databases, *db)
		}

		instance.Setting = setting

		feature.Instances = append(feature.Instances, *instance)
	}

	return &feature, merr
}
