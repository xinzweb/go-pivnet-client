package utils

import (
	"fmt"
	"github.com/baotingfang/go-pivnet-client/vlog"
	"github.com/pivotal-cf/go-pivnet/v4"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func ExistsPath(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func GetAllFiles(dirPath string) []string {
	var files []string

	if !ExistsPath(dirPath) {
		return files
	}

	if !IsDir(dirPath) {
		return files
	}

	_ = filepath.Walk(dirPath,
		func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				files = append(files, path)
			}
			return nil
		},
	)
	return files
}

func Empty(s interface{}) bool {
	switch v := s.(type) {
	case string:
		return len(strings.TrimSpace(s.(string))) == 0
	case pivnet.Release:
		return s.(pivnet.Release) == pivnet.Release{}
	case pivnet.ReleaseType:
		return strings.TrimSpace(string(s.(pivnet.ReleaseType))) == ""
	default:
		vlog.Fatal("Empty() doesn't support this type: %T", v)
	}
	return false
}

type Date struct {
	time.Time
}

func (d Date) String() string {
	return d.Time.Format("2006-01-02")
}

func (d Date) LastDayOfCurrentMonth() Date {
	year, month, _ := d.Date()
	return Date{time.Date(year, month+1, 0, 0, 0, 0, 0, d.Location())}
}

func (d Date) Offset(offsetExpression string) Date {
	if Empty(offsetExpression) {
		return d
	}
	dayOffset := 0
	monthOffset := 0
	yearOffset := 0

	matcher := func(patter string) int {
		dayOffsetMatcher := regexp.MustCompile(patter)
		r := dayOffsetMatcher.FindStringSubmatch(strings.ToLower(offsetExpression))
		if len(r) > 0 {
			offset, _ := strconv.Atoi(r[1])
			return offset
		}
		return 0
	}

	dayOffset = matcher(`\+(\d*)d`)
	monthOffset = matcher(`\+(\d*)m`)
	yearOffset = matcher(`\+(\d*)y`)

	return Date{d.Time.AddDate(yearOffset, monthOffset, dayOffset)}
}

func MustParseDateFrom(value string) Date {
	t, err := time.Parse("2006-01-02", value)
	if err != nil {
		errMessage := fmt.Sprintf("can not parse %s to Date. err: %s", value, err.Error())
		vlog.Error(errMessage)
		vlog.Fatal(errMessage)
	}
	return Date{t}
}

func ParseDateFrom(value string) (Date, error) {
	if Empty(value) {
		return Date{}, fmt.Errorf("can not parse empty string")
	}

	t, err := time.Parse("2006-01-02", value)
	if err != nil {
		return Date{}, fmt.Errorf("parse Date failed: %s, err: %s", value, err.Error())
	}

	return Date{t}, nil
}

func IsDate(value string) bool {
	_, err := ParseDateFrom(value)
	return err == nil
}
