package gremlin

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGoGremlin(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GoGremlin Suite")
}
