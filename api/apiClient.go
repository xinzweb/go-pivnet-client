package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/baotingfang/go-pivnet-client/config"
	"github.com/baotingfang/go-pivnet-client/utils"
	"github.com/baotingfang/go-pivnet-client/vlog"
	semver "github.com/cppforlife/go-semi-semantic/version"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"sort"
	"strconv"
	"time"
)

// Struct
var DefaultClient = newDefaultClient()

const (
	BaseUrlKey  = "PIVNET_ENDPOINT"
	SlugKey     = "PIVNET_PRODUCT_SLUG"
	UaaTokenKey = "PIVNET_REFRESH_TOKEN"
)

func newDefaultClient() AccessInterface {
	baseUrl, ex := os.LookupEnv(BaseUrlKey)
	if !ex {
		vlog.Fatal("The env variable %s is not set.\n", BaseUrlKey)
	}
	slug, ex := os.LookupEnv(SlugKey)
	if !ex {
		vlog.Fatal("The env variable %s is not set.\n", SlugKey)
	}
	uaaToken, ex := os.LookupEnv(UaaTokenKey)
	if !ex {
		vlog.Fatal("The env variable %s is not set.\n", UaaTokenKey)
	}

	return NewApiClient(baseUrl, slug, uaaToken)
}

type Request struct {
	EndPoint string
	Payload  io.Reader
}

type ReleaseResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Release struct {
		Id                     int64      `json:"id"`
		Version                string     `json:"version"`
		ReleaseType            string     `json:"release_type"`
		ReleaseDate            utils.Date `json:"release_date"`
		ReleaseNotesUrl        string     `json:"release_notes_url"`
		Availability           string     `json:"availability"`
		Description            string     `json:"description"`
		EndOfSupportDate       utils.Date `json:"end_of_support_date"`
		EndOfGuidanceDate      utils.Date `json:"end_of_guidance_date"`
		EndOfAvailabilityDate  utils.Date `json:"end_of_availability_date"`
		Eccn                   string     `json:"eccn"`
		LicenseException       string     `json:"license_exception"`
		UpdatedAt              time.Time  `json:"updated_at"`
		SoftwareFilesUpdatedAt time.Time  `json:"software_files_updated_at"`
		BecameGaAt             time.Time  `json:"became_ga_at"`
		Eula                   struct {
			Id   int64  `json:"Id"`
			Slug string `json:"slug"`
			Name string `json:"name"`
		} `json:"eula"`
		ProductFiles []struct {
			Id           int64  `json:"id"`
			AwsObjectKey string `json:"aws_object_key"`
			FileType     string `json:"file_type"`
			FileVersion  string `json:"file_version"`
			Md5          string `json:"md5"`
			Sha256       string `json:"sha256"`
			Name         string `json:"name"`
		} `json:"product_files"`
		FileGroups []struct {
			Id           int64  `json:"id"`
			Name         string `json:"name"`
			ProductFiles []struct {
				Id           int64  `json:"id"`
				AwsObjectKey string `json:"aws_object_key"`
				FileType     string `json:"file_type"`
				FileVersion  string `json:"file_version"`
				Md5          string `json:"md5"`
				Sha256       string `json:"sha256"`
				Name         string `json:"name"`
			}
		} `json:"file_groups"`
	} `json:"release"`
}

type ReleaseItem struct {
	Id                     int64      `json:"id"`
	Version                string     `json:"version"`
	ReleaseType            string     `json:"release_type"`
	ReleaseDate            utils.Date `json:"release_date"`
	ReleaseNotesUrl        string     `json:"release_notes_url"`
	Availability           string     `json:"availability"`
	Description            string     `json:"description"`
	EndOfSupportDate       utils.Date `json:"end_of_support_date"`
	EndOfGuidanceDate      utils.Date `json:"end_of_guidance_date"`
	EndOfAvailabilityDate  utils.Date `json:"end_of_availability_date"`
	Eccn                   string     `json:"eccn"`
	LicenseException       string     `json:"license_exception"`
	UpdatedAt              time.Time  `json:"updated_at"`
	SoftwareFilesUpdatedAt time.Time  `json:"software_files_updated_at"`
	BecameGaAt             time.Time  `json:"became_ga_at"`
	Eula                   struct {
		Id   int64  `json:"Id"`
		Slug string `json:"slug"`
		Name string `json:"name"`
	} `json:"eula"`
}

type ReleaseArrayResponse struct {
	Releases []ReleaseItem `json:"releases"`
}

type GroupResponse struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type FederationTokenResponse struct {
	AccessKeyId     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
	SessionToken    string `json:"session_token"`
	Bucket          string `json:"bucket"`
	Region          string `json:"region"`
}

type ProductFileResponse struct {
	ProductFile struct {
		Id                 int64    `json:"id"`
		AwsObjectKey       string   `json:"aws_object_key"`
		Description        string   `json:"description"`
		DocsUrl            string   `json:"docs_url"`
		FileTransferStatus string   `json:"file_transfer_status"`
		FileType           string   `json:"file_type"`
		FileVersion        string   `json:"file_version"`
		IncludeFiles       []string `json:"include_files"`
		Md5                string   `json:"md5"`
		Sha256             string   `json:"sha256"`
		Name               string   `json:"name"`
		ReadyToServe       bool     `json:"ready_to_serve"`
		ReleasedAt         string   `json:"released_at"`
		Size               int64    `json:"Size"`
		SystemRequirements []string `json:"system_requirements"`
	} `json:"product_file"`
}

func MetadataFrom(reader io.Reader, gpdbVersion string) (*config.Metadata, error) {

	var metadata config.Metadata
	if err := yaml.NewDecoder(reader).Decode(&metadata); err != nil {
		return &config.Metadata{}, err
	}

	v, err := semver.NewVersionFromString(gpdbVersion)
	if err != nil {
		return &config.Metadata{}, err
	}

	metadata.Release.Version = v
	vlog.Info("GPDB Version: %s", v.Release.AsString())

	majorVersion, err := strconv.Atoi(v.Release.Components[0].AsString())
	if err != nil {
		vlog.Error("covert gpdb major version failed")
		return &config.Metadata{}, err
	}
	metadata.Release.MajorVersion = majorVersion

	minorVersion, err := strconv.Atoi(v.Release.Components[1].AsString())
	if err != nil {
		vlog.Error("covert gpdb minor version failed")
		return &config.Metadata{}, err
	}
	metadata.Release.MinorVersion = minorVersion

	patchVersion, err := strconv.Atoi(v.Release.Components[2].AsString())
	if err != nil {
		vlog.Error("covert gpdb patch version failed")
		return &config.Metadata{}, err
	}
	metadata.Release.PatchVersion = patchVersion

	r := &metadata.Release
	r.Id = -1

	releaseType, err := r.ComputeReleaseType()
	if err != nil {
		return &config.Metadata{}, err
	}
	metadata.Release.ReleaseType = releaseType

	previousMinorRelease, err := DefaultClient.GetPreviousRelease(r.MajorVersion, config.MinorRelease.String())
	if err != nil {
		vlog.Info("can not find previous minor release: %s", err)
		previousMinorRelease = nil
	}

	previousMajorRelease, err := DefaultClient.GetPreviousRelease(r.MajorVersion, config.MajorRelease.String())
	if err != nil {
		vlog.Info("can not find previous major release: %s", err)
		previousMajorRelease = nil
	}
	r.PreviousMinorRelease = previousMinorRelease
	r.PreviousMajorRelease = previousMajorRelease

	endOfSupportDate, err := r.ComputeEndOfSupportDate()
	if err != nil {
		return &config.Metadata{}, err
	}
	metadata.Release.EndOfSupportDate = endOfSupportDate

	metadata.Release.EndOfGuidanceDate = r.ComputeEndOfGuidanceDate()
	metadata.Release.EndOfAvailabilityDate = r.ComputeEndOfAvailabilityDate()

	return &metadata, nil
}

// Api interface
//go:generate counterfeiter . AccessInterface

type AccessInterface interface {
	CreateRelease(release config.Release) (resp *ReleaseResponse, err error)
	CreateFileGroup(fileGroup config.FileGroup) (resp *GroupResponse, err error)
	CreateFederationToken() (resp *FederationTokenResponse, err error)
	CreateProductFile(productFile config.ProductFile) (resp *ProductFileResponse, err error)
	DeleteProductFile(productFileId int64) (resp *ProductFileResponse, err error)
	AddProductFileToFileGroup(productFileId, fileGroupId int64) (err error)
	AddProductFileToRelease(productFileId, releaseId int64) (err error)
	GetAllReleases() (resp *ReleaseArrayResponse, err error)
	GetPreviousRelease(gpdbMajorVersion int, releaseType string) (release *config.Release, err error)
	UpdateRelease(release config.Release) (resp *ReleaseResponse, err error)
	GetProductFile(productFileId int64) (resp *ProductFileResponse, err error)
	IsProductFileTransferInProgress(productFileId int64) bool
}

type Client struct {
	ProductSlug   string
	UaaFreshToken string
	pivnetClient  utils.HttpClient
}

func NewApiClient(pivnetBaseUrl, productSlug, uaaFreshToken string) AccessInterface {
	return &Client{
		ProductSlug:   productSlug,
		UaaFreshToken: uaaFreshToken,
		pivnetClient:  utils.NewPivnetHttpClient(pivnetBaseUrl, uaaFreshToken),
	}
}

func (c *Client) CreateRelease(release config.Release) (resp *ReleaseResponse, err error) {
	payload, err := PayloadFromRelease(release)
	if err != nil {
		return &ReleaseResponse{}, err
	}

	endPoint := fmt.Sprintf("/products/%s/releases", c.ProductSlug)
	request, err := CreateRequest(payload, endPoint)
	if err != nil {
		return &ReleaseResponse{}, err
	}

	responseData, err := c.pivnetClient.Post(request.EndPoint, request.Payload)
	if err != nil {
		return &ReleaseResponse{}, err
	}

	resp = &ReleaseResponse{}
	err = json.Unmarshal(responseData, resp)
	return
}

func (c *Client) CreateFileGroup(fileGroup config.FileGroup) (resp *GroupResponse, err error) {
	payload := struct {
		FileGroup struct {
			Name string `json:"name"`
		} `json:"file_group"`
	}{}
	payload.FileGroup.Name = fileGroup.Name

	endPoint := fmt.Sprintf("/products/%s/file_groups", c.ProductSlug)
	request, err := CreateRequest(payload, endPoint)
	if err != nil {
		return &GroupResponse{}, err
	}

	responseData, err := c.pivnetClient.Post(request.EndPoint, request.Payload)
	if err != nil {
		return &GroupResponse{}, err
	}

	resp = &GroupResponse{}
	err = json.Unmarshal(responseData, resp)
	return
}

func (c *Client) CreateFederationToken() (resp *FederationTokenResponse, err error) {
	payload := struct {
		ProductId string `json:"product_id"`
	}{}
	payload.ProductId = c.ProductSlug

	endPoint := "/federation_token"
	request, err := CreateRequest(payload, endPoint)
	if err != nil {
		return &FederationTokenResponse{}, err
	}

	responseData, err := c.pivnetClient.Post(request.EndPoint, request.Payload)
	if err != nil {
		return &FederationTokenResponse{}, err
	}

	resp = &FederationTokenResponse{}
	err = json.Unmarshal(responseData, resp)
	return
}

func (c *Client) CreateProductFile(productFile config.ProductFile) (resp *ProductFileResponse, err error) {
	payload := struct {
		ProductFile struct {
			AwsObjectKey       string   `json:"aws_object_key"`
			Description        string   `json:"description"`
			DocsUrl            string   `json:"docs_url"`
			FileType           string   `json:"file_type"`
			FileVersion        string   `json:"file_version"`
			IncludeFiles       []string `json:"include_files"`
			Sha256             string   `json:"sha256"`
			Name               string   `json:"name"`
			ReleasedAt         string   `json:"released_at"`
			SystemRequirements []string `json:"system_requirements"`
		} `json:"product_file"`
	}{}

	pf := &payload.ProductFile
	// TODO
	pf.AwsObjectKey = ""
	// TODO
	pf.Description = productFile.Description
	pf.DocsUrl = productFile.DocsUrl
	pf.FileType = productFile.FileType
	// TODO
	pf.FileVersion = productFile.FileVersion
	pf.IncludeFiles = productFile.IncludedFiles
	// TODO
	pf.Sha256 = ""
	// TODO
	pf.Name = productFile.UploadAs
	pf.ReleasedAt = time.Now().Format("2006/01/02")
	pf.SystemRequirements = productFile.SystemRequirements

	endPoint := fmt.Sprintf("/products/%s/product_files", c.ProductSlug)
	request, err := CreateRequest(payload, endPoint)
	if err != nil {
		return &ProductFileResponse{}, err
	}

	responseData, err := c.pivnetClient.Post(request.EndPoint, request.Payload)
	if err != nil {
		return &ProductFileResponse{}, err
	}

	resp = &ProductFileResponse{}
	err = json.Unmarshal(responseData, resp)
	return
}

func (c *Client) DeleteProductFile(productFileId int64) (resp *ProductFileResponse, err error) {
	endPoint := fmt.Sprintf("/products/%s/product_files/%d", c.ProductSlug, productFileId)
	request, err := CreateRequest(nil, endPoint)
	if err != nil {
		return &ProductFileResponse{}, err
	}

	responseData, err := c.pivnetClient.Delete(request.EndPoint)
	if err != nil {
		return &ProductFileResponse{}, err
	}

	resp = &ProductFileResponse{}
	err = json.Unmarshal(responseData, resp)
	return
}

func (c *Client) AddProductFileToFileGroup(productFileId, fileGroupId int64) (err error) {
	payload := struct {
		ProductFile struct {
			Id int64 `json:"id"`
		} `json:"product_file"`
	}{}
	payload.ProductFile.Id = productFileId

	endPoint := fmt.Sprintf("/products/%s/file_groups/%d/add_product_file", c.ProductSlug, fileGroupId)
	request, err := CreateRequest(payload, endPoint)
	if err != nil {
		return err
	}

	_, err = c.pivnetClient.Patch(request.EndPoint, request.Payload)
	return
}

func (c *Client) AddProductFileToRelease(productFileId, releaseId int64) (err error) {
	payload := struct {
		ProductFile struct {
			Id int64 `json:"id"`
		} `json:"product_file"`
	}{}
	payload.ProductFile.Id = productFileId

	endPoint := fmt.Sprintf("/products/%s/releases/%d/add_product_file", c.ProductSlug, releaseId)
	request, err := CreateRequest(payload, endPoint)
	if err != nil {
		return err
	}

	_, err = c.pivnetClient.Patch(request.EndPoint, request.Payload)
	return
}

func (c *Client) AddFileGroupToRelease(fileGroupId, releaseId int64) (err error) {
	payload := struct {
		FileGroup struct {
			Id int64 `json:"id"`
		} `json:"file_group"`
	}{}
	payload.FileGroup.Id = fileGroupId

	endPoint := fmt.Sprintf("/products/%s/releases/%d/add_file_group", c.ProductSlug, releaseId)
	request, err := CreateRequest(payload, endPoint)
	if err != nil {
		return err
	}

	_, err = c.pivnetClient.Patch(request.EndPoint, request.Payload)
	return
}

func (c *Client) GetAllReleases() (resp *ReleaseArrayResponse, err error) {
	endPoint := fmt.Sprintf("products/%s/releases", c.ProductSlug)
	responseData, err := c.pivnetClient.Get(endPoint)
	if err != nil {
		return &ReleaseArrayResponse{}, err
	}

	resp = &ReleaseArrayResponse{}
	err = json.Unmarshal(responseData, resp)
	return

}

func (c *Client) GetPreviousRelease(gpdbMajorVersion int, releaseType string) (release *config.Release, err error) {
	allReleasesResponse, err := c.GetAllReleases()
	if err != nil {
		return &config.Release{}, err
	}
	var versions []semver.Version
	versionMap := make(map[string]ReleaseItem)

	for _, release := range allReleasesResponse.Releases {

		if release.Availability != "All Users" {
			continue
		}

		version := semver.MustNewVersionFromString(release.Version)
		majorVersion, err := strconv.Atoi(version.Release.Components[0].AsString())
		if err != nil {
			return &config.Release{}, err
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
			majorVersion, _ := strconv.Atoi(v.Release.Components[0].AsString())
			minorVersion, _ := strconv.Atoi(v.Release.Components[1].AsString())
			patchVersion, _ := strconv.Atoi(v.Release.Components[2].AsString())
			release := config.Release{
				ReleaseType:           previousRelease.ReleaseType,
				EulaSlug:              previousRelease.Eula.Slug,
				Description:           previousRelease.Description,
				ReleaseNotesUrl:       previousRelease.ReleaseNotesUrl,
				Availability:          previousRelease.Availability,
				Controlled:            false,
				Eccn:                  previousRelease.Eccn,
				LicenseException:      previousRelease.LicenseException,
				ReleaseDate:           previousRelease.ReleaseDate,
				EndOfSupportDate:      previousRelease.EndOfSupportDate,
				EndOfGuidanceDate:     previousRelease.EndOfGuidanceDate,
				EndOfAvailabilityDate: previousRelease.EndOfAvailabilityDate,
				Id:                    previousRelease.Id,
				Version:               v,
				MajorVersion:          majorVersion,
				MinorVersion:          minorVersion,
				PatchVersion:          patchVersion,
			}
			return &release, nil
		}
	}

	return &config.Release{},
		fmt.Errorf("can not found previous release. major version: %d, release type: %s",
			gpdbMajorVersion, releaseType)
}

func (c *Client) UpdateRelease(release config.Release) (resp *ReleaseResponse, err error) {
	if release.Id <= 0 {
		return &ReleaseResponse{}, fmt.Errorf("invalid release id: %d", release.Id)
	}

	payload, err := PayloadFromRelease(release)
	if err != nil {
		return &ReleaseResponse{}, err
	}

	endPoint := fmt.Sprintf("/products/%s/releases/%d", c.ProductSlug, release.Id)
	request, err := CreateRequest(payload, endPoint)
	if err != nil {
		return &ReleaseResponse{}, err
	}

	responseData, err := c.pivnetClient.Post(request.EndPoint, request.Payload)
	if err != nil {
		return &ReleaseResponse{}, err
	}

	resp = &ReleaseResponse{}
	err = json.Unmarshal(responseData, resp)
	return
}

func (c *Client) GetProductFile(productFileId int64) (resp *ProductFileResponse, err error) {
	endPoint := fmt.Sprintf("products/%s/product_files/%d", c.ProductSlug, productFileId)
	responseData, err := c.pivnetClient.Get(endPoint)
	if err != nil {
		return &ProductFileResponse{}, err
	}

	resp = &ProductFileResponse{}
	err = json.Unmarshal(responseData, resp)
	return
}

func (c *Client) IsProductFileTransferInProgress(productFileId int64) bool {
	productFile, err := c.GetProductFile(productFileId)
	if err != nil {
		vlog.Fatal("get product file (id=%d) failed in %s slug.\n%v", productFileId, c.ProductSlug, err)
	}
	return productFile.ProductFile.FileTransferStatus == "in_progress"
}

// internal usage: utils

func PayloadFromRelease(r config.Release) ([]byte, error) {
	req := struct {
		CopyMetadata bool `json:"copy_metadata,omitempty"`
		Release      struct {
			Version               string     `json:"version,omitempty"`
			ReleaseNotesUrl       string     `json:"release_notes_url,omitempty"`
			Description           string     `json:"description"`
			ReleaseType           string     `json:"release_type,omitempty"`
			Availability          string     `json:"availability"`
			OssCompliant          string     `json:"oss_compliant,omitempty"`
			ReleaseDate           utils.Date `json:"release_date,omitempty"`
			EndOfSupportDate      utils.Date `json:"end_of_support_date,omitempty"`
			EndOfGuidanceDate     utils.Date `json:"end_of_guidance_date,omitempty"`
			EndOfAvailabilityDate utils.Date `json:"end_of_availability_date,omitempty"`
			Eccn                  string     `json:"eccn"`
			LicenseException      string     `json:"license_exception"`
			Eula                  struct {
				Slug string `json:"slug,omitempty"`
			} `json:"eula,omitempty"`
		} `json:"release,omitempty"`
	}{}

	req.CopyMetadata = false

	release := &req.Release
	release.Version = r.Version.AsString()
	release.ReleaseNotesUrl = r.ReleaseNotesUrl
	release.Description = r.Description
	release.ReleaseDate = r.ReleaseDate
	release.ReleaseType = r.ReleaseType
	release.Availability = r.Availability
	release.Eula.Slug = r.EulaSlug
	release.OssCompliant = "confirm"
	release.EndOfSupportDate = r.EndOfSupportDate
	release.EndOfGuidanceDate = r.EndOfGuidanceDate
	release.EndOfAvailabilityDate = r.EndOfAvailabilityDate
	release.Eccn = r.Eccn
	release.LicenseException = r.LicenseException

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

func CreateRequest(data interface{}, endPoint string) (*Request, error) {
	if utils.IsEmpty(endPoint) {
		return &Request{}, fmt.Errorf("EndPoint is empty")
	}

	if data == nil {
		return &Request{
			EndPoint: endPoint,
			Payload:  nil,
		}, nil
	}
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return &Request{
		EndPoint: endPoint,
		Payload:  bytes.NewReader(payload),
	}, nil
}
