package api

import (
	"fmt"
	"github.com/baotingfang/go-pivnet-client/gp"
	"github.com/baotingfang/go-pivnet-client/vlog"
	"github.com/baotingfang/go-pivnet-client/wrapper"
	semver "github.com/cppforlife/go-semi-semantic/version"
	"github.com/pivotal-cf/go-pivnet/v4"
	"sort"
	"strconv"
)

//go:generate counterfeiter . AccessClient

type AccessClient interface {
	GetAllReleases() ([]pivnet.Release, error)
	GetLatestPublicReleaseByReleaseType(gpdbMajorVersion int, releaseType pivnet.ReleaseType) (release pivnet.Release, err error)
	CreateRelease(releaseConfig pivnet.CreateReleaseConfig) (pivnet.Release, error)
	CreateFileGroup(groupName string) (pivnet.FileGroup, error)
	CreateFederationToken() (pivnet.FederationToken, error)
	CreateProductFile(productFileConfig pivnet.CreateProductFileConfig) (pivnet.ProductFile, error)
	DeleteProductFile(productFileId int) (pivnet.ProductFile, error)
	AddProductFileToFileGroup(productFileId, fileGroupId int) error
	AddProductFileToRelease(productFileId, releaseId int) error
	AddFileGroupToRelease(fileGroupId, releaseId int) error
	UpdateRelease(release pivnet.Release) (pivnet.Release, error)
	FileTransferStatusInProgress(productFileId int) bool
}

type Client struct {
	ProductSlug   string
	UaaFreshToken string
	pivnetClient  wrapper.PivnetClient
}

func NewApiClient(context gp.Context) AccessClient {
	return &Client{
		ProductSlug:   context.Slug,
		UaaFreshToken: context.UaaFreshToken,
		pivnetClient:  context.Client,
	}
}

func (c Client) CreateRelease(releaseConfig pivnet.CreateReleaseConfig) (pivnet.Release, error) {
	return c.pivnetClient.CreateRelease(releaseConfig)
}

func (c Client) CreateFileGroup(groupName string) (pivnet.FileGroup, error) {
	return c.pivnetClient.CreateFileGroup(c.ProductSlug, groupName)
}

func (c Client) CreateFederationToken() (pivnet.FederationToken, error) {
	return c.pivnetClient.CreateFederationToken(c.ProductSlug)
}

func (c Client) CreateProductFile(productFileConfig pivnet.CreateProductFileConfig) (pivnet.ProductFile, error) {
	return c.pivnetClient.CreateProductFile(productFileConfig)
}

func (c Client) DeleteProductFile(productFileId int) (pivnet.ProductFile, error) {
	return c.pivnetClient.DeleteProductFile(c.ProductSlug, productFileId)
}

func (c Client) AddProductFileToFileGroup(productFileId, fileGroupId int) error {
	return c.pivnetClient.AddProductFileToFileGroup(c.ProductSlug, productFileId, fileGroupId)
}

func (c Client) AddProductFileToRelease(productFileId, releaseId int) error {
	return c.pivnetClient.AddProductFileToRelease(c.ProductSlug, productFileId, releaseId)
}

func (c Client) AddFileGroupToRelease(fileGroupId, releaseId int) error {
	return c.pivnetClient.AddFileGroupToRelease(c.ProductSlug, fileGroupId, releaseId)
}

func (c Client) UpdateRelease(release pivnet.Release) (pivnet.Release, error) {
	return c.pivnetClient.UpdateRelease(c.ProductSlug, release)
}

func (c Client) GetAllReleases() ([]pivnet.Release, error) {
	return c.pivnetClient.GetAllReleases(c.ProductSlug)
}

func (c Client) GetLatestPublicReleaseByReleaseType(gpdbMajorVersion int, releaseType pivnet.ReleaseType) (release pivnet.Release, err error) {
	allReleases, err := c.GetAllReleases()
	if err != nil {
		return pivnet.Release{}, err
	}

	var versions []semver.Version
	versionMap := make(map[string]pivnet.Release)

	for _, release := range allReleases {
		if release.Availability != "All Users" {
			continue
		}

		version := semver.MustNewVersionFromString(release.Version)
		majorVersion, err := strconv.Atoi(version.Release.Components[0].AsString())
		if err != nil {
			return pivnet.Release{}, err
		}

		if majorVersion == gpdbMajorVersion {
			versions = append(versions, version)
			versionMap[version.AsString()] = release
		}
	}

	sort.Sort(sort.Reverse(semver.AscSorting(versions)))

	for _, v := range versions {
		if releaseType == versionMap[v.AsString()].ReleaseType {
			previousRelease := versionMap[v.AsString()]
			return previousRelease, nil
		}
	}

	return pivnet.Release{},
		fmt.Errorf("can not found previous release. major version: %d, release type: %s",
			gpdbMajorVersion, releaseType)
}

func (c Client) FileTransferStatusInProgress(productFileId int) bool {
	pf, err := c.pivnetClient.GetProductFile(c.ProductSlug, productFileId)
	if err != nil {
		vlog.Error("can not find product file. id=%d", productFileId)
		vlog.Fatal(err.Error())
	}
	return pf.FileTransferStatus == "in_progress"
}
