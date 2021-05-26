module github.com/ercole-io/ercole-agent/v2

go 1.16

require (
	github.com/aws/aws-sdk-go v1.38.43 // indirect
	github.com/ercole-io/ercole/v2 v2.0.0-20210520083642-444c57e56dfa
	github.com/fatih/color v1.12.0
	github.com/golang/mock v1.5.0
	github.com/google/go-cmp v0.5.4 // indirect
	github.com/hashicorp/go-version v1.2.1
	github.com/kardianos/service v1.2.0
	github.com/kr/text v0.2.0 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/tidwall/pretty v1.0.2 // indirect
	golang.org/x/sys v0.0.0-20210525143221-35b2ab0089ea
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
)

// replace gopkg.in/robfig/cron.v3 => github.com/robfig/cron/v3 v3.0.1
// replace github.com/ercole-io/ercole/v2 => ../ercole
