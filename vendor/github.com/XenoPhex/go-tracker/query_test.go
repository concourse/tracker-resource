package tracker_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/xoebus/go-tracker"
)

var _ = Describe("Queries", func() {
	queryString := func(query tracker.Query) string {
		return query.Query().Encode()
	}

	Describe("StoriesQuery", func() {
		It("only has date_format by default", func() {
			query := tracker.StoriesQuery{}
			Ω(queryString(query)).Should(Equal(""))
		})

		It("can query by story state", func() {
			query := tracker.StoriesQuery{
				State: tracker.StoryStateRejected,
			}
			Ω(queryString(query)).Should(Equal("with_state=rejected"))
		})

		It("can query by story labels", func() {
			query := tracker.StoriesQuery{
				Label: "blocked",
			}
			Ω(queryString(query)).Should(Equal("with_label=blocked"))
		})

		It("can limit the numer of results", func() {
			query := tracker.StoriesQuery{
				Limit: 33,
			}
			Ω(queryString(query)).Should(Equal("limit=33"))
		})
	})
})
