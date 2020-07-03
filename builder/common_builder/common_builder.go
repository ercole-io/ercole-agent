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

	"github.com/ercole-io/ercole-agent/config"
	"github.com/ercole-io/ercole-agent/fetcher"
	"github.com/ercole-io/ercole-agent/logger"
	"github.com/ercole-io/ercole-agent/utils"
	"github.com/ercole-io/ercole/model"
)

// CommonBuilder for Linux and Windows hosts
type CommonBuilder struct {
	fetcher       fetcher.Fetcher
	configuration config.Configuration
	log           logger.Logger
}

// NewCommonBuilder initialize an appropriate builder for Linux or Windows
func NewCommonBuilder(configuration config.Configuration, log logger.Logger) CommonBuilder {
	var specializedFetcher fetcher.SpecializedFetcher

	log.Debugf("runtime.GOOS: [%v]", runtime.GOOS)

	if runtime.GOOS == "windows" {
		wf := fetcher.NewWindowsFetcherImpl(configuration, log)
		specializedFetcher = &wf
	} else {
		if runtime.GOOS != "linux" {
			log.Errorf("Unknow runtime.GOOS: [%v], I'll try with linux\n", runtime.GOOS)
		}

		lf := fetcher.NewLinuxFetcherImpl(configuration, log)
		specializedFetcher = &lf
	}

	fetcherImpl := &fetcher.CommonFetcherImpl{
		SpecializedFetcher: specializedFetcher,
		Configuration:      configuration,
		Log:                log,
	}

	builder := CommonBuilder{
		fetcher:       fetcherImpl,
		configuration: configuration,
		log:           log,
	}

	return builder
}

// Run fill hostData
func (b *CommonBuilder) Run(hostData *model.HostData) {
	if b.configuration.Features.Exadata.Enabled {
		b.log.Debug("Exadata mode enabled")

		b.checksToRunExadata()

		b.setExadataFetchersUser()
	}

	hostData.Info = b.fetcher.GetHost()

	hostData.Hostname = hostData.Info.Hostname
	if b.configuration.Hostname != "default" {
		hostData.Hostname = b.configuration.Hostname
	}

	hostData.ClusterMembershipStatus = b.fetcher.GetClustersMembershipStatus()

	if b.configuration.Features.Databases.Enabled {
		b.log.Debug("Databases mode enabled")
		hostData.Filesystems = b.fetcher.GetFilesystems()

		hostData.Features.Oracle = new(model.OracleFeature)
		hostData.Features.Oracle.Database = new(model.OracleDatabaseFeature)
		hostData.Features.Oracle.Database.Databases = b.getOracleDBs(hostData.Info.HardwareAbstractionTechnology,
			hostData.Info.CPUCores, hostData.Info.CPUSockets)
	}

	if b.configuration.Features.Virtualization.Enabled {
		b.log.Debug("Virtualization mode enabled")
		hostData.Clusters = b.getClustersInfos()
	}

	if b.configuration.Features.Exadata.Enabled {
		b.log.Debug("Exadata mode enabled")

		if err := b.fetcher.SetUserAsCurrent(); err != nil {
			b.log.Panicf("Can't set current user for fetcher, err: [%v]", err)
		}

		if hostData.Features.Oracle == nil {
			hostData.Features.Oracle = new(model.OracleFeature)
		}

		hostData.Features.Oracle.Exadata = new(model.OracleExadataFeature)
		hostData.Features.Oracle.Exadata.Components = b.getOracleExadataComponents()
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

func (b *CommonBuilder) setExadataFetchersUser() {
	if strings.TrimSpace(b.configuration.Features.Exadata.FetchersUser) == "" {
		b.log.Warn("You didn't set FetchersUser, but you have exadata mode enabled, using current user")
		return
	}

	if err := b.fetcher.SetUser(b.configuration.Features.Exadata.FetchersUser); err != nil {
		b.log.Panicf("Can't set user [%s] for fetcher, err: [%v]", b.configuration.Features.Exadata.FetchersUser, err)
	}
}
