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

	"github.com/ercole-io/ercole-agent/v2/agentmodel"
	"github.com/ercole-io/ercole-agent/v2/utils"
	"github.com/ercole-io/ercole/v2/model"
)

func (b *CommonBuilder) checksToRunExadata() {
	if runtime.GOOS != "linux" {
		b.log.Panicf("Can't run exadata mode if os is different from linux, current os: [%v]", runtime.GOOS)
	}

	if !utils.IsRunnigAsRootInLinux() {
		b.log.Panicf("You must be root to run in exadata mode")
	}
}

func (b *CommonBuilder) getOracleExadataFeature() (*model.OracleExadataFeature, []error) {
	oracleExadataFeature := new(model.OracleExadataFeature)
	var errs []error
	oracleExadataFeature.Components, errs = b.getOracleExadataComponents()
	if errs != nil {
		return nil, errs
	}

	return oracleExadataFeature, nil
}

func (b *CommonBuilder) getOracleExadataComponents() ([]model.OracleExadataComponent, []error) {
	exadataDevices, errs := b.fetcher.GetOracleExadataComponents()
	if errs != nil {
		return nil, errs
	}

	exadataCellDisks, errs := b.fetcher.GetOracleExadataCellDisks()
	if errs != nil {
		return nil, errs
	}

	for i := range exadataDevices {
		cellDisks := exadataCellDisks[agentmodel.StorageServerName(exadataDevices[i].Hostname)]
		exadataDevices[i].CellDisks = &cellDisks
	}

	return exadataDevices, nil
}
