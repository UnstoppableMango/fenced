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
