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
	"github.com/ercole-io/ercole-agent/config"
	"github.com/ercole-io/ercole-agent/model"
	"github.com/sirupsen/logrus"
)

// BuildData will build HostData
func BuildData(configuration config.Configuration, version string, hostDataSchemaVersion int, log *logrus.Logger) *model.HostData {
	hostData := new(model.HostData)

	hostData.Environment = configuration.Envtype
	hostData.Location = configuration.Location
	hostData.HostType = configuration.HostType
	hostData.Version = version
	hostData.HostDataSchemaVersion = hostDataSchemaVersion

	builder := NewCommonBuilder(configuration, log)

	builder.Run(hostData)

	return hostData
}
