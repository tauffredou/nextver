package provider

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"github.com/tauffredou/nextver/model"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

type ProviderSuite struct {
	suite.Suite
	provider Provider
	gitPath  string
}

func (suite *ProviderSuite) SetupSuite() {
	suite.gitPath = fmt.Sprintf(filepath.Join(os.TempDir(), "nextver-%d"), time.Now().UnixNano())
	suite.provider = NewGitProvider(suite.gitPath, "vSEMVER")
	cloneRepo("https://github.com/tauffredou/test-semver.git", suite.gitPath)
}

func cloneRepo(url string, directory string) {
	var err error

	_, err = git.PlainClone(directory, false, &git.CloneOptions{
		URL:               url,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (suite *ProviderSuite) TearDownSuite() {
	cleanRepo(suite.gitPath)
}

func cleanRepo(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		log.Fatal(err)
	}
}

func TestProviderSuite(t *testing.T) {
	suite.Run(t, new(ProviderSuite))
}

func (suite *ProviderSuite) TestGitProvider_GetReleases() {
	p := suite.provider
	actual, err := p.GetReleases()
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), 2, len(actual))
	assert.Equal(suite.T(), "v1.1.0", actual[0].CurrentVersion)
	assert.Equal(suite.T(), "v1.0.1", actual[1].CurrentVersion)
}

func (suite *ProviderSuite) TestGitProvider_GetRelease_empty() {
	p := suite.provider
	actual, err := p.GetRelease("")
	require.NoError(suite.T(), err)
	require.Nil(suite.T(), actual)
}

func (suite *ProviderSuite) TestGitProvider_getPreviousRelease() {
	p := NewGitProvider(suite.gitPath, "vSEMVER")

	tests := []struct {
		name string
		want *model.Release
	}{
		{"v1.1.0", &model.Release{CurrentVersion: "v1.0.1"}},
		{"v1.0.1", nil},
	}
	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, p.getPreviousRelease(tt.name))
		})
	}
}

/* other test */

func TestGitProvider_tagFilter(t *testing.T) {

	tests := []struct {
		tag  string
		want bool
	}{
		{"refs/tags/someTag", false},
		{"refs/tags/v1.0.0", true},
	}

	p := NewGitProvider("", "vSEMVER")

	for _, tt := range tests {
		t.Run(tt.tag, func(t *testing.T) {
			reference := plumbing.NewReferenceFromStrings(tt.tag, "a")
			assert.Equal(t, p.tagFilter(reference), tt.want)
		})
	}
}
