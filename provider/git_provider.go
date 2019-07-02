package provider

import (
	"fmt"
	"github.com/tauffredou/nextver/model"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"regexp"
)

type GitProvider struct {
	path           string
	versionPattern string
	versionRegexp  *regexp.Regexp
}

func NewGitProvider(path string, versionPattern string) *GitProvider {
	provider := GitProvider{path: path, versionPattern: versionPattern}
	provider.versionRegexp = GetVersionRegexp(versionPattern)
	return &provider
}

func (p *GitProvider) GetReleases() []model.Release {
	repo, _ := git.PlainOpen(p.path)
	r := make([]model.Release, 0)
	tags, _ := repo.Tags()
	_ = tags.ForEach(func(reference *plumbing.Reference) error {
		fmt.Println(reference)
		if p.tagFilter(reference) {
			tag := p.tagMapper(reference)
			r = append([]model.Release{tag}, r...)
		}
		return nil
	})
	return r
}

func (*GitProvider) GetNextRelease() *model.Release {

	panic("implement me")
}

func (*GitProvider) GetRelease(name string) (*model.Release, error) {
	panic("implement me")
}

func (p *GitProvider) tagFilter(reference *plumbing.Reference) bool {
	s := reference.Name().Short()
	return p.versionRegexp.MatchString(s)
}

func (p *GitProvider) tagMapper(reference *plumbing.Reference) model.Release {
	return model.Release{
		CurrentVersion: reference.Name().Short(),
	}
}
