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
	"github.com/ercole-io/ercole/v2/model"
	ercutils "github.com/ercole-io/ercole/v2/utils"
	"github.com/hashicorp/go-multierror"
)

func (b *CommonBuilder) getMySQLFeature() (*model.MySQLFeature, error) {
	var merr error

	mysql := &model.MySQLFeature{
		Instances: []model.MySQLInstance{},
	}

	for _, conf := range b.configuration.Features.MySQL.Instances {
		instance, err := b.fetcher.GetMySQLInstance(conf)
		if err != nil {
			b.log.Errorf("Can't get MySQL instance: %s", conf.Host)

			merr = multierror.Append(merr, ercutils.NewError(err))

			continue
		}

		if instance.HighAvailability, err = b.fetcher.GetMySQLHighAvailability(conf); err != nil {
			b.log.Errorf("Can't get MySQL HighAvailability: %s", conf.Host)

			merr = multierror.Append(merr, ercutils.NewError(err))

			continue
		}

		if instance.UUID, err = b.fetcher.GetMySQLUUID(); err != nil {
			b.log.Errorf("Can't get MySQL UUID: %s", conf.Host)

			merr = multierror.Append(merr, ercutils.NewError(err))

			continue
		}

		if instance.IsMaster, instance.SlaveUUIDs, err = b.fetcher.GetMySQLSlaveHosts(conf); err != nil {
			b.log.Errorf("Can't get MySQL slave hosts: %s", conf.Host)

			merr = multierror.Append(merr, ercutils.NewError(err))

			continue
		}

		if instance.IsSlave, instance.MasterUUID, err = b.fetcher.GetMySQLSlaveStatus(conf); err != nil {
			b.log.Errorf("Can't get MySQL slave status: %s", conf.Host)

			merr = multierror.Append(merr, ercutils.NewError(err))

			continue
		}

		if instance.Databases, err = b.fetcher.GetMySQLDatabases(conf); err != nil {
			b.log.Errorf("Can't get MySQL databases: %s", conf.Host)

			merr = multierror.Append(merr, ercutils.NewError(err))

			continue
		}

		if instance.TableSchemas, err = b.fetcher.GetMySQLTableSchemas(conf); err != nil {
			b.log.Errorf("Can't get MySQL table schemas: %s", conf.Host)

			merr = multierror.Append(merr, ercutils.NewError(err))

			continue
		}

		if instance.SegmentAdvisors, err = b.fetcher.GetMySQLSegmentAdvisors(conf); err != nil {
			b.log.Errorf("Can't get MySQL segment advisors: %s", conf.Host)

			merr = multierror.Append(merr, ercutils.NewError(err))

			continue
		}

		mysql.Instances = append(mysql.Instances, *instance)
	}

	return mysql, merr
}
