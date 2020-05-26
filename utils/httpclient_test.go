package utils_test

import (
	"github.com/baotingfang/go-pivnet-client/vlog"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"strings"

	. "github.com/baotingfang/go-pivnet-client/utils"
)

var _ = Describe("Httpclient", func() {
	It("HttpClient api version", func() {
		client := NewPivnetHttpClient("https://www.example.com", "123")
		Expect(client.BaseUrl).To(Equal("https://www.example.com/api/v2"))
	})

	It("HttpClient log init", func() {
		vlog.InitLog("test", vlog.DebugLevel)
		client := NewPivnetHttpClient("https://www.example.com", "123")
		Expect(client.BaseUrl).To(Equal("https://www.example.com/api/v2"))
		Expect(client.UaaToken).To(Equal("123"))
		logger, r := client.HttpClient.Logger.(*vlog.Logger)
		Expect(r).To(BeTrue())
		Expect(logger.LogLevel).To(Equal(vlog.DebugLevel))
	})

	It("HttpClient log init", func() {
		client := NewPivnetHttpClient("https://www.example.com", "123")
		Expect(client.BaseUrl).To(Equal("https://www.example.com/api/v2"))
		Expect(client.UaaToken).To(Equal("123"))
		logger, r := client.HttpClient.Logger.(*vlog.Logger)
		Expect(r).To(BeTrue())
		Expect(logger.LogLevel).To(Equal(vlog.DebugLevel))
	})

	It("Generate http request", func() {
		client := NewPivnetHttpClient("https://www.example.com", "123")
		client.AccessToken = "XYZ"
		r, err := client.GenerateRequest("abc/def", "POST", strings.NewReader("payload"))
		Expect(err).NotTo(HaveOccurred())
		Expect(r.URL.String()).To(Equal("https://www.example.com/api/v2/abc/def"))

		payload, _ := ioutil.ReadAll(r.Body)
		Expect(string(payload)).To(Equal("payload"))
		Expect(r.Method).To(Equal("POST"))
		Expect(r.Header.Get("content-type")).To(Equal(JsonContentType))
		Expect(r.Header.Get("Authorization")).To(Equal("Bearer XYZ"))
	})
})
