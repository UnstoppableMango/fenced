package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestFenced(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Fenced Suite")
}
