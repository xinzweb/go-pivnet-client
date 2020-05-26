package service_test

import (
	semver "github.com/cppforlife/go-semi-semantic/version"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/baotingfang/go-pivnet-client/service"
)

var _ = Describe("Uploader", func() {
	Context("NewVersionReplacer", func() {
		It("NewVersionReplacer: empty string", func() {
			resolvedFile := ResolvedFile{
				LocalFilePath:   "/tmp/path/file-1.0.0.txt",
				LocalFileName:   "file.txt",
				ResolvedVersion: semver.Version{},
			}
			vr := NewVersionReplacer(resolvedFile)
			Expect(vr.Replace("")).To(Equal(""))
		})

		It("NewVersionReplacer: expression doesn't include ${VERSION_REGEX}", func() {
			resolvedFile := ResolvedFile{
				LocalFilePath:   "/tmp/path/file-1.0.0.txt",
				LocalFileName:   "file.txt",
				ResolvedVersion: semver.Version{},
			}
			vr := NewVersionReplacer(resolvedFile)
			Expect(vr.Replace("1.2.0")).To(Equal("1.2.0"))
		})

		It("NewVersionReplacer: resolvedFile is empty", func() {
			resolvedFile := ResolvedFile{
				LocalFilePath:   "/tmp/path/file-1.0.0.txt",
				LocalFileName:   "file.txt",
				ResolvedVersion: semver.Version{},
			}
			vr := NewVersionReplacer(resolvedFile)
			f := func() {
				vr.Replace("abc-${VERSION_REGEX}")
			}
			Expect(f).To(PanicWith(`[Default Logger][FATAL] resolved version is empty`))
		})

		It("NewVersionReplacer: correct replace action", func() {
			resolvedFile := ResolvedFile{
				LocalFilePath:   "/tmp/path/file-1.0.0.txt",
				LocalFileName:   "file.txt",
				ResolvedVersion: semver.MustNewVersionFromString("1.0.0"),
			}
			vr := NewVersionReplacer(resolvedFile)
			Expect(vr.Replace("abc-${VERSION_REGEX}-DEF")).To(Equal("abc-1.0.0-DEF"))
		})

	})
})
