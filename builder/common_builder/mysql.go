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
	"strings"

	"github.com/ercole-io/ercole/v2/model"
	ercutils "github.com/ercole-io/ercole/v2/utils"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-version"
)

func (b *CommonBuilder) getMySQLFeature() (*model.MySQLFeature, error) {
	var merr error

	var instance *model.MySQLInstance

	var errInstance error

	mysql := &model.MySQLFeature{
		Instances: []model.MySQLInstance{},
	}

	for i, conf := range b.configuration.Features.MySQL.Instances {
		version, err := b.fetcher.GetMySQLVersion(conf)
		if err != nil {
			b.log.Errorf("Can't get MySQL version: %s", conf.Host)

			merr = multierror.Append(merr, ercutils.NewError(err))

			continue
		}

		isold := b.isOldVersion(version)

		if isold {
			instance, errInstance = b.fetcher.GetMySQLOldInstance(conf)
			if errInstance != nil {
				b.log.Errorf("Can't get MySQL old instance: %s", conf.Host)

				merr = multierror.Append(merr, ercutils.NewError(err))

				continue
			}

			if strings.Contains(strings.ToLower(version), "enterprise") {
				instance.Edition = model.MySQLEditionEnterprise
			}
		} else {
			instance, errInstance = b.fetcher.GetMySQLInstance(conf)
			if errInstance != nil {
				b.log.Errorf("Can't get MySQL instance: %s", conf.Host)

				merr = multierror.Append(merr, ercutils.NewError(err))

				continue
			}
		}

		if instance.HighAvailability, err = b.fetcher.GetMySQLHighAvailability(conf); err != nil {
			b.log.Errorf("Can't get MySQL HighAvailability: %s", conf.Host)

			merr = multierror.Append(merr, ercutils.NewError(err))

			continue
		}

		if instance.UUID, err = b.fetcher.GetMySQLUUID(); err != nil {
			b.log.Warnf("Can't get MySQL UUID: %s", conf.Host)
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

		if len(b.configuration.Features.MySQL.Instances) > 1 {
			instance.Name = b.configuration.Features.MySQL.Instances[i].Host + ":" + b.configuration.Features.MySQL.Instances[i].Port
		}

		mysql.Instances = append(mysql.Instances, *instance)
	}

	return mysql, merr
}

func (b *CommonBuilder) isOldVersion(dbversion string) bool {
	re := regexp.MustCompile(`\b\d+(\.\d+)*\b`)
	matches := re.FindStringSubmatch(dbversion)

	if len(matches) == 0 {
		b.log.Error("cannot find matches in db version")
		return false
	}

	v, err := version.NewVersion(matches[0])
	if err != nil {
		b.log.Error(err)
		return false
	}

	constraints, err := version.NewConstraint("< 5.7")
	if err != nil {
		b.log.Error(err)
		return false
	}

	return constraints.Check(v)
}
