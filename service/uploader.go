package service

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/baotingfang/go-pivnet-client/api"
	"github.com/baotingfang/go-pivnet-client/config"
	"github.com/baotingfang/go-pivnet-client/gp"
	. "github.com/baotingfang/go-pivnet-client/utils"
	"github.com/baotingfang/go-pivnet-client/vlog"
	"github.com/pivotal-cf/go-pivnet/v4"
	"io"
	"os"
	"path"
	"strings"
)

type Uploader struct {
	GpdbVersion     string
	Metadata        config.Metadata
	SearchPath      string
	AwsObjectPrefix string

	Context  gp.Context
	Client   api.AccessClient
	Resolver Resolver
}

func NewUploader(context gp.Context, gpdbVersion string, metadataReader io.Reader, searchPath string) (Uploader, error) {
	metadata, err := config.MetadataFrom(metadataReader, gpdbVersion)
	if err != nil {
		return Uploader{}, err
	}
	metadata.Release.Version = gpdbVersion
	client := api.NewApiClient(context)

	return Uploader{
		GpdbVersion: gpdbVersion,
		Metadata:    metadata,
		SearchPath:  searchPath,

		Context:  context,
		Client:   client,
		Resolver: NewResourceResolver(searchPath, FilesWalker{}),
	}, nil
}

func (u Uploader) Run() error {
	r := u.Metadata.Release

	mv := NewMetaDataValidator(u.Metadata)
	if !mv.Validate() {
		messages := mv.errorMessages
		fmt.Println(strings.Join(messages, "\n"))
		return fmt.Errorf("validate metata data failed")
	}

	crc, err := u.NewCreateReleaseConfig(r)
	if err != nil {
		return err
	}

	release, err := u.Client.CreateRelease(crc)
	if err != nil {
		return err
	}

	federationToken, err := u.Client.CreateFederationToken()
	if err != nil {
		return err
	}

	err = u.HandleFileGroups(release, federationToken)
	if err != nil {
		return err
	}

	err = u.HandleProductFiles(release, federationToken)
	if err != nil {
		return err
	}

	return nil
}

func (u Uploader) NewCreateReleaseConfig(r config.Release) (pivnet.CreateReleaseConfig, error) {
	previousMajorRelease, err := u.Client.GetLatestPublicReleaseByReleaseType(
		r.GpdbMajorVersion(),
		config.MajorReleaseType,
	)
	if err != nil {
		vlog.Warn(err.Error())
	}

	previousMinorRelease, err := u.Client.GetLatestPublicReleaseByReleaseType(
		r.GpdbMajorVersion(),
		config.MinorReleaseType,
	)
	if err != nil {
		vlog.Warn(err.Error())
	}

	releaseType, err := r.ComputeReleaseType()
	if err != nil {
		return pivnet.CreateReleaseConfig{}, err
	}

	releaseNotesURL, err := r.ComputeReleaseNotesUrl()
	if err != nil {
		return pivnet.CreateReleaseConfig{}, err
	}

	endOfSupportDate, err := r.ComputeEndOfSupportDate(previousMajorRelease, previousMinorRelease)
	if err != nil {
		return pivnet.CreateReleaseConfig{}, err
	}
	endOfGuidanceDate, err := r.ComputeEndOfGuidanceDate(previousMajorRelease, previousMinorRelease)
	if err != nil {
		return pivnet.CreateReleaseConfig{}, err
	}

	return pivnet.CreateReleaseConfig{
		ProductSlug:           u.Context.Slug,
		Version:               r.Version,
		ReleaseType:           string(releaseType),
		ReleaseDate:           r.ReleaseDate,
		EULASlug:              r.EulaSlug,
		Description:           r.Description,
		ReleaseNotesURL:       releaseNotesURL,
		ECCN:                  r.ECCN,
		LicenseException:      r.LicenseException,
		EndOfSupportDate:      endOfSupportDate.String(),
		EndOfGuidanceDate:     endOfGuidanceDate.String(),
		EndOfAvailabilityDate: r.ComputeEndOfAvailabilityDate().String(),
		CopyMetadata:          false,
	}, nil
}

func (u Uploader) NewCreateProductFileConfig(f config.ProductFile) (pivnet.CreateProductFileConfig, error) {
	resolvedFile, err := u.Resolver.Resolve(f.File)
	if err != nil {
		return pivnet.CreateProductFileConfig{}, err
	}

	versionReplacer := NewVersionReplacer(resolvedFile)

	description := versionReplacer.Replace(f.Description)
	if Empty(description) {
		description = versionReplacer.Replace(f.UploadAs)
	}

	return pivnet.CreateProductFileConfig{
		ProductSlug:   u.Context.Slug,
		AWSObjectKey:  f.AWSObjectKey,
		Description:   description,
		DocsURL:       f.DocsURL,
		FileType:      f.FileType,
		FileVersion:   versionReplacer.Replace(f.FileVersion),
		IncludedFiles: f.IncludedFiles,
		SHA256:        f.SHA256,
		//MD5:                f.MD5,
		Name:               f.Name,
		Platforms:          f.Platforms,
		ReleasedAt:         f.ReleasedAt,
		SystemRequirements: f.SystemRequirements,
	}, nil
}

func (u Uploader) HandleFileGroups(release pivnet.Release, federationToken pivnet.FederationToken) error {
	fileGroups := u.Metadata.FileGroups
	for _, group := range fileGroups {
		g, err := u.Client.CreateFileGroup(group.Name)
		if err != nil {
			return err
		}
		err = u.Client.AddFileGroupToRelease(g.ID, release.ID)
		if err != nil {
			return err
		}

		for _, productFile := range group.ProductFiles {
			updatedProductFile, err := u.uploadToS3(productFile, federationToken)
			if err != nil {
				return err
			}
			cpfc, err := u.NewCreateProductFileConfig(updatedProductFile)
			if err != nil {
				return err
			}
			pf, err := u.Client.CreateProductFile(cpfc)
			if err != nil {
				return err
			}
			err = u.Client.AddProductFileToFileGroup(pf.ID, g.ID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (u Uploader) HandleProductFiles(release pivnet.Release, federationToken pivnet.FederationToken) error {
	for _, f := range u.Metadata.ProductFiles {
		updatedPf, err := u.uploadToS3(f, federationToken)
		if err != nil {
			return err
		}

		cpfc, err := u.NewCreateProductFileConfig(updatedPf)
		if err != nil {
			return err
		}

		pf, err := u.Client.CreateProductFile(cpfc)
		if err != nil {
			return err
		}

		err = u.Client.AddProductFileToRelease(pf.ID, release.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u Uploader) uploadToS3(productFile config.ProductFile, federationToken pivnet.FederationToken) (config.ProductFile, error) {
	rv, err := u.Resolver.Resolve(productFile.File)
	if err != nil {
		return config.ProductFile{}, err
	}

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(federationToken.Region),
		Credentials: credentials.NewStaticCredentials(
			federationToken.AccessKeyID,
			federationToken.SecretAccessKey,
			federationToken.SessionToken,
		),
	}))

	uploader := s3manager.NewUploader(sess)

	f, err := os.Open(rv.LocalFilePath)
	if err != nil {
		return config.ProductFile{}, fmt.Errorf("failed to open file %q, %v", rv.LocalFilePath, err)
	}

	AwsObjectKey := path.Join(u.AwsObjectPrefix, rv.LocalFileName)
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(federationToken.Bucket),
		Key:    aws.String(AwsObjectKey),
		Body:   f,
	})
	if err != nil {
		return config.ProductFile{}, fmt.Errorf("failed to upload file, %v", err)
	}

	vlog.Info("file uploaded to, %s\n", result.Location)

	productFile.AWSObjectKey = AwsObjectKey

	return productFile, nil
}

type VersionReplacer struct {
	resolvedFile ResolvedFile
}

func NewVersionReplacer(resolvedFile ResolvedFile) VersionReplacer {
	return VersionReplacer{
		resolvedFile: resolvedFile,
	}
}

func (vr VersionReplacer) Replace(expression string) string {
	if !strings.Contains(expression, `${VERSION_REGEX}`) {
		return expression
	}

	if vr.resolvedFile.ResolvedVersion.Empty() {
		vlog.Fatal("resolved version is empty")
	}

	return strings.ReplaceAll(expression, `${VERSION_REGEX}`, vr.resolvedFile.ResolvedVersion.String())
}
