package provider

import (
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
	return &GitProvider{
		path:           path,
		versionPattern: versionPattern,
		versionRegexp:  GetVersionRegexp(versionPattern),
	}
}

func (p *GitProvider) GetReleases() ([]model.Release, error) {
	repo, err := git.PlainOpen(p.path)
	if err != nil {
		return nil, err
	}
	r := make([]model.Release, 0)
	t, err := repo.Tags()
	if err != nil {
		return nil, err
	}
	_ = t.ForEach(func(reference *plumbing.Reference) error {
		if p.tagFilter(reference) {
			tag := p.tagMapper(reference)
			r = append([]model.Release{tag}, r...)
		}
		return nil
	})
	return r, nil
}

func (p *GitProvider) GetNextRelease() *model.Release {
	return nil
}

func (*GitProvider) GetRelease(name string) (*model.Release, error) {
	return nil, nil
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

func (p *GitProvider) getPreviousRelease(release string) *model.Release {
	releases, _ := p.GetReleases()
	for i := 0; i < len(releases)-1; i++ {
		it := &releases[i]

		if release == "" {
			return it
		}

		if it.CurrentVersion == release {
			return &releases[i+1]
		}
	}

	return nil
}
