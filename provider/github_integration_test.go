package provider

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	intConfig = &GithubProviderConfig{
		Branch:  "master",
		Pattern: "vSEMVER",
	}
)

func TestGithubProvider_GetLatestRelease(t *testing.T) {
	// using a public repository to test integration
	p, err := NewGithubProvider("tauffredou", "test-semver", "a", intConfig)
	assert.NoError(t, err)

	var q latestReleasesQuery

	_ = p.queryLatestRelease(&q)
	assert.NotNil(t, q.Repository)
	//assert.NotEmpty(t, q.Repository.Refs.Edges)
	//assert.Equal(t, "v1.1.0", p.getFirstTag(q))
}
