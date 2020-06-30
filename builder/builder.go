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

package builder

import (
	common "github.com/ercole-io/ercole-agent/builder/common_builder"
	"github.com/ercole-io/ercole-agent/config"
	"github.com/ercole-io/ercole-agent/logger"
	"github.com/ercole-io/ercole/model"
)

// BuildData will build HostData
func BuildData(configuration config.Configuration, log logger.Logger) *model.HostData {
	hostData := new(model.HostData)

	hostData.Location = configuration.Location
	hostData.Environment = configuration.Envtype

	builder := common.NewCommonBuilder(configuration, log)

	builder.Run(hostData)

	return hostData
}
