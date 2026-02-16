// Package cmd implements the command-line interface for fenced.
package cmd

import (
	"errors"
	"io"
	"os"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/unmango/go/cli"
	fenced "github.com/unstoppablemango/fenced/pkg"
	"golang.org/x/term"
)

var rootCmd = &cobra.Command{
	Use:   "fenced [PATH]",
	Short: "Parse code fences from anywhere",
	Args:  cobra.MaximumNArgs(1),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if d := os.Getenv("DEBUG"); d != "" {
			log.SetLevel(log.DebugLevel)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		in, err := Open(cmd, args)
		if err != nil {
			cli.Fail(err)
		}

		blocks, err := fenced.Parse(in)
		if err != nil {
			cli.Fail(err)
		}

		out := cmd.OutOrStdout()
		for _, b := range blocks {
			if _, err := io.WriteString(out, b.String()); err != nil {
				cli.Fail(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

// Open returns a reader for the input source (file or stdin).
func Open(cmd *cobra.Command, args []string) (io.Reader, error) {
	if len(args) == 0 {
		log.Debug("Choosing stdin")
		if term.IsTerminal(int(os.Stdin.Fd())) {
			return nil, errors.New("stdin is a terminal; please provide a file path or pipe input")
		}
		return cmd.InOrStdin(), nil
	}

	log.Debug("Opening file", "path", args[0])
	if file, err := os.Open(args[0]); err != nil {
		return nil, err
	} else {
		return file, nil
	}
}
