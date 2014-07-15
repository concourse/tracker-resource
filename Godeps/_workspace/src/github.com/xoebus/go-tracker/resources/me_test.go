package resources_test

import (
	"encoding/json"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/xoebus/go-tracker/resources"
)

var _ = Describe("Me", func() {
	It("has attributes", func() {
		var me resources.Me
		reader := strings.NewReader(Fixture("me.json"))
		err := json.NewDecoder(reader).Decode(&me)
		Ω(err).ToNot(HaveOccurred())

		Ω(me.Username).To(Equal("vader"))
		Ω(me.Name).To(Equal("Darth Vader"))
		Ω(me.Initials).To(Equal("DV"))
		Ω(me.ID).To(Equal(101))
		Ω(me.Email).To(Equal("vader@deathstar.mil"))
	})
})
