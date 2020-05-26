package config

import (
	"gopkg.in/yaml.v2"
	"io"
	"time"
)

type Release struct {
	ReleaseType      string    `yaml:"release_type,omitempty"`
	EulaSlug         string    `yaml:"eula_slug"`
	Description      string    `yaml:"description"`
	ReleaseNotesUrl  string    `yaml:"release_notes_url,omitempty"`
	Availability     string    `yaml:"availability"`
	Controlled       bool      `yaml:"controlled"`
	Eccn             string    `yaml:"eccn"`
	LicenseException string    `yaml:"license_exception"`
	ReleaseDate      time.Time `yaml:"release_date,omitempty"`
}

type FileGroup struct {
	Name         string        `yaml:"name"`
	ProductFiles []ProductFile `yaml:"product_files"`
}

type ProductFile struct {
	File               string `yaml:"file"`
	UploadAs           string `yaml:"upload_as"`
	Description        string `yaml:"description"`
	FileType           string `yaml:"file_type"`
	DocsUrl            string `yaml:"docs_url,omitempty"`
	SystemRequirements string `yaml:"system_requirements,omitempty"`
	Platforms          string `yaml:"platforms,omitempty"`
	IncludedFiles      string `yaml:"included_files,omitempty"`
	FileVersion        string `yaml:"file_version"`
}

type Metadata struct {
	Release      Release       `yaml:"release"`
	FileGroups   []FileGroup   `yaml:"file_groups"`
	ProductFiles []ProductFile `yaml:"product_files"`
}

func MetadataFrom(reader io.Reader) (Metadata, error) {

	var metadata Metadata
	if err := yaml.NewDecoder(reader).Decode(&metadata); err != nil {
		return Metadata{}, err
	}

	return metadata, nil
}
