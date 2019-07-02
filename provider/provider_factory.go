package provider

import (
	"fmt"
	"os"
	"regexp"
)

type ProviderFactory struct {
	Token string
}

func (f *ProviderFactory) CreateProvider(repo string) (Provider, error) {
	r, err := ParseRepo(repo)
	if err != nil {
		return nil, err
	}
	gr := r.(GithubRepository)
	provider, err := NewGithubProvider(gr.Repo, gr.Owner, f.Token, nil)
	if err != nil {
		return nil, err
	}

	return provider, nil
}

func ParseRepo(repo string) (interface{}, error) {
	if repo == "" {
		return nil, &InvalidRepositoryError{repo: repo}
	}

	re := regexp.MustCompile(`^(https://|git@)?github.com[:/]([a-zA-Z0-9-]+)/([a-zA-Z0-9-]+)(\.git)?$`)
	if re.MatchString(repo) {
		v := re.FindStringSubmatch(repo)
		return GithubRepository{Owner: v[2], Repo: v[3]}, nil
	}

	if _, err := os.Stat(repo); os.IsNotExist(err) {
		return nil, &InvalidRepositoryError{repo: repo}
	}

	return GitRepository{path: repo}, nil
}

type GithubRepository struct {
	Owner string
	Repo  string
}

type GitRepository struct {
	path string
}

type InvalidRepositoryError struct{ repo string }

func (e InvalidRepositoryError) Error() string { return fmt.Sprintf("Invalid repository %s", e.repo) }
