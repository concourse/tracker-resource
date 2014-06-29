package out_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/gexec"

	"github.com/concourse/tracker-resource/out"
)

var _ = Describe("In", func() {
	var tmpdir string

	var outCmd *exec.Cmd

	BeforeEach(func() {
		var err error

		tmpdir, err = ioutil.TempDir("", "out-tmp")
		Ω(err).ShouldNot(HaveOccurred())
		err = os.MkdirAll(tmpdir, 0755)
		Ω(err).ShouldNot(HaveOccurred())

		setupTestEnvironment(tmpdir)

		outCmd = exec.Command(outPath, tmpdir)
	})

	AfterEach(func() {
		os.RemoveAll(tmpdir)
	})

	Context("when executed", func() {
		var request out.OutRequest
		var response out.OutResponse

		BeforeEach(func() {
			request = out.OutRequest{}

			response = out.OutResponse{}
		})

		JustBeforeEach(func() {
			stdin, err := outCmd.StdinPipe()
			Ω(err).ShouldNot(HaveOccurred())

			session, err := gexec.Start(outCmd, GinkgoWriter, GinkgoWriter)
			Ω(err).ShouldNot(HaveOccurred())

			err = json.NewEncoder(stdin).Encode(request)
			Ω(err).ShouldNot(HaveOccurred())

			Eventually(session).Should(gexec.Exit(0))

			err = json.Unmarshal(session.Out.Contents(), &response)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("finds finished stories that are mentioned in recent git commits", func() {
		})
	})
})

func setupTestEnvironment(path string) {
	cmd := exec.Command(filepath.Join("scripts/setup.sh"), path)
	cmd.Stdout = GinkgoWriter
	cmd.Stderr = GinkgoWriter

	err := cmd.Run()
	Ω(err).ShouldNot(HaveOccurred())
}
