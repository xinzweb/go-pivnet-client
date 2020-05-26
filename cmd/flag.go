package cmd

type FlagName int

//go:generate stringer -type FlagName -linecomment -output flag_string.go
const (
	FlagNameMetaFilePath FlagName = iota // metadata
	FlagNameSearchPath                   // search-path
	FlagNameVerbose                      // verbose
	FlagNameGpdbVersion                  // gpdb-version
)
