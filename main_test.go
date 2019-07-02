package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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
		{name: "current dir", repo: "doc/", want: GitRepository{path: "doc/"}},

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
