package in_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/gexec"

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
		Expect(err).NotTo(HaveOccurred())

		tmpDir, err = ioutil.TempDir("", "tracker_resource_in")

		stdin := &bytes.Buffer{}
		err = json.NewEncoder(stdin).Encode(request)
		Expect(err).NotTo(HaveOccurred())

		cmd := exec.Command(binPath, tmpDir)
		cmd.Stdin = stdin
		cmd.Dir = tmpDir

		session, err := gexec.Start(
			cmd,
			GinkgoWriter,
			GinkgoWriter,
		)
		Expect(err).NotTo(HaveOccurred())

		Eventually(session).Should(gexec.Exit(0))

		err = json.Unmarshal(session.Out.Contents(), &response)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		err := os.RemoveAll(tmpDir)
		Expect(err).NotTo(HaveOccurred())
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
			Expect(response.Version.Time).To(BeTemporally("~", request.Version.Time, time.Second))
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
			Expect(response.Version.Time).To(BeTemporally("~", time.Now(), time.Second))
		})
	})
})
