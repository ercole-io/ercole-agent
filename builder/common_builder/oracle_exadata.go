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
	ercutils "github.com/ercole-io/ercole/v2/utils"
)

func (b *CommonBuilder) checksToRunExadata() error {
	if runtime.GOOS != "linux" {
		err := ercutils.NewErrorf("Can't run exadata mode if os is different from linux, current os: [%v]", runtime.GOOS)

		b.log.Error(err)

		return err
	}

	if !utils.IsRunnigAsRootInLinux() {
		err := ercutils.NewErrorf("You must be root to run in exadata mode")

		b.log.Error(err)

		return err
	}

	return nil
}

func (b *CommonBuilder) getOracleExadataFeature() (*model.OracleExadataFeature, error) {
	oracleExadataFeature := new(model.OracleExadataFeature)

	var err error

	oracleExadataFeature.Components, err = b.getOracleExadataComponents()
	if err != nil {
		b.log.Error(err)
		return nil, err
	}

	return oracleExadataFeature, nil
}

func (b *CommonBuilder) getOracleExadataComponents() ([]model.OracleExadataComponent, error) {
	exadataDevices, err := b.fetcher.GetOracleExadataComponents()
	if err != nil {
		return nil, err
	}

	exadataCellDisks, err := b.fetcher.GetOracleExadataCellDisks()
	if err != nil {
		return nil, err
	}

	for i := range exadataDevices {
		cellDisks := exadataCellDisks[agentmodel.StorageServerName(exadataDevices[i].Hostname)]
		exadataDevices[i].CellDisks = &cellDisks
	}

	return exadataDevices, nil
}
