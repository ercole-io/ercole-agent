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
	"runtime"
	"strings"

	"github.com/ercole-io/ercole-agent/v2/config"
	"github.com/ercole-io/ercole-agent/v2/fetcher"
	"github.com/ercole-io/ercole-agent/v2/logger"
	"github.com/ercole-io/ercole-agent/v2/utils"
	"github.com/ercole-io/ercole/v2/model"
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
	// build data about host info
	hostData.Info = b.fetcher.GetHost()
	hostData.Filesystems = b.fetcher.GetFilesystems()
	hostData.Hostname = hostData.Info.Hostname
	if b.configuration.Hostname != "default" {
		hostData.Hostname = b.configuration.Hostname
	}
	hostData.ClusterMembershipStatus = b.fetcher.GetClustersMembershipStatus()

	// build data about Oracle/Database
	if b.configuration.Features.OracleDatabase.Enabled {
		b.log.Debugf("Oracle/Database mode enabled (user='%s')", b.configuration.Features.OracleDatabase.FetcherUser)
		b.setOrResetFetcherUser(b.configuration.Features.OracleDatabase.FetcherUser)

		lazyInitOracleFeature(&hostData.Features)
		hostData.Features.Oracle.Database = b.getOracleDatabaseFeature(hostData.Info)
	}

	// build data about Oracle/Exadata
	if b.configuration.Features.OracleExadata.Enabled {
		b.log.Debugf("Oracle/Exadata mode enabled (user='%s')", b.configuration.Features.OracleExadata.FetcherUser)
		b.setOrResetFetcherUser(b.configuration.Features.OracleExadata.FetcherUser)
		b.checksToRunExadata()

		lazyInitOracleFeature(&hostData.Features)
		hostData.Features.Oracle.Exadata = b.getOracleExadataFeature()
	}

	// build data about Oracle/Database
	if b.configuration.Features.MicrosoftSQLServer.Enabled {
		b.log.Debugf("Microsoft/SQLServer mode enabled (user='%s')", b.configuration.Features.MicrosoftSQLServer.FetcherUser)
		b.setOrResetFetcherUser(b.configuration.Features.MicrosoftSQLServer.FetcherUser)

		lazyInitMicrosoftFeature(&hostData.Features)
		hostData.Features.Microsoft.SQLServer = b.getMicrosoftSQLServerFeature()
	}

	// build data about Virtualization
	if b.configuration.Features.Virtualization.Enabled {
		b.log.Debugf("Virtualization mode enabled (user='%s')", b.configuration.Features.Virtualization.FetcherUser)
		b.setOrResetFetcherUser(b.configuration.Features.Virtualization.FetcherUser)

		hostData.Clusters = b.getClustersInfos()
	}
}

func (b *CommonBuilder) checksToRunExadata() {
	if runtime.GOOS != "linux" {
		b.log.Panicf("Can't run exadata mode if os is different from linux, current os: [%v]", runtime.GOOS)
	}

	if !utils.IsRunnigAsRootInLinux() {
		b.log.Panicf("You must be root to run in exadata mode")
	}
}

func (b *CommonBuilder) setOrResetFetcherUser(user string) {
	if runtime.GOOS != "linux" {
		if strings.TrimSpace(user) != "" {
			b.log.Errorf("Can't set user [%s] for fetcher because it is not supported")
		}

		return
	}

	if strings.TrimSpace(user) == "" {
		if err := b.fetcher.SetUserAsCurrent(); err != nil {
			b.log.Panicf("Can't set current user for fetchers, err: [%v]", user, err)
		}
	} else {
		if err := b.fetcher.SetUser(user); err != nil {
			b.log.Panicf("Can't set user [%s] for fetchers, err: [%v]", user, err)
		}
	}
}

func lazyInitOracleFeature(fs *model.Features) {
	if fs.Oracle == nil {
		fs.Oracle = new(model.OracleFeature)
	}
}

func lazyInitMicrosoftFeature(fs *model.Features) {
	if fs.Microsoft == nil {
		fs.Microsoft = new(model.MicrosoftFeature)
	}
}
