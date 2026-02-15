package main_test

import (
	"embed"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

//go:embed testdata/*
var testdata embed.FS

func TestFenced(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Fenced Suite")
}
