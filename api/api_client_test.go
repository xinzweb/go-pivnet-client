package api_test

import (
	"errors"
	. "github.com/baotingfang/go-pivnet-client/api"
	"github.com/baotingfang/go-pivnet-client/config"
	"github.com/baotingfang/go-pivnet-client/gp"
	"github.com/baotingfang/go-pivnet-client/wrapper/wrapperfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/go-pivnet/v4"
)

var _ = Describe("ApiClient", func() {
	var (
		context          gp.Context
		apiClient        AccessClient
		fakePivnetClient wrapperfakes.FakePivnetClient
	)

	BeforeEach(func() {
		context = gp.NewContext("http://fakesite/", "fakeslug", "faketoken", true, true)
		fakePivnetClient = wrapperfakes.FakePivnetClient{}
		context.Client = &fakePivnetClient
		apiClient = NewApiClient(context)
	})

	Context("GetLatestPublicReleaseByReleaseType", func() {
		It("Can not found previous release, due to GetAllReleases failed", func() {
			fakePivnetClient.GetAllReleasesReturns([]pivnet.Release{}, errors.New("failed get all releases"))
			context.Client = &fakePivnetClient

			r, err := apiClient.GetLatestPublicReleaseByReleaseType(4, config.MajorReleaseType)
			Expect(err).To(HaveOccurred())
			Expect(r.Version).To(BeEmpty())
			Expect(err.Error()).To(Equal("failed get all releases"))

		})

		It("Can not found previous release", func() {
			context.Client = &fakePivnetClient

			r, err := apiClient.GetLatestPublicReleaseByReleaseType(4, config.MajorReleaseType)
			Expect(err).To(HaveOccurred())
			Expect(r.Version).To(BeEmpty())
			Expect(err.Error()).To(Equal("can not found previous release. major version: 4, release type: Major Release"))

		})

		It("for gpdb4/5/6", func() {
			fakePivnetClient.GetAllReleasesReturns([]pivnet.Release{
				{
					Version:      "4.3.2.1",
					Availability: "All Users",
					ReleaseType:  config.MaintenanceReleaseType,
				},
				{
					Version:      "4.3.30.4",
					Availability: "All Users",
					ReleaseType:  config.MaintenanceReleaseType,
				},
				{
					Version:      "4.3.26.0",
					Availability: "All Users",
					ReleaseType:  config.MinorReleaseType,
				},
				{
					Version:      "4.3.27.0",
					Availability: "All Users",
					ReleaseType:  config.MinorReleaseType,
				},
				{
					Version:      "5.0.0",
					Availability: "All Users",
					ReleaseType:  config.MajorReleaseType,
				},
				{
					Version:      "5.26.0",
					Availability: "All Users",
					ReleaseType:  config.MinorReleaseType,
				},
				{
					Version:      "5.27.0",
					Availability: "All Users",
					ReleaseType:  config.MinorReleaseType,
				},
				{
					Version:      "5.27.0",
					Availability: "Admin Only",
					ReleaseType:  config.MinorReleaseType,
				},
				{
					Version:      "5.26.1",
					Availability: "All Users",
					ReleaseType:  config.MaintenanceReleaseType,
				},
				{
					Version:      "5.27.1",
					Availability: "All Users",
					ReleaseType:  config.MaintenanceReleaseType,
				},
				{
					Version:      "5.27.2",
					Availability: "Admin Only",
					ReleaseType:  config.MaintenanceReleaseType,
				},
				{
					Version:      "6.0.0",
					Availability: "All Users",
					ReleaseType:  config.MajorReleaseType,
				},
				{
					Version:      "6.6.0",
					Availability: "All Users",
					ReleaseType:  config.MinorReleaseType,
				},
				{
					Version:      "6.7.0",
					Availability: "All Users",
					ReleaseType:  config.MinorReleaseType,
				},
				{
					Version:      "6.8.0",
					Availability: "Admin Only",
					ReleaseType:  config.MinorReleaseType,
				},
				{
					Version:      "6.6.1",
					Availability: "All Users",
					ReleaseType:  config.MaintenanceReleaseType,
				},
				{
					Version:      "6.7.1",
					Availability: "All Users",
					ReleaseType:  config.MaintenanceReleaseType,
				},
				{
					Version:      "6.7.2",
					Availability: "Admin Only",
					ReleaseType:  config.MaintenanceReleaseType,
				},
			}, nil)

			By("for gpdb4")
			r, err := apiClient.GetLatestPublicReleaseByReleaseType(4, config.MinorReleaseType)
			Expect(err).NotTo(HaveOccurred())
			Expect(r.Version).To(Equal("4.3.27.0"))

			r, err = apiClient.GetLatestPublicReleaseByReleaseType(4, config.MaintenanceReleaseType)
			Expect(err).NotTo(HaveOccurred())
			Expect(r.Version).To(Equal("4.3.30.4"))

			By("for gpdb5")
			r, err = apiClient.GetLatestPublicReleaseByReleaseType(5, config.MajorReleaseType)
			Expect(err).NotTo(HaveOccurred())
			Expect(r.Version).To(Equal("5.0.0"))

			r, err = apiClient.GetLatestPublicReleaseByReleaseType(5, config.MinorReleaseType)
			Expect(err).NotTo(HaveOccurred())
			Expect(r.Version).To(Equal("5.27.0"))

			r, err = apiClient.GetLatestPublicReleaseByReleaseType(5, config.MaintenanceReleaseType)
			Expect(err).NotTo(HaveOccurred())
			Expect(r.Version).To(Equal("5.27.1"))

			By("for gpdb6")
			r, err = apiClient.GetLatestPublicReleaseByReleaseType(6, config.MajorReleaseType)
			Expect(err).NotTo(HaveOccurred())
			Expect(r.Version).To(Equal("6.0.0"))

			r, err = apiClient.GetLatestPublicReleaseByReleaseType(6, config.MinorReleaseType)
			Expect(err).NotTo(HaveOccurred())
			Expect(r.Version).To(Equal("6.7.0"))

			r, err = apiClient.GetLatestPublicReleaseByReleaseType(6, config.MaintenanceReleaseType)
			Expect(err).NotTo(HaveOccurred())
			Expect(r.Version).To(Equal("6.7.1"))
		})
	})

	Context("FileTransferStatusInProgress", func() {
		It("Test in progress", func() {
			fakePivnetClient.GetProductFileReturns(pivnet.ProductFile{
				FileTransferStatus: "in_progress",
			}, nil)

			result := apiClient.FileTransferStatusInProgress(1)
			Expect(result).To(BeTrue())
		})

		It("Test not in progress", func() {
			fakePivnetClient.GetProductFileReturns(pivnet.ProductFile{
				FileTransferStatus: "not_in_progress",
			}, nil)

			result := apiClient.FileTransferStatusInProgress(1)
			Expect(result).To(BeFalse())
		})

		It("Test it is failed when get product file", func() {
			fakePivnetClient.GetProductFileReturns(
				pivnet.ProductFile{}, errors.New("can not find product"))

			f := func() {
				apiClient.FileTransferStatusInProgress(1)
			}
			Expect(f).To(PanicWith(`[Default Logger][FATAL] can not find product`))
		})
	})
})
