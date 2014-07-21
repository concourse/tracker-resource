package out_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"

	"github.com/onsi/gomega/ghttp"

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

		var server *ghttp.Server

		BeforeEach(func() {
			server = ghttp.NewServer()

			request = out.OutRequest{
				Source: out.Source{
					Token:      "abc",
					TrackerURL: server.URL(),
					ProjectID:  1234,
				},
				Params: out.Params{
					Repos: []string{
						"git",
						"middle/git2",
					},
				},
			}
			response = out.OutResponse{}
		})

		AfterEach(func() {
			server.Close()
		})

		It("finds finished stories that are mentioned in recent git commits", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/services/v5/projects/1234/stories"),
					ghttp.VerifyHeaderKV("X-TrackerToken", "abc"),
					ghttp.RespondWith(http.StatusOK, Fixture("stories.json")),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("PUT", "/services/v5/projects/1234/stories/123456"),
					ghttp.VerifyHeaderKV("X-TrackerToken", "abc"),
					ghttp.VerifyJSON(`{"current_state":"delivered"}`),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("PUT", "/services/v5/projects/1234/stories/123457"),
					ghttp.VerifyHeaderKV("X-TrackerToken", "abc"),
					ghttp.VerifyJSON(`{"current_state":"delivered"}`),
				),
			)

			stdin, err := outCmd.StdinPipe()
			Ω(err).ShouldNot(HaveOccurred())

			session, err := Start(outCmd, GinkgoWriter, GinkgoWriter)
			Ω(err).ShouldNot(HaveOccurred())
			err = json.NewEncoder(stdin).Encode(request)
			Ω(err).ShouldNot(HaveOccurred())
			Eventually(session).Should(Exit(0))

			Ω(session.Err).Should(Say("could not find story for delivery: 565"))
			Ω(session.Err).Should(Say("delivering it!: 123456"))
			Ω(session.Err).Should(Say("delivering it!: 123457"))

			// Output
			err = json.Unmarshal(session.Out.Contents(), &response)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(response.Version.Time).Should(BeTemporally("~", time.Now(), time.Second))
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
