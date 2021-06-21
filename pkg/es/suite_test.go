package es

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestEs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Es Suite")
}
