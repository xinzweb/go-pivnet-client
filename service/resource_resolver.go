package service

import (
	"errors"
	"fmt"
	"github.com/baotingfang/go-pivnet-client/utils"
	semver "github.com/cppforlife/go-semi-semantic/version"
	"net/url"
	"path"
	"regexp"
)

const (
	FileSchema = "file"
)

type ResolvedFile struct {
	LocalFilePath   string
	LocalFileName   string
	ResolvedVersion semver.Version
}

//go:generate counterfeiter . Resolver

type Resolver interface {
	Resolve(file string) (ResolvedFile, error)
}

//go:generate counterfeiter . Walker

type Walker interface {
	GetAllFiles(dirPath string) []string
}

type ResourceResolver struct {
	searchPath string
	walker     Walker
}

func NewResourceResolver(searchPath string, walker Walker) Resolver {
	return ResourceResolver{
		searchPath: searchPath,
		walker:     walker,
	}
}

func (r ResourceResolver) Resolve(file string) (ResolvedFile, error) {
	resourceUrl, err := url.Parse(file)
	if err != nil {
		return ResolvedFile{}, err
	}
	schema := resourceUrl.Scheme
	if schema == FileSchema {
		return resolveLocalFile(r, resourceUrl)
	} else {
		return ResolvedFile{}, errors.New("not support schema:" + schema)
	}
}

func resolveLocalFile(resolver ResourceResolver, resource *url.URL) (ResolvedFile, error) {
	rawFilePath := path.Join(resolver.searchPath, resource.Hostname(), resource.Path)
	rawFileName := path.Base(rawFilePath)

	if !isContainsGroupExpression(rawFilePath) {
		return ResolvedFile{
			LocalFilePath: rawFilePath,
			LocalFileName: rawFileName,
		}, nil
	}

	files := resolver.walker.GetAllFiles(path.Dir(rawFilePath))

	fileNameRegexp, err := regexp.Compile(rawFileName)
	if err != nil {
		return ResolvedFile{}, err
	}

	var matchedFiles []ResolvedFile
	for _, f := range files {
		if !fileNameRegexp.Match([]byte(f)) {
			continue
		}
		matches := fileNameRegexp.FindStringSubmatch(f)
		version, err := semver.NewVersionFromString(matches[1])
		if err != nil {
			return ResolvedFile{}, err
		}
		matchedFiles = append(matchedFiles, ResolvedFile{
			LocalFilePath:   f,
			LocalFileName:   path.Base(f),
			ResolvedVersion: version,
		})
	}

	if len(matchedFiles) == 0 {
		return ResolvedFile{}, fmt.Errorf("can not match file")
	}

	if len(matchedFiles) > 1 {
		return ResolvedFile{}, fmt.Errorf("match multiple files")
	}

	return matchedFiles[0], nil
}

func isContainsGroupExpression(value string) bool {
	r, _ := regexp.Match(`\(.*\)`, []byte(value))
	return r
}

type FilesWalker struct {
}

func (fw FilesWalker) GetAllFiles(dirPath string) []string {
	return utils.GetAllFiles(dirPath)
}
