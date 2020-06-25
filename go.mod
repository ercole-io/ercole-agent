module github.com/ercole-io/ercole-agent

go 1.14

require (
	github.com/ercole-io/ercole v0.0.0-20200617121441-a788422d1a00
	github.com/kardianos/service v1.1.0
	github.com/sirupsen/logrus v1.6.0
	github.com/stretchr/testify v1.5.1
)

replace gopkg.in/robfig/cron.v3 => github.com/robfig/cron/v3 v3.0.1

replace github.com/ercole-io/ercole => github.com/amreo/ercole v0.0.0-20200618144554-fd2960c28b08
