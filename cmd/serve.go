// Copyright (c) 2022 Sorint.lab S.p.A.
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
package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/ercole-io/ercole-agent/v2/builder"
	"github.com/ercole-io/ercole-agent/v2/client"
	"github.com/ercole-io/ercole-agent/v2/config"
	"github.com/ercole-io/ercole-agent/v2/logger"
	"github.com/ercole-io/ercole-agent/v2/queue"
	ercutils "github.com/ercole-io/ercole/v2/utils"
	"github.com/go-co-op/gocron"
	"github.com/kardianos/service"
	"github.com/shirou/gopsutil/host"
	"github.com/spf13/cobra"
)

var version = "latest"
var hostDataSchemaVersion = 1

const maxSecondsToWait = 1800

type program struct {
	log logger.Logger
}

var serviceLogger service.Logger

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run ercole-agent services",
	Long:  `Run ercole-agent services`,
	Run: func(cmd *cobra.Command, args []string) {
		serve(new(program))
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func (p *program) run() {
	opts := make([]logger.LoggerOption, 0)
	if configuration.Verbose {
		opts = append(opts, logger.LogVerbosely(configuration.Verbose))
	}

	if len(configuration.LogDirectory) > 0 {
		opts = append(opts, logger.LogDirectory(configuration.LogDirectory))
	}

	var err error

	p.log, err = logger.NewLogger("AGENT", opts...)
	if err != nil {
		log.Fatal("Can't initialize AGENT logger: ", err)
	}

	clients := make([]*client.Client, 0, len(configuration.Queue.DataServices))

	// old configuration must be supported
	if configuration.DataserviceURL != "" {
		configuration.Queue.DataServices = append(configuration.Queue.DataServices, config.DataService{
			Url:                    configuration.DataserviceURL,
			AgentUser:              configuration.AgentUser,
			AgentPassword:          configuration.AgentPassword,
			EnableServerValidation: configuration.EnableServerValidation,
		})
	}

	for _, dataservice := range configuration.Queue.DataServices {
		client, err := client.NewClient(
			client.EnableServerValidation(dataservice.EnableServerValidation),
			client.SetAuthentication(dataservice.AgentUser, dataservice.AgentPassword),
			client.SetBaseUrl(dataservice.Url),
		)
		if err != nil {
			log.Fatal("Can't initialize AGENT client: ", err)
			continue
		}

		ping(p.log, client)

		clients = append(clients, client)
	}

	uptime(p.log)

	wk := queue.NewWorker("api-worker")

	waitingTimeDuration, err := time.ParseDuration(fmt.Sprintf("%dm", configuration.Queue.WaitingTime))
	if err != nil {
		log.Fatal(err)
	}

	tk := queue.NewConsumer("host-sender", waitingTimeDuration, configuration.Queue.RetryLimit, sendData)

	scheduler := gocron.NewScheduler(time.UTC)

	_, err = scheduler.Every(int(configuration.Period)).Hour().Do(func() {
		for _, c := range clients {
			doBuildAndSend(p.log, c, configuration, wk, tk)
		}
	})
	if err != nil {
		p.log.Fatal("Error scheduling Ercole agent", err)
	}

	scheduler.StartBlocking()
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
		body, err := io.ReadAll(resp.Body)
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

func doBuildAndSend(log logger.Logger, client *client.Client, configuration config.Configuration, wk queue.Worker, tk queue.Consumer) {
	hostData := builder.BuildData(configuration, log)

	hostData.AgentVersion = version
	hostData.SchemaVersion = hostDataSchemaVersion
	hostData.Period = configuration.Period
	hostData.Tags = []string{}

	tMessageHost := tk.NewMessage(context.Background(), log, client, configuration, hostData, "hosts")

	if err := wk.Add(tMessageHost); err != nil {
		log.Error(err)
	}

	if configuration.Features.OracleExadata.Enabled {
		exadata := builder.BuildExadata(configuration, log)

		tMessageExa := tk.NewMessage(context.Background(), log, client, configuration, exadata, "exadatas")
		if err := wk.Add(tMessageExa); err != nil {
			log.Error(err)
		}
	}
}

func sendData(log logger.Logger, client *client.Client, configuration config.Configuration, data interface{}, endopoint string) error {
	log.Info("Sending data...")

	dataBytes, err := json.Marshal(data)
	if err != nil {
		log.Error(err)
		return err
	}

	log.Debugf("Data: %v", string(dataBytes))

	if configuration.Verbose {
		writeDataOnTmpFile(data, log)
	}

	resp, err := client.DoRequest("POST", fmt.Sprintf("/%s", endopoint), dataBytes)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Errorf("Error sending data: unable to reach ercole server in %d seconds", client.Timeout())
		} else {
			log.Error("Error sending data: ", err)
		}

		log.Warn("Sending result: FAILED")

		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 300 {
		log.Info("Response status: ", resp.Status)
		log.Info("Sending result: SUCCESS")
	} else {
		log.Warnf("Response status: %s", resp.Status)
		logResponseBody(log, resp.Body)
		log.Warn("Sending result: FAILED")

		return errors.New("failed send data to ercole data-service")
	}

	return nil
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

func writeDataOnTmpFile(data interface{}, log logger.Logger) {
	dataBytes, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		log.Error(err)
		return
	}

	filePath := fmt.Sprintf("%s/ercole-agent-data-%s.json", os.TempDir(), time.Now().Local().Format("06-01-02_15-04-05"))

	tmpFile, err := os.Create(filePath)
	if err != nil {
		log.Debugf("Can't create data file: %v", os.TempDir()+filePath)
		return
	}

	defer tmpFile.Close()

	if _, err := tmpFile.Write(dataBytes); err != nil {
		log.Debugf("Can't write data in file: %v", os.TempDir()+filePath)
		return
	}

	log.Debugf("Data pretty-printed on file: %v", filePath)
}

func serve(prg *program) {
	svcConfig := &service.Config{
		Name:        "ErcoleAgent",
		DisplayName: "The Ercole Agent",
		Description: "Asset management agent from the Ercole project.",
	}

	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	serviceLogger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}

	err = s.Run()
	if err != nil {
		if err := serviceLogger.Error(err); err != nil {
			log.Fatal(err)
		}
	}
}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) Stop(s service.Service) error {
	return nil
}
