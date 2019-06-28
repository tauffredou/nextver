// +build integration

package provider

import (
	"github.com/stretchr/testify/assert"
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
		token, err = ReadHubToken(DEFAULT_HUB_CONFIG)
		if err != nil {
			return ""
		}
	}
	return token
}

func TestIntergration_GithubProvider_GetLatestRelease(t *testing.T) {
	// using a public repository to test integration
	p, err := NewGithubProvider("tauffredou", "test-semver", getToken(), intConfig)
	assert.NoError(t, err)
	r := p.GetLatestRelease()
	assert.Equal(t, "v1.1.0", r.CurrentVersion)

}

func TestIntegration_GithubProvider_GetRelease(t *testing.T) {
	// using a public repository to test integration
	p, err := NewGithubProvider("tauffredou", "test-semver", getToken(), intConfig)
	assert.NoError(t, err)
	_, err = p.GetRelease("v1.1.0")
	assert.NoError(t, err)
	//assert.Equal(t, "", r.CurrentVersion)

}
