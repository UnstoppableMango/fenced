package fenced_test

import (
	"bytes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	fenced "github.com/unstoppablemango/fenced/pkg"
)

var _ = Describe("Write", func() {
	It("should write the block content", func() {
		block := fenced.Block{Content: "package main\n"}
		var buf bytes.Buffer

		_, err := fenced.Write(&buf, block)

		Expect(err).NotTo(HaveOccurred())
		Expect(buf.String()).To(Equal("package main\n"))
	})

	It("should return the number of bytes written", func() {
		block := fenced.Block{Content: "abc\n"}
		var buf bytes.Buffer

		n, err := fenced.Write(&buf, block)

		Expect(err).NotTo(HaveOccurred())
		Expect(n).To(Equal(4))
	})
})

var _ = Describe("WriteAll", func() {
	It("should write a single block", func() {
		blocks := []fenced.Block{{Content: "package main\n"}}
		var buf bytes.Buffer

		_, err := fenced.WriteAll(&buf, blocks)

		Expect(err).NotTo(HaveOccurred())
		Expect(buf.String()).To(Equal("package main\n"))
	})

	It("should write multiple blocks without delimiter", func() {
		blocks := []fenced.Block{
			{Content: "package main\n"},
			{Content: "func main() {}\n"},
		}
		var buf bytes.Buffer

		_, err := fenced.WriteAll(&buf, blocks)

		Expect(err).NotTo(HaveOccurred())
		Expect(buf.String()).To(Equal("package main\nfunc main() {}\n"))
	})

	It("should write nothing for empty blocks", func() {
		var buf bytes.Buffer

		_, err := fenced.WriteAll(&buf, []fenced.Block{})

		Expect(err).NotTo(HaveOccurred())
		Expect(buf.String()).To(BeEmpty())
	})

	It("should return the number of bytes written", func() {
		blocks := []fenced.Block{{Content: "abc\n"}}
		var buf bytes.Buffer

		n, err := fenced.WriteAll(&buf, blocks)

		Expect(err).NotTo(HaveOccurred())
		Expect(n).To(Equal(4))
	})

	Describe("WithDelimiter", func() {
		It("should not write delimiter before the first block", func() {
			blocks := []fenced.Block{{Content: "package main\n"}}
			var buf bytes.Buffer

			_, err := fenced.WriteAll(&buf, blocks, fenced.WithDelimiter("---"))

			Expect(err).NotTo(HaveOccurred())
			Expect(buf.String()).To(Equal("package main\n"))
		})

		It("should insert delimiter between multiple blocks", func() {
			blocks := []fenced.Block{
				{Content: "package main\n"},
				{Content: "func main() {}\n"},
			}
			var buf bytes.Buffer

			_, err := fenced.WriteAll(&buf, blocks, fenced.WithDelimiter("---"))

			Expect(err).NotTo(HaveOccurred())
			Expect(buf.String()).To(Equal("package main\n---\nfunc main() {}\n"))
		})

		It("should insert delimiter between three blocks", func() {
			blocks := []fenced.Block{
				{Content: "a\n"},
				{Content: "b\n"},
				{Content: "c\n"},
			}
			var buf bytes.Buffer

			_, err := fenced.WriteAll(&buf, blocks, fenced.WithDelimiter("==="))

			Expect(err).NotTo(HaveOccurred())
			Expect(buf.String()).To(Equal("a\n===\nb\n===\nc\n"))
		})
	})

	Describe("WithNoImplicitNewline", func() {
		It("should not append newline after delimiter", func() {
			blocks := []fenced.Block{
				{Content: "package main\n"},
				{Content: "func main() {}\n"},
			}
			var buf bytes.Buffer

			_, err := fenced.WriteAll(&buf, blocks, fenced.WithDelimiter("---"), fenced.WithNoImplicitNewline)

			Expect(err).NotTo(HaveOccurred())
			Expect(buf.String()).To(Equal("package main\n---func main() {}\n"))
		})

		It("should have no effect without a delimiter", func() {
			blocks := []fenced.Block{
				{Content: "package main\n"},
				{Content: "func main() {}\n"},
			}
			var buf bytes.Buffer

			_, err := fenced.WriteAll(&buf, blocks, fenced.WithNoImplicitNewline)

			Expect(err).NotTo(HaveOccurred())
			Expect(buf.String()).To(Equal("package main\nfunc main() {}\n"))
		})
	})
})
