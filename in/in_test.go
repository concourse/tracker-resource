package in_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	. "github.com/concourse/tracker-resource/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/concourse/tracker-resource/Godeps/_workspace/src/github.com/onsi/gomega"

	"github.com/concourse/tracker-resource/Godeps/_workspace/src/github.com/onsi/gomega/gexec"

	"github.com/concourse/tracker-resource"
	"github.com/concourse/tracker-resource/in"
)

var _ = Describe("In", func() {
	var (
		tmpDir   string
		request  in.InRequest
		response in.InResponse
	)

	JustBeforeEach(func() {
		binPath, err := gexec.Build("github.com/concourse/tracker-resource/in/cmd/in")
		Ω(err).ShouldNot(HaveOccurred())

		tmpDir, err = ioutil.TempDir("", "tracker_resource_in")

		stdin := &bytes.Buffer{}
		err = json.NewEncoder(stdin).Encode(request)
		Ω(err).ShouldNot(HaveOccurred())

		cmd := exec.Command(binPath, tmpDir)
		cmd.Stdin = stdin
		cmd.Dir = tmpDir

		session, err := gexec.Start(
			cmd,
			GinkgoWriter,
			GinkgoWriter,
		)
		Ω(err).ShouldNot(HaveOccurred())

		Eventually(session).Should(gexec.Exit(0))

		err = json.Unmarshal(session.Out.Contents(), &response)
		Ω(err).ShouldNot(HaveOccurred())
	})

	AfterEach(func() {
		err := os.RemoveAll(tmpDir)
		Ω(err).ShouldNot(HaveOccurred())
	})

	Context("when a version is given to the executable", func() {
		BeforeEach(func() {
			request = in.InRequest{
				Source: resource.Source{
					Token:     "abc",
					ProjectID: "1234",
				},
				Version: resource.Version{
					Time: time.Now().Add(332 * time.Hour),
				},
			}
		})

		It("outputs that version", func() {
			Ω(response.Version.Time).Should(BeTemporally("~", request.Version.Time, time.Second))
		})
	})

	Context("when a version is not given to the executable", func() {
		BeforeEach(func() {
			request = in.InRequest{
				Source: resource.Source{
					Token:     "abc",
					ProjectID: "1234",
				},
			}
		})

		It("generates a 'fake' current version", func() {
			Ω(response.Version.Time).Should(BeTemporally("~", time.Now(), time.Second))
		})
	})
})
