package provider

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	changeLog = []ReleaseItem{
		NewReleaseItem("feat(feature-1234): commit message"),
		NewReleaseItem("fix(feature-1234): commit message"),
		NewReleaseItem("empty"),
	}
)

func TestRelease_NextVersion_withSemver_withPrefix(t *testing.T) {
	r := Release{
		Changelog:      changeLog,
		CurrentVersion: "v2.0.1",
		VersionPattern: "vSEMVER",
	}

	v, _ := r.NextVersion()

	assert.Equal(t, "v2.1.0", v)

}

func TestRelease_NextVersion_withSemver_withoutPrefix(t *testing.T) {
	r := Release{
		Changelog:      changeLog,
		CurrentVersion: "2.0.1",
		VersionPattern: "SEMVER",
	}

	v, _ := r.NextVersion()

	assert.Equal(t, "2.1.0", v)

}

func TestRelease_NextVersion_withDate(t *testing.T) {
	r := Release{
		Changelog:      changeLog,
		CurrentVersion: "release-2019-03-29-161011 ",
		VersionPattern: "release-DATE",
	}

	v, _ := r.NextVersion()

	assert.Regexp(t, `release-\d{4}-\d{2}-\d{2}-\d{6}`, v)

}
