package utils_test

import (
	. "github.com/baotingfang/go-pivnet-client/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/go-pivnet/v4"
)

var _ = Describe("Tools", func() {
	It("Test Offset", func() {
		baseDate, _ := ParseDateFrom("2013-05-19")
		d := baseDate.Offset("+3d")
		Expect(d.String()).To(Equal("2013-05-22"))

		d = baseDate.Offset("+18m")
		Expect(d.String()).To(Equal("2014-11-19"))

		d = baseDate.Offset("+1y+3d")
		Expect(d.String()).To(Equal("2014-05-22"))

		d = baseDate.Offset("+1y+2m+3d")
		Expect(d.String()).To(Equal("2014-07-22"))

		d = baseDate.Offset("")
		Expect(d.String()).To(Equal("2013-05-19"))
	})

	It("Test LastDayOfCurrentMonth", func() {
		inputs := []string{
			"2013-05-19",
			"1983-02-11",
			"1984-04-06",
			"2008-08-08",
		}

		results := []string{
			"2013-05-31",
			"1983-02-28",
			"1984-04-30",
			"2008-08-31",
		}
		for i := 0; i < 4; i++ {
			d1, _ := ParseDateFrom(inputs[i])
			d := d1.LastDayOfCurrentMonth()
			Expect(d.String()).To(Equal(results[i]))
		}
	})

	It("Test Empty", func() {
		s := "     "
		Expect(len(s)).To(BeNumerically(">", 0))
		Expect(Empty(s)).To(BeTrue())

		// 2 tab 4 whitespaces
		s = "	     	"
		Expect(len(s)).To(BeNumerically(">", 0))
		Expect(Empty(s)).To(BeTrue())

		s = "abc"
		Expect(len(s)).To(BeNumerically(">", 0))
		Expect(Empty(s)).To(BeFalse())

		r := pivnet.Release{}
		Expect(Empty(r)).To(BeTrue())

		r = pivnet.Release{
			ID:      1000,
			Version: "4.6.0",
		}
		Expect(Empty(r)).To(BeFalse())
	})

	It("Test ParseDateFrom", func() {
		By("Test empty string")
		d, err := ParseDateFrom("")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("can not parse empty string"))
		Expect(d.IsZero()).To(BeTrue())

		By("Test correct format")
		d, err = ParseDateFrom("2013-05-19")
		Expect(err).NotTo(HaveOccurred())
		Expect(d.String()).To(Equal("2013-05-19"))

		By("Test wrong format")
		d, err = ParseDateFrom("2013/05/19")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal(`parse Date failed: 2013/05/19, err: parsing time "2013/05/19" as "2006-01-02": cannot parse "/05/19" as "-"`))
		Expect(d.IsZero()).To(BeTrue())
	})

	It("Test MustParseDateFrom", func() {
		d := MustParseDateFrom("2013-05-19")
		Expect(d.String()).To(Equal("2013-05-19"))
	})
})
