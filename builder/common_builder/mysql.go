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

import "github.com/ercole-io/ercole/v2/model"

func (b *CommonBuilder) getMySQLFeature() *model.MySQLFeature {
	mysql := &model.MySQLFeature{
		Instances: []model.MySQLInstance{},
	}

	for _, conf := range b.configuration.Features.MySQL.Instances {
		instance := b.fetcher.GetMySQLInstance(conf)

		if instance == nil {
			b.log.Warnf("Can't get instance: %s\n", conf.Host)
			continue
		}

		instance.HighAvailability = b.fetcher.GetMySQLHighAvailability(conf)

		instance.UUID = b.fetcher.GetMySQLUUID()
		instance.IsMaster, instance.SlaveUUIDs = b.fetcher.GetMySQLSlaveHosts(conf)
		instance.IsSlave, instance.MasterUUID = b.fetcher.GetMySQLSlaveStatus(conf)

		instance.Databases = b.fetcher.GetMySQLDatabases(conf)
		instance.TableSchemas = b.fetcher.GetMySQLTableSchemas(conf)
		instance.SegmentAdvisors = b.fetcher.GetMySQLSegmentAdvisors(conf)

		mysql.Instances = append(mysql.Instances, *instance)
	}

	return mysql
}
