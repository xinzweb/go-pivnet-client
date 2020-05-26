package utils

import (
	"github.com/baotingfang/go-pivnet-client/vlog"
	"net/url"
	"path"
)

func UrlJoin(baseUrl string, paths ...string) string {
	u, err := url.Parse(baseUrl)
	if err != nil {
		vlog.Fatal(baseUrl)
	}
	u.Path = path.Join(u.Path, path.Join(paths...))
	return u.String()
}
