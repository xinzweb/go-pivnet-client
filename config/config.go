package config

import (
	"fmt"
	"github.com/baotingfang/go-pivnet-client/utils"
	semver "github.com/cppforlife/go-semi-semantic/version"
	"strings"
	"time"
)

const (
	COMPUTED = "<COMPUTED>"
)

type Release struct {
	ReleaseType                 string     `yaml:"release_type,omitempty"`
	EulaSlug                    string     `yaml:"eula_slug"`
	Description                 string     `yaml:"description"`
	ReleaseNotesUrl             string     `yaml:"release_notes_url,omitempty"`
	Availability                string     `yaml:"availability"`
	Controlled                  bool       `yaml:"controlled"`
	Eccn                        string     `yaml:"eccn"`
	LicenseException            string     `yaml:"license_exception"`
	ReleaseDate                 utils.Date `yaml:"release_date,omitempty"`
	EndOfSupportDate            utils.Date `yaml:"end_of_support_date,omitempty"`
	EndOfGuidanceDate           utils.Date `yaml:"end_of_guidance_date,omitempty"`
	EndOfAvailabilityDate       utils.Date `yaml:"end_of_availability_date,omitempty"`
	EndOfAvailabilityDateOffset string     `yaml:"end_of_availability_date_offset"`

	Id           int64
	Version      semver.Version
	MajorVersion int
	MinorVersion int
	PatchVersion int

	PreviousMinorRelease *Release
	PreviousMajorRelease *Release
}

func (r *Release) ComputeReleaseType() (string, error) {
	if r.ReleaseType != COMPUTED {
		return r.ReleaseType, nil
	}

	if r.Version.PreRelease.Components != nil {
		for _, p := range r.Version.PreRelease.Components {
			if strings.Contains(p.AsString(), "alpha") {
				return AlphaRelease.String(), nil
			}
			if strings.Contains(p.AsString(), "beta") {
				return BetaRelease.String(), nil
			}
		}
	}

	// gpdb4
	if r.MajorVersion == 4 {
		if len(r.Version.Release.Components) != 4 {
			return "", fmt.Errorf("invalid release version for gpdb4: %s", r.Version)
		}
		if r.PatchVersion == 0 {
			return MinorRelease.String(), nil
		} else {
			return MaintenanceRelease.String(), nil
		}
	}

	// gpdb5 gpdb6
	if r.MajorVersion == 5 || r.MajorVersion == 6 {
		if len(r.Version.Release.Components) != 3 {
			return "", fmt.Errorf("invalid release version for gpdb%d: %s", r.MajorVersion, r.Version)
		}

		if r.MinorVersion == 0 && r.PatchVersion == 0 {
			return MajorRelease.String(), nil
		}

		if r.MinorVersion != 0 && r.PatchVersion == 0 {
			return MinorRelease.String(), nil
		}

		if r.PatchVersion != 0 {
			return MaintenanceRelease.String(), nil
		}
	}
	return "", fmt.Errorf("invalid gpdb release version: %s", r.Version.String())
}

func (r *Release) ComputeEndOfSupportDate() (utils.Date, error) {
	if !r.EndOfSupportDate.IsZero() {
		return r.EndOfSupportDate, nil
	}

	const (
		OffsetFromMajorReleaseMonths        = 36
		OffsetFromCurrentMinorReleaseMonths = 18
	)
	initReleaseDate := utils.Date{Time: time.Now()}
	if !r.ReleaseDate.IsZero() {
		initReleaseDate = r.ReleaseDate
	}

	if r.ReleaseType == MaintenanceRelease.String() {
		if r.PreviousMajorRelease == nil {
			return utils.Date{},
				fmt.Errorf("current release type: %s, Can not find the previous minor release", r.ReleaseType)
		}
		return utils.EndDayOfCurrentMonth(r.PreviousMinorRelease.ReleaseDate), nil
	}

	if r.ReleaseType == MinorRelease.String() {
		if r.PreviousMajorRelease == nil {
			return utils.Date{},
				fmt.Errorf("current release type: %s, Can not find the previous major release", r.ReleaseType)
		}

		t1 := utils.Date{
			Time: r.PreviousMajorRelease.ReleaseDate.AddDate(0, OffsetFromMajorReleaseMonths, 0),
		}
		t2 := utils.Date{
			Time: initReleaseDate.AddDate(0, OffsetFromCurrentMinorReleaseMonths, 0),
		}

		if t1.Time.Before(t2.Time) {
			return utils.EndDayOfCurrentMonth(t2), nil
		} else {
			return utils.EndDayOfCurrentMonth(t1), nil
		}
	}

	if r.ReleaseType == MajorRelease.String() ||
		r.ReleaseType == AlphaRelease.String() ||
		r.ReleaseType == BetaRelease.String() {
		t1 := utils.Date{
			Time: initReleaseDate.AddDate(0, OffsetFromMajorReleaseMonths, 0),
		}
		return utils.EndDayOfCurrentMonth(t1), nil
	}
	return utils.Date{}, fmt.Errorf("invalid release type: %s", r.ReleaseType)
}

func (r *Release) ComputeEndOfGuidanceDate() utils.Date {
	if !r.EndOfGuidanceDate.IsZero() {
		return r.EndOfGuidanceDate
	}
	const OffsetFromEndOfSupportDate = 12
	return utils.Date{Time: r.EndOfSupportDate.AddDate(0, OffsetFromEndOfSupportDate, 0)}
}

func (r *Release) ComputeEndOfAvailabilityDate() utils.Date {
	if !r.EndOfAvailabilityDate.IsZero() {
		return r.EndOfAvailabilityDate
	}

	initReleaseDate := utils.Date{Time: time.Now()}
	if !r.ReleaseDate.IsZero() {
		initReleaseDate = r.ReleaseDate
	}

	return utils.ComputeFromOffset(initReleaseDate, r.EndOfAvailabilityDateOffset)
}

func (r *Release) ComputeReleaseNotesUrl() (string, error) {
	if r.ReleaseNotesUrl != COMPUTED || utils.IsEmpty(r.ReleaseNotesUrl) {
		return r.ReleaseNotesUrl, nil
	}

	if r.MajorVersion == 6 {
		return generateGpdb6ReleaseNotesUrl(r.Version), nil
	}

	if r.MajorVersion == 5 {
		return generateGpdb5ReleaseNotesUrl(r.Version), nil
	}

	if r.MajorVersion == 4 {
		return generateGpdb4ReleaseNotesUrl(r.Version), nil
	}
	return "", fmt.Errorf("compute release notes url faild. only support gpdb4/5/6")

}

func generateGpdb6ReleaseNotesUrl(v semver.Version) string {
	releaseComponents := make([]string, len(v.Release.Components))
	for _, component := range v.Release.Components {
		releaseComponents = append(releaseComponents, component.AsString())
	}

	preRelease := v.PreRelease
	if !preRelease.Empty() {
		preReleaseStr := strings.ToUpper(preRelease.Components[0].AsString())

		return fmt.Sprintf(
			"https://gpdb.docs.pivotal.io/%s%s/main/index.html",
			strings.Join(releaseComponents[:1], "-"),
			preReleaseStr,
		)
	}

	return fmt.Sprintf(
		"https://gpdb.docs.pivotal.io/%s/main/index.html",
		strings.Join(releaseComponents[:1], "-"),
	)
}

func generateGpdb5ReleaseNotesUrl(v semver.Version) string {
	releaseComponents := make([]string, len(v.Release.Components))
	for _, component := range v.Release.Components {
		releaseComponents = append(releaseComponents, component.AsString())
	}

	preRelease := v.PreRelease
	if !preRelease.Empty() {
		preReleaseStr := strings.ToUpper(preRelease.Components[0].AsString())

		return fmt.Sprintf(
			"https://gpdb.docs.pivotal.io/%s%s/main/index.html",
			strings.Join(releaseComponents, ""),
			preReleaseStr,
		)
	}

	p1 := strings.Join(releaseComponents[:1], "") + "0"
	p2 := strings.Join(releaseComponents, "")

	return fmt.Sprintf(
		"https://gpdb.docs.pivotal.io/%s/relnotes/GPDB_%s_README.html", p1, p2)
}

func generateGpdb4ReleaseNotesUrl(v semver.Version) string {
	releaseComponents := make([]string, len(v.Release.Components))
	for _, component := range v.Release.Components {
		releaseComponents = append(releaseComponents, component.AsString())
	}

	preRelease := v.PreRelease
	if !preRelease.Empty() {
		preReleaseStr := strings.ToUpper(preRelease.Components[0].AsString())

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

type FileGroup struct {
	Name         string        `yaml:"name"`
	ProductFiles []ProductFile `yaml:"product_files"`
}

type ProductFile struct {
	File               string   `yaml:"file"`
	UploadAs           string   `yaml:"upload_as"`
	Description        string   `yaml:"description"`
	FileType           string   `yaml:"file_type"`
	DocsUrl            string   `yaml:"docs_url,omitempty"`
	SystemRequirements []string `yaml:"system_requirements,omitempty"`
	Platforms          []string `yaml:"platforms,omitempty"`
	IncludedFiles      []string `yaml:"included_files,omitempty"`
	FileVersion        string   `yaml:"file_version"`
}

type Metadata struct {
	Release      Release       `yaml:"release"`
	FileGroups   []FileGroup   `yaml:"file_groups"`
	ProductFiles []ProductFile `yaml:"product_files"`
}
