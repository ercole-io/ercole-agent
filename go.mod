module github.com/ercole-io/ercole-agent

go 1.14

require (
	github.com/ercole-io/ercole v0.0.0-20200916082827-baa822e04562
	github.com/hashicorp/go-version v1.2.1
	github.com/kardianos/service v1.1.0
	github.com/sirupsen/logrus v1.6.0
	github.com/stretchr/testify v1.6.1
)

replace gopkg.in/robfig/cron.v3 => github.com/robfig/cron/v3 v3.0.1
