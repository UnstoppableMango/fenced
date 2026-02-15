package main_test

import (
	"os"
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/unstoppablemango/fenced/cmd"
)

var _ = Describe("E2e", func() {
	When("no args are provided", func() {
		var ses *gexec.Session

		BeforeEach(func() {
			ses = run()
		})

		It("should exit", func() {
			Eventually(ses).Should(gexec.Exit())
		})
	})

	When("a filepath is provided", func() {
		var ses *gexec.Session

		BeforeEach(func() {
			wd, err := os.Getwd()
			Expect(err).NotTo(HaveOccurred())
			ses = run(filepath.Join(wd, "testdata", "markdown.md"))
		})

		It("should exit", func() {
			Eventually(ses).Should(gexec.Exit())
		})

		It("should print the fenced code", func() {
			expected := "import \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}\n"

			Eventually(ses).Should(gexec.Exit(0))
			Expect(ses.Out.Contents()).To(Equal([]byte(expected)))
		})
	})

	Describe("version", func() {
		It("should print the version", func() {
			ses := run("version")

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

func run(args ...string) *gexec.Session {
	cmd := exec.Command(binPath, args...)
	ses, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())

	return ses
}
