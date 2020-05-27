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
