module github.com/ercole-io/ercole-agent/v2

go 1.16

require (
	github.com/aws/aws-sdk-go v1.38.67 // indirect
	github.com/ercole-io/ercole/v2 v2.0.0-20211209152515-ad4278654157
	github.com/fatih/color v1.12.0
	github.com/felixge/httpsnoop v1.0.2 // indirect
	github.com/golang/mock v1.5.0
	github.com/hashicorp/go-multierror v1.1.1
	github.com/hashicorp/go-version v1.3.0
	github.com/kardianos/service v1.2.0
	github.com/klauspost/compress v1.13.1 // indirect
	github.com/shirou/gopsutil v3.21.11+incompatible
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/tidwall/pretty v1.0.2 // indirect
	github.com/tklauser/go-sysconf v0.3.9 // indirect
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
	go.mongodb.org/mongo-driver v1.5.3 // indirect
	golang.org/x/crypto v0.0.0-20210616213533-5ff15b29337e // indirect
	golang.org/x/sys v0.0.0-20210816074244-15123e1e1f71
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)

// replace github.com/ercole-io/ercole/v2 => ../ercole-aws
