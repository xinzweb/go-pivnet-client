package config_test

import (
	"fmt"
	"github.com/baotingfang/go-pivnet-client/api"
	"github.com/baotingfang/go-pivnet-client/api/apifakes"
	"github.com/baotingfang/go-pivnet-client/config"
	"github.com/baotingfang/go-pivnet-client/utils"
	semver "github.com/cppforlife/go-semi-semantic/version"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strings"
)

var metadataYaml = `
---
release:
  release_type: "Major Release"
  eula_slug: pivotal_software_eula
  description: "test description"
  release_notes_url: "http://example.com/notes/url"
  availability: Admins Only
  controlled: false
  eccn: 5D002
  license_exception: TSU
  release_date: 2013-05-19


file_groups:
- name: Greenplum Database Server
  product_files:
  - file: file://server-rhel6/greenplum-db-(6\..*)-rhel6-x86_64.rpm
    upload_as: Greenplum Database ${VERSION_REGEX} Installer for RHEL 6
    description:
    file_type: Software
    docs_url:
    system_requirements:
    platforms:
    included_files:
    file_version: ${VERSION_REGEX}
  - file: file://server-rhel7/greenplum-db-(6\..*)-rhel7-x86_64.rpm
    upload_as: Greenplum Database ${VERSION_REGEX} Installer for RHEL 7
    description:
    file_type: Software
    docs_url:
    system_requirements:
    platforms:
    included_files:
    file_version: ${VERSION_REGEX}
product_files:
- file: file://gpdb-osl/open_source_license_pivotal-gpdb-([0-9]+\.[0-9]+\.[0-9]+)-(.*).txt
  upload_as: Open Source Licenses for GPDB 6.x
  description:
  file_type: Open Source License
  docs_url:
  system_requirements:
  platforms:
  included_files:
  file_version: ${VERSION_REGEX}
- file: file://pl-extensions-osl/open_source_license_pivotal-gpdb-pl-extensions-([0-9]+\.[0-9]+\.[0-9]+)-(.*).txt
  upload_as: Open Source Licenses for Greenplum 6.x Procedural Language Extensions
  description:
  file_type: Open Source License
  docs_url:
  system_requirements:
  platforms:
  included_files:
  file_version: ${VERSION_REGEX}
`

var _ = Describe("Config", func() {
	BeforeEach(func() {
		api.DefaultClient = &apifakes.FakeAccessInterface{}
	})

	Context("Metadata config", func() {
		It("Decode metadata config", func() {
			metadataReader := strings.NewReader(metadataYaml)
			metaData, err := api.MetadataFrom(metadataReader, "6.6.0")
			Expect(err).NotTo(HaveOccurred())

			Expect(metaData.Release.ReleaseType).To(Equal("Major Release"))
			Expect(metaData.Release.EulaSlug).To(Equal("pivotal_software_eula"))
			Expect(len(metaData.FileGroups)).To(Equal(1))
			Expect(len(metaData.FileGroups[0].ProductFiles)).To(Equal(2))
			Expect(len(metaData.ProductFiles)).To(Equal(2))
		})
	})

	Context("Release business logic", func() {
		It("ComputeReleaseType: User provided", func() {
			inputReleaseType := []config.ReleaseType{
				config.MinorRelease,
				config.MinorRelease,
				config.AlphaRelease,
				config.BetaRelease,
				config.MaintenanceRelease,
			}

			for _, rt := range inputReleaseType {
				r := &config.Release{
					ReleaseType: rt.String(),
				}
				t, err := r.ComputeReleaseType()
				Expect(err).NotTo(HaveOccurred())
				Expect(t).To(Equal(rt.String()))
			}
		})

		It("ComputeReleaseType: gpdb4 error version", func() {
			inputVersions := []string{
				"4.3.3",
				"4.3",
				"4.3.37.5.6",
			}
			for _, version := range inputVersions {
				r := &config.Release{
					ReleaseType: config.COMPUTED,
					Version:     semver.MustNewVersionFromString(version),
				}
				t, err := r.ComputeReleaseType()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("invalid release version for gpdb4: " + version))
				Expect(t).To(Equal(""))
			}
		})

		It("ComputeReleaseType: gpdb5 error version", func() {
			inputVersions := []string{
				"5.2",
				"5.27.1.2",
			}
			for _, version := range inputVersions {
				r := &config.Release{
					ReleaseType: config.COMPUTED,
					Version:     semver.MustNewVersionFromString(version),
				}
				t, err := r.ComputeReleaseType()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("invalid release version for gpdb5: " + version))
				Expect(t).To(Equal(""))
			}
		})

		It("ComputeReleaseType: gpdb6 error version", func() {
			inputVersions := []string{
				"6.7.2.1",
				"6.7",
			}
			for _, version := range inputVersions {
				r := &config.Release{
					ReleaseType: config.COMPUTED,
					Version:     semver.MustNewVersionFromString(version),
				}
				t, err := r.ComputeReleaseType()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("invalid release version for gpdb6: " + version))
				Expect(t).To(Equal(""))
			}
		})

		It("ComputeReleaseType: error version not gpdb4/5/6", func() {
			inputVersions := []string{
				"1.0.0",
				"7.1.0",
				"7.0.0",
			}
			for _, version := range inputVersions {
				r := &config.Release{
					ReleaseType: config.COMPUTED,
					Version:     semver.MustNewVersionFromString(version),
				}
				t, err := r.ComputeReleaseType()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("invalid gpdb release version: " + version))
				Expect(t).To(Equal(""))
			}
		})

		It("Test Major version", func() {
			r := &config.Release{
				ReleaseType: config.COMPUTED,
			}

			defer utils.TestHelper(`\[Default Logger\]\[FATAL\] can not found gpdb major version\.`)
			_ = r.MajorVersion()
		})

		It("Test Minor version", func() {
			r := &config.Release{
				ReleaseType: config.COMPUTED,
			}

			defer utils.TestHelper(`\[Default Logger\]\[FATAL\] can not found gpdb minor version\.`)
			_ = r.MinorVersion()
		})

		It("Test Patch version", func() {
			r := &config.Release{
				ReleaseType: config.COMPUTED,
			}

			defer utils.TestHelper(`\[Default Logger\]\[FATAL\] can not found gpdb patch version\.`)
			_ = r.PatchVersion()
		})

		It("ComputeReleaseType: correct version", func() {
			inputVersions := []string{
				"4.3.3.0",
				"4.3.3.1",
				"4.3.3.0-alpha.1",
				"4.3.3.0-beta.1",

				"5.0.0",
				"5.1.0",
				"5.1.1",
				"5.1.1-alpha.1",
				"5.1.1-beta.1",

				"6.0.0",
				"6.1.0",
				"6.1.1",
				"6.1.1-alpha.1",
				"6.1.1-beta.1",
			}

			releaseTypes := []config.ReleaseType{
				config.MinorRelease,
				config.MaintenanceRelease,
				config.AlphaRelease,
				config.BetaRelease,

				config.MajorRelease,
				config.MinorRelease,
				config.MaintenanceRelease,
				config.AlphaRelease,
				config.BetaRelease,

				config.MajorRelease,
				config.MinorRelease,
				config.MaintenanceRelease,
				config.AlphaRelease,
				config.BetaRelease,
			}

			for i, v := range inputVersions {
				fmt.Printf("%d: %s\n", i, v)
				r := &config.Release{
					ReleaseType: config.COMPUTED,
					Version:     semver.MustNewVersionFromString(inputVersions[i]),
				}
				t, err := r.ComputeReleaseType()
				Expect(err).NotTo(HaveOccurred())
				Expect(t).To(Equal(releaseTypes[i].String()))
			}
		})
	})
})
