module github.com/concourse/tracker-resource

go 1.16

require (
	github.com/golang/protobuf v0.0.0-20161117033126-8ee79997227b // indirect
	github.com/mitchellh/colorstring v0.0.0-20150917214807-8631ce90f286
	github.com/onsi/ginkgo v1.2.1-0.20170126062008-bb93381d543b
	github.com/onsi/gomega v0.0.0-20170214000320-c463cd2a8578
	github.com/xoebus/go-tracker v0.0.0-00010101000000-000000000000
	golang.org/x/sync v0.0.0-20201207232520-09787c993a3a // indirect
	golang.org/x/sys v0.0.0-20160717071931-a646d33e2ee3 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v2 v2.0.0-20170208141851-a3f3340b5840 // indirect
)

replace github.com/xoebus/go-tracker => ./go-tracker/
