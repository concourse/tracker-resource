package out_test

import (
	"encoding/json"
	"fmt"
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

	"github.com/concourse/tracker-resource"
	"github.com/concourse/tracker-resource/out"
)

var _ = Describe("Out", func() {
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

		trackerToken := "abc"
		projectId := "1234"

		BeforeEach(func() {
			server = ghttp.NewServer()

			request = out.OutRequest{
				Source: resource.Source{
					Token:      trackerToken,
					TrackerURL: server.URL(),
					ProjectID:  projectId,
				},
				Params: out.Params{
					Repos: []string{
						"git",
						"middle/git2",
					},
				},
			}
			response = out.OutResponse{}

			comment := "Delivered by Concourse"

			server.AppendHandlers(
				listStoriesHandler(),
				deliverStoryHandler(trackerToken, projectId, 123456, comment),
				deliverStoryHandler(trackerToken, projectId, 123457, comment),
			)
		})

		AfterEach(func() {
			server.Close()
		})

		It("does not output credentials", func() {
			session := runCommand(outCmd, request)

			Ω(session.Err).ShouldNot(Say(trackerToken))
		})

		It("finds finished stories that are mentioned in recent git commits", func() {
			session := runCommand(outCmd, request)

			Ω(session.Err).Should(Say("Checking for finished story: .*#565"))
			Ω(session.Err).Should(Say("git.*... .*SKIPPING"))
			Ω(session.Err).Should(Say("middle/git2.*... .*SKIPPING"))

			Ω(session.Err).Should(Say("Checking for finished story: .*#123456"))
			Ω(session.Err).Should(Say("git.*... .*DELIVERING"))
			Ω(session.Err).Should(Say("middle/git2.*... .*SKIPPING"))

			Ω(session.Err).Should(Say("Checking for finished story: .*#123457"))
			Ω(session.Err).Should(Say("git.*... .*SKIPPING"))
			Ω(session.Err).Should(Say("middle/git2.*... .*DELIVERING"))
		})

		It("outputs the current time", func() {
			session := runCommand(outCmd, request)

			err := json.Unmarshal(session.Out.Contents(), &response)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(response.Version.Time).Should(BeTemporally("~", time.Now(), time.Second))
		})
	})
})

func runCommand(outCmd *exec.Cmd, request out.OutRequest) *Session {
	stdin, err := outCmd.StdinPipe()
	Ω(err).ShouldNot(HaveOccurred())

	session, err := Start(outCmd, GinkgoWriter, GinkgoWriter)
	Ω(err).ShouldNot(HaveOccurred())
	err = json.NewEncoder(stdin).Encode(request)
	Ω(err).ShouldNot(HaveOccurred())
	Eventually(session).Should(Exit(0))

	return session
}

func listStoriesHandler() http.HandlerFunc {
	return ghttp.CombineHandlers(
		ghttp.VerifyRequest("GET", "/services/v5/projects/1234/stories"),
		ghttp.VerifyHeaderKV("X-TrackerToken", "abc"),
		ghttp.RespondWith(http.StatusOK, Fixture("stories.json")),
	)
}

func deliverStoryHandler(token string, projectId string, storyId int, comment string) http.HandlerFunc {
	body := fmt.Sprintf(`{"current_state":"delivered", "comment":"%s"}`, comment)
	return ghttp.CombineHandlers(
		ghttp.VerifyRequest("PUT", fmt.Sprintf("/services/v5/projects/%s/stories/%d", projectId, storyId)),
		ghttp.VerifyHeaderKV("X-TrackerToken", token),
		ghttp.VerifyJSON(body),
	)
}

func setupTestEnvironment(path string) {
	cmd := exec.Command(filepath.Join("scripts/setup.sh"), path)
	cmd.Stdout = GinkgoWriter
	cmd.Stderr = GinkgoWriter

	err := cmd.Run()
	Ω(err).ShouldNot(HaveOccurred())
}
