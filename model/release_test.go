package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRelease_NextVersion_date(t *testing.T) {
	r := &Release{VersionPattern: "DATE"}

	actual, err := r.NextVersion()
	assert.NoError(t, err)
	assert.Regexp(t, DateRegexp, actual)
}

func TestRelease_NextVersion_semver(t *testing.T) {
	r := &Release{
		VersionPattern: "SEMVER",
		CurrentVersion: "1.0",
		Changelog: []ReleaseItem{
			NewReleaseItem("abc", "Picsou", time.Now(), "feat: gain more money"),
		},
	}

	actual, err := r.NextVersion()
	assert.NoError(t, err)
	assert.Equal(t, "1.1.0", actual)
}

func TestRelease_NextVersion_unknown(t *testing.T) {
	r := &Release{VersionPattern: "bad"}

	_, err := r.NextVersion()
	assert.Error(t, err)
}
