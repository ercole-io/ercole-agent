module github.com/ercole-io/ercole-agent

go 1.14

require (
	github.com/ercole-io/ercole v0.0.0-20200721075312-7a9582b28223
	github.com/kardianos/service v1.1.0
	github.com/plandem/xlsx v1.0.4 // indirect
	github.com/sirupsen/logrus v1.6.0
	github.com/stretchr/testify v1.6.1
)

replace gopkg.in/robfig/cron.v3 => github.com/robfig/cron/v3 v3.0.1
