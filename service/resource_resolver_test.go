package service_test

import (
	"github.com/baotingfang/go-pivnet-client/service/servicefakes"
	semver "github.com/cppforlife/go-semi-semantic/version"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/baotingfang/go-pivnet-client/service"
)

var _ = Describe("ResourceResolver", func() {
	var fakeWalker servicefakes.FakeWalker
	BeforeEach(func() {
		fakeWalker = servicefakes.FakeWalker{}
	})
	It("search path is dot (current dir)", func() {
		fakeWalker.GetAllFilesReturns([]string{
			"server-sles11sp4/greenplum-db-6.6.7-sles11-x86_64.zip",
		})
		resolver := NewResourceResolver(".", &fakeWalker)
		rv, err := resolver.Resolve("file://server-sles11sp4/greenplum-db-(.*)-sles11-x86_64.zip")
		Expect(err).NotTo(HaveOccurred())
		Expect(rv).To(Equal(ResolvedFile{
			LocalFilePath:   "server-sles11sp4/greenplum-db-6.6.7-sles11-x86_64.zip",
			LocalFileName:   "greenplum-db-6.6.7-sles11-x86_64.zip",
			ResolvedVersion: semver.MustNewVersionFromString("6.6.7"),
		}))
	})

	It("search path is empty (current dir)", func() {
		fakeWalker.GetAllFilesReturns([]string{
			"server-sles11sp4/greenplum-db-6.6.7-sles11-x86_64.zip",
		})
		resolver := NewResourceResolver("", &fakeWalker)
		rv, err := resolver.Resolve("file://server-sles11sp4/greenplum-db-(.*)-sles11-x86_64.zip")
		Expect(err).NotTo(HaveOccurred())
		Expect(rv).To(Equal(ResolvedFile{
			LocalFilePath:   "server-sles11sp4/greenplum-db-6.6.7-sles11-x86_64.zip",
			LocalFileName:   "greenplum-db-6.6.7-sles11-x86_64.zip",
			ResolvedVersion: semver.MustNewVersionFromString("6.6.7"),
		}))
	})

	It("search path is not empty", func() {
		fakeWalker.GetAllFilesReturns([]string{
			"/tmp/server-sles11sp4/greenplum-db-6.6.7-sles11-x86_64.zip",
		})
		resolver := NewResourceResolver("/tmp", &fakeWalker)
		rv, err := resolver.Resolve("file://server-sles11sp4/greenplum-db-(.*)-sles11-x86_64.zip")
		Expect(err).NotTo(HaveOccurred())
		Expect(rv).To(Equal(ResolvedFile{
			LocalFilePath:   "/tmp/server-sles11sp4/greenplum-db-6.6.7-sles11-x86_64.zip",
			LocalFileName:   "greenplum-db-6.6.7-sles11-x86_64.zip",
			ResolvedVersion: semver.MustNewVersionFromString("6.6.7"),
		}))
	})

	It("only support file schema", func() {
		resolver := NewResourceResolver("/tmp", FilesWalker{})
		rv, err := resolver.Resolve("gs://server-sles11sp4/greenplum-db-(.*)-sles11-x86_64.zip")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("not support schema:gs"))
		Expect(rv).To(Equal(ResolvedFile{}))
	})

	It("file name doesn't include regexp group", func() {
		resolver := NewResourceResolver("/tmp", FilesWalker{})
		rv, err := resolver.Resolve("file://server-sles11sp4/greenplum-db-6.6.7-sles11-x86_64.zip")
		Expect(err).NotTo(HaveOccurred())
		Expect(rv).To(Equal(ResolvedFile{
			LocalFilePath: "/tmp/server-sles11sp4/greenplum-db-6.6.7-sles11-x86_64.zip",
			LocalFileName: "greenplum-db-6.6.7-sles11-x86_64.zip",
		}))
	})

	It("file name doesn't include regexp group, search path is empty", func() {
		resolver := NewResourceResolver("", FilesWalker{})
		rv, err := resolver.Resolve("file://server-sles11sp4/greenplum-db-6.6.7-sles11-x86_64.zip")
		Expect(err).NotTo(HaveOccurred())
		Expect(rv).To(Equal(ResolvedFile{
			LocalFilePath: "server-sles11sp4/greenplum-db-6.6.7-sles11-x86_64.zip",
			LocalFileName: "greenplum-db-6.6.7-sles11-x86_64.zip",
		}))
	})

	It("no match", func() {
		fakeWalker.GetAllFilesReturns([]string{
			"/tmp/server-sles11sp4/no_file.zip",
		})
		resolver := NewResourceResolver("", &fakeWalker)
		rv, err := resolver.Resolve("file://server-sles11sp4/greenplum-db-(.*)-sles11-x86_64.zip")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("can not match file"))
		Expect(rv).To(Equal(ResolvedFile{}))
	})

	It("multiple match", func() {
		fakeWalker.GetAllFilesReturns([]string{
			"/tmp/server-sles11sp4/greenplum-db-6.6.0-sles11-x86_64.zip",
			"/tmp/server-sles11sp4/greenplum-db-6.7.0-sles11-x86_64.zip",
		})
		resolver := NewResourceResolver("", &fakeWalker)
		rv, err := resolver.Resolve("file://server-sles11sp4/greenplum-db-(.*)-sles11-x86_64.zip")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("match multiple files"))
		Expect(rv).To(Equal(ResolvedFile{}))
	})
})
