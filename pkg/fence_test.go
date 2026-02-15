package fenced_test

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	fenced "github.com/unstoppablemango/fenced/pkg"
)

var _ = Describe("Fence", func() {
	Describe("Parse", func() {
		It("should parse a single code block", func() {
			input := "```\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}\n```"
			expected := []string{"import \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}\n"}

			codeBlocks, err := fenced.Parse(strings.NewReader(input))

			Expect(err).NotTo(HaveOccurred())
			Expect(codeBlocks).To(Equal(expected))
		})
	})
})
