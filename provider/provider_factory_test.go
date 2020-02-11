package provider

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateProvider_github(t *testing.T) {
	f := ProviderFactory{
		TokenReader: func() string { return "aToken" },
		Pattern:     "vSEMVER",
	}

	p, err := f.CreateProvider("github.com/test/test-rep")

	assert.NoError(t, err)
	assert.NotNil(t, p)
	assert.IsType(t, &GithubProvider{}, p)
	gp := p.(*GithubProvider)
	assert.Equal(t, "test", gp.Owner)
	assert.Equal(t, "test-rep", gp.Repo)
}

func TestCreateProvider_git(t *testing.T) {
	f := ProviderFactory{
		Pattern: "vSEMVER",
	}

	p, err := f.CreateProvider(".")

	assert.NoError(t, err)
	assert.NotNil(t, p)
	assert.IsType(t, &GitProvider{}, p)
	gp := p.(*GitProvider)
	assert.Equal(t, ".", gp.path)
}

func TestRepoParam(t *testing.T) {
	tests := []struct {
		name    string
		repo    string
		want    interface{}
		wantErr error
	}{
		{name: "empty", repo: "", wantErr: &InvalidRepositoryError{repo: ""}},
		{name: "non existing dir", repo: "non/existing/dir", wantErr: &InvalidRepositoryError{repo: "non/existing/dir"}},

		{name: "current dir", repo: ".", want: GitRepository{path: "."}},
		{name: "existing dir", repo: "../doc/", want: GitRepository{path: "../doc/"}},

		{name: "github short", repo: "github.com/test/test-rep", want: GithubRepository{Owner: "test", Repo: "test-rep"}},
		{name: "github https", repo: "https://github.com/test/test-rep", want: GithubRepository{Owner: "test", Repo: "test-rep"}},
		{name: "github https .git", repo: "https://github.com/test/test-rep.git", want: GithubRepository{Owner: "test", Repo: "test-rep"}},
		{name: "github git", repo: "git@github.com:test/test-rep.git", want: GithubRepository{Owner: "test", Repo: "test-rep"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parseRepo, err := ParseRepo(test.repo)
			if test.wantErr != nil {
				assert.Equal(t, test.wantErr, err)
			} else {
				assert.Equal(t, test.want, parseRepo)
			}
		})
	}
}
