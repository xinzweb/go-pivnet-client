package service_test

import (
	"fmt"
	"github.com/baotingfang/go-pivnet-client/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strings"

	. "github.com/baotingfang/go-pivnet-client/service"
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

  end_of_support_date: 2025-01-01
  end_of_guidance_date: 2028-02-01
  end_of_availability_date: 2030-03-03


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

var _ = Describe("Validator", func() {
	Context("MetaDataValidator", func() {
		var metadata config.Metadata

		BeforeEach(func() {
			metadataReader := strings.NewReader(metadataYaml)
			metadata, _ = config.MetadataFrom(metadataReader, "6.6.0")
		})

		It("MetaDataValidator: set EndOfAvailabilityDate and EndOfAvailabilityDateOffset at same time",
			func() {
				metadata.Release.EndOfAvailabilityDateOffset = "+1y+3m"
				mv := NewMetaDataValidator(metadata)
				Expect(mv.Validate()).To(BeFalse())
				Expect(mv.GetErrorMessages()).To(Equal([]string{
					"can not specify both end_of_availability_date and end_of_availability_date_offset",
				}))
			})

		It("MetaDataValidator: date settings with wrong format", func() {
			metadata.Release.EndOfSupportDate = "2013/05/19"
			metadata.Release.EndOfGuidanceDate = "2013/05/19"
			metadata.Release.EndOfAvailabilityDate = "2013/05/19"
			mv := NewMetaDataValidator(metadata)
			Expect(mv.Validate()).To(BeFalse())
			Expect(mv.GetErrorMessages()).To(Equal([]string{
				"end_of_support_date must be a valid date of the format \"YYYY-MM-DD\"",
				"end_of_guidance_date must be a valid date of the format \"YYYY-MM-DD\"",
				"end_of_availability_date must be a valid date of the format \"YYYY-MM-DD\"",
			}))
		})

		It("MetaDataValidator: end_of_availability_date_offset setting with wrong format", func() {
			// can not specify both end_of_availability_date and end_of_availability_date_offset
			metadata.Release.EndOfAvailabilityDate = ""
			metadata.Release.EndOfAvailabilityDateOffset = "1y3m"
			mv := NewMetaDataValidator(metadata)
			Expect(mv.Validate()).To(BeFalse())
			Expect(mv.GetErrorMessages()).To(Equal([]string{
				"end_of_availability_date_offset must be a valid offset of the form \"(+\\d+[mdyMDY])+\"",
			}))
		})

		It("MetaDataValidator: test file groups", func() {
			metadata.FileGroups[0].ProductFiles[0].File = ""
			mv := NewMetaDataValidator(metadata)
			Expect(mv.Validate()).To(BeFalse())
			Expect(mv.GetErrorMessages()).To(Equal([]string{
				"One of settings is empty.[file, upload_as, file_type, file_version]",
				"value is empty, index=0 (|  | Greenplum Database ${VERSION_REGEX} Installer for RHEL 6 | Software | ${VERSION_REGEX} |)",
			}))
		})

		It("MetaDataValidator: test product files", func() {
			metadata.ProductFiles[0].File = ""
			mv := NewMetaDataValidator(metadata)
			Expect(mv.Validate()).To(BeFalse())
			Expect(mv.GetErrorMessages()).To(Equal([]string{
				"One of settings is empty.[file, upload_as, file_type, file_version]",
				"value is empty, index=0 (|  | Open Source Licenses for GPDB 6.x | Open Source License | ${VERSION_REGEX} |)",
			}))
		})
	})

	Context("OffsetValidator", func() {

		It("Test OffsetValidator: correct pattern", func() {
			correctPattern := []string{
				"+1d",
				"+1D",
				"+2d",
				"+2D",
				"+3m",
				"+3M",
				"+4y",
				"+4Y",

				"+10d",
				"+10D",
				"+20d",
				"+20D",
				"+30m",
				"+30M",
				"+40y",
				"+40Y",

				"+1y+2m+3d",
				"+10y+20m+30d",
			}

			for _, pattern := range correctPattern {
				ov := NewOffsetValidator(pattern)
				Expect(ov.Validate()).To(BeTrue())
			}

		})

		It("Test OffsetValidator: incorrect pattern", func() {
			incorrectPattern := []string{
				"1d",
				"1day",
				"-1d",
				"-1day",
				"1y2m3d",
			}

			for _, pattern := range incorrectPattern {
				ov := NewOffsetValidator(pattern)
				Expect(ov.Validate()).To(BeFalse())
				Expect(ov.GetErrorMessages()[0]).To(Equal(fmt.Sprintf("\"%s\" is not a valid offset value", pattern)))
			}
		})
	})

	Context("RequiredValidator", func() {
		It("Test RequiredValidator: single value", func() {
			rv := NewRequiredValidator("abc")
			Expect(rv.Validate()).To(BeTrue())
			Expect(len(rv.GetErrorMessages())).To(Equal(0))
		})

		It("Test RequiredValidator: empty string", func() {
			rv := NewRequiredValidator("")
			Expect(rv.Validate()).To(BeFalse())
			Expect(rv.GetErrorMessages()[0]).To(Equal("value is empty, index=0 (|  |)"))
		})

		It("Test RequiredValidator: multiple string", func() {
			rv := NewRequiredValidator("abc", "def", "123")
			Expect(rv.Validate()).To(BeTrue())
			Expect(len(rv.GetErrorMessages())).To(Equal(0))
		})

		It("Test RequiredValidator: multiple string that include empty string", func() {
			rv := NewRequiredValidator("abc", "def", "123", "", "\t", "xyz")
			Expect(rv.Validate()).To(BeFalse())
			Expect(len(rv.GetErrorMessages())).To(Equal(2))
			Expect(rv.GetErrorMessages()).To(Equal([]string{
				"value is empty, index=3 (| abc | def | 123 |  | \t | xyz |)",
				"value is empty, index=4 (| abc | def | 123 |  | \t | xyz |)",
			}))
		})

		It("Test RequiredValidator: multiple string that include non-string type", func() {
			rv := NewRequiredValidator("abc", "def", "123", "", "\t", 123, "xyz", true, 12.4)
			Expect(rv.Validate()).To(BeFalse())
			Expect(len(rv.GetErrorMessages())).To(Equal(5))
			Expect(rv.GetErrorMessages()).To(Equal([]string{
				"value is empty, index=3 (| abc | def | 123 |  | \t | 123 | xyz | true | 12.4 |)",
				"value is empty, index=4 (| abc | def | 123 |  | \t | 123 | xyz | true | 12.4 |)",
				"RequiredValidator only support string type, and not support this type: int",
				"RequiredValidator only support string type, and not support this type: bool",
				"RequiredValidator only support string type, and not support this type: float64",
			}))
		})
	})

	Context("ProductFileValidator", func() {
		It("ProductFileValidator: empty product file", func() {
			pv := NewProductFileValidator(config.ProductFile{})
			Expect(pv.Validate()).To(BeFalse())
			Expect(pv.GetErrorMessages()).To(Equal([]string{
				"One of settings is empty.[file, upload_as, file_type, file_version]",
				"value is empty, index=0 (|  |  |  |  |)",
				"value is empty, index=1 (|  |  |  |  |)",
				"value is empty, index=2 (|  |  |  |  |)",
				"value is empty, index=3 (|  |  |  |  |)",
			}))
		})

		It("ProductFileValidator: empty product file", func() {
			pf := config.ProductFile{}
			pf.File = "file://path/to/file"
			pf.FileVersion = "3.4.2"

			pv := NewProductFileValidator(pf)
			Expect(pv.Validate()).To(BeFalse())
			Expect(pv.GetErrorMessages()).To(Equal([]string{
				"One of settings is empty.[file, upload_as, file_type, file_version]",
				"value is empty, index=1 (| file://path/to/file |  |  | 3.4.2 |)",
				"value is empty, index=2 (| file://path/to/file |  |  | 3.4.2 |)",
			}))
		})
	})
})
