package fenced

import (
	"io"

	"github.com/unmango/go/fopt"
)

// Writer writes fenced code blocks to an underlying io.Writer.
type Writer struct {
	w                 io.Writer
	delimiter         string
	noImplicitNewline bool
}

// NewWriter creates a new Writer that writes to w with the given options applied.
func NewWriter(w io.Writer, options ...Option) *Writer {
	fw := &Writer{w: w}
	fopt.ApplyAll(fw, options)
	return fw
}

func (w *Writer) implicitNewline() bool {
	return !w.noImplicitNewline
}

func (w *Writer) Write(b Block) (n int, err error) {
	return w.write(b, true)
}

func (w *Writer) write(b Block, first bool) (n int, err error) {
	if !first && w.delimiter != "" {
		if nn, err := io.WriteString(w.w, w.delimiter); err != nil {
			return 0, err
		} else {
			n += nn
		}
		if w.implicitNewline() {
			if nn, err := io.WriteString(w.w, "\n"); err != nil {
				return 0, err
			} else {
				n += nn
			}
		}
	}

	if nn, err := io.WriteString(w.w, b.String()); err != nil {
		return 0, err
	} else {
		n += nn
	}

	return
}

// Option configures a Writer.
type Option func(*Writer)

// WithDelimiter sets the string written between blocks.
func WithDelimiter(delim string) Option {
	return func(o *Writer) {
		o.delimiter = delim
	}
}

// WithNoImplicitNewline disables the implicit newline written after the delimiter.
func WithNoImplicitNewline(o *Writer) {
	o.noImplicitNewline = true
}

// WriteAll writes all blocks to w using the given options.
func WriteAll(w io.Writer, blocks []Block, options ...Option) (n int, err error) {
	fs := NewWriter(w, options...)
	for i, b := range blocks {
		if nn, err := fs.write(b, i == 0); err != nil {
			return 0, err
		} else {
			n += nn
		}
	}

	return
}

// Write writes a single block to w using the given options.
func Write(w io.Writer, b Block, options ...Option) (n int, err error) {
	return NewWriter(w, options...).Write(b)
}
