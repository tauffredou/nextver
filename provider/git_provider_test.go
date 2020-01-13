package provider

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"github.com/tauffredou/nextver/model"
	"gopkg.in/src-d/go-billy.v4/osfs"
	"gopkg.in/src-d/go-git.v4/plumbing/cache"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
	"io/ioutil"
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

func (suite *ProviderSuite) TestGitProvider_getPreviousRelease() {
	p := NewGitProvider(suite.gitPath, "vSEMVER")

	tests := []struct {
		name string
		want *model.Release
	}{
		{"v1.1.0", &model.Release{CurrentVersion: "v1.0.1", VersionPattern: "vSEMVER"}},
		{"v1.0.1", nil},
		{"bad", nil},
		{"", &model.Release{CurrentVersion: "v1.1.0", VersionPattern: "vSEMVER"}},
	}
	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, p.getPreviousRelease(tt.name))
		})
	}
}

func TestGitProvider_getPreviousRelease_noPreviousRelease(t *testing.T) {

	outputDir, _ := ioutil.TempDir("", "nextver-test-")
	defer os.RemoveAll(outputDir)

	fs := osfs.New(outputDir)
	_, _ = git.Init(filesystem.NewStorage(fs, cache.NewObjectLRUDefault()), fs)

	p := NewGitProvider(outputDir, "vSEMVER")
	assert.Nil(t, p.getPreviousRelease("any"))
}

func (suite *ProviderSuite) TestGitProvider_GetReleases_badPath() {
	p := NewGitProvider("badPath", "vSEMVER")
	_, err := p.GetReleases()
	assert.Error(suite.T(), err)
}

func (suite *ProviderSuite) TestGitProvider_GetNextRelease() {
	actual := suite.provider.GetNextRelease()

	require.NotNil(suite.T(), actual)
	assert.Equal(suite.T(), "v1.2.0", actual.MustNextVersion())
}

func (suite *ProviderSuite) TestGitProvider_GetRelease_empty() {
	p := suite.provider
	actual, err := p.GetRelease("")
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), actual)
}

func (suite *ProviderSuite) TestGitProvider_GetRelease_last() {
	p := suite.provider
	actual, err := p.GetRelease("v1.0.1")
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), actual)

	assert.Equal(suite.T(), "v1.0.1", actual.CurrentVersion)
	require.Len(suite.T(), actual.Changelog, 5)
	assert.Equal(suite.T(), "change f1", actual.Changelog[0].Title)
	assert.Equal(suite.T(), "Initial commit", actual.Changelog[4].Title)
}

func (suite *ProviderSuite) TestProvider_GetRelease_withBoundaries() {
	p := suite.provider
	r, err := p.GetRelease("v1.1.0")

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), r)
	assert.Equal(suite.T(), "v1.1.0", r.CurrentVersion)
	require.Len(suite.T(), r.Changelog, 2)
	assert.Equal(suite.T(), "feature 3", r.Changelog[0].Title)
	assert.Equal(suite.T(), "feature 2", r.Changelog[1].Title)
}

func (suite *ProviderSuite) TestProvider_GetReleases() {
	p := suite.provider

	actual, err := p.GetReleases()
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), 2, len(actual))
	assert.Equal(suite.T(), "v1.1.0", actual[0].CurrentVersion)
	assert.Equal(suite.T(), "v1.0.1", actual[1].CurrentVersion)
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

func TestGitProvider_mustGetPattern_default(t *testing.T) {
	p := GitProvider{}
	actual := p.VersionPattern()
	assert.Equal(t, model.DefaultPattern, actual)
}

func TestGitProvider_mustGetPattern_fromFile(t *testing.T) {
	p := GitProvider{path: "../fixtures/local"}
	actual := p.VersionPattern()
	assert.Equal(t, "testSEMVER", actual)
}

func TestGitProvider_mustGetPattern_override(t *testing.T) {
	p := GitProvider{versionPattern: "overridePattern"}
	actual := p.VersionPattern()
	assert.Equal(t, "overridePattern", actual)
}

func TestGitProvider_GetVersionRegexp_semver(t *testing.T) {
	tests := []struct {
		pattern  string
		match    string
		expected bool
	}{
		{"SEMVER", "bad-v1", false},
		{"SEMVER", "v1", true},
		{"SEMVER", "1", true},
		{"SEMVER", "1.0", true},
		{"SEMVER", "1.0.1", true},
		{"SEMVER", "1.0.1.0", false},
		{"SEMVER", "1.0.1-bad", false},
		{"SEMVER", "v1", true},
		{"SEMVER", "prefix-1.0", false},
		{"prefix-SEMVER", "prefix-1.0", true},
		{"prefix-SEMVER-suffix", "prefix-1.0-suffix", true},
	}
	for _, test := range tests {
		t.Run(test.pattern+"_"+test.match, func(t *testing.T) {
			p := GitProvider{versionPattern: test.pattern}
			assert.Equal(t, test.expected, p.VersionRegexp().MatchString(test.match))
		})
	}
}

func TestGitProvider_GetVersionRegexp_date(t *testing.T) {
	tests := []struct {
		pattern  string
		match    string
		expected bool
	}{
		{"DATE", "2019-06-12-095000", true},
		{"prefix-DATE", "prefix-2019-06-12-095000", true},
		{"DATE-suffix", "2019-06-12-095000-suffix", true},
	}
	for _, test := range tests {
		t.Run(test.pattern+"_"+test.match, func(t *testing.T) {
			p := GitProvider{versionPattern: test.pattern}
			assert.Equal(t, test.expected, p.VersionRegexp().MatchString(test.match))
		})
	}
}
