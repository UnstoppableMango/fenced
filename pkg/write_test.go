package fenced_test

import (
	"bytes"
	"errors"
	"io"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	fenced "github.com/unstoppablemango/fenced/pkg"
)

var _ = Describe("Writer", func() {
	Describe("Write", func() {
		It("should insert delimiter between blocks written in separate Write calls", func() {
			var buf bytes.Buffer
			w := fenced.NewWriter(&buf, fenced.WithDelimiter("---"))

			_, err := w.Write(fenced.Block{Content: "a\n"})
			Expect(err).NotTo(HaveOccurred())

			_, err = w.Write(fenced.Block{Content: "b\n"})
			Expect(err).NotTo(HaveOccurred())

			Expect(buf.String()).To(Equal("a\n---\nb\n"))
		})
	})
})

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

	It("should return an error when the writer fails", func() {
		pr, pw := io.Pipe()
		_ = pr.Close()

		_, err := fenced.Write(pw, fenced.Block{Content: "x"})

		Expect(err).To(MatchError(io.ErrClosedPipe))
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

	It("should return an error when the writer fails", func() {
		pr, pw := io.Pipe()
		_ = pr.Close()

		_, err := fenced.WriteAll(pw, []fenced.Block{{Content: "x"}})

		Expect(err).To(MatchError(io.ErrClosedPipe))
	})

	It("should return bytes written for successful blocks before a failure", func() {
		pr, pw := io.Pipe()
		go func() {
			buf := make([]byte, 2)
			_, _ = pr.Read(buf) // drain "a\n"
			pr.CloseWithError(errors.New("write error"))
		}()

		blocks := []fenced.Block{{Content: "a\n"}, {Content: "b\n"}}
		n, err := fenced.WriteAll(pw, blocks)
		_ = pw.Close()

		Expect(err).To(HaveOccurred())
		Expect(n).To(Equal(2))
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

		It("should return an error when writing delimiter fails", func() {
			pr, pw := io.Pipe()
			go func() {
				buf := make([]byte, 100)
				_, _ = pr.Read(buf) // drain first block write
				pr.CloseWithError(errors.New("delimiter error"))
			}()

			blocks := []fenced.Block{{Content: "a\n"}, {Content: "b\n"}}
			_, err := fenced.WriteAll(pw, blocks, fenced.WithDelimiter("---"))
			_ = pw.Close()

			Expect(err).To(MatchError("delimiter error"))
		})

		It("should return an error when writing implicit newline fails", func() {
			pr, pw := io.Pipe()
			go func() {
				buf := make([]byte, 100)
				_, _ = pr.Read(buf) // drain first block write
				_, _ = pr.Read(buf) // drain delimiter write
				pr.CloseWithError(errors.New("newline error"))
			}()

			blocks := []fenced.Block{{Content: "a\n"}, {Content: "b\n"}}
			_, err := fenced.WriteAll(pw, blocks, fenced.WithDelimiter("---"))
			_ = pw.Close()

			Expect(err).To(MatchError("newline error"))
		})

		It("should return bytes written before delimiter write failure", func() {
			pr, pw := io.Pipe()
			go func() {
				buf := make([]byte, 2)
				_, _ = pr.Read(buf) // drain "a\n"
				pr.CloseWithError(errors.New("delimiter error"))
			}()

			blocks := []fenced.Block{{Content: "a\n"}, {Content: "b\n"}}
			n, err := fenced.WriteAll(pw, blocks, fenced.WithDelimiter("---"))
			_ = pw.Close()

			Expect(err).To(HaveOccurred())
			Expect(n).To(Equal(2))
		})

		It("should return bytes written before block content write failure", func() {
			pr, pw := io.Pipe()
			go func() {
				buf := make([]byte, 100)
				_, _ = pr.Read(buf) // drain "a\n"
				_, _ = pr.Read(buf) // drain "---"
				_, _ = pr.Read(buf) // drain "\n"
				pr.CloseWithError(errors.New("content error"))
			}()

			blocks := []fenced.Block{{Content: "a\n"}, {Content: "b\n"}}
			n, err := fenced.WriteAll(pw, blocks, fenced.WithDelimiter("---"))
			_ = pw.Close()

			Expect(err).To(HaveOccurred())
			Expect(n).To(Equal(6)) // "a\n"(2) + "---"(3) + "\n"(1)
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
