// Package cmd implements the command-line interface for fenced.
package cmd

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/unmango/go/cli"
	fenced "github.com/unstoppablemango/fenced/pkg"
	"golang.org/x/term"
)

var rootCmd = &cobra.Command{
	Use:   "fenced [PATH...]",
	Short: "Parses code fences",
	Args:  cobra.ArbitraryArgs,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if d := os.Getenv("DEBUG"); d != "" {
			log.SetLevel(log.DebugLevel)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		readers, err := OpenAll(cmd, args)
		if err != nil {
			cli.Fail(err)
		}

		delimiter, err := cmd.Flags().GetString("delimiter")
		if err != nil {
			cli.Fail(err)
		}

		out := cmd.OutOrStdout()
		first := true
		for _, in := range readers {
			blocks, err := fenced.Parse(in)
			if err != nil {
				cli.Fail(err)
			}

			for _, b := range blocks {
				if !first && delimiter != "" {
					if _, err := io.WriteString(out, delimiter); err != nil {
						cli.Fail(err)
					}
				}
				first = false

				if _, err := io.WriteString(out, b.String()); err != nil {
					cli.Fail(err)
				}
			}

			if err := in.Close(); err != nil {
				log.Warn("Failed to close reader", "error", err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.Flags().StringP("delimiter", "d", "", "delimiter to insert between code blocks")
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

// Open returns a reader for the given input (file or stdin).
func Open(cmd *cobra.Command, path string) (io.ReadCloser, error) {
	if path == "-" {
		log.Debug("Reading from stdin")
		return io.NopCloser(cmd.InOrStdin()), nil
	} else {
		log.Debug("Opening file", "path", path)
		return os.Open(filepath.Clean(path))
	}
}

// OpenAll returns readers for the input sources (files or stdin).
func OpenAll(cmd *cobra.Command, args []string) ([]io.ReadCloser, error) {
	if len(args) == 0 {
		if term.IsTerminal(int(os.Stdin.Fd())) {
			return nil, errors.New("stdin is a terminal; provide a file path or pipe input")
		}
		in := io.NopCloser(cmd.InOrStdin())
		return []io.ReadCloser{in}, nil
	}

	readers := make([]io.ReadCloser, 0, len(args))
	for _, path := range args {
		if r, err := Open(cmd, path); err != nil {
			return nil, err
		} else {
			readers = append(readers, r)
		}
	}

	return readers, nil
}
