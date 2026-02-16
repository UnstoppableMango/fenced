package main_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/unstoppablemango/fenced/cmd"

	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("E2e", func() {
	When("no args are provided and stdin is empty", func() {
		It("should exit successfully", func() {
			ses := run(exec.Command(binPath))

			Eventually(ses).Should(gexec.Exit(0))
		})
	})

	When("stdin has piped content", func() {
		It("should parse fenced blocks from stdin", func() {
			cmd := exec.Command(binPath)
			cmd.Stdin = bytes.NewBufferString("```go\nfmt.Println(\"test\")\n```\n")

			ses := run(cmd)

			Eventually(ses).Should(gexec.Exit(0))
			Expect(ses.Out).Should(gbytes.Say("fmt.Println"))
		})
	})

	When("a filepath is provided", func() {
		var testdata string

		BeforeEach(func() {
			wd, err := os.Getwd()
			Expect(err).NotTo(HaveOccurred())
			testdata = filepath.Join(wd, "testdata")
		})

		It("should exit", func() {
			cmd := exec.Command(binPath, filepath.Join(testdata, "markdown.md"))

			ses := run(cmd)

			Eventually(ses).Should(gexec.Exit())
		})

		It("should print the fenced code", func() {
			cmd := exec.Command(binPath, filepath.Join(testdata, "markdown.md"))
			expected := "import \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}\n"

			ses := run(cmd)

			Eventually(ses).Should(gexec.Exit(0))
			Expect(ses.Out.Contents()).To(Equal([]byte(expected)))
		})
	})

	Describe("version", func() {
		It("should print the version", func() {
			ses := run(exec.Command(binPath, "version"))

			Eventually(ses).Should(gexec.Exit(0))
			Expect(ses.Out).Should(gbytes.Say(cmd.Version))
		})
	})
})

var binPath string

var _ = BeforeSuite(func() {
	var err error
	binPath, err = gexec.Build("main.go")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})

func run(cmd *exec.Cmd) *gexec.Session {
	ses, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	return ses
}
