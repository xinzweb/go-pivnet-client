package utils_test

import (
	. "github.com/baotingfang/go-pivnet-client/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tools", func() {
	It("UrlJoin", func() {
		baseUrl := "http://www.example.com/"
		p := UrlJoin(baseUrl, "/abc", "/def/", "/xyz")
		Expect(p).To(Equal("http://www.example.com/abc/def/xyz"))
	})
})
