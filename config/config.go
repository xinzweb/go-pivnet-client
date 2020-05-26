package config

import (
	"fmt"
	. "github.com/baotingfang/go-pivnet-client/utils"
	"github.com/baotingfang/go-pivnet-client/vlog"
	semver "github.com/cppforlife/go-semi-semantic/version"
	"github.com/pivotal-cf/go-pivnet/v4"
	"gopkg.in/yaml.v2"
	"io"
	"strconv"
	"strings"
	"time"
)

const (
	COMPUTED = "<COMPUTED>"

	AlphaReleaseType       = "Alpha Release"
	BetaReleaseType        = "Beta Release"
	MajorReleaseType       = "Major Release"
	MinorReleaseType       = "Minor Release"
	MaintenanceReleaseType = "Maintenance Release"
)

type Version struct {
	version semver.Version
}

func NewVersion(v semver.Version) Version {
	return Version{version: v}
}

func (v Version) MajorVersion() int {
	return v.versionAt(0)
}

func (v Version) MinorVersion() int {
	return v.versionAt(1)
}

func (v Version) PatchVersion() int {
	return v.versionAt(2)
}

func (v Version) versionAt(position int) int {
	ver := v.version
	if len(ver.Release.Components) > position {
		patchVersion, err := strconv.Atoi(ver.Release.Components[position].AsString())
		if err != nil {
			vlog.Fatal(err.Error())
			return -1
		}
		return patchVersion
	}
	return -1
}

type Release struct {
	pivnet.Release              `json:",inline" yaml:",inline"`
	EulaSlug                    string `json:"eula_slug,omitempty" yaml:"eula_slug,omitempty"`
	EndOfAvailabilityDateOffset string `json:"end_of_availability_date_offset,omitempty" yaml:"end_of_availability_date_offset,omitempty"`
}

func (r Release) GpdbVersion() semver.Version {
	if !Empty(r.Version) {
		v, err := semver.NewVersionFromString(r.Version)
		if err != nil {
			vlog.Fatal("invalid gpdb version: %s", r.Version)
		}
		return v
	}
	vlog.Info("gpdb version is empty")
	return semver.Version{}
}

func (r Release) GpdbMajorVersion() int {
	return NewVersion(r.GpdbVersion()).MajorVersion()
}

func (r Release) GpdbMinorVersion() int {
	return NewVersion(r.GpdbVersion()).MinorVersion()
}

func (r Release) GpdbPatchVersion() int {
	return NewVersion(r.GpdbVersion()).PatchVersion()
}

func (r Release) Empty() bool {
	return r == Release{}
}

func (r *Release) ComputeReleaseType() (pivnet.ReleaseType, error) {
	if r.ReleaseType != COMPUTED && !Empty(r.ReleaseType) {
		return r.ReleaseType, nil
	}

	version := r.GpdbVersion()

	if version.PreRelease.Components != nil {
		for _, p := range version.PreRelease.Components {
			if strings.Contains(p.AsString(), "alpha") {
				r.ReleaseType = AlphaReleaseType
				return r.ReleaseType, nil
			}
			if strings.Contains(p.AsString(), "beta") {
				r.ReleaseType = BetaReleaseType
				return r.ReleaseType, nil
			}
		}
	}

	// gpdb4
	if r.GpdbMajorVersion() == 4 {
		if len(version.Release.Components) != 4 {
			return "", fmt.Errorf("invalid release version for gpdb4: %s", r.Version)
		}

		// 4.3.33.0, not Patch Version
		if version.Release.Components[3].AsString() == "0" {
			r.ReleaseType = MinorReleaseType
			return r.ReleaseType, nil
		} else {
			r.ReleaseType = MaintenanceReleaseType
			return r.ReleaseType, nil
		}
	}

	// gpdb5 gpdb6
	if r.GpdbMajorVersion() == 5 || r.GpdbMajorVersion() == 6 {
		if len(version.Release.Components) != 3 {
			return "", fmt.Errorf("invalid release version for gpdb%d: %s", r.GpdbMajorVersion(), r.Version)
		}

		if r.GpdbMinorVersion() == 0 && r.GpdbPatchVersion() == 0 {
			r.ReleaseType = MajorReleaseType
			return r.ReleaseType, nil
		}

		if r.GpdbMinorVersion() != 0 && r.GpdbPatchVersion() == 0 {
			r.ReleaseType = MinorReleaseType
			return r.ReleaseType, nil
		}

		if r.GpdbPatchVersion() != 0 {
			r.ReleaseType = MaintenanceReleaseType
			return r.ReleaseType, nil
		}
	}
	return "", fmt.Errorf("invalid gpdb release version: %s", r.Version)
}

func (r *Release) ComputeEndOfSupportDate(previousMajorRelease, previousMinorRelease pivnet.Release) (Date, error) {
	const (
		OffsetFromMajorReleaseMonths        = 36
		OffsetFromCurrentMinorReleaseMonths = 18
	)

	if !Empty(r.EndOfSupportDate) && r.EndOfSupportDate != COMPUTED {
		return MustParseDateFrom(r.EndOfSupportDate), nil
	}

	if !Empty(previousMajorRelease) &&
		previousMajorRelease.ReleaseType != MajorReleaseType {
		return Date{},
			fmt.Errorf("previous major release type is wrong. actual release type:%s", previousMajorRelease.ReleaseType)
	}

	if !Empty(previousMinorRelease) &&
		previousMinorRelease.ReleaseType != MinorReleaseType {
		return Date{},
			fmt.Errorf("previous minor release type is wrong. actual release type:%s", previousMinorRelease.ReleaseType)
	}

	initReleaseDate := Date{Time: time.Now()}
	releaseDate := MustParseDateFrom(r.ReleaseDate)
	if !releaseDate.IsZero() {
		initReleaseDate = releaseDate
	}

	if r.ReleaseType == MaintenanceReleaseType {
		if Empty(previousMinorRelease) {
			return Date{},
				fmt.Errorf("current release type: %s, Can not find the previous minor release", r.ReleaseType)
		}
		previousMinorReleaseDate := MustParseDateFrom(previousMinorRelease.ReleaseDate)
		r.EndOfSupportDate = previousMinorReleaseDate.LastDayOfCurrentMonth().String()
		return MustParseDateFrom(r.EndOfSupportDate), nil
	}

	if r.ReleaseType == MinorReleaseType {
		if Empty(previousMajorRelease) {
			return Date{},
				fmt.Errorf("current release type: %s, Can not find the previous major release", r.ReleaseType)
		}

		previousMajorReleaseDate := MustParseDateFrom(previousMajorRelease.ReleaseDate)
		t1 := Date{
			Time: previousMajorReleaseDate.AddDate(0, OffsetFromMajorReleaseMonths, 0),
		}
		t2 := Date{
			Time: initReleaseDate.AddDate(0, OffsetFromCurrentMinorReleaseMonths, 0),
		}

		if t1.Time.Before(t2.Time) {
			r.EndOfSupportDate = t2.LastDayOfCurrentMonth().String()
			return MustParseDateFrom(r.EndOfSupportDate), nil
		} else {
			r.EndOfSupportDate = t1.LastDayOfCurrentMonth().String()
			return MustParseDateFrom(r.EndOfSupportDate), nil
		}
	}

	if r.ReleaseType == MajorReleaseType ||
		r.ReleaseType == AlphaReleaseType ||
		r.ReleaseType == BetaReleaseType {
		t1 := Date{
			Time: initReleaseDate.AddDate(0, OffsetFromMajorReleaseMonths, 0),
		}
		r.EndOfSupportDate = t1.LastDayOfCurrentMonth().String()
		return MustParseDateFrom(r.EndOfSupportDate), nil
	}

	return Date{}, fmt.Errorf("invalid release type: %s", r.ReleaseType)
}

func (r *Release) ComputeEndOfGuidanceDate(previousMajorRelease, previousMinorRelease pivnet.Release) (Date, error) {

	if !Empty(r.EndOfGuidanceDate) && r.EndOfGuidanceDate != COMPUTED {
		return MustParseDateFrom(r.EndOfGuidanceDate), nil
	}

	var endOfSupportDate Date
	if !Empty(r.EndOfSupportDate) && r.EndOfSupportDate != COMPUTED {
		endOfSupportDate = MustParseDateFrom(r.EndOfSupportDate)
	}

	if endOfSupportDate.IsZero() {
		eod, err := r.ComputeEndOfSupportDate(previousMajorRelease, previousMinorRelease)
		if err != nil {
			return Date{}, err
		}
		endOfSupportDate = eod
	}

	r.EndOfSupportDate = endOfSupportDate.String()

	const OffsetFromEndOfSupportDate = 12
	r.EndOfGuidanceDate = Date{
		Time: endOfSupportDate.AddDate(0, OffsetFromEndOfSupportDate, 0),
	}.String()

	return MustParseDateFrom(r.EndOfGuidanceDate), nil
}

func (r *Release) ComputeEndOfAvailabilityDate() Date {
	if !Empty(r.EndOfAvailabilityDate) && r.EndOfAvailabilityDate != COMPUTED {
		return MustParseDateFrom(r.EndOfAvailabilityDate)
	}

	initReleaseDate := Date{Time: time.Now()}
	if !Empty(r.ReleaseDate) && r.ReleaseDate != COMPUTED {
		initReleaseDate = MustParseDateFrom(r.ReleaseDate)
	}

	endOfAvailabilityDate := initReleaseDate.Offset(r.EndOfAvailabilityDateOffset)
	r.EndOfAvailabilityDate = endOfAvailabilityDate.String()
	return endOfAvailabilityDate
}

func (r *Release) ComputeReleaseNotesUrl() (string, error) {
	if !Empty(r.ReleaseNotesURL) && r.ReleaseNotesURL != COMPUTED {
		return r.ReleaseNotesURL, nil
	}

	if r.GpdbMajorVersion() == 6 {
		return generateGpdb6ReleaseNotesUrl(r.GpdbVersion()), nil
	}

	if r.GpdbMajorVersion() == 5 {
		return generateGpdb5ReleaseNotesUrl(r.GpdbVersion()), nil
	}

	if r.GpdbMajorVersion() == 4 {
		return generateGpdb4ReleaseNotesUrl(r.GpdbVersion()), nil
	}

	return "", fmt.Errorf("compute release notes url failed. only support gpdb4/5/6")

}

func generateGpdb6ReleaseNotesUrl(v semver.Version) string {
	releaseComponents := make([]string, 0)
	for _, component := range v.Release.Components {
		releaseComponents = append(releaseComponents, component.AsString())
	}

	preRelease := v.PreRelease
	if !preRelease.Empty() {
		preReleaseStr := strings.Title(preRelease.Components[0].AsString())

		return fmt.Sprintf(
			"https://gpdb.docs.pivotal.io/%s%s/main/index.html",
			strings.Join(releaseComponents[:2], "-"),
			preReleaseStr,
		)
	}
	fmt.Println(releaseComponents)
	return fmt.Sprintf(
		"https://gpdb.docs.pivotal.io/%s/main/index.html",
		strings.Join(releaseComponents[:2], "-"),
	)
}

func generateGpdb5ReleaseNotesUrl(v semver.Version) string {
	releaseComponents := make([]string, 0)
	for _, component := range v.Release.Components {
		releaseComponents = append(releaseComponents, component.AsString())
	}

	preRelease := v.PreRelease
	if !preRelease.Empty() {
		preReleaseStr := strings.Title(preRelease.Components[0].AsString())

		return fmt.Sprintf(
			"https://gpdb.docs.pivotal.io/%s%s/main/index.html",
			strings.Join(releaseComponents, ""),
			preReleaseStr,
		)
	}

	p1 := strings.Join(releaseComponents[:2], "") + "0"
	p2 := strings.Join(releaseComponents, "")

	return fmt.Sprintf(
		"https://gpdb.docs.pivotal.io/%s/relnotes/GPDB_%s_README.html", p1, p2)
}

func generateGpdb4ReleaseNotesUrl(v semver.Version) string {
	releaseComponents := make([]string, 0)
	for _, component := range v.Release.Components {
		releaseComponents = append(releaseComponents, component.AsString())
	}

	preRelease := v.PreRelease
	if !preRelease.Empty() {
		preReleaseStr := strings.Title(preRelease.Components[0].AsString())

		return fmt.Sprintf(
			"https://gpdb.docs.pivotal.io/%s%s/main/index.html",
			strings.Join(releaseComponents, ""),
			preReleaseStr,
		)
	}

	p1 := strings.Join(releaseComponents[:3], "") + "0"
	p2 := strings.Join(releaseComponents, "")

	return fmt.Sprintf(
		"https://gpdb.docs.pivotal.io/%s/relnotes/GPDB_%s_README.html", p1, p2)
}

type FileGroup struct {
	Name         string        `json:"name,omitempty" yaml:"name,omitempty"`
	ProductFiles []ProductFile `json:"product_files,omitempty" yaml:"product_files,omitempty"`
}

type ProductFile struct {
	pivnet.ProductFile `json:",inline" yaml:",inline"`
	File               string `json:"file,omitempty" yaml:"file,omitempty"`
	UploadAs           string `json:"upload_as,omitempty" yaml:"upload_as,omitempty"`
}

type Metadata struct {
	Release      Release       `json:"release,omitempty" yaml:"release,omitempty"`
	FileGroups   []FileGroup   `json:"file_groups,omitempty" yaml:"file_groups,omitempty"`
	ProductFiles []ProductFile `json:"product_file,omitempty" yaml:"product_files,omitempty"`
}

func MetadataFrom(reader io.Reader, gpdbVersion string) (Metadata, error) {
	_, err := semver.NewVersionFromString(gpdbVersion)
	if err != nil {
		return Metadata{}, err
	}

	var metadata Metadata
	if err := yaml.NewDecoder(reader).Decode(&metadata); err != nil {
		return Metadata{}, err
	}
	metadata.Release.Version = gpdbVersion
	vlog.Info("GPDB Version: %s", metadata.Release.GpdbVersion().String())

	return metadata, nil
}
