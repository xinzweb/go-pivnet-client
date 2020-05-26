package config_test

import (
	"fmt"
	"github.com/baotingfang/go-pivnet-client/config"
	"github.com/baotingfang/go-pivnet-client/vlog"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/pivotal-cf/go-pivnet/v4"
	"log"
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
	Context("Metadata config", func() {
		It("Decode metadata config", func() {
			metadataReader := strings.NewReader(metadataYaml)
			metaData, err := config.MetadataFrom(metadataReader, "6.6.0")
			Expect(err).NotTo(HaveOccurred())

			Expect(string(metaData.Release.ReleaseType)).To(Equal("Major Release"))
			Expect(metaData.Release.EulaSlug).To(Equal("pivotal_software_eula"))
			Expect(len(metaData.FileGroups)).To(Equal(1))
			Expect(len(metaData.FileGroups[0].ProductFiles)).To(Equal(2))
			Expect(len(metaData.ProductFiles)).To(Equal(2))
		})

		It("gpdb version is not valid", func() {
			metadataReader := strings.NewReader(metadataYaml)
			metaData, err := config.MetadataFrom(metadataReader, "invalid version")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Expected version 'invalid version' to match version format"))
			Expect(metaData).To(Equal(config.Metadata{}))
		})

		It("metadata yaml is not valid", func() {
			metadataReader := strings.NewReader("invalid yaml content")
			metaData, err := config.MetadataFrom(metadataReader, "6.7.0")
			Expect(err).To(HaveOccurred())
			fmt.Println(err.Error())
			Expect(err.Error()).To(HavePrefix("yaml: unmarshal errors:"))
			Expect(metaData).To(Equal(config.Metadata{}))
		})
	})

	Context("Test major/minor/patch versions", func() {
		var (
			outLog *gbytes.Buffer
			errLog *gbytes.Buffer
		)
		BeforeEach(func() {
			outLog = gbytes.NewBuffer()
			errLog = gbytes.NewBuffer()
			vlog.Log = &vlog.Logger{
				OutLogger: log.New(outLog, "", log.LstdFlags),
				ErrLogger: log.New(errLog, "", log.LstdFlags),
				LogLevel:  vlog.DebugLevel,
				Prefix:    "Test",
			}
		})

		It("Test Major version", func() {
			r := &config.Release{}

			r.ReleaseType = config.COMPUTED
			majorVersion := r.GpdbMajorVersion()
			Expect(majorVersion).To(Equal(-1))
			Expect(string(outLog.Contents())).To(ContainSubstring(
				"gpdb version is empty"))
		})

		It("Test Major version with panic", func() {
			r := &config.Release{}

			r.ReleaseType = config.COMPUTED
			r.Version = "beta-1.2.3.build.1"
			f := func() {
				r.GpdbMajorVersion()
			}
			Expect(f).To(PanicWith(MatchRegexp(`strconv.Atoi: parsing "beta": invalid syntax`)))
		})

		It("Test Major version with panic. invalid version", func() {
			r := &config.Release{}

			r.ReleaseType = config.COMPUTED
			r.Version = "invalid version"
			f := func() {
				r.GpdbMajorVersion()
			}
			Expect(f).To(PanicWith(MatchRegexp(`invalid gpdb version: invalid version`)))
		})

		It("Test Minor version", func() {
			r := &config.Release{}

			r.ReleaseType = config.COMPUTED
			minorVersion := r.GpdbMinorVersion()
			Expect(minorVersion).To(Equal(-1))
			Expect(string(outLog.Contents())).To(ContainSubstring("gpdb version is empty"))
		})

		It("Test Patch version", func() {
			r := &config.Release{}
			r.ReleaseType = config.COMPUTED

			patchVersion := r.GpdbPatchVersion()
			Expect(patchVersion).To(Equal(-1))
			Expect(string(outLog.Contents())).To(ContainSubstring("gpdb version is empty"))
		})
	})

	Context("Test ComputeReleaseType", func() {
		It("ComputeReleaseType: User provided", func() {
			inputReleaseType := []pivnet.ReleaseType{
				config.MinorReleaseType,
				config.MinorReleaseType,
				config.AlphaReleaseType,
				config.BetaReleaseType,
				config.MaintenanceReleaseType,
			}

			for _, rt := range inputReleaseType {
				r := &config.Release{}
				r.ReleaseType = rt
				t, err := r.ComputeReleaseType()
				Expect(err).NotTo(HaveOccurred())
				Expect(t).To(Equal(rt))
			}
		})

		It("ComputeReleaseType: gpdb4 error version", func() {
			inputVersions := []string{
				"4.3.3",
				"4.3",
				"4.3.37.5.6",
			}
			for _, version := range inputVersions {
				r := &config.Release{}
				r.ReleaseType = config.COMPUTED
				r.Version = version
				t, err := r.ComputeReleaseType()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("invalid release version for gpdb4: " + version))
				Expect(string(t)).To(Equal(""))
			}
		})

		It("ComputeReleaseType: gpdb5 error version", func() {
			inputVersions := []string{
				"5.2",
				"5.27.1.2",
			}
			for _, version := range inputVersions {
				r := &config.Release{}
				r.ReleaseType = config.COMPUTED
				r.Version = version
				t, err := r.ComputeReleaseType()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("invalid release version for gpdb5: " + version))
				Expect(string(t)).To(Equal(""))
			}
		})

		It("ComputeReleaseType: gpdb6 error version", func() {
			inputVersions := []string{
				"6.7.2.1",
				"6.7",
			}
			for _, version := range inputVersions {
				r := &config.Release{}
				r.ReleaseType = config.COMPUTED
				r.Version = version
				t, err := r.ComputeReleaseType()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("invalid release version for gpdb6: " + version))
				Expect(string(t)).To(Equal(""))
			}
		})

		It("ComputeReleaseType: error version not gpdb4/5/6", func() {
			inputVersions := []string{
				"1.0.0",
				"7.1.0",
				"7.0.0",
			}
			for _, version := range inputVersions {
				r := &config.Release{}
				r.ReleaseType = config.COMPUTED
				r.Version = version
				t, err := r.ComputeReleaseType()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("invalid gpdb release version: " + version))
				Expect(string(t)).To(Equal(""))
			}
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

			releaseTypes := []pivnet.ReleaseType{
				config.MinorReleaseType,
				config.MaintenanceReleaseType,
				config.AlphaReleaseType,
				config.BetaReleaseType,

				config.MajorReleaseType,
				config.MinorReleaseType,
				config.MaintenanceReleaseType,
				config.AlphaReleaseType,
				config.BetaReleaseType,

				config.MajorReleaseType,
				config.MinorReleaseType,
				config.MaintenanceReleaseType,
				config.AlphaReleaseType,
				config.BetaReleaseType,
			}

			for i, version := range inputVersions {
				fmt.Printf("%d: %s\n", i, version)
				r := &config.Release{}
				r.ReleaseType = config.COMPUTED
				r.Version = version
				t, err := r.ComputeReleaseType()
				Expect(err).NotTo(HaveOccurred())
				Expect(t).To(Equal(releaseTypes[i]))
			}
		})
	})

	Context("Test ComputeEndOfSupportDate", func() {
		It("ComputeEndOfSupportDate: user provide ", func() {
			release := config.Release{
				Release: pivnet.Release{
					ReleaseDate:      "2008-08-18",
					EndOfSupportDate: "2015-08-18",
					ReleaseType:      config.MaintenanceReleaseType,
				},
			}

			d, err := release.ComputeEndOfSupportDate(pivnet.Release{}, pivnet.Release{})
			Expect(err).NotTo(HaveOccurred())
			Expect(d.String()).To(Equal("2015-08-18"))
		})

		It("ComputeEndOfSupportDate: previousMajorRelease with error release type", func() {
			release := config.Release{}
			preMajorRelease := pivnet.Release{
				ReleaseType: config.MinorReleaseType,
			}

			d, err := release.ComputeEndOfSupportDate(preMajorRelease, pivnet.Release{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("previous major release type is wrong. actual release type:" + config.MinorReleaseType))
			Expect(d.IsZero()).To(BeTrue())
		})

		It("ComputeEndOfSupportDate: previousMinorRelease with error release type", func() {
			release := config.Release{}
			preMinorRelease := pivnet.Release{
				ReleaseType: config.MajorReleaseType,
			}
			d, err := release.ComputeEndOfSupportDate(pivnet.Release{}, preMinorRelease)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("previous minor release type is wrong. actual release type:" + config.MajorReleaseType))
			Expect(d.IsZero()).To(BeTrue())
		})

		It("ComputeEndOfSupportDate: maintenanceRelease type", func() {
			release := config.Release{
				Release: pivnet.Release{
					ReleaseDate:      "2008-08-18",
					EndOfSupportDate: "<COMPUTED>",
					ReleaseType:      config.MaintenanceReleaseType,
				},
			}

			By("previousMinorRelease is nil")
			d, err := release.ComputeEndOfSupportDate(pivnet.Release{}, pivnet.Release{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(fmt.Sprintf("current release type: %s, Can not find the previous minor release", release.ReleaseType)))
			Expect(d.IsZero()).To(BeTrue())

			By("Correct path")
			preMinorRelease := pivnet.Release{
				ReleaseDate: "2007-07-01",
				ReleaseType: config.MinorReleaseType,
			}

			d, err = release.ComputeEndOfSupportDate(pivnet.Release{}, preMinorRelease)
			Expect(err).NotTo(HaveOccurred())
			Expect(d.String()).To(Equal("2007-07-31"))
		})

		It("ComputeEndOfSupportDate: minorRelease type", func() {
			By("previousMajorRelease is nil")
			release := config.Release{
				Release: pivnet.Release{
					ReleaseDate:      "2008-08-18",
					EndOfSupportDate: "<COMPUTED>",
					ReleaseType:      config.MinorReleaseType,
				},
			}
			d, err := release.ComputeEndOfSupportDate(pivnet.Release{}, pivnet.Release{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(fmt.Sprintf("current release type: %s, Can not find the previous major release", release.ReleaseType)))
			Expect(d.IsZero()).To(BeTrue())

			By("Correct path1: release minor release during previous major release eod")
			release = config.Release{
				Release: pivnet.Release{
					ReleaseDate:      "2008-08-18",
					EndOfSupportDate: "<COMPUTED>",
					ReleaseType:      config.MinorReleaseType,
				},
			}

			preMajorRelease := pivnet.Release{
				ReleaseDate: "2007-08-18",
				ReleaseType: config.MajorReleaseType,
			}

			d, err = release.ComputeEndOfSupportDate(preMajorRelease, pivnet.Release{})
			Expect(err).NotTo(HaveOccurred())
			Expect(d.String()).To(Equal("2010-08-31"))

			By("Correct path2: release minor release after previous major release eod")
			release = config.Release{
				Release: pivnet.Release{
					ReleaseDate:      "2009-08-18",
					EndOfSupportDate: "<COMPUTED>",
					ReleaseType:      config.MinorReleaseType,
				},
			}

			preMajorRelease = pivnet.Release{
				ReleaseDate: "2007-08-18",
				ReleaseType: config.MajorReleaseType,
			}
			d, err = release.ComputeEndOfSupportDate(preMajorRelease, pivnet.Release{})
			Expect(err).NotTo(HaveOccurred())
			Expect(d.String()).To(Equal("2011-02-28"))

		})

		It("ComputeEndOfSupportDate: MajorRelease, AlphaRelease, BetaRelease type", func() {
			release := config.Release{
				Release: pivnet.Release{
					ReleaseDate:      "2008-08-18",
					EndOfSupportDate: "<COMPUTED>",
					ReleaseType:      config.MajorReleaseType,
				},
			}

			d, err := release.ComputeEndOfSupportDate(pivnet.Release{}, pivnet.Release{})
			Expect(err).NotTo(HaveOccurred())
			Expect(d.String()).To(Equal("2011-08-31"))

			release = config.Release{
				Release: pivnet.Release{
					ReleaseDate:      "2008-08-18",
					EndOfSupportDate: "<COMPUTED>",
					ReleaseType:      config.AlphaReleaseType,
				},
			}

			d, err = release.ComputeEndOfSupportDate(pivnet.Release{}, pivnet.Release{})
			Expect(err).NotTo(HaveOccurred())
			Expect(d.String()).To(Equal("2011-08-31"))

			release = config.Release{
				Release: pivnet.Release{
					ReleaseDate:      "2008-08-18",
					EndOfSupportDate: "<COMPUTED>",
					ReleaseType:      config.BetaReleaseType,
				},
			}

			d, err = release.ComputeEndOfSupportDate(pivnet.Release{}, pivnet.Release{})
			Expect(err).NotTo(HaveOccurred())
			Expect(d.String()).To(Equal("2011-08-31"))

		})

		It("ComputeEndOfSupportDate: With incorrect release type", func() {
			release := config.Release{
				Release: pivnet.Release{
					ReleaseDate:      "2008-08-18",
					EndOfSupportDate: "<COMPUTED>",
					ReleaseType:      pivnet.ReleaseType("INVALID_RELEASE_TYPE")},
			}
			d, err := release.ComputeEndOfSupportDate(pivnet.Release{}, pivnet.Release{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("invalid release type: INVALID_RELEASE_TYPE"))
			Expect(d.IsZero()).To(BeTrue())
		})
	})

	Context("Test ComputeEndOfGuidanceDate", func() {
		It("ComputeEndOfGuidanceDate: user provide EndOfGuidanceDate ", func() {
			release := config.Release{
				Release: pivnet.Release{
					ReleaseDate:       "2008-08-18",
					EndOfGuidanceDate: "2015-08-18",
					ReleaseType:       config.MaintenanceReleaseType,
				},
			}

			d, err := release.ComputeEndOfGuidanceDate(pivnet.Release{}, pivnet.Release{})
			Expect(err).NotTo(HaveOccurred())
			Expect(d.String()).To(Equal("2015-08-18"))
		})

		It("ComputeEndOfGuidanceDate: user provide EndOfSupportDate ", func() {
			release := config.Release{
				Release: pivnet.Release{
					ReleaseDate:      "2008-08-18",
					EndOfSupportDate: "2016-08-18",
					ReleaseType:      config.MaintenanceReleaseType,
				},
			}

			d, err := release.ComputeEndOfGuidanceDate(pivnet.Release{}, pivnet.Release{})
			Expect(err).NotTo(HaveOccurred())
			Expect(d.String()).To(Equal("2017-08-18"))
		})

		It("ComputeEndOfGuidanceDate: user provide both EndOfSupportDate and EndOfGuidanceDate", func() {
			release := config.Release{
				Release: pivnet.Release{
					ReleaseDate:       "2008-08-18",
					EndOfGuidanceDate: "2015-08-18",
					EndOfSupportDate:  "2016-08-18",
					ReleaseType:       config.MaintenanceReleaseType},
			}

			d, err := release.ComputeEndOfGuidanceDate(pivnet.Release{}, pivnet.Release{})
			Expect(err).NotTo(HaveOccurred())
			Expect(d.String()).To(Equal("2015-08-18"))
		})

		It("ComputeEndOfGuidanceDate: previousMajorRelease with error release type", func() {
			release := config.Release{}
			preMajorRelease := pivnet.Release{
				ReleaseType: config.MinorReleaseType,
			}

			d, err := release.ComputeEndOfGuidanceDate(preMajorRelease, pivnet.Release{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("previous major release type is wrong. actual release type:" + config.MinorReleaseType))
			Expect(d.IsZero()).To(BeTrue())
		})

		It("ComputeEndOfGuidanceDate: previousMinorRelease with error release type", func() {
			release := config.Release{}
			preMinorRelease := pivnet.Release{
				ReleaseType: config.MajorReleaseType,
			}
			d, err := release.ComputeEndOfGuidanceDate(pivnet.Release{}, preMinorRelease)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("previous minor release type is wrong. actual release type:" + config.MajorReleaseType))
			Expect(d.IsZero()).To(BeTrue())
		})

		It("ComputeEndOfGuidanceDate: maintenanceRelease type", func() {
			release := config.Release{
				Release: pivnet.Release{
					ReleaseDate:      "2008-08-18",
					EndOfSupportDate: "<COMPUTED>",
					ReleaseType:      config.MaintenanceReleaseType,
				},
			}

			By("previousMinorRelease is nil")
			d, err := release.ComputeEndOfGuidanceDate(pivnet.Release{}, pivnet.Release{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(fmt.Sprintf("current release type: %s, Can not find the previous minor release", release.ReleaseType)))
			Expect(d.IsZero()).To(BeTrue())

			By("Correct path")
			preMinorRelease := pivnet.Release{
				ReleaseDate: "2007-07-01",
				ReleaseType: config.MinorReleaseType,
			}
			d, err = release.ComputeEndOfGuidanceDate(pivnet.Release{}, preMinorRelease)
			Expect(err).NotTo(HaveOccurred())
			Expect(d.String()).To(Equal("2008-07-31"))
		})

		It("ComputeEndOfGuidanceDate: minorRelease type", func() {
			By("previousMajorRelease is nil")
			release := config.Release{
				Release: pivnet.Release{
					ReleaseDate:      "2008-08-18",
					EndOfSupportDate: "<COMPUTED>",
					ReleaseType:      config.MinorReleaseType,
				},
			}
			d, err := release.ComputeEndOfGuidanceDate(pivnet.Release{}, pivnet.Release{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(fmt.Sprintf("current release type: %s, Can not find the previous major release", release.ReleaseType)))
			Expect(d.IsZero()).To(BeTrue())

			By("Correct path1: release minor release during previous major release eod")
			release = config.Release{
				Release: pivnet.Release{
					ReleaseDate:      "2008-08-18",
					EndOfSupportDate: "<COMPUTED>",
					ReleaseType:      config.MinorReleaseType,
				},
			}

			preMajorRelease := pivnet.Release{
				ReleaseDate: "2007-08-18",
				ReleaseType: config.MajorReleaseType,
			}

			d, err = release.ComputeEndOfGuidanceDate(preMajorRelease, pivnet.Release{})
			Expect(err).NotTo(HaveOccurred())
			Expect(d.String()).To(Equal("2011-08-31"))

			By("Correct path2: release minor release after previous major release eod")
			release = config.Release{
				Release: pivnet.Release{
					ReleaseDate:      "2009-08-18",
					EndOfSupportDate: "<COMPUTED>",
					ReleaseType:      config.MinorReleaseType,
				},
			}

			preMajorRelease = pivnet.Release{
				ReleaseDate: "2007-08-18",
				ReleaseType: config.MajorReleaseType,
			}

			d, err = release.ComputeEndOfGuidanceDate(preMajorRelease, pivnet.Release{})
			Expect(err).NotTo(HaveOccurred())
			Expect(d.String()).To(Equal("2012-02-28"))

		})

		It("ComputeEndOfGuidanceDate: MajorRelease, AlphaRelease, BetaRelease type", func() {
			release := config.Release{
				Release: pivnet.Release{
					ReleaseDate:      "2008-08-18",
					EndOfSupportDate: "<COMPUTED>",
					ReleaseType:      config.MajorReleaseType,
				},
			}

			d, err := release.ComputeEndOfGuidanceDate(pivnet.Release{}, pivnet.Release{})
			Expect(err).NotTo(HaveOccurred())
			Expect(d.String()).To(Equal("2012-08-31"))

			release = config.Release{
				Release: pivnet.Release{
					ReleaseDate:      "2008-08-18",
					EndOfSupportDate: "<COMPUTED>",
					ReleaseType:      config.AlphaReleaseType,
				},
			}

			d, err = release.ComputeEndOfGuidanceDate(pivnet.Release{}, pivnet.Release{})
			Expect(err).NotTo(HaveOccurred())
			Expect(d.String()).To(Equal("2012-08-31"))

			release = config.Release{
				Release: pivnet.Release{
					ReleaseDate:      "2008-08-18",
					EndOfSupportDate: "<COMPUTED>",
					ReleaseType:      config.BetaReleaseType,
				},
			}

			d, err = release.ComputeEndOfGuidanceDate(pivnet.Release{}, pivnet.Release{})
			Expect(err).NotTo(HaveOccurred())
			Expect(d.String()).To(Equal("2012-08-31"))

		})

		It("ComputeEndOfGuidanceDate: With incorrect release type", func() {
			release := config.Release{
				Release: pivnet.Release{
					ReleaseDate:      "2008-08-18",
					EndOfSupportDate: "<COMPUTED>",
					ReleaseType:      pivnet.ReleaseType("INVALID_RELEASE_TYPE")},
			}
			d, err := release.ComputeEndOfGuidanceDate(pivnet.Release{}, pivnet.Release{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("invalid release type: INVALID_RELEASE_TYPE"))
			Expect(d.IsZero()).To(BeTrue())
		})
	})

	Context("Test ComputeEndOfAvailabilityDate", func() {
		It("ComputeEndOfAvailabilityDate: user provide EndOfAvailabilityDate", func() {
			release := config.Release{
				Release: pivnet.Release{
					ReleaseDate:           "2008-08-18",
					EndOfAvailabilityDate: "2012-12-10",
				},
			}
			d := release.ComputeEndOfAvailabilityDate()
			Expect(d.String()).To(Equal("2012-12-10"))
		})

		It("ComputeEndOfAvailabilityDate: user provide release date", func() {
			release := config.Release{
				Release: pivnet.Release{
					ReleaseDate: "2008-08-18",
				},
			}
			d := release.ComputeEndOfAvailabilityDate()
			Expect(d.String()).To(Equal("2008-08-18"))
		})

		It("ComputeEndOfAvailabilityDate: user provide release date and offset", func() {
			release := config.Release{
				Release: pivnet.Release{
					ReleaseDate: "2008-08-18",
				},
				EndOfAvailabilityDateOffset: "+1y+3m+10d",
			}
			d := release.ComputeEndOfAvailabilityDate()
			Expect(d.String()).To(Equal("2009-11-28"))
		})
	})

	Context("Test ComputeReleaseNotesUrl", func() {
		It("ComputeReleaseNotesUrl: user provide", func() {
			release := config.Release{
				Release: pivnet.Release{
					ReleaseNotesURL: "https://www.example/gpdb/docs/index.html",
				},
			}
			url, err := release.ComputeReleaseNotesUrl()
			Expect(err).NotTo(HaveOccurred())
			Expect(url).To(Equal("https://www.example/gpdb/docs/index.html"))
		})

		It("ComputeReleaseNotesUrl: gpdb6", func() {
			versions := []string{
				"6.7.0",
				"6.7.0-alpha.1",
				"6.7.0-beta.1",

				"5.27.1",
				"5.27.1-alpha.1",
				"5.27.1-beta.1",

				"4.3.33.7",
				"4.3.33.7-alpha.1",
				"4.3.33.7-beta.1",
			}

			urls := []string{
				"https://gpdb.docs.pivotal.io/6-7/main/index.html",
				"https://gpdb.docs.pivotal.io/6-7Alpha/main/index.html",
				"https://gpdb.docs.pivotal.io/6-7Beta/main/index.html",

				"https://gpdb.docs.pivotal.io/5270/relnotes/GPDB_5271_README.html",
				"https://gpdb.docs.pivotal.io/5271Alpha/main/index.html",
				"https://gpdb.docs.pivotal.io/5271Beta/main/index.html",

				"https://gpdb.docs.pivotal.io/43330/relnotes/GPDB_43337_README.html",
				"https://gpdb.docs.pivotal.io/43337Alpha/main/index.html",
				"https://gpdb.docs.pivotal.io/43337Beta/main/index.html",
			}

			for i := range versions {
				release := config.Release{
					Release: pivnet.Release{
						Version: versions[i],
					},
				}
				url, err := release.ComputeReleaseNotesUrl()
				Expect(err).NotTo(HaveOccurred())
				Expect(url).To(Equal(urls[i]))
			}
		})

		It("ComputeReleaseNotesUrl: not gpdb4/5/6", func() {
			release := config.Release{
				Release: pivnet.Release{
					Version: "7.0.0",
				},
			}
			url, err := release.ComputeReleaseNotesUrl()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("compute release notes url failed. only support gpdb4/5/6"))
			Expect(url).To(Equal(""))
		})
	})

	It("Release empty()", func() {
		r := config.Release{}
		Expect(r.Empty()).To(BeTrue())

		metadataReader := strings.NewReader(metadataYaml)
		metaData, err := config.MetadataFrom(metadataReader, "6.6.0")
		Expect(err).NotTo(HaveOccurred())
		Expect(metaData.Release.Empty()).To(BeFalse())
	})
})
