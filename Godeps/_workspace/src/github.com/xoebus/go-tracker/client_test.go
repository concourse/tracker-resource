package tracker_test

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/ghttp"

	"github.com/xoebus/go-tracker"
)

var _ = Describe("Tracker Client", func() {
	var (
		server *ghttp.Server
		client *tracker.Client
	)

	BeforeEach(func() {
		server = ghttp.NewServer()
		tracker.DefaultURL = server.URL()
		client = tracker.NewClient("api-token")
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("getting information about the current user", func() {
		var statusCode int

		It("works if everything goes to plan", func() {
			statusCode = http.StatusOK

			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/services/v5/me"),
				verifyTrackerToken(),

				ghttp.RespondWith(statusCode, Fixture("me.json")),
			))

			me, err := client.Me()

			Ω(err).ToNot(HaveOccurred())
			Ω(me.Username).To(Equal("vader"))
		})

		It("returns an error if the response is not successful", func() {
			statusCode = http.StatusInternalServerError

			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.RespondWith(statusCode, ""),
			))

			client := tracker.NewClient("api-token")
			_, err := client.Me()
			Ω(err).To(MatchError("request failed (500)"))
		})

		It("returns a helpful error if the token is invalid", func() {
			statusCode = http.StatusUnauthorized

			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.RespondWith(statusCode, ""),
			))

			client := tracker.NewClient("api-token")
			_, err := client.Me()
			Ω(err).To(MatchError("invalid token"))
		})

		It("returns an error if the request fails", func() {
			server.Close()

			client := tracker.NewClient("api-token")
			_, err := client.Me()

			Ω(err).To(HaveOccurred())
			Ω(err.Error()).To(MatchRegexp("failed to make request"))
			server = ghttp.NewServer()
		})

		It("returns an error if the request can't be created", func() {
			tracker.DefaultURL = "aaaaa)#Q&%*(*"

			client := tracker.NewClient("api-token")
			_, err := client.Me()

			Ω(err).To(HaveOccurred())
			Ω(err.Error()).To(MatchRegexp("failed to create request"))
		})

		It("returns an error if the response JSON is broken", func() {
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.RespondWith(http.StatusOK, `{"`),
			))

			client := tracker.NewClient("api-token")
			_, err := client.Me()

			Ω(err).To(HaveOccurred())
			Ω(err.Error()).To(MatchRegexp("invalid json response"))
		})
	})

	Describe("listing stories", func() {
		It("gets all the stories by default", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/services/v5/projects/99/stories", "date_format=millis"),
					verifyTrackerToken(),

					ghttp.RespondWith(http.StatusOK, Fixture("stories.json")),
				),
			)

			client := tracker.NewClient("api-token")

			stories, err := client.InProject(99).Stories(tracker.StoriesQuery{})
			Ω(stories).Should(HaveLen(4))
			Ω(err).ToNot(HaveOccurred())
		})

		It("allows different queries to be made", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/services/v5/projects/99/stories", "date_format=millis&with_state=finished"),
					verifyTrackerToken(),

					ghttp.RespondWith(http.StatusOK, Fixture("stories.json")),
				),
			)

			client := tracker.NewClient("api-token")

			query := tracker.StoriesQuery{
				State: tracker.StateFinished,
			}
			stories, err := client.InProject(99).Stories(query)
			Ω(stories).Should(HaveLen(4))
			Ω(err).ToNot(HaveOccurred())
		})
	})

	Describe("delivering a story", func() {
		It("HTTP PUTs it in its place", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("PUT", "/services/v5/projects/99/stories/15225523"),
					ghttp.VerifyJSON(`{"current_state":"delivered"}`),
					verifyTrackerToken(),

					ghttp.RespondWith(http.StatusOK, ""),
				),
			)

			client := tracker.NewClient("api-token")

			err := client.InProject(99).DeliverStory(15225523)
			Ω(err).ShouldNot(HaveOccurred())
		})
	})
})

func verifyTrackerToken() http.HandlerFunc {
	headers := http.Header{
		"X-TrackerToken": {"api-token"},
	}

	return ghttp.VerifyHeader(headers)
}
