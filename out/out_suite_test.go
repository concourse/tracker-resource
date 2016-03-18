package out_test

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	. "github.com/concourse/tracker-resource/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/concourse/tracker-resource/Godeps/_workspace/src/github.com/onsi/gomega"
	"github.com/concourse/tracker-resource/Godeps/_workspace/src/github.com/onsi/gomega/gexec"
)

var (
	outPath string
)

var _ = BeforeSuite(func() {
	var err error

	outPath, err = gexec.Build("github.com/concourse/tracker-resource/out/cmd/out")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})

func TestOut(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Out Suite")
}

func Fixture(filename string) string {
	path := filepath.Join("fixtures", filename)
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return string(contents)
}
