package provider

import (
	"errors"
	"log"
	"strings"
)

type Release struct {
	Project           string        `json:"project"`
	Ref               string        `json:"ref"`
	CurrentVersion    string        `json:"current_version"`
	Changelog         []ReleaseItem `json:"changelog"`
	versionCalculator func(*Release) (string, error)
	VersionPattern    string `json:"version_pattern"`
}

//NextVersion calculates next semver version from commits
func (r *Release) NextVersion() (string, error) {
	switch {
	case strings.Contains(r.VersionPattern, "SEMVER"):
		return SemverCalculator(r)
	case strings.Contains(r.VersionPattern, "DATE"):
		return DateVersionCalculator(r)
	}
	return "", errors.New("unknown version calculator")
}

func (r *Release) MustNextVersion() string {
	version, err := r.NextVersion()
	if err != nil {
		log.Fatal(err)
	}
	return version
}
