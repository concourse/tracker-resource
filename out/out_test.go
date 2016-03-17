package out_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	. "github.com/concourse/tracker-resource/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/concourse/tracker-resource/Godeps/_workspace/src/github.com/onsi/gomega"
	. "github.com/concourse/tracker-resource/Godeps/_workspace/src/github.com/onsi/gomega/gbytes"
	. "github.com/concourse/tracker-resource/Godeps/_workspace/src/github.com/onsi/gomega/gexec"
	"github.com/concourse/tracker-resource/Godeps/_workspace/src/github.com/xoebus/go-tracker"

	"github.com/concourse/tracker-resource/Godeps/_workspace/src/github.com/onsi/gomega/ghttp"

	"github.com/concourse/tracker-resource"
	"github.com/concourse/tracker-resource/out"
)

var _ = Describe("Out", func() {
	var (
		outCmd *exec.Cmd
		tmpdir string
	)

	BeforeEach(func() {
		var err error

		tmpdir, err = ioutil.TempDir("", "out-tmp")
		Ω(err).ShouldNot(HaveOccurred())
		err = os.MkdirAll(tmpdir, 0755)
		Ω(err).ShouldNot(HaveOccurred())

		outCmd = exec.Command(outPath, tmpdir)
	})

	AfterEach(func() {
		os.RemoveAll(tmpdir)
	})

	Describe("integration with the real Tracker API", func() {
		var (
			request            out.OutRequest
			storyId            string
			actualTrackerToken string
		)

		projectId := "1412996"

		BeforeEach(func() {
			actualTrackerToken = os.Getenv("TRACKER_TOKEN")
			if actualTrackerToken == "" {
				Skip("TRACKER_TOKEN must be provided.")
			}

			storyId = createActualStory(projectId, actualTrackerToken)
			setupTestEnvironmentWithActualStoryID(tmpdir, storyId)

			request = out.OutRequest{
				Source: resource.Source{
					Token:      actualTrackerToken,
					TrackerURL: "https://www.pivotaltracker.com",
					ProjectID:  projectId,
				},
				Params: out.Params{
					Repos: []string{
						"middle/git3",
					},
				},
			}
		})

		AfterEach(func() {
			deleteActualStory(projectId, actualTrackerToken, storyId)
		})

		It("finds finished stories that are mentioned in recent git commits", func() {
			session := runCommand(outCmd, request)

			Ω(session.Err).Should(Say(fmt.Sprintf("Checking for finished story: .*#%s", storyId)))
			Ω(session.Err).Should(Say("middle/git3.*... .*DELIVERING"))
		})
	})

	Context("when executed against a mock URL", func() {
		var request out.OutRequest
		var response out.OutResponse

		var server *ghttp.Server

		trackerToken := "abc"
		projectId := "1234"

		BeforeEach(func() {
			setupTestEnvironment(tmpdir)

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

		})

		AfterEach(func() {
			server.Close()
		})

		Context("without a comment file specified", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					listStoriesHandler(trackerToken),
					deliverStoryHandler(trackerToken, projectId, 123456),
					deliverStoryHandler(trackerToken, projectId, 123457),
					deliverStoryHandler(trackerToken, projectId, 223456),
					deliverStoryHandler(trackerToken, projectId, 323456),
					deliverStoryHandler(trackerToken, projectId, 423456),
					deliverStoryHandler(trackerToken, projectId, 523456),
					deliverStoryHandler(trackerToken, projectId, 789456),
					deliverStoryHandler(trackerToken, projectId, 223457),
					deliverStoryHandler(trackerToken, projectId, 323457),
					deliverStoryHandler(trackerToken, projectId, 423457),
				)
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

				Ω(session.Err).Should(Say("Checking for finished story: .*#223456"))
				Ω(session.Err).Should(Say("git.*... .*DELIVERING"))
				Ω(session.Err).Should(Say("middle/git2.*... .*SKIPPING"))

				Ω(session.Err).Should(Say("Checking for finished story: .*#323456"))
				Ω(session.Err).Should(Say("git.*... .*DELIVERING"))
				Ω(session.Err).Should(Say("middle/git2.*... .*SKIPPING"))

				Ω(session.Err).Should(Say("Checking for finished story: .*#423456"))
				Ω(session.Err).Should(Say("git.*... .*DELIVERING"))
				Ω(session.Err).Should(Say("middle/git2.*... .*SKIPPING"))

				Ω(session.Err).Should(Say("Checking for finished story: .*#523456"))
				Ω(session.Err).Should(Say("git.*... .*DELIVERING"))
				Ω(session.Err).Should(Say("middle/git2.*... .*SKIPPING"))

				Ω(session.Err).Should(Say("Checking for finished story: .*#789456"))
				Ω(session.Err).Should(Say("git.*... .*DELIVERING"))
				Ω(session.Err).Should(Say("middle/git2.*... .*SKIPPING"))

				Ω(session.Err).Should(Say("Checking for finished story: .*#223457"))
				Ω(session.Err).Should(Say("git.*... .*SKIPPING"))
				Ω(session.Err).Should(Say("middle/git2.*... .*DELIVERING"))

				Ω(session.Err).Should(Say("Checking for finished story: .*#323457"))
				Ω(session.Err).Should(Say("git.*... .*SKIPPING"))
				Ω(session.Err).Should(Say("middle/git2.*... .*DELIVERING"))

				Ω(session.Err).Should(Say("Checking for finished story: .*#423457"))
				Ω(session.Err).Should(Say("git.*... .*SKIPPING"))
				Ω(session.Err).Should(Say("middle/git2.*... .*DELIVERING"))
			})

			// It("outputs the current time", func() {
			// 	session := runCommand(outCmd, request)

			// 	err := json.Unmarshal(session.Out.Contents(), &response)
			// 	Ω(err).ShouldNot(HaveOccurred())
			// 	Ω(response.Version.Time).Should(BeTemporally("~", time.Now(), time.Second))
			// })
		})

		Context("when a comment file is specified", func() {
			BeforeEach(func() {
				commentPath := "tracker-resource-comment"
				request.Params.CommentPath = commentPath
				err := ioutil.WriteFile(filepath.Join(tmpdir, commentPath), []byte("some custom comment"), os.ModePerm)
				Ω(err).ShouldNot(HaveOccurred())

				server.AppendHandlers(
					listStoriesHandler(trackerToken),
					deliverStoryHandler(trackerToken, projectId, 123456),
					deliverStoryCommentHandler(trackerToken, projectId, 123456, "some custom comment"),
					deliverStoryHandler(trackerToken, projectId, 123457),
					deliverStoryCommentHandler(trackerToken, projectId, 123457, "some custom comment"),
					deliverStoryHandler(trackerToken, projectId, 223456),
					deliverStoryCommentHandler(trackerToken, projectId, 223456, "some custom comment"),
					deliverStoryHandler(trackerToken, projectId, 323456),
					deliverStoryCommentHandler(trackerToken, projectId, 323456, "some custom comment"),
					deliverStoryHandler(trackerToken, projectId, 423456),
					deliverStoryCommentHandler(trackerToken, projectId, 423456, "some custom comment"),
					deliverStoryHandler(trackerToken, projectId, 523456),
					deliverStoryCommentHandler(trackerToken, projectId, 523456, "some custom comment"),
					deliverStoryHandler(trackerToken, projectId, 789456),
					deliverStoryCommentHandler(trackerToken, projectId, 789456, "some custom comment"),
					deliverStoryHandler(trackerToken, projectId, 223457),
					deliverStoryCommentHandler(trackerToken, projectId, 223457, "some custom comment"),
					deliverStoryHandler(trackerToken, projectId, 323457),
					deliverStoryCommentHandler(trackerToken, projectId, 323457, "some custom comment"),
					deliverStoryHandler(trackerToken, projectId, 423457),
					deliverStoryCommentHandler(trackerToken, projectId, 423457, "some custom comment"),
				)
			})

			It("should make a comment with the file's contents", func() {
				session := runCommand(outCmd, request)
				Ω(session.Err).Should(Say("Checking for finished story: .*#123456"))
				Ω(session.Err).Should(Say("Checking for finished story: .*#123457"))

				os.Remove(request.Params.CommentPath)
			})
		})

		Context("when the comment file does not exist", func() {
			It("should return a fatal error", func() {
				request.Params.CommentPath = "some-non-existent-file"

				session := runCommandExpectingStatus(outCmd, request, 1)
				Ω(session.Err).Should(Say("error reading comment file: open"))
			})
		})
	})
})

func runCommand(outCmd *exec.Cmd, request out.OutRequest) *Session {
	return runCommandExpectingStatus(outCmd, request, 0)
}

func runCommandExpectingStatus(outCmd *exec.Cmd, request out.OutRequest, status int) *Session {
	timeout := 10 * time.Second
	stdin, err := outCmd.StdinPipe()
	Ω(err).ShouldNot(HaveOccurred())

	session, err := Start(outCmd, GinkgoWriter, GinkgoWriter)
	Ω(err).ShouldNot(HaveOccurred())
	err = json.NewEncoder(stdin).Encode(request)
	Ω(err).ShouldNot(HaveOccurred())
	Eventually(session, timeout).Should(Exit(status))

	return session
}

func listStoriesHandler(trackerToken string) http.HandlerFunc {
	return ghttp.CombineHandlers(
		ghttp.VerifyRequest("GET", "/services/v5/projects/1234/stories"),
		ghttp.VerifyHeaderKV("X-TrackerToken", trackerToken),
		ghttp.RespondWith(http.StatusOK, Fixture("stories.json")),
	)
}

func deliverStoryCommentHandler(token string, projectId string, storyId int, comment string) http.HandlerFunc {
	body := fmt.Sprintf(`{"text":"%s"}`, comment)
	return ghttp.CombineHandlers(
		ghttp.VerifyRequest(
			"POST",
			fmt.Sprintf("/services/v5/projects/%s/stories/%d/comments", projectId, storyId),
		), ghttp.VerifyHeaderKV("X-TrackerToken", token),
		ghttp.VerifyJSON(body),
	)
}

func deliverStoryHandler(token string, projectId string, storyId int) http.HandlerFunc {
	body := `{"current_state":"delivered"}`
	return ghttp.CombineHandlers(
		ghttp.VerifyRequest(
			"PUT",
			fmt.Sprintf("/services/v5/projects/%s/stories/%d", projectId, storyId),
		), ghttp.VerifyHeaderKV("X-TrackerToken", token),
		ghttp.VerifyJSON(body),
	)
}

func setupTestEnvironment(path string) {
	setupTestEnvironmentWithActualStoryID(path, "")
}

func setupTestEnvironmentWithActualStoryID(path string, storyId string) {
	cmd := exec.Command(filepath.Join("scripts/setup.sh"), path, storyId)
	cmd.Stdout = GinkgoWriter
	cmd.Stderr = GinkgoWriter

	err := cmd.Run()
	Ω(err).ShouldNot(HaveOccurred())
}

func createActualStory(projectID string, trackerToken string) string {
	projectIDInt, err := strconv.Atoi(projectID)
	Ω(err).NotTo(HaveOccurred())

	client := tracker.NewClient(trackerToken).InProject(projectIDInt)
	story := tracker.Story{
		Name:  "concourse test story",
		Type:  tracker.StoryTypeBug,
		State: tracker.StoryStateFinished,
	}
	story, err = client.CreateStory(story)
	Ω(err).NotTo(HaveOccurred())
	return strconv.Itoa(story.ID)
}

func deleteActualStory(projectID string, trackerToken string, storyId string) {
	projectIDInt, err := strconv.Atoi(projectID)
	Ω(err).NotTo(HaveOccurred())

	storyIDInt, err := strconv.Atoi(storyId)
	Ω(err).NotTo(HaveOccurred())

	client := tracker.NewClient(trackerToken).InProject(projectIDInt)
	err = client.DeleteStory(storyIDInt)
	Ω(err).NotTo(HaveOccurred())
}
