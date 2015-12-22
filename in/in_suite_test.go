package in_test

import (
	. "github.com/concourse/tracker-resource/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/concourse/tracker-resource/Godeps/_workspace/src/github.com/onsi/gomega"

	"testing"
)

func TestIn(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "In Suite")
}
