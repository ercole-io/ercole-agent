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
	b64 "encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/ercole-io/ercole-agent/builder"
	"github.com/ercole-io/ercole-agent/config"
	"github.com/ercole-io/ercole-agent/model"
	"github.com/ercole-io/ercole-agent/scheduler"
	"github.com/ercole-io/ercole-agent/scheduler/storage"

	"github.com/kardianos/service"
)

var logger service.Logger
var version = "latest"
var hostDataSchemaVersion = 4

type program struct{}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) run() {
	configuration := config.ReadConfig()

	doBuildAndSend(configuration)

	memStorage := storage.NewMemoryStorage()
	scheduler := scheduler.New(memStorage)

	_, err := scheduler.RunEvery(time.Duration(configuration.Frequency)*time.Hour, doBuildAndSend, configuration)

	if err != nil {
		log.Fatal("Error sending data", err)
	}

	scheduler.Start()
	scheduler.Wait()
}

func doBuildAndSend(configuration config.Configuration) {
	hostData := builder.BuildData(configuration, version, hostDataSchemaVersion)
	sendData(hostData, configuration)
}

func sendData(data *model.HostData, configuration config.Configuration) {
	log.Println("Sending data...")

	b, _ := json.Marshal(data)
	s := string(b)

	log.Println("Data:", s)

	client := &http.Client{}

	//Disable certificate validation if enableServerValidation is false
	if configuration.EnableServerValidation == false {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	req, err := http.NewRequest("POST", configuration.Serverurl, bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	auth := configuration.Serverusr + ":" + configuration.Serverpsw
	authEnc := b64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Add("Authorization", "Basic "+authEnc)
	resp, err := client.Do(req)

	sendResult := "FAILED"

	if err != nil {
		log.Println("Error sending data", err)
	} else {
		log.Println("Response status:", resp.Status)
		if resp.StatusCode == 200 {
			sendResult = "SUCCESS"
		}
		defer resp.Body.Close()
	}

	log.Println("Sending result:", sendResult)
}

func (p *program) Stop(s service.Service) error {
	return nil
}

func main() {
	svcConfig := &service.Config{
		Name:        "ErcoleAgent",
		DisplayName: "The Ercole Agent",
		Description: "Asset management agent from the Ercole project.",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}
