package config

type ReleaseType int

//go:generate stringer -type ReleaseType -linecomment -output release_type_string.go
const (
	AlphaRelease       ReleaseType = iota // Alpha Release
	BetaRelease                           // Beta Release
	MajorRelease                          // Major Release
	MinorRelease                          // Minor Release
	MaintenanceRelease                    // Maintenance Release
)
