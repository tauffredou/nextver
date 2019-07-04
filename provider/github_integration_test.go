// +build integration

package provider

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

var (
	intConfig = &GithubProviderConfig{
		Branch:  "master",
		Pattern: "vSEMVER",
	}
)

func getToken() string {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		var err error
		token, err = ReadHubToken(DefaultHubConfig)
		if err != nil {
			return ""
		}
	}
	return token
}

func TestIntegration_GithubProvider_GetNextRelease(t *testing.T) {
	// using a public repository to test integration
	p, err := NewGithubProvider("tauffredou", "test-semver", getToken(), intConfig)
	assert.NoError(t, err)
	r := p.GetNextRelease()
	assert.Equal(t, "v1.1.0", r.CurrentVersion)
	assert.Equal(t, "v1.2.0", r.MustNextVersion())
	assert.Len(t, r.Changelog, 2)
	assert.Equal(t, "feature 5", r.Changelog[0].Title)
	assert.Equal(t, "feature 4", r.Changelog[1].Title)
}

func TestIntegration_GithubProvider_GetRelease(t *testing.T) {
	// using a public repository to test integration
	p, err := NewGithubProvider("tauffredou", "test-semver", getToken(), intConfig)
	assert.NoError(t, err)
	r, err := p.GetRelease("v1.1.0")

	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, "v1.1.0", r.CurrentVersion)
	assert.Len(t, r.Changelog, 2)
	assert.Equal(t, "feature 3", r.Changelog[0].Title)
	assert.Equal(t, "feature 2", r.Changelog[1].Title)
}

func TestIntegration_GithubProvider_GetReleases(t *testing.T) {
	// using a public repository to test integration
	p, err := NewGithubProvider("tauffredou", "test-semver", getToken(), intConfig)
	assert.NoError(t, err)

	actual, err := p.GetReleases()
	require.NoError(t, err)
	assert.Equal(t, 2, len(actual))
	assert.Equal(t, "v1.1.0", actual[0].CurrentVersion)
	assert.Equal(t, "v1.0.1", actual[1].CurrentVersion)
}
