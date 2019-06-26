package provider

import (
	"testing"

	"github.com/shurcooL/githubv4"
	"github.com/stretchr/testify/assert"
)

var (
	testConfig = &GithubProviderConfig{
		Branch:  "master",
		Pattern: "vSEMVER",
	}
)

func TestNewGithubProvider_emptyConfig(t *testing.T) {
	_, err := NewGithubProvider("owner", "repo", "token", nil)
	if assert.Error(t, err) {
		assert.Equal(t, &ConfigurationError{}, err)
	}

}

func TestNewGithubProvider_emptyOwner(t *testing.T) {
	_, err := NewGithubProvider("", "repo", "token", testConfig)
	if assert.Error(t, err) {
		assert.Equal(t, &ConfigurationError{}, err)
	}
}

func TestNewGithubProvider_emptyRepo(t *testing.T) {
	_, err := NewGithubProvider("owner", "", "token", testConfig)
	if assert.Error(t, err) {
		assert.Equal(t, &ConfigurationError{}, err)
	}
}

func TestNewGithubProvider_emptyToken(t *testing.T) {
	_, err := NewGithubProvider("owner", "repo", "", testConfig)
	if assert.Error(t, err) {
		assert.Equal(t, &ConfigurationError{}, err)
	}
}

func TestNewGithubProvider_constructor(t *testing.T) {
	p, err := NewGithubProvider("owner", "repo", "token", testConfig)
	assert.NoError(t, err)
	assert.IsType(t, &githubv4.Client{}, p.client)
	assert.Equal(t, "owner", p.Owner)
	assert.Equal(t, "repo", p.Repo)
}

func TestNewGithubProvider_obfuscateToken(t *testing.T) {
	tests := []struct {
		token    string
		expected string
	}{
		{"a", "a"},
		{"ab", "ab"},
		{"abcd", "abcd"},
		{"abcde", "ab*de"},
		{"abcdefabcdefabcdefabcdefabcdef", "ab**************************ef"},
	}
	for _, test := range tests {
		t.Run(test.token, func(t *testing.T) {
			assert.Equal(t, test.expected, obfuscateToken(test.token))
		})
	}
}

func TestConfigurationError_Error(t *testing.T) {
	err := &ConfigurationError{}
	assert.Error(t, err)
	assert.Equal(t, "Invalid configuration", err.Error())
}

func TestGithubProvider_GetVersionRegexp_createMatcherOnce(t *testing.T) {
	p1, _ := NewGithubProvider("owner", "repo", "token", &GithubProviderConfig{Pattern: "first"})
	assert.Equal(t, true, p1.GetVersionRegexp().MatchString("first"))
	p1.config.Pattern = "second"

	p2, _ := NewGithubProvider("owner", "repo", "token", &GithubProviderConfig{Pattern: "second"})
	p2.VersionRegexp = p1.GetVersionRegexp()
	assert.Equal(t, true, p2.GetVersionRegexp().MatchString("first"))
	assert.Equal(t, false, p1.GetVersionRegexp().MatchString("second"))
}

func TestGithubProvider_GetVersionRegexp_semver(t *testing.T) {
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
			p, _ := NewGithubProvider("owner", "repo", "token", &GithubProviderConfig{
				Pattern: test.pattern,
			})
			assert.Equal(t, test.expected, p.GetVersionRegexp().MatchString(test.match))
		})
	}
}

func TestGithubProvider_GetVersionRegexp_date(t *testing.T) {
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
			p, _ := NewGithubProvider("owner", "repo", "token", &GithubProviderConfig{
				Pattern: test.pattern,
			})
			assert.Equal(t, test.expected, p.GetVersionRegexp().MatchString(test.match))
		})
	}
}
