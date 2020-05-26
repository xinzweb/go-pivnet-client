package utils_test

import (
	"bytes"
	. "github.com/baotingfang/go-pivnet-client/utils"
	"github.com/baotingfang/go-pivnet-client/utils/utilsfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var _ = Describe("Httpclient", func() {
	It("HttpClient api version", func() {
		client := NewPivnetHttpClient("https://www.example.com", "123")
		Expect(client.BaseUrl).To(Equal("https://www.example.com"))
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

var _ = Describe("Test http methods", func() {
	var retryClient *utilsfakes.FakeRetryHttpClient
	var client *PivnetHttpClient

	BeforeEach(func() {
		retryClient = &utilsfakes.FakeRetryHttpClient{}
		client = &PivnetHttpClient{
			BaseUrl:                "https://www.fakepivnet.pivotal.io",
			UaaToken:               "123456",
			HttpClient:             retryClient,
			AccessToken:            "abcdef",
			AccessTokenExpiredTime: time.Time{},
		}
	})

	It("Test Post Method", func() {
		responseForAccessToken := &http.Response{
			Status:     "200 OK",
			StatusCode: 200,
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Body:       ioutil.NopCloser(bytes.NewBufferString(`{"access_token":"19491001"}`)),
		}

		responseForPost := &http.Response{
			Status:     "200 OK",
			StatusCode: 200,
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Body:       ioutil.NopCloser(bytes.NewBufferString(`{"data":"abc"}`)),
		}

		retryClient.DoReturnsOnCall(0, responseForAccessToken, nil)
		retryClient.DoReturnsOnCall(1, responseForPost, nil)

		resp, err := client.Post("/post/url", bytes.NewReader([]byte(`{"a":"1""}`)))
		Expect(err).NotTo(HaveOccurred())
		Expect(string(resp)).To(Equal(`{"data":"abc"}`))
		Expect(client.AccessToken).To(Equal("19491001"))
	})

	It("Test Get Method", func() {
		responseForAccessToken := &http.Response{
			Status:     "200 OK",
			StatusCode: 200,
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Body:       ioutil.NopCloser(bytes.NewBufferString(`{"access_token":"123456"}`)),
		}

		responseForPost := &http.Response{
			Status:     "200 OK",
			StatusCode: 200,
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Body:       ioutil.NopCloser(bytes.NewBufferString(`{"data":"abc"}`)),
		}

		retryClient.DoReturnsOnCall(0, responseForAccessToken, nil)
		retryClient.DoReturnsOnCall(1, responseForPost, nil)

		resp, err := client.Get("/get/url")
		Expect(err).NotTo(HaveOccurred())
		Expect(string(resp)).To(Equal(`{"data":"abc"}`))
		Expect(client.AccessToken).To(Equal("123456"))
	})

	It("Test Delete Method", func() {
		responseForAccessToken := &http.Response{
			Status:     "200 OK",
			StatusCode: 200,
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Body:       ioutil.NopCloser(bytes.NewBufferString(`{"access_token":"123456"}`)),
		}

		responseForPost := &http.Response{
			Status:     "200 OK",
			StatusCode: 200,
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Body:       ioutil.NopCloser(bytes.NewBufferString(`{"data":"abc"}`)),
		}

		retryClient.DoReturnsOnCall(0, responseForAccessToken, nil)
		retryClient.DoReturnsOnCall(1, responseForPost, nil)

		resp, err := client.Delete("/delete/url")
		Expect(err).NotTo(HaveOccurred())
		Expect(string(resp)).To(Equal(`{"data":"abc"}`))
		Expect(client.AccessToken).To(Equal("123456"))
	})

	It("Test Patch Method", func() {
		responseForAccessToken := &http.Response{
			Status:     "200 OK",
			StatusCode: 200,
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Body:       ioutil.NopCloser(bytes.NewBufferString(`{"access_token":"123456"}`)),
		}

		responseForPost := &http.Response{
			Status:     "200 OK",
			StatusCode: 200,
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Body:       ioutil.NopCloser(bytes.NewBufferString(`{"data":"abc"}`)),
		}

		retryClient.DoReturnsOnCall(0, responseForAccessToken, nil)
		retryClient.DoReturnsOnCall(1, responseForPost, nil)

		resp, err := client.Patch("/patch/url", bytes.NewReader([]byte(`{"a":"1""}`)))
		Expect(err).NotTo(HaveOccurred())
		Expect(string(resp)).To(Equal(`{"data":"abc"}`))
		Expect(client.AccessToken).To(Equal("123456"))
	})

	It("Test Do Method", func() {
		responseForPost := &http.Response{
			Status:     "200 OK",
			StatusCode: 200,
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Body:       ioutil.NopCloser(bytes.NewBufferString(`{"data":"abc"}`)),
		}

		retryClient.DoReturnsOnCall(0, responseForPost, nil)

		req, err := client.GenerateRequest("/do/url", "POST", bytes.NewReader([]byte(`{"a":"1""}`)))
		Expect(err).NotTo(HaveOccurred())
		Expect(req.Host).To(Equal("www.fakepivnet.pivotal.io"))
		Expect(req.URL.String()).To(Equal("https://www.fakepivnet.pivotal.io/api/v2/do/url"))

		resp, err := client.Do(req)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(resp)).To(Equal(`{"data":"abc"}`))
	})

	It("Test Refresh access token Method", func() {
		responseForAccessToken1 := &http.Response{
			Status:     "200 OK",
			StatusCode: 200,
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Body:       ioutil.NopCloser(bytes.NewBufferString(`{"access_token":"123456"}`)),
		}

		responseForAccessToken2 := &http.Response{
			Status:     "200 OK",
			StatusCode: 200,
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Body:       ioutil.NopCloser(bytes.NewBufferString(`{"access_token":"7890"}`)),
		}

		retryClient.DoReturnsOnCall(0, responseForAccessToken1, nil)
		retryClient.DoReturnsOnCall(1, responseForAccessToken2, nil)

		err := client.RefreshAccessToken(false)
		Expect(err).NotTo(HaveOccurred())
		Expect(client.AccessToken).To(Equal("123456"))
		Expect(client.AccessTokenExpiredTime.IsZero()).To(BeFalse())

		firstExpiredTime := client.AccessTokenExpiredTime
		err = client.RefreshAccessToken(false)
		Expect(err).NotTo(HaveOccurred())
		Expect(client.AccessToken).To(Equal("123456"))
		Expect(client.AccessTokenExpiredTime).To(Equal(firstExpiredTime))

		err = client.RefreshAccessToken(true)
		Expect(err).NotTo(HaveOccurred())
		Expect(client.AccessToken).To(Equal("7890"))
		Expect(client.AccessTokenExpiredTime).NotTo(Equal(firstExpiredTime))
	})

	It("Client error", func() {
		responseForPost := &http.Response{
			Status:     "401 ERROR",
			StatusCode: 401,
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Body:       ioutil.NopCloser(bytes.NewBufferString(`{"data":"abc"}`)),
		}

		retryClient.DoReturnsOnCall(0, responseForPost, nil)

		req, err := client.GenerateRequest("/do/url", "POST", bytes.NewReader([]byte(`{"a":"1""}`)))
		Expect(err).NotTo(HaveOccurred())
		Expect(req.Host).To(Equal("www.fakepivnet.pivotal.io"))
		Expect(req.URL.String()).To(Equal("https://www.fakepivnet.pivotal.io/api/v2/do/url"))

		resp, err := client.Do(req)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal(`invoke api failed! error code: 401

request: 
	https://www.fakepivnet.pivotal.io/api/v2/do/url

response:
	{"data":"abc"}
`))
		Expect(resp).To(BeNil())
	})

	It("Server error", func() {
		responseForPost := &http.Response{
			Status:     "500 ERROR",
			StatusCode: 500,
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Body:       ioutil.NopCloser(bytes.NewBufferString(`{"data":"abc"}`)),
		}

		retryClient.DoReturnsOnCall(0, responseForPost, nil)

		req, err := client.GenerateRequest("/do/url", "POST", bytes.NewReader([]byte(`{"a":"1""}`)))
		Expect(err).NotTo(HaveOccurred())
		Expect(req.Host).To(Equal("www.fakepivnet.pivotal.io"))
		Expect(req.URL.String()).To(Equal("https://www.fakepivnet.pivotal.io/api/v2/do/url"))

		resp, err := client.Do(req)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal(`invoke api failed! error code: 500

request: 
	https://www.fakepivnet.pivotal.io/api/v2/do/url

response:
	{"data":"abc"}
`))
		Expect(resp).To(BeNil())
	})
})
