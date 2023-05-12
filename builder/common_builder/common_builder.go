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
	"runtime"
	"strings"

	"github.com/ercole-io/ercole-agent/v2/config"
	"github.com/ercole-io/ercole-agent/v2/fetcher"
	"github.com/ercole-io/ercole-agent/v2/logger"
	"github.com/ercole-io/ercole/v2/model"
	ercutils "github.com/ercole-io/ercole/v2/utils"
)

// CommonBuilder for Linux and Windows hosts
type CommonBuilder struct {
	fetcher       fetcher.Fetcher
	configuration config.Configuration
	log           logger.Logger
}

// NewCommonBuilder initialize an appropriate builder for Linux or Windows
func NewCommonBuilder(configuration config.Configuration, log logger.Logger) CommonBuilder {
	var f fetcher.Fetcher

	log.Debugf("runtime.GOOS: [%v]", runtime.GOOS)

	if runtime.GOOS == "windows" {
		f = fetcher.NewWindowsFetcherImpl(configuration, log)
	} else {
		if runtime.GOOS != "linux" {
			log.Errorf("Unknow runtime.GOOS: [%v], I'll try with linux\n", runtime.GOOS)
		}

		f = fetcher.NewLinuxFetcherImpl(configuration, log)
	}

	builder := CommonBuilder{
		fetcher:       f,
		configuration: configuration,
		log:           log,
	}

	return builder
}

// Run fill hostData
func (b *CommonBuilder) Run(hostData *model.HostData) {
	if host, err := b.fetcher.GetHost(); err != nil {
		b.log.Error(err)
		hostData.AddErrors(err)

		return
	} else {
		hostData.Info = *host
	}

	var err error

	if hostData.Filesystems, err = b.fetcher.GetFilesystems(); err != nil {
		hostData.Filesystems = []model.Filesystem{}

		b.log.Error(err)
		hostData.AddErrors(err)
	}

	hostData.Hostname = hostData.Info.Hostname
	if b.configuration.Hostname != "default" {
		hostData.Hostname = b.configuration.Hostname
	}

	if cms, err := b.fetcher.GetClustersMembershipStatus(); err != nil {
		b.log.Error(err)
		hostData.AddErrors(err)

		return
	} else {
		hostData.ClusterMembershipStatus = *cms
	}

	b.runCloud(hostData)

	b.runOracleDatabase(hostData)

	b.runMicrosoftSQLServer(hostData)

	b.runVirtualization(hostData)

	b.runMySQL(hostData)

	b.runPostgreSQL(hostData)

	b.runMongoDB(hostData)
}

func (b *CommonBuilder) runCloud(hostData *model.HostData) {
	hostData.Cloud = model.Cloud{
		Membership: model.CloudMembershipUnknown,
	}

	membership, err := b.fetcher.GetCloudMembership()
	if err != nil {
		b.log.Error(err)
		hostData.AddErrors(err)

		return
	}

	hostData.Cloud.Membership = membership
}

func (b *CommonBuilder) runOracleDatabase(hostData *model.HostData) {
	if !b.configuration.Features.OracleDatabase.Enabled {
		return
	}

	b.log.Debugf("Oracle/Database mode enabled (user='%s')", b.configuration.Features.OracleDatabase.FetcherUser)

	if err := b.setOrResetFetcherUser(b.configuration.Features.OracleDatabase.FetcherUser); err != nil {
		b.log.Error(err)
		hostData.AddErrors(err)

		return
	}

	var err error

	database, err := b.getOracleDatabaseFeature(hostData.Info, CoreFactor(*hostData))
	if err != nil {
		hostData.AddErrors(err)
	}

	hostData.Features.Oracle = &model.OracleFeature{Database: database}
}

func (b *CommonBuilder) RunExadata(exadata *model.OracleExadataInstance) {
	if !b.configuration.Features.OracleExadata.Enabled {
		return
	}

	b.log.Debugf("Oracle/Exadata mode enabled (user='%s')", b.configuration.Features.OracleExadata.FetcherUser)

	if err := b.setOrResetFetcherUser(b.configuration.Features.OracleExadata.FetcherUser); err != nil {
		b.log.Error(err)
		return
	}

	if err := b.checksToRunExadata(); err != nil {
		b.log.Error(err)

		return
	}

	var err error

	if exadata.Components, err = b.getOracleExadataComponents(); err != nil {
		b.log.Error(err)
		return
	}

	if len(exadata.Components) > 0 {
		exadata.RackID = exadata.Components[0].RackID
	}
}

func (b *CommonBuilder) runMicrosoftSQLServer(hostData *model.HostData) {
	if !b.configuration.Features.MicrosoftSQLServer.Enabled {
		return
	}

	b.log.Debugf("Microsoft/SQLServer mode enabled (user='%s')", b.configuration.Features.MicrosoftSQLServer.FetcherUser)

	if err := b.setOrResetFetcherUser(b.configuration.Features.MicrosoftSQLServer.FetcherUser); err != nil {
		b.log.Error(err)
		hostData.AddErrors(err)

		return
	}

	lazyInitMicrosoftFeature(&hostData.Features)

	var err error

	if hostData.Features.Microsoft.SQLServer, err = b.getMicrosoftSQLServerFeature(); err != nil {
		hostData.AddErrors(err)
	}
}

func (b *CommonBuilder) runVirtualization(hostData *model.HostData) {
	if !b.configuration.Features.Virtualization.Enabled {
		return
	}

	b.log.Debugf("Virtualization mode enabled (user='%s')", b.configuration.Features.Virtualization.FetcherUser)

	if err := b.setOrResetFetcherUser(b.configuration.Features.Virtualization.FetcherUser); err != nil {
		b.log.Error(err)
		hostData.AddErrors(err)

		return
	}

	var err error
	if hostData.Clusters, err = b.getClustersInfos(); err != nil {
		b.log.Error(err)
		hostData.AddErrors(err)
	}
}

func (b *CommonBuilder) runMySQL(hostData *model.HostData) {
	if !b.configuration.Features.MySQL.Enabled {
		return
	}

	b.log.Debugf("MySQL mode enabled")

	if err := b.setOrResetFetcherUser(b.configuration.Features.MySQL.FetcherUser); err != nil {
		b.log.Error(err)
		hostData.AddErrors(err)

		return
	}

	var err error
	if hostData.Features.MySQL, err = b.getMySQLFeature(); err != nil {
		b.log.Error(err)
		hostData.AddErrors(err)
	}
}

func (b *CommonBuilder) runPostgreSQL(hostData *model.HostData) {
	if !b.configuration.Features.PostgreSQL.Enabled {
		return
	}

	b.log.Debugf("PostgrSQL mode enabled (user='%s')", b.configuration.Features.PostgreSQL.FetcherUser)

	if err := b.setOrResetFetcherUser(b.configuration.Features.PostgreSQL.FetcherUser); err != nil {
		b.log.Error(err)
		hostData.AddErrors(err)

		return
	}

	lazyInitPostgreSQLFeature(&hostData.Features)

	var err error

	if hostData.Features.PostgreSQL, err = b.getPostgreSQLFeature(hostData.Hostname); err != nil {
		hostData.AddErrors(err)
	}
}

func (b *CommonBuilder) runMongoDB(hostData *model.HostData) {
	if !b.configuration.Features.MongoDB.Enabled {
		return
	}

	b.log.Debugf("MongoDB mode enabled")

	if err := b.setOrResetFetcherUser(b.configuration.Features.MongoDB.FetcherUser); err != nil {
		b.log.Error(err)
		hostData.AddErrors(err)

		return
	}

	lazyInitMongoDBFeature(&hostData.Features)

	var err error
	if hostData.Features.MongoDB, err = b.getMongoDBFeature(hostData.Hostname); err != nil {
		b.log.Error(err)
		hostData.AddErrors(err)
	}
}

func (b *CommonBuilder) setOrResetFetcherUser(user string) error {
	user = strings.TrimSpace(user)

	if runtime.GOOS != "linux" {
		if user == "" {
			return nil
		}

		err := ercutils.NewErrorf("Can't set user [%s] for fetcher because it is not supported", user)

		b.log.Error(err)

		return err
	}

	if user == "" {
		if err := b.fetcher.SetUserAsCurrent(); err != nil {
			err := ercutils.NewErrorf("Can't set current user for fetchers, err: [%v]", err)
			b.log.Error(err)

			return err
		}

		return nil
	}

	if err := b.fetcher.SetUser(user); err != nil {
		err = ercutils.NewErrorf("Can't set user [%s] for fetchers, err: [%v]", user, err)

		b.log.Error(err)

		return err
	}

	return nil
}

func lazyInitMicrosoftFeature(fs *model.Features) {
	if fs.Microsoft == nil {
		fs.Microsoft = new(model.MicrosoftFeature)
	}
}

func lazyInitPostgreSQLFeature(fs *model.Features) {
	if fs.PostgreSQL == nil {
		fs.PostgreSQL = new(model.PostgreSQLFeature)
	}
}

func lazyInitMongoDBFeature(fs *model.Features) {
	if fs.MongoDB == nil {
		fs.MongoDB = new(model.MongoDBFeature)
	}
}

func CoreFactor(v model.HostData) float64 {
	if v.Cloud.Membership == model.CloudMembershipAws {
		return 1
	}

	return 0.5
}
