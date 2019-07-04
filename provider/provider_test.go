package provider

import (
	"github.com/stretchr/testify/assert"
	"github.com/tauffredou/nextver/model"
	"log"
	"testing"
	"time"
)

var (
	changeLog = []model.ReleaseItem{
		model.NewReleaseItem("tauf", testDate, "feat(feature-1234): commit message"),
		model.NewReleaseItem("tauf", testDate, "fix(feature-1234): commit message"),
		model.NewReleaseItem("tauf", testDate, "empty"),
	}
)

var (
	testDate = MustParse(time.RFC3339, "2010-01-02T12:34:00Z")
)

func MustParse(layout, value string) time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		log.Fatal(err)
	}
	return t
}

func TestRelease_NextVersion_withSemver_empty(t *testing.T) {
	r := model.Release{
		Changelog:      changeLog,
		CurrentVersion: "",
		VersionPattern: "vSEMVER",
	}

	v, _ := r.NextVersion()

	assert.Equal(t, "v0.1.0", v)

}

func TestRelease_NextVersion_withSemver_withPrefix(t *testing.T) {
	r := model.Release{
		Changelog:      changeLog,
		CurrentVersion: "v2.0.1",
		VersionPattern: "vSEMVER",
	}

	v, _ := r.NextVersion()

	assert.Equal(t, "v2.1.0", v)

}

func TestRelease_NextVersion_withSemver_withoutPrefix(t *testing.T) {
	r := model.Release{
		Changelog:      changeLog,
		CurrentVersion: "2.0.1",
		VersionPattern: "SEMVER",
	}

	v, _ := r.NextVersion()

	assert.Equal(t, "2.1.0", v)

}

func TestRelease_NextVersion_withDate(t *testing.T) {
	r := model.Release{
		Changelog:      changeLog,
		CurrentVersion: "release-2019-03-29-161011 ",
		VersionPattern: "release-DATE",
	}

	v, _ := r.NextVersion()

	assert.Regexp(t, `release-\d{4}-\d{2}-\d{2}-\d{6}`, v)

}

func TestRelease_NextVersion_withEmptyChangelog(t *testing.T) {
	changeLog = []model.ReleaseItem{}

	r := model.Release{
		Changelog:      changeLog,
		CurrentVersion: "v1.0.0 ",
		VersionPattern: "vSEMVER",
	}

	v, _ := r.NextVersion()

	assert.Regexp(t, "v1.0.0", v)

}
