package utils

import (
	"fmt"
	"github.com/baotingfang/go-pivnet-client/vlog"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func UrlJoin(baseUrl string, paths ...string) string {
	u, err := url.Parse(baseUrl)
	if err != nil {
		vlog.Fatal(baseUrl)
	}
	u.Path = path.Join(u.Path, path.Join(paths...))
	return u.String()
}

func IsEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func EndDayOfCurrentMonth(t Date) Date {
	year, month, _ := t.Date()
	return Date{time.Date(year, month+1, 0, 0, 0, 0, 0, t.Location())}
}

func ComputeFromOffset(base Date, offsetExpression string) Date {
	if IsEmpty(offsetExpression) {
		return base
	}
	dayOffset := 0
	monthOffset := 0
	yearOffset := 0

	matcher := func(patter, value string) int {
		dayOffsetMatcher := regexp.MustCompile(patter)
		r := dayOffsetMatcher.FindStringSubmatch(offsetExpression)
		if len(r) > 0 {
			offset, _ := strconv.Atoi(r[1])
			return offset
		}
		return 0
	}

	dayOffset = matcher(`\+(\d*)d`, offsetExpression)
	monthOffset = matcher(`\+(\d*)m`, offsetExpression)
	yearOffset = matcher(`\+(\d*)y`, offsetExpression)

	return Date{base.Time.AddDate(yearOffset, monthOffset, dayOffset)}
}

type Date struct {
	time.Time
}

func (d Date) MarshalYAML() (interface{}, error) {
	if d.IsZero() {
		return nil, nil
	}

	return d.Time.Format("2006-01-02"), nil
}

func (d *Date) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var data time.Time
	err := unmarshal(&data)
	if err != nil {
		return err
	}
	d.Time = data
	return nil
}

func (d *Date) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		d.Time = time.Time{}
		return
	}
	d.Time, err = time.Parse("2006-01-02", s)
	return
}

func (d Date) MarshalJSON() ([]byte, error) {
	if d.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", d.Time.Format("2006-01-02"))), nil
}

func (d Date) String() string {
	return d.Time.Format("2006-01-02")
}

func TestHelper(expectedStringRegexp string) {
	if err := recover(); err != nil {
		message := fmt.Sprintf("%v", err)
		match, err := regexp.Match(expectedStringRegexp, []byte(message))
		if !match || err != nil {
			panic("Test failed in test helper function: " + err.Error())
		}
	}
}
