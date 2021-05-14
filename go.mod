module github.com/ercole-io/ercole-agent/v2

go 1.16

require (
	github.com/ercole-io/ercole/v2 v2.0.0-20210412140120-4da40b4eb822
	github.com/golang/mock v1.5.0
	github.com/google/go-cmp v0.5.4 // indirect
	github.com/hashicorp/go-version v1.2.1
	github.com/kardianos/service v1.2.0
	github.com/kr/text v0.2.0 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/tidwall/pretty v1.0.2 // indirect
	golang.org/x/net v0.0.0-20200202094626-16171245cfb2
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
)

// replace gopkg.in/robfig/cron.v3 => github.com/robfig/cron/v3 v3.0.1
// replace github.com/pkg/sftp => github.com/amreo/sftp v1.10.2-0.20200107133605-5981645e4b3b
// replace github.com/ercole-io/ercole/v2 => ../ercole
