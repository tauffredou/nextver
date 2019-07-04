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

	err = t.ForEach(func(reference *plumbing.Reference) error {
		if p.tagFilter(reference) {
			tag := p.tagMapper(reference)
			r = append([]model.Release{tag}, r...) //preprend to reverse order
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (p *GitProvider) GetNextRelease() *model.Release {
	tag := p.getPreviousRelease("").CurrentVersion
	release, _ := p.GetRelease(tag)
	return release
}

func (p *GitProvider) GetRelease(name string) (*model.Release, error) {
	var (
		err error
		ref *plumbing.Reference
	)

	release := model.Release{
		CurrentVersion: name,
		Changelog:      make([]model.ReleaseItem, 0),
		VersionPattern: p.versionPattern,
	}

	repo, err := git.PlainOpen(p.path)
	if err != nil {
		return nil, err
	}

	if name != "" {
		ref, _ = repo.Tag(name)
	}

	previousRelease := p.getPreviousRelease(name)
	prev, _ := repo.Tag(previousRelease.CurrentVersion)
	prevObject, err := repo.TagObject(prev.Hash())
	if err != nil {
		return nil, err
	}

	var options git.LogOptions
	if ref != nil {
		refObject, err := repo.TagObject(ref.Hash())
		if err != nil {
			return nil, err
		}

		options = git.LogOptions{
			From: refObject.Target,
		}
	} else {
		options = git.LogOptions{
			All: true,
		}
	}

	it, err := repo.Log(&options)
	if err != nil {
		return nil, err
	}

	for {
		commit, err := it.Next()
		if err != nil {
			break
		}

		if commit.Hash == prevObject.Target {
			break
		}
		item := model.NewReleaseItem(commit.Author.Name, commit.Author.When, commit.Message)
		release.Changelog = append(release.Changelog, item)
	}
	it.Close()
	return &release, nil
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
