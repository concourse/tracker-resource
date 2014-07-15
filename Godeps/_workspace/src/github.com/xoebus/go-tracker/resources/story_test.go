package resources_test

import (
	"encoding/json"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/xoebus/go-tracker/resources"
)

var _ = Describe("Story", func() {
	It("has attributes", func() {
		var stories []resources.Story
		reader := strings.NewReader(Fixture("stories.json"))
		err := json.NewDecoder(reader).Decode(&stories)
		Ω(err).ToNot(HaveOccurred())
		story := stories[0]

		Ω(story.ID).Should(Equal(560))
		Ω(story.Name).Should(Equal("Tractor beam loses power intermittently"))
	})
})
