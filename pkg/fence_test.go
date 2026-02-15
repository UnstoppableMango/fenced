package fenced_test

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	fenced "github.com/unstoppablemango/fenced/pkg"
)

var _ = Describe("Fence", func() {
	Describe("Block", func() {
		It("should return the content as a string", func() {
			block := fenced.Block{Content: "import \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}\n"}

			Expect(block.String()).To(Equal("import \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}\n"))
		})
	})

	Describe("Parse", func() {
		It("should parse a single code block", func() {
			input := "```\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}\n```"
			expected := []fenced.Block{{
				Content: "import \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}\n",
			}}

			codeBlocks, err := fenced.Parse(strings.NewReader(input))

			Expect(err).NotTo(HaveOccurred())
			Expect(codeBlocks).To(Equal(expected))
		})

		It("should parse a single code block with language hint", func() {
			input := "```go\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}\n```"
			expected := []fenced.Block{{
				Content: "import \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}\n",
				Lang:    "go",
			}}

			codeBlocks, err := fenced.Parse(strings.NewReader(input))

			Expect(err).NotTo(HaveOccurred())
			Expect(codeBlocks).To(Equal(expected))
		})

		It("should parse a multiple code blocks", func() {
			input := "```\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}\n```\n```\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}\n```"
			expected := []fenced.Block{
				{Content: "import \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}\n"},
				{Content: "import \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}\n"},
			}

			codeBlocks, err := fenced.Parse(strings.NewReader(input))

			Expect(err).NotTo(HaveOccurred())
			Expect(codeBlocks).To(Equal(expected))
		})

		It("should parse a single code block with asterisks", func() {
			input := "***\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}\n***"
			expected := []fenced.Block{{
				Content: "import \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}\n",
			}}

			codeBlocks, err := fenced.Parse(strings.NewReader(input))

			Expect(err).NotTo(HaveOccurred())
			Expect(codeBlocks).To(Equal(expected))
		})

		It("should parse an unclosed code block", func() {
			input := "```\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}\n"

			codeBlocks, err := fenced.Parse(strings.NewReader(input))

			Expect(err).NotTo(HaveOccurred())
			Expect(codeBlocks).To(BeEmpty())
		})
	})
})
