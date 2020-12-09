module github.com/ercole-io/ercole-agent

go 1.14

require (
	github.com/ercole-io/ercole v0.0.0-20201209152328-aeb7f1615f3a
	github.com/google/go-cmp v0.5.4 // indirect
	github.com/hashicorp/go-version v1.2.1
	github.com/kardianos/service v1.2.0
	github.com/kr/text v0.2.0 // indirect
	github.com/sirupsen/logrus v1.7.0
	github.com/stretchr/testify v1.6.1
	github.com/tidwall/pretty v1.0.2 // indirect
	go.mongodb.org/mongo-driver v1.4.4 // indirect
	golang.org/x/sys v0.0.0-20201202213521-69691e467435 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
)

replace gopkg.in/robfig/cron.v3 => github.com/robfig/cron/v3 v3.0.1

// replace github.com/ercole-io/ercole => ../ercole
