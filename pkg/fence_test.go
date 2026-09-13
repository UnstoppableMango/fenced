package fenced_test

import (
	"errors"
	"fmt"
	"strings"
	"testing/iotest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	fenced "github.com/unstoppablemango/fenced/pkg"
)

const helloProgram = "import \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}\n"

var _ = Describe("Fence", func() {
	Describe("Block", func() {
		It("should return the content as a string", func() {
			block := fenced.Block{Content: helloProgram}

			Expect(block.String()).To(Equal(helloProgram))
		})
	})

	Describe("Parse", func() {
		It("should parse a single code block", func() {
			input := fmt.Sprint("```\n", helloProgram, "```")
			expected := []fenced.Block{{
				Content: helloProgram,
			}}

			codeBlocks, err := fenced.Parse(strings.NewReader(input))

			Expect(err).NotTo(HaveOccurred())
			Expect(codeBlocks).To(Equal(expected))
		})

		It("should parse a single code block with language hint", func() {
			input := fmt.Sprint("```go\n", helloProgram, "```")
			expected := []fenced.Block{{
				Content: helloProgram,
				Lang:    "go",
			}}

			codeBlocks, err := fenced.Parse(strings.NewReader(input))

			Expect(err).NotTo(HaveOccurred())
			Expect(codeBlocks).To(Equal(expected))
		})

		It("should parse a multiple code blocks", func() {
			input := fmt.Sprint("```\n", helloProgram, "```\n```\n", helloProgram, "```")
			expected := []fenced.Block{
				{Content: helloProgram},
				{Content: helloProgram},
			}

			codeBlocks, err := fenced.Parse(strings.NewReader(input))

			Expect(err).NotTo(HaveOccurred())
			Expect(codeBlocks).To(Equal(expected))
		})

		It("should parse a single code block with tildes", func() {
			input := fmt.Sprint("~~~\n", helloProgram, "~~~")
			expected := []fenced.Block{{
				Content: helloProgram,
			}}

			codeBlocks, err := fenced.Parse(strings.NewReader(input))

			Expect(err).NotTo(HaveOccurred())
			Expect(codeBlocks).To(Equal(expected))
		})

		It("should parse an unclosed code block", func() {
			// input := "```\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}\n"
			input := fmt.Sprint("```\n", helloProgram)

			codeBlocks, err := fenced.Parse(strings.NewReader(input))

			Expect(err).NotTo(HaveOccurred())
			Expect(codeBlocks).To(BeEmpty())
		})

		It("should return an error when the reader fails", func() {
			r := iotest.ErrReader(errors.New("read error"))

			_, err := fenced.Parse(r)

			Expect(err).To(MatchError("read error"))
		})
	})
})
