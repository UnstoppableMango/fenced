package cmd_test

import (
	"os"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/spf13/cobra"
	"github.com/unstoppablemango/fenced/cmd"
)

var _ = Describe("Execute", func() {
	var origArgs []string

	BeforeEach(func() {
		origArgs = os.Args
	})

	AfterEach(func() {
		os.Args = origArgs
	})

	It("should run the version subcommand", func() {
		os.Args = []string{"fenced", "version"}

		Expect(cmd.Execute()).To(Succeed())
	})

	It("should parse a file", func() {
		wd, err := os.Getwd()
		Expect(err).NotTo(HaveOccurred())
		os.Args = []string{"fenced", filepath.Join(wd, "..", "testdata", "markdown.md")}

		Expect(cmd.Execute()).To(Succeed())
	})

	It("should enable debug logging when DEBUG is set", func() {
		DeferCleanup(os.Unsetenv, "DEBUG")
		Expect(os.Setenv("DEBUG", "1")).To(Succeed())
		os.Args = []string{"fenced", "version"}

		Expect(cmd.Execute()).To(Succeed())
	})

	It("should apply no-implicit-newline flag", func() {
		wd, err := os.Getwd()
		Expect(err).NotTo(HaveOccurred())
		os.Args = []string{"fenced", "-N", filepath.Join(wd, "..", "testdata", "markdown.md")}

		Expect(cmd.Execute()).To(Succeed())
	})
})

var _ = Describe("Open", func() {
	It("should return stdin when path is '-'", func() {
		c := &cobra.Command{}
		c.SetIn(strings.NewReader("hello"))

		rc, err := cmd.Open(c, "-")

		Expect(err).NotTo(HaveOccurred())
		Expect(rc).NotTo(BeNil())
	})

	It("should open a file when path is a file", func() {
		wd, err := os.Getwd()
		Expect(err).NotTo(HaveOccurred())
		path := filepath.Join(wd, "..", "testdata", "markdown.md")

		rc, err := cmd.Open(&cobra.Command{}, path)

		Expect(err).NotTo(HaveOccurred())
		Expect(rc).NotTo(BeNil())
		Expect(rc.Close()).To(Succeed())
	})

	It("should return an error for a non-existent file", func() {
		_, err := cmd.Open(&cobra.Command{}, "/does/not/exist.md")

		Expect(err).To(HaveOccurred())
	})
})

var _ = Describe("OpenAll", func() {
	It("should return stdin when no args are provided", func() {
		c := &cobra.Command{}
		c.SetIn(strings.NewReader("hello"))

		readers, err := cmd.OpenAll(c, []string{})

		Expect(err).NotTo(HaveOccurred())
		Expect(readers).To(HaveLen(1))
	})

	It("should use cmd.InOrStdin when os.Stdin is replaced with a non-terminal", func() {
		r, w, err := os.Pipe()
		Expect(err).NotTo(HaveOccurred())
		defer func() { Expect(r.Close()).To(Succeed()) }()
		defer func() { Expect(w.Close()).To(Succeed()) }()

		old := os.Stdin
		os.Stdin = r
		defer func() { os.Stdin = old }()

		c := &cobra.Command{}
		c.SetIn(strings.NewReader("hello"))

		readers, err := cmd.OpenAll(c, []string{})

		Expect(err).NotTo(HaveOccurred())
		Expect(readers).To(HaveLen(1))
	})

	It("should open files for each path argument", func() {
		wd, err := os.Getwd()
		Expect(err).NotTo(HaveOccurred())
		path := filepath.Join(wd, "..", "testdata", "markdown.md")

		readers, err := cmd.OpenAll(&cobra.Command{}, []string{path})

		Expect(err).NotTo(HaveOccurred())
		Expect(readers).To(HaveLen(1))
		for _, r := range readers {
			Expect(r.Close()).To(Succeed())
		}
	})

	It("should return an error when a path cannot be opened", func() {
		_, err := cmd.OpenAll(&cobra.Command{}, []string{"/does/not/exist.md"})

		Expect(err).To(HaveOccurred())
	})

	It("should open multiple files", func() {
		wd, err := os.Getwd()
		Expect(err).NotTo(HaveOccurred())
		testdata := filepath.Join(wd, "..", "testdata")

		readers, err := cmd.OpenAll(&cobra.Command{}, []string{
			filepath.Join(testdata, "markdown.md"),
			filepath.Join(testdata, "python.md"),
		})

		Expect(err).NotTo(HaveOccurred())
		Expect(readers).To(HaveLen(2))
		for _, r := range readers {
			Expect(r.Close()).To(Succeed())
		}
	})
})
