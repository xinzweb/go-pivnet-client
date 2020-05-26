package wrapper

import (
	"github.com/pivotal-cf/go-pivnet/v4"
	"github.com/pivotal-cf/go-pivnet/v4/logger"
)

//go:generate counterfeiter . AccessTokenService

type AccessTokenService interface {
	AccessToken() (string, error)
}

//go:generate counterfeiter . PivnetClient

type PivnetClient interface {
	GetAllReleases(productSlug string) ([]pivnet.Release, error)
	CreateRelease(releaseConfig pivnet.CreateReleaseConfig) (pivnet.Release, error)
	CreateFileGroup(productSlug, groupName string) (pivnet.FileGroup, error)
	CreateFederationToken(productSlug string) (pivnet.FederationToken, error)
	GetProductFile(productSlug string, productFileId int) (pivnet.ProductFile, error)
	CreateProductFile(productFileConfig pivnet.CreateProductFileConfig) (pivnet.ProductFile, error)
	DeleteProductFile(productSlug string, productFileId int) (pivnet.ProductFile, error)
	AddProductFileToFileGroup(productSlug string, productFileId, fileGroupId int) error
	AddProductFileToRelease(productSlug string, productFileId, releaseId int) error
	AddFileGroupToRelease(productSlug string, fileGroupId, releaseId int) error
	UpdateRelease(productSlug string, release pivnet.Release) (pivnet.Release, error)
}

type Client struct {
	client pivnet.Client
}

func NewClient(token AccessTokenService, config pivnet.ClientConfig, logger logger.Logger) PivnetClient {
	return &Client{
		client: pivnet.NewClient(token, config, logger),
	}
}

func (c Client) GetAllReleases(productSlug string) ([]pivnet.Release, error) {
	return c.client.Releases.List(productSlug)
}

func (c Client) CreateRelease(releaseConfig pivnet.CreateReleaseConfig) (pivnet.Release, error) {
	return c.client.Releases.Create(releaseConfig)
}

func (c Client) CreateFileGroup(productSlug, groupName string) (pivnet.FileGroup, error) {
	return c.client.FileGroups.Create(pivnet.CreateFileGroupConfig{ProductSlug: productSlug, Name: groupName})
}

func (c Client) CreateFederationToken(productSlug string) (pivnet.FederationToken, error) {
	return c.client.FederationToken.GenerateFederationToken(productSlug)
}

func (c Client) GetProductFile(productSlug string, productFileId int) (pivnet.ProductFile, error) {
	return c.client.ProductFiles.Get(productSlug, productFileId)
}

func (c Client) CreateProductFile(productFileConfig pivnet.CreateProductFileConfig) (pivnet.ProductFile, error) {
	return c.client.ProductFiles.Create(productFileConfig)
}

func (c Client) DeleteProductFile(productSlug string, productFileId int) (pivnet.ProductFile, error) {
	return c.client.ProductFiles.Delete(productSlug, productFileId)
}

func (c Client) AddProductFileToFileGroup(productSlug string, productFileId, fileGroupId int) error {
	return c.client.ProductFiles.AddToFileGroup(productSlug, fileGroupId, productFileId)
}

func (c Client) AddProductFileToRelease(productSlug string, productFileId, releaseId int) error {
	return c.client.ProductFiles.AddToRelease(productSlug, releaseId, productFileId)
}

func (c Client) AddFileGroupToRelease(productSlug string, fileGroupId, releaseId int) error {
	return c.client.FileGroups.AddToRelease(productSlug, releaseId, fileGroupId)
}

func (c Client) UpdateRelease(productSlug string, release pivnet.Release) (pivnet.Release, error) {
	return c.client.Releases.Update(productSlug, release)
}
