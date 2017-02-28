package matcho_test

import (
	"github.com/vlad-stoian/matcho/matchers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Matcho", func() {
	Context("YAML Matcher", func() {
		It("Should return default values", func() {
			yamlToMatch := matcho.MatchYAMLMatcher{YAMLToMatch: "abc"}
			success, err := yamlToMatch.Match("abc")

			Expect(success).Should(BeTrue())
			Expect(err).To(BeNil())

		})
	})

})
