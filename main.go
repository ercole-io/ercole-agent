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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/shirou/gopsutil/host"

	"github.com/ercole-io/ercole-agent/v2/builder"
	"github.com/ercole-io/ercole-agent/v2/client"
	"github.com/ercole-io/ercole-agent/v2/config"
	"github.com/ercole-io/ercole-agent/v2/logger"
	"github.com/ercole-io/ercole-agent/v2/scheduler"
	"github.com/ercole-io/ercole-agent/v2/scheduler/storage"
	"github.com/ercole-io/ercole/v2/model"
	ercutils "github.com/ercole-io/ercole/v2/utils"
)

var version = "latest"
var hostDataSchemaVersion = 1

const maxSecondsToWait = 1800

type program struct {
	log logger.Logger
}

func (p *program) run() {
	confLog, err := logger.NewLogger("CONFIG", logger.LogLevel(logger.DebugLevel))
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

	client, err := client.NewClient(
		client.EnableServerValidation(configuration.EnableServerValidation),
		client.SetAuthentication(configuration.AgentUser, configuration.AgentPassword),
		client.SetBaseUrl(configuration.DataserviceURL),
	)
	if err != nil {
		log.Fatal("Can't initialize AGENT client: ", err)
	}

	ping(p.log, client)

	uptime(p.log)

	doBuildAndSend(p.log, client, configuration)

	memStorage := storage.NewMemoryStorage()
	scheduler := scheduler.New(memStorage)

	_, err = scheduler.RunEvery(time.Duration(configuration.Period)*time.Hour, func() {
		doBuildAndSend(p.log, client, configuration)
	})
	if err != nil {
		p.log.Fatal("Error scheduling Ercole agent", err)
	}

	if err := scheduler.Start(); err != nil {
		p.log.Fatal("Error starting Ercole agent scheduler", err)
	}

	scheduler.Wait()
}

func uptime(log logger.Logger) {
	log.Debug("Uptime...")

	uptime, err := host.Uptime()
	if err != nil {
		log.Error(err)
		return
	}

	if uptime < maxSecondsToWait {
		secondsToWait := time.Duration(maxSecondsToWait - uptime)
		log.Debugf("Seconds to wait: %d", secondsToWait)
		time.Sleep(secondsToWait * time.Second)
	}

	log.Debug("Uptime OK")
}

func ping(log logger.Logger, client *client.Client) {
	log.Debug("Ping...")

	resp, err := client.DoRequest("GET", "/ping", nil)
	if err != nil {
		log.Warn("Can't ping ercole data-service: " + err.Error())
		time.Sleep(3 * time.Second)

		return
	}

	if resp.StatusCode != 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Error(err)
			return
		}

		defer resp.Body.Close()

		log.Warn("Can't ping ercole data-service: " + resp.Status)
		log.Debug("Responde body: " + string(body))
		time.Sleep(3 * time.Second)

		return
	}

	log.Debug("Ping OK")
}

func doBuildAndSend(log logger.Logger, client *client.Client, configuration config.Configuration) {
	hostData := builder.BuildData(configuration, log)

	hostData.AgentVersion = version
	hostData.SchemaVersion = hostDataSchemaVersion
	hostData.Tags = []string{}

	sendData(log, client, configuration, hostData)
}

func sendData(log logger.Logger, client *client.Client, configuration config.Configuration, data *model.HostData) {
	log.Info("Sending data...")

	dataBytes, err := json.Marshal(data)
	if err != nil {
		log.Error(err)
		return
	}

	log.Debugf("Hostdata: %v", string(dataBytes))

	if configuration.Verbose {
		writeHostDataOnTmpFile(data, log)
	}

	resp, err := client.DoRequest("POST", "/hosts", dataBytes)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Errorf("Error sending data: unable to reach ercole server in %d seconds", client.Timeout())
		} else {
			log.Error("Error sending data: ", err)
		}

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

	var errFE ercutils.ErrorResponseFE

	err = json.Unmarshal(bytes, &errFE)
	if err != nil {
		return
	}

	log.Warnf("%s\n%s\n", errFE.Message, errFE.Error)
}

func writeHostDataOnTmpFile(data *model.HostData, log logger.Logger) {
	dataBytes, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		log.Error(err)
		return
	}

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
