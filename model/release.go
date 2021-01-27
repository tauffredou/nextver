package model

import (
	"errors"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	UNDEFINED = 0
	PATCH     = 1
	MINOR     = PATCH << 1
	MAJOR     = MINOR << 1
)
const (
	SemverRegex              = `v?(\d+)(\.(\d+))?(\.(\d+))?`
	DateRegexp               = `\d{4}-\d{2}-\d{2}-\d{6}`
	ConventionalCommitRegexp = `^([a-zA-Z-_]+)(\(([^\):]+)\))?[ ]?: ?(.*)$`
	FirstVersion             = "0.0.0"
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
	log.WithField("VersionPattern", r.VersionPattern).Debug("NextVersion")

	if strings.Contains(r.VersionPattern, "SEMVER") {
		return SemverCalculator(r)
	}

	if strings.Contains(r.VersionPattern, "DATE") {
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
