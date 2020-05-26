package gp

import (
	"github.com/baotingfang/go-pivnet-client/wrapper"
	"github.com/pivotal-cf/go-pivnet/v4"
	"github.com/pivotal-cf/go-pivnet/v4/logshim"
	"log"
	"os"
)

type Context struct {
	BaseUrl           string
	Slug              string
	UaaFreshToken     string
	UserAgent         string
	SkipSSLValidation bool
	Verbose           bool
	Client            wrapper.PivnetClient
}

func NewContext(pivnetBaseUrl, productSlug, uaaFreshToken string, skipSSLValidation bool, verbose bool) Context {
	userAgent := "RelEng Release Tools"

	tokenService := pivnet.NewAccessTokenOrLegacyToken(uaaFreshToken, pivnetBaseUrl, skipSSLValidation, userAgent)

	pivnetConfig := pivnet.ClientConfig{
		Host:              pivnetBaseUrl,
		UserAgent:         userAgent,
		SkipSSLValidation: skipSSLValidation,
	}

	stdoutLogger := log.New(os.Stdout, "[apiClient]", log.LstdFlags)
	stderrLogger := log.New(os.Stderr, "[apiClient]", log.LstdFlags)
	logger := logshim.NewLogShim(stdoutLogger, stderrLogger, verbose)

	pivnetClient := wrapper.NewClient(tokenService, pivnetConfig, logger)

	return Context{
		BaseUrl:           pivnetBaseUrl,
		Slug:              productSlug,
		UaaFreshToken:     uaaFreshToken,
		UserAgent:         userAgent,
		SkipSSLValidation: skipSSLValidation,
		Verbose:           verbose,

		Client: pivnetClient,
	}
}
