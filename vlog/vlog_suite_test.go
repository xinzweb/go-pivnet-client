package vlog_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestVlog(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Vlog Suite")
}
