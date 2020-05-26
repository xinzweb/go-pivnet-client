package utils_test

import (
	"encoding/json"
	. "github.com/baotingfang/go-pivnet-client/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"
	"time"
)

var _ = Describe("Tools", func() {
	It("UrlJoin", func() {
		baseUrl := "http://www.example.com/"
		p := UrlJoin(baseUrl, "/abc", "/def/", "/xyz")
		Expect(p).To(Equal("http://www.example.com/abc/def/xyz"))
	})

	It("Test Date struct unmarshal json", func() {
		jsonStr := `{"create_at":"2013-05-19"}`
		jsonData := struct {
			CreateAt Date `json:"create_at"`
		}{}

		err := json.Unmarshal([]byte(jsonStr), &jsonData)
		Expect(err).NotTo(HaveOccurred())
		year, month, day := jsonData.CreateAt.Time.Date()
		Expect(year).To(Equal(2013))
		Expect(month).To(Equal(time.May))
		Expect(day).To(Equal(19))

		output, err := json.Marshal(jsonData)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(output)).To(Equal(`{"create_at":"2013-05-19"}`))
	})

	It("Test Date struct unmarshal json: null", func() {
		jsonStr := `{"create_at": null}`
		jsonData := struct {
			CreateAt Date `json:"create_at"`
		}{}

		err := json.Unmarshal([]byte(jsonStr), &jsonData)
		Expect(err).NotTo(HaveOccurred())
		Expect(jsonData.CreateAt.IsZero()).To(BeTrue())

		output, err := json.Marshal(jsonData)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(output)).To(Equal(`{"create_at":null}`))
	})

	It("Test Date struct unmarshal yaml", func() {
		yamlStr := `create_at: 2013-05-19`
		yamlData := struct {
			CreateAt Date `yaml:"create_at"`
		}{}

		err := yaml.Unmarshal([]byte(yamlStr), &yamlData)
		Expect(err).NotTo(HaveOccurred())
		year, month, day := yamlData.CreateAt.Time.Date()
		Expect(year).To(Equal(2013))
		Expect(month).To(Equal(time.May))
		Expect(day).To(Equal(19))

		output, err := yaml.Marshal(&yamlData)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(output)).To(Equal("create_at: \"2013-05-19\"\n"))
	})

	It("Test Date struct unmarshal yaml: null", func() {
		yamlStr := `create_at: `
		yamlData := struct {
			CreateAt Date `yaml:"create_at"`
		}{}

		err := yaml.Unmarshal([]byte(yamlStr), &yamlData)
		Expect(err).NotTo(HaveOccurred())
		Expect(yamlData.CreateAt.IsZero()).To(BeTrue())

		output, err := yaml.Marshal(&yamlData)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(output)).To(Equal("create_at: null\n"))
	})

	It("Test ComputeFromOffset", func() {
		baseDate, _ := time.Parse("2006-01-02", "2013-05-19")
		d := ComputeFromOffset(Date{Time: baseDate}, "+3d")
		Expect(d.String()).To(Equal("2013-05-22"))

		d = ComputeFromOffset(Date{Time: baseDate}, "+18m")
		Expect(d.String()).To(Equal("2014-11-19"))

		d = ComputeFromOffset(Date{Time: baseDate}, "+1y+3d")
		Expect(d.String()).To(Equal("2014-05-22"))

		d = ComputeFromOffset(Date{Time: baseDate}, "+1y+2m+3d")
		Expect(d.String()).To(Equal("2014-07-22"))

		d = ComputeFromOffset(Date{Time: baseDate}, "")
		Expect(d.String()).To(Equal("2013-05-19"))
	})

	It("Test EndDayOfCurrentMonth", func() {
		generateDate := func(date string) Date {
			baseDate, _ := time.Parse("2006-01-02", date)
			return Date{Time: baseDate}
		}

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
			d1 := generateDate(inputs[i])
			d := EndDayOfCurrentMonth(d1)
			Expect(d.String()).To(Equal(results[i]))
		}
	})

	It("Test IsEmpty", func() {
		s := "     "
		Expect(len(s)).To(BeNumerically(">", 0))
		Expect(IsEmpty(s)).To(BeTrue())

		// 2 tab 4 whitespaces
		s = "	     	"
		Expect(len(s)).To(BeNumerically(">", 0))
		Expect(IsEmpty(s)).To(BeTrue())

		s = "abc"
		Expect(len(s)).To(BeNumerically(">", 0))
		Expect(IsEmpty(s)).To(BeFalse())
	})
})
