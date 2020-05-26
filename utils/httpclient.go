package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/baotingfang/go-pivnet-client/vlog"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	JsonContentType = "application/json"
	ApiVersion      = "/api/v2"
)

type PivnetHttpClient struct {
	BaseUrl    string
	UaaToken   string
	HttpClient *retryablehttp.Client

	AccessToken            string
	AccessTokenExpiredTime time.Time
}

func NewPivnetHttpClient(baseUrl, UaaToken string) *PivnetHttpClient {
	httpClient := retryablehttp.NewClient()

	httpClient.RetryWaitMin = 1 * time.Second
	httpClient.RetryWaitMax = 10 * time.Second
	httpClient.RetryMax = 10
	httpClient.Logger = vlog.Log

	pivnetClient := &PivnetHttpClient{
		BaseUrl:    UrlJoin(baseUrl, ApiVersion),
		UaaToken:   UaaToken,
		HttpClient: httpClient,
	}

	if err := pivnetClient.RefreshAccessToken(true); err != nil {
		vlog.Fatal(err)
	}
	return pivnetClient
}

func (p *PivnetHttpClient) Post(apiEndPoint string, body io.Reader) (response *http.Response, error error) {
	request, err := p.generateRequest(apiEndPoint, "POST", body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("content-type", JsonContentType)
	response, error = p.Do(request)
	return
}

func (p *PivnetHttpClient) Get(apiEndPoint string) (response *http.Response, error error) {
	request, err := p.generateRequest(apiEndPoint, "GET", nil)
	if err != nil {
		return nil, err
	}
	response, error = p.Do(request)
	return
}

func (p *PivnetHttpClient) Delete(apiEndPoint string) (response *http.Response, error error) {
	request, err := p.generateRequest(apiEndPoint, "DELETE", nil)
	if err != nil {
		return nil, err
	}
	response, error = p.Do(request)
	return
}

func (p *PivnetHttpClient) Patch(apiEndPoint string, body io.Reader) (response *http.Response, error error) {
	request, err := p.generateRequest(apiEndPoint, "PATCH", body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("content-type", JsonContentType)
	response, error = p.Do(request)
	return
}

func (p *PivnetHttpClient) Do(request *http.Request) (response *http.Response, error error) {
	if err := p.RefreshAccessToken(false); err != nil {
		vlog.Fatal(err)
	}
	req, err := retryablehttp.FromRequest(request)
	if err != nil {
		return nil, err
	}
	response, error = p.HttpClient.Do(req)
	return
}

func (p *PivnetHttpClient) generateRequest(apiEndPoint string, httpMethod string, body io.Reader) (request *http.Request, err error) {
	endPointUrl := UrlJoin(p.BaseUrl, apiEndPoint)
	request, err = http.NewRequest(httpMethod, endPointUrl, body)
	if err != nil {
		return nil, err
	}
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

	resp, err := p.HttpClient.Post(UrlJoin(p.BaseUrl, "/authentication/access_tokens"), JsonContentType, bytes.NewReader(payload))
	if err != nil {
		return err
	}

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusOK {
		r := make(map[string]string)
		err = json.Unmarshal(responseData, &r)
		if err != nil {
			return err
		}
		// Pivnet access token will be expired after 1 hour
		p.AccessTokenExpiredTime = time.Now().Add(59 * time.Minute)
		p.AccessToken = r["access_token"]
	} else {
		return errors.New(fmt.Sprintf("Error Code: %d Error Message%s", resp.StatusCode, responseData))
	}

	return nil
}
