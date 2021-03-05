// Copyright (c) 2019 Sorint.lab S.p.A.
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
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

// Configuration holds the agent configuration options
type Configuration struct {
	Hostname               string
	Environment            string
	Location               string
	DataserviceURL         string
	AgentUser              string
	AgentPassword          string
	EnableServerValidation bool
	ForcePwshVersion       string
	Period                 int
	Verbose                bool
	ParallelizeRequests    bool
	Features               Features
}

// Features holds features params
type Features struct {
	OracleDatabase     OracleDatabaseFeature
	Virtualization     VirtualizationFeature
	OracleExadata      OracleExadataFeature
	MicrosoftSQLServer MicrosoftSQLServerFeature
	MySQL              MySQLFeature
}

// OracleDatabaseFeature holds oracle database feature params
type OracleDatabaseFeature struct {
	Enabled     bool
	FetcherUser string
	Oratab      string
	AWR         int
	Forcestats  bool
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
	Enabled   bool
	Instances []MySQLInstanceConnection
}

type MySQLInstanceConnection struct {
	Host     string
	User     string
	Password string
}

// ReadConfig reads the configuration file from the current dir
// or /opt/ercole-agent
func ReadConfig() Configuration {
	baseDir := GetBaseDir()
	configFile := ""
	if runtime.GOOS == "windows" {
		configFile = baseDir + "\\config.json"
	} else {
		configFile = baseDir + "/config.json"
	}
	ex := exists(configFile)
	if !ex {
		if runtime.GOOS == "windows" {
			configFile = "C:\\ErcoleAgent\\config.json"
		} else {
			configFile = "/opt/ercole-agent/config.json"
		}
	}
	raw, err := ioutil.ReadFile(configFile)

	if err != nil {
		log.Fatal("Unable to read configuration file", err)
	}

	var conf Configuration
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()

	err = decoder.Decode(&conf)
	if err != nil {
		log.Fatal("Unable to parse configuration file", err)
	}

	if conf.Features.OracleDatabase.Oratab == "" {
		conf.Features.OracleDatabase.Oratab = "/etc/oratab"
	}

	return conf
}

func exists(name string) bool {
	_, err := os.Stat(name)

	return err == nil
}

// GetBaseDir return executable base directory, os independant
func GetBaseDir() string {
	var s string
	if runtime.GOOS == "windows" {
		s, _ = os.Executable()
		s = filepath.Dir(s)
	} else {
		s, _ = os.Readlink("/proc/self/exe")
		s = filepath.Dir(s)
	}

	return s
}
