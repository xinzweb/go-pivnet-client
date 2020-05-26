package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/baotingfang/go-pivnet-client/vlog"
	"github.com/hashicorp/go-retryablehttp"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	JsonContentType = "application/json"
	ApiVersion      = "/api/v2"
)

//go:generate counterfeiter . RetryHttpClient

type RetryHttpClient interface {
	Do(req *retryablehttp.Request) (*http.Response, error)
}

//go:generate counterfeiter . HttpClient

type HttpClient interface {
	Post(apiEndPoint string, body io.Reader) (response []byte, error error)
	Get(apiEndPoint string) (response []byte, error error)
	Delete(apiEndPoint string) (response []byte, error error)
	Patch(apiEndPoint string, body io.Reader) (response []byte, error error)
	Do(request *http.Request) (response []byte, error error)
	RefreshAccessToken(force bool) error
}

type PivnetHttpClient struct {
	BaseUrl    string
	UaaToken   string
	HttpClient RetryHttpClient

	AccessToken            string
	AccessTokenExpiredTime time.Time
}

func NewPivnetHttpClient(baseUrl, UaaToken string) *PivnetHttpClient {
	httpClient := retryablehttp.NewClient()

	httpClient.RetryWaitMin = 1 * time.Second
	httpClient.RetryWaitMax = 10 * time.Second
	httpClient.RetryMax = 10
	httpClient.Logger = vlog.Log

	return &PivnetHttpClient{
		BaseUrl:    baseUrl,
		UaaToken:   UaaToken,
		HttpClient: httpClient,
	}
}

func (p *PivnetHttpClient) Post(apiEndPoint string, body io.Reader) (response []byte, error error) {
	if err := p.RefreshAccessToken(false); err != nil {
		vlog.Fatal(err.Error())
	}
	request, err := p.GenerateRequest(apiEndPoint, "POST", body)
	if err != nil {
		return nil, err
	}
	response, error = p.Do(request)
	return
}

func (p *PivnetHttpClient) Get(apiEndPoint string) (response []byte, error error) {
	if err := p.RefreshAccessToken(false); err != nil {
		vlog.Fatal(err.Error())
	}
	request, err := p.GenerateRequest(apiEndPoint, "GET", nil)
	if err != nil {
		return nil, err
	}
	response, error = p.Do(request)
	return
}

func (p *PivnetHttpClient) Delete(apiEndPoint string) (response []byte, error error) {
	if err := p.RefreshAccessToken(false); err != nil {
		vlog.Fatal(err.Error())
	}
	request, err := p.GenerateRequest(apiEndPoint, "DELETE", nil)
	if err != nil {
		return nil, err
	}
	response, error = p.Do(request)
	return
}

func (p *PivnetHttpClient) Patch(apiEndPoint string, body io.Reader) (response []byte, error error) {
	if err := p.RefreshAccessToken(false); err != nil {
		vlog.Fatal(err.Error())
	}
	request, err := p.GenerateRequest(apiEndPoint, "PATCH", body)
	if err != nil {
		return nil, err
	}
	response, error = p.Do(request)
	return
}

func (p *PivnetHttpClient) Do(request *http.Request) (response []byte, error error) {
	req, err := retryablehttp.FromRequest(request)
	if err != nil {
		return nil, err
	}
	resp, err := p.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("invoke api failed! error code: %d\n\nrequest: \n\t%v\n\nresponse:\n\t%v\n", resp.StatusCode, request.URL, string(responseData))
	}

	return responseData, nil
}

func (p *PivnetHttpClient) GenerateRequest(apiEndPoint string, httpMethod string, body io.Reader) (request *http.Request, err error) {
	endPointUrl := UrlJoin(p.BaseUrl, ApiVersion, apiEndPoint)
	request, err = http.NewRequest(httpMethod, endPointUrl, body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("content-type", JsonContentType)
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", p.AccessToken))
	return
}

func (p *PivnetHttpClient) RefreshAccessToken(force bool) error {
	if p.AccessTokenExpiredTime.After(time.Now()) && len(p.AccessToken) != 0 && !force {
		return nil
	}
	payload, err := json.Marshal(map[string]string{
		"refresh_token": p.UaaToken,
	})
	if err != nil {
		return err
	}

	endPoint := "/authentication/access_tokens"
	request, err := p.GenerateRequest(endPoint, "POST", bytes.NewReader(payload))
	if err != nil {
		return err
	}
	responseData, err := p.Do(request)
	if err != nil {
		return err
	}

	r := make(map[string]string)
	err = json.Unmarshal(responseData, &r)
	if err != nil {
		return err
	}
	// Pivnet access token will be expired after 1 hour
	p.AccessTokenExpiredTime = time.Now().Add(59 * time.Minute)
	p.AccessToken = r["access_token"]

	return nil
}
