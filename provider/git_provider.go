package provider

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tauffredou/nextver/model"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
)

type GitProvider struct {
	path           string
	versionPattern string
	versionRegexp  *regexp.Regexp
}

func (p *GitProvider) String() string {
	return fmt.Sprintf("path:%s, versionPattern:%s, versionRegexp:%s", p.path, p.VersionPattern(), p.VersionRegexp().String())
}

func NewGitProvider(path string, versionPattern string) *GitProvider {
	provider := GitProvider{
		path:           path,
		versionPattern: versionPattern,
	}
	return &provider
}

func (p *GitProvider) VersionRegexp() *regexp.Regexp {
	if p.versionRegexp == nil {
		p.versionRegexp = GetVersionRegexp(p.VersionPattern())
	}
	return p.versionRegexp
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
			r = append(r, tag)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.SliceStable(r, func(i, j int) bool {
		return r[i].CurrentVersion > r[j].CurrentVersion
	})
	return r, nil
}

func (p *GitProvider) GetNextRelease() *model.Release {
	previousRelease := p.getPreviousRelease("")
	var tag string
	if previousRelease == nil {
		tag = ""
	} else {
		tag = previousRelease.CurrentVersion
	}
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
		VersionPattern: p.VersionPattern(),
	}

	repo, err := git.PlainOpen(p.path)
	if err != nil {
		return nil, err
	}

	if name != "" {
		ref, _ = repo.Tag(name)
	}

	previousRelease := p.getPreviousRelease(name)
	var prevObject *object.Tag
	if previousRelease != nil {
		prev, _ := repo.Tag(previousRelease.CurrentVersion)
		prevObject, err = repo.TagObject(prev.Hash())
		if err != nil {
			return nil, err
		}
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
			Order: git.LogOrderCommitterTime,
		}
	}

	it, err := repo.Log(&options)
	if err != nil {
		return nil, err
	}

	release.Changelog = mapChangelog(it, prevObject)
	it.Close()
	return &release, nil
}

func mapChangelog(it object.CommitIter, prevObject *object.Tag) []model.ReleaseItem {
	changelog := make([]model.ReleaseItem, 0)
	for {
		commit, err := it.Next()
		if err != nil {
			break
		}

		if prevObject != nil && commit.Hash == prevObject.Target {
			break
		}
		item := model.NewReleaseItem(commit.Author.Name, commit.Author.When, commit.Message)
		changelog = append(changelog, item)
	}
	return changelog
}

func (p *GitProvider) tagFilter(reference *plumbing.Reference) bool {
	s := reference.Name().Short()
	return p.VersionRegexp().MatchString(s)
}

func (p *GitProvider) tagMapper(reference *plumbing.Reference) model.Release {
	return model.Release{
		CurrentVersion: reference.Name().Short(),
		VersionPattern: p.VersionPattern(),
	}
}

func (p *GitProvider) getPreviousRelease(release string) *model.Release {
	releases, err := p.GetReleases()
	if err != nil {
		log.Fatal(err)
	}

	if len(releases) == 0 {
		return nil
	}

	for i := 0; i < len(releases)-1; i++ {
		it := releases[i]
		it.VersionPattern = p.VersionPattern()

		if release == "" {
			return &it
		}

		if it.CurrentVersion == release {
			r := releases[i+1]
			r.VersionPattern = p.VersionPattern()
			return &r
		}
	}

	return nil
}

//VersionPattern tries to fetch the config file
func (p *GitProvider) VersionPattern() string {
	if p.versionPattern != "" {
		return p.versionPattern
	}

	c, err := p.ReadConfigFile()
	if err == nil && c.Pattern != "" {
		return c.Pattern
	}

	return model.DefaultPattern
}

func (p *GitProvider) ReadConfigFile() (*model.Config, error) {
	//log.Debug("provider.GitProvider::ReadConfigFile")
	f := filepath.Join(p.path, model.DefaultConfigFile)

	if _, err := os.Stat(f); os.IsNotExist(err) {
		log.WithField("filename", f).Debug("Cannot read configuration")
		return nil, err
	}
	bytes, err := ioutil.ReadFile(f)
	var c model.Config
	err = yaml.Unmarshal(bytes, &c)
	if err != nil {
		log.Fatal(err)
	}
	return &c, nil
}
