// Copyright (c) 2023 Sorint.lab S.p.A.
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

package config

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/goraz/onion"
	"github.com/goraz/onion/onionwriter"

	"github.com/ercole-io/ercole-agent/v2/logger"
	"github.com/ercole-io/ercole/v2/model"
)

// Configuration holds the agent configuration options
type Configuration struct {
	Hostname            string
	ExadataName         string
	Environment         string
	Location            string
	Queue               Queue
	ForcePwshVersion    string
	Period              uint
	Verbose             bool
	ParallelizeRequests bool
	LogDirectory        string
	Features            Features
}

// Features holds features params
type Features struct {
	OracleDatabase     OracleDatabaseFeature
	Virtualization     VirtualizationFeature
	OracleExadata      OracleExadataFeature
	MicrosoftSQLServer MicrosoftSQLServerFeature
	MySQL              MySQLFeature
	PostgreSQL         PostgreSQLFeature
	MongoDB            MongoDBFeature
}

// OracleDatabaseFeature holds oracle database feature params
type OracleDatabaseFeature struct {
	Enabled     bool
	FetcherUser string
	Oratab      string
	AWR         int
	Forcestats  bool
	OracleUser  OracleUser
}

func (odf *OracleDatabaseFeature) IsOracleUser() bool {
	return odf.OracleUser != OracleUser{}
}

type OracleUser struct {
	Username string
	Password string
}

// VirtualizationFeature holds virtualization feature params
type VirtualizationFeature struct {
	Enabled     bool
	FetcherUser string
	Hypervisors []Hypervisor
}

// Hypervisor holds the parameters used to connect to an hypervisor
type Hypervisor struct {
	Type       string
	Endpoint   string
	Username   string
	Password   string
	OvmUserKey string
	OvmControl string
}

// OracleExadataFeature holds oracle exadata feature params
type OracleExadataFeature struct {
	Enabled     bool
	FetcherUser string
}

// MicrosoftSQLServerFeature holds microsoft sql server feature params
type MicrosoftSQLServerFeature struct {
	Enabled     bool
	FetcherUser string
}

type MySQLFeature struct {
	Enabled     bool
	FetcherUser string
	Instances   []MySQLInstanceConnection
}

type MySQLInstanceConnection struct {
	Host          string
	Port          string
	User          string
	Password      string
	DataDirectory string
	Socket        string
}

type PostgreSQLFeature struct {
	Enabled     bool
	FetcherUser string
	Instances   []PostgreSQLInstanceConnection
}

type PostgreSQLInstanceConnection struct {
	Port     string
	User     string
	Password string
}

type MongoDBFeature struct {
	Enabled     bool
	FetcherUser string
	Instances   []MongoDBInstanceConnection
}

type MongoDBInstanceConnection struct {
	Host             string
	Port             string
	User             string
	Password         string
	DirectConnection bool
}

type Queue struct {
	DataServices []DataService
	WaitingTime  int
	RetryLimit   int
}

type DataService struct {
	Url                    string
	AgentUser              string
	AgentPassword          string
	EnableServerValidation bool
}

// ReadConfig reads the configuration file from the current dir
// or /opt/ercole-agent
func ReadConfig(log logger.Logger, extraConfigFile string) Configuration {
	baseDir, err := GetBaseDir(log)
	if err != nil {
		log.Fatal("Unable to get base directory: ", err)
	}

	configFile := ""

	layers := make([]onion.Layer, 0)

	if runtime.GOOS == "windows" {
		configFile = baseDir + "\\config.json"
		if !exists(configFile) {
			configFile = "C:\\ErcoleAgent\\config.json"
		}

		layers = addFileLayers(log, layers, configFile)
	} else {
		layers = addFileLayers(log, layers, "/opt/ercole-agent/config.json")
		layers = addFileLayers(log, layers, "/usr/share/ercole-agent/config.json")
		layers = addFileLayers(log, layers, "/etc/ercole-agent/ercole-agent.json")
		layers = addFileLayers(log, layers, "/etc/ercole-agent/conf.d/*.json")
		layers = addFileLayers(log, layers, "./config.json")
		layers = addFileLayers(log, layers, extraConfigFile)
	}

	configOnion := onion.New(layers...)

	var conf Configuration

	err = onionwriter.DecodeOnion(configOnion, &conf)
	if err != nil {
		log.Fatal("something went wrong while reading your configuration files")
	}

	checkConfiguration(log, &conf)

	return conf
}

func exists(name string) bool {
	_, err := os.Stat(name)

	return err == nil
}

func checkConfiguration(log logger.Logger, config *Configuration) {
	checkPeriod(log, config)
	checkLogDirectory(log, config)

	checkFeatureOracleDatabase(log, config)
	checkFeatureVirtualization(log, config)
	checkMySQLFeature(log, config)
}

func checkPeriod(log logger.Logger, config *Configuration) {
	if config.Period == 0 {
		defaultPeriod := uint(24)
		log.Warnf("Period has invalid value [%d], set to default value [%d]", config.Period, defaultPeriod)
		config.Period = defaultPeriod
	}
}

func checkLogDirectory(log logger.Logger, config *Configuration) {
	path := config.LogDirectory
	if path == "" {
		return
	}

	if err := checkDirectoryIsWritable(path); err != nil {
		log.Fatalf("LogDirectory \"%s\" is not valid: %s", path, err)
	}
}

func checkFeatureOracleDatabase(log logger.Logger, config *Configuration) {
	if !config.Features.OracleDatabase.Enabled {
		return
	}

	if runtime.GOOS == "windows" {
		return
	}

	if config.Features.OracleDatabase.Oratab == "" {
		config.Features.OracleDatabase.Oratab = "/etc/oratab"
	}

	_, err := ioutil.ReadFile(config.Features.OracleDatabase.Oratab)
	if err != nil {
		log.Fatalf("Oracle Database: oratab file \"%s\" can't be opened: %s", config.Features.OracleDatabase.Oratab, err)
	}
}

func checkFeatureVirtualization(log logger.Logger, config *Configuration) {
	if config.Features.Virtualization.Hypervisors == nil {
		return
	}

	hypervisorTypes := map[string]string{
		"ovm":    model.TechnologyOracleVM,
		"vmware": model.TechnologyVMWare,
		"olvm":   model.TechnologyOracleLVM,
	}

	for i := range config.Features.Virtualization.Hypervisors {
		hv := &config.Features.Virtualization.Hypervisors[i]

		correctType, found := hypervisorTypes[hv.Type]
		if !found {
			log.Errorf("Hypervisor type not supported: %v", hv.Type)
			log.Errorf("Hypervisor types supported are:")

			for k, v := range hypervisorTypes {
				log.Errorf("\t\"%v\" for %v", k, v)
			}

			log.Fatalf("Fix you configuration file")
		}

		hv.Type = correctType
	}
}

func checkMySQLFeature(log logger.Logger, config *Configuration) {
	for i, instance := range config.Features.MySQL.Instances {
		if instance.DataDirectory == "" {
			config.Features.MySQL.Instances[i].DataDirectory = "/var/lib/mysql"
		}
	}
}

// GetBaseDir return executable base directory, os independant
func GetBaseDir(log logger.Logger) (string, error) {
	var s string

	if runtime.GOOS == "windows" {
		execString, err := os.Executable()
		if err != nil {
			return s, err
		}

		s = filepath.Dir(execString)
	} else {
		execString, err := os.Readlink("/proc/self/exe")
		if err != nil {
			return s, err
		}

		s = filepath.Dir(execString)
	}

	return s, nil
}

func addFileLayers(log logger.Logger, layers []onion.Layer, configFiles ...string) []onion.Layer {
	for _, file := range configFiles {
		layer, err := onion.NewFileLayer(file, nil)

		var pathErr *os.PathError

		if err == nil {
			log.Debugf("Read file for conf: %s", file)

			layers = append(layers, layer)
		} else if !errors.As(err, &pathErr) {
			log.Warnf("error reading file [%s]: [%s]", file, err)
		}
	}

	return layers
}
