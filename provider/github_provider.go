package provider

import (
	"fmt"
	"github.com/shurcooL/githubv4"
	log "github.com/sirupsen/logrus"
	"github.com/tauffredou/nextver/model"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"
)

const (
	FirstCommit      = ""
	DefaultHubConfig = "$HOME/.config/hub"
)

type GithubProvider struct {
	client        *githubv4.Client
	VersionRegexp *regexp.Regexp
	config        *GithubProviderConfig
	Owner         string
	Repo          string
	pattern       string
	Branch        string
}

type GithubProviderConfig struct {
	Branch    string
	Pattern   string
	BeforeRef string
}

func NewGithubProvider(owner string, repo string, token string, config *GithubProviderConfig) (*GithubProvider, error) {
	if owner == "" || repo == "" || token == "" || config == nil {
		return nil, &ConfigurationError{}
	}

	log.WithField("token", obfuscateToken(token)).Debug("Init github provider")

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	return &GithubProvider{
		Owner:  owner,
		Repo:   repo,
		client: githubv4.NewClient(httpClient),
		config: config,
	}, nil
}

type ConfigurationError struct{}

func (e *ConfigurationError) Error() string {
	return "Invalid configuration"
}

func (p *GithubProvider) GetNextRelease() *model.Release {

	release := model.Release{
		Project:        fmt.Sprintf("%s/%s", p.Owner, p.Repo),
		VersionPattern: p.MustGetPattern(),
	}

	previousTag := p.getLastReleaseTag()
	if previousTag != nil {
		release.CurrentVersion = previousTag.getId()
		release.Ref = previousTag.getCommitId()
		release.Changelog = p.getHistory("HEAD", previousTag.getCommitId())
	} else {
		release.CurrentVersion = model.FirstVersion
		release.Ref = FirstCommit
		release.Changelog = p.getHistory("HEAD", FirstCommit)
	}

	return &release
}

//GetReleases returns the list of tags matching the release pattern
func (p *GithubProvider) GetReleases() []model.Release {
	log.Debug("Getting release")

	query := p.mustQueryReleases()

	r := make([]model.Release, 0)
	tags := query.Repository.Refs.TagNodes

	// reverse order
	for i := len(tags) - 1; i >= 0; i-- {
		v := tags[i]
		if p.tagFilter(v.Target.TagInfo) {
			tag := p.tagMapper(v.Target.TagInfo, nil)
			r = append(r, tag)
		}
	}

	return r
}

func (p *GithubProvider) GetRelease(name string) (*model.Release, error) {

	from, to, _ := p.getReleaseBoundary(name)

	r := model.Release{
		CurrentVersion: name,
		Changelog:      p.getHistory(from, to),
		VersionPattern: p.MustGetPattern(),
	}
	return &r, nil
}

//MustGetPattern tries to fetch the config file
func (p *GithubProvider) MustGetPattern() string {
	log.Debug("get pattern")
	if p.pattern != "" {
		return p.pattern
	}

	if p.config.Pattern != "" {
		p.pattern = p.config.Pattern
		return p.pattern
	}

	query := p.mustQueryConfigFile()

	if query.hasFile() {
		p.pattern = query.mustGetConfig().Pattern
		log.WithField("pattern", p.pattern).Debug("got pattern from github")
	} else {
		p.pattern = model.DefaultPattern
		log.WithField("pattern", p.pattern).Debug("got pattern from default")
	}
	return p.pattern
}

// readHubToken read token form hub config when available
// default location is ~/.config/hub
func ReadHubToken(f string) (string, error) {

	var v struct {
		Github []struct {
			Token string `yaml:"oauth_token"`
		} `yaml:"github.com,flow"`
	}
	f = os.ExpandEnv(f)
	if _, err := os.Stat(f); err == nil {
		bytes, _ := ioutil.ReadFile(f)
		err := yaml.Unmarshal(bytes, &v)
		if err != nil {
			return "", err
		}

		return v.Github[0].Token, nil
	} else {
		return "", err
	}

}

func (p *GithubProvider) getLastReleaseTag() *TagNode {
	tagsQuery := p.mustQueryReleases()
	tags := tagsQuery.GetTags()

	// reverse order
	for i := len(tags) - 1; i >= 0; i-- {
		tag := tags[i]

		if p.GetVersionRegexp().MatchString(tag.getId()) {
			return &tag
		}
	}

	return nil
}

func (p *GithubProvider) getHistory(fromRef string, toRef string) []model.ReleaseItem {

	variables := p.defaultVariables()

	variables["release"] = githubv4.String(fromRef)
	variables["itemsCount"] = githubv4.Int(50)
	ts, _ := time.Parse(time.RFC3339, "1900-01-01T00:00:00Z")
	variables["since"] = githubv4.GitTimestamp{Time: ts}

	query := p.mustGetHistory(variables)

	result := make([]model.ReleaseItem, 0)

	commits := query.getCommits()

	for _, c := range commits {
		if c.Oid == toRef {
			break
		}
		ri := model.NewReleaseItem(c.Author.Name, c.Author.Date, c.Message)
		result = append(result, ri)
	}

	return result

}

func (p *GithubProvider) defaultVariables() map[string]interface{} {
	return map[string]interface{}{
		"owner": githubv4.String(p.Owner),
		"name":  githubv4.String(p.Repo),
	}
}

func obfuscateToken(token string) string {
	strlen := len(token)
	var sb strings.Builder
	for pos, char := range token {
		if pos < 2 || pos > strlen-3 {
			sb.WriteRune(char)
		} else {
			sb.WriteRune('*')
		}
	}
	return sb.String()
}

func (p *GithubProvider) mustGetHistory(variables map[string]interface{}) *historyQuery {
	var query historyQuery

	err := p.client.Query(context.Background(), &query, variables)
	if err != nil {
		log.Fatal(err)
	}
	return &query
}

func (p *GithubProvider) mustGetDefaultBranch() string {
	var query defaultBranchQuery
	err := p.client.Query(context.Background(), &query, p.defaultVariables())
	if err != nil {
		log.Fatal(err)
	}
	return query.Repository.DefaultBranchRef.Name
}

func (p *GithubProvider) tagFilter(v TagInfo) bool {
	return p.GetVersionRegexp().MatchString(v.getId())
}

func (p *GithubProvider) tagMapper(tag TagInfo, changeLog []model.ReleaseItem) model.Release {
	return model.Release{
		Project:        fmt.Sprintf("%s/%s", p.Owner, p.Repo),
		CurrentVersion: tag.getId(),
		Ref:            tag.getCommitId(),
		Changelog:      changeLog,
		VersionPattern: p.MustGetPattern(),
	}
}

func (p *GithubProvider) mustQueryReleases() *tagsQuery {

	var query tagsQuery

	variables := map[string]interface{}{
		"owner": githubv4.String(p.Owner),
		"name":  githubv4.String(p.Repo),
	}
	err := p.client.Query(context.Background(), &query, variables)
	if err != nil {
		log.WithError(err).Fatal("cannot get last tags")
	}
	return &query
}

func (p *GithubProvider) mustQueryConfigFile() *configFileQuery {
	var query configFileQuery
	variables := p.defaultVariables()
	variables["file"] = githubv4.String(p.mustGetBranch() + ":" + model.DefaultConfigFile)

	err := p.client.Query(context.Background(), &query, variables)
	if err != nil {
		log.Fatal(err)
	}
	return &query
}

// mustGetBranch get the target branch by order from:
// 1. local param
// 2. repository default branch
func (p *GithubProvider) mustGetBranch() string {
	if p.Branch != "" {
		return p.Branch
	}

	if p.config.Branch != "" {
		p.Branch = p.config.Branch
		return p.Branch
	}

	variables := p.defaultVariables()
	var query defaultBranchQuery

	err := p.client.Query(context.Background(), &query, variables)
	if err != nil {
		log.Fatal(err)
	}
	p.Branch = query.Repository.DefaultBranchRef.Name
	return p.Branch
}

func (p *GithubProvider) GetVersionRegexp() *regexp.Regexp {
	if p.VersionRegexp != nil {
		return p.VersionRegexp
	}

	replacer := strings.NewReplacer(
		"SEMVER", model.SemverRegex,
		"DATE", model.DateRegexp,
	)

	p.VersionRegexp = regexp.MustCompile("^" + replacer.Replace(p.MustGetPattern()) + "$")
	return p.VersionRegexp
}

func (p *GithubProvider) getReleaseBoundary(release string) (string, string, error) {
	var first, last string

	tags := p.mustQueryReleases()

	TagNodes := tags.Repository.Refs.TagNodes
	for i, t := range TagNodes {
		if t.getId() == release {
			first = t.getCommitId()
			if i != 0 {
				last = TagNodes[i-1].getCommitId()
			}
		}
	}

	return first, last, nil
}
