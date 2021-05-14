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
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ercole-io/ercole-agent/v2/builder"
	"github.com/ercole-io/ercole-agent/v2/config"
	"github.com/ercole-io/ercole-agent/v2/logger"
	"github.com/ercole-io/ercole-agent/v2/scheduler"
	"github.com/ercole-io/ercole-agent/v2/scheduler/storage"
	"github.com/ercole-io/ercole/v2/model"
	ercole_utils "github.com/ercole-io/ercole/v2/utils"
	"golang.org/x/net/context"
)

var version = "latest"
var hostDataSchemaVersion = 1

type program struct {
	log logger.Logger
}

func (p *program) run() {
	confLog, err := logger.NewLogger("CONFIG")
	if err != nil {
		log.Fatal("Can't initialize CONFIG logger: ", err)
	}
	configuration := config.ReadConfig(confLog)

	opts := make([]logger.LoggerOption, 0)
	if configuration.Verbose {
		opts = append(opts, logger.LogLevel(logger.DebugLevel))
	}
	if len(configuration.LogDirectory) > 0 {
		opts = append(opts, logger.LogDirectory(configuration.LogDirectory))
	}

	p.log, err = logger.NewLogger("AGENT", opts...)
	if err != nil {
		log.Fatal("Can't initialize AGENT logger: ", err)
	}

	doBuildAndSend(configuration, p.log)

	memStorage := storage.NewMemoryStorage()
	scheduler := scheduler.New(memStorage)

	_, err = scheduler.RunEvery(time.Duration(configuration.Period)*time.Hour, func() {
		doBuildAndSend(configuration, p.log)
	})
	if err != nil {
		p.log.Fatal("Error scheduling Ercole agent", err)
	}

	if err := scheduler.Start(); err != nil {
		p.log.Fatal("Error starting Ercole agent scheduler", err)
	}

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
	log.Debugf("Hostdata: %v", string(dataBytes))

	if configuration.Verbose {
		writeHostDataOnTmpFile(data, log)
	}

	client := &http.Client{}
	if !configuration.EnableServerValidation {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "POST", configuration.DataserviceURL+"/hosts", bytes.NewReader(dataBytes))
	if err != nil {
		log.Error("Error creating request: ", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(configuration.AgentUser, configuration.AgentPassword)
	resp, err := client.Do(req)
	if err != nil {
		log.Error("Error sending data: ", err)
		log.Warn("Sending result: FAILED")
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode < 300 {
		log.Info("Response status: ", resp.Status)
		log.Info("Sending result: SUCCESS")
	} else {
		log.Warnf("Response status: %s", resp.Status)
		logResponseBody(log, resp.Body)
		log.Warn("Sending result: FAILED")
	}
}

func logResponseBody(log logger.Logger, body io.ReadCloser) {
	bytes, err := io.ReadAll(body)
	if err != nil {
		return
	}

	var errFE ercole_utils.ErrorResponseFE
	err = json.Unmarshal(bytes, &errFE)
	if err != nil {
		return
	}

	log.Warnf("%s\n%s\n", errFE.Message, errFE.Error)
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
	serve(new(program))
}
