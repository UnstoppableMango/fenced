// Package main is the entry point for the fenced CLI application.
package main

import (
	"github.com/unmango/go/cli"
	"github.com/unstoppablemango/fenced/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		cli.Fail(err)
	}
}
