package config_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strings"

	. "github.com/baotingfang/go-pivnet-client/config"
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
			metaData, err := MetadataFrom(metadataReader)
			Expect(err).NotTo(HaveOccurred())

			Expect(metaData.Release.ReleaseType).To(Equal("Major Release"))
			Expect(metaData.Release.EulaSlug).To(Equal("pivotal_software_eula"))
			Expect(len(metaData.FileGroups)).To(Equal(1))
			Expect(len(metaData.FileGroups[0].ProductFiles)).To(Equal(2))
			Expect(len(metaData.ProductFiles)).To(Equal(2))
		})
	})
})
