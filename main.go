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

package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/ercole-io/ercole-agent/builder"
	"github.com/ercole-io/ercole-agent/config"
	"github.com/ercole-io/ercole-agent/logger"
	"github.com/ercole-io/ercole-agent/scheduler"
	"github.com/ercole-io/ercole-agent/scheduler/storage"
	"github.com/ercole-io/ercole/model"
)

var version = "latest"
var hostDataSchemaVersion = 1

type program struct {
	log logger.Logger
}

func (p *program) run() {
	configuration := config.ReadConfig()

	if configuration.Verbose == true {
		p.log.SetLevel(logger.DebugLevel)
	}

	doBuildAndSend(configuration, p.log)

	memStorage := storage.NewMemoryStorage()
	scheduler := scheduler.New(memStorage)

	_, err := scheduler.RunEvery(time.Duration(configuration.Period)*time.Hour, doBuildAndSend, configuration, p.log)

	if err != nil {
		p.log.Fatal("Error sending data", err)
	}

	scheduler.Start()
	scheduler.Wait()
}

func doBuildAndSend(configuration config.Configuration, log logger.Logger) {
	hostData := builder.BuildData(configuration, log)

	hostData.AgentVersion = version
	hostData.SchemaVersion = hostDataSchemaVersion
	hostData.Tags = []string{}

	sendData(hostData, configuration, log)
}

func sendData(data *model.HostData, configuration config.Configuration, log logger.Logger) {
	log.Info("Sending data...")

	dataBytes, _ := json.Marshal(data)
	log.Infof("Hostdata: %v", string(dataBytes))

	if configuration.Verbose {
		writeHostDataOnTmpFile(data, log)
	}

	client := &http.Client{}

	//Disable certificate validation if enableServerValidation is false
	if configuration.EnableServerValidation == false {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	req, err := http.NewRequest("POST", configuration.DataserviceURL, bytes.NewReader(dataBytes))
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(configuration.AgentUser, configuration.AgentPassword)
	resp, err := client.Do(req)

	sendResult := "FAILED"

	if err != nil {
		log.Error("Error sending data", err)
	} else {
		log.Info("Response status:", resp.Status)
		if resp.StatusCode == 200 {
			sendResult = "SUCCESS"
		}
		defer resp.Body.Close()
	}

	log.Info("Sending result:", sendResult)
}

func writeHostDataOnTmpFile(data *model.HostData, log logger.Logger) {
	dataBytes, _ := json.MarshalIndent(data, "", "    ")

	filePath := fmt.Sprintf("%s/ercole-agent-hostdata-%s.json", os.TempDir(), time.Now().Local().Format("06-01-02-15:04:05"))

	tmpFile, err := os.Create(filePath)
	if err != nil {
		log.Debugf("Can't create hostdata file: %v", os.TempDir()+filePath)
		return
	}

	defer tmpFile.Close()

	if _, err := tmpFile.Write(dataBytes); err != nil {
		log.Debugf("Can't write hostdata in file: %v", os.TempDir()+filePath)
		return
	}

	log.Debugf("Hostdata pretty-printed on file: %v", filePath)
}

func main() {
	log := logger.NewLogger("AGENT")
	prg := &program{log}

	serve(prg)
}
