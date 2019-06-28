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

type GithubProvider struct {
	client        *githubv4.Client
	VersionRegexp *regexp.Regexp
	config        *GithubProviderConfig
	Owner         string
	Repo          string
	Pattern       string
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

type ConfigurationError struct{}

func (e *ConfigurationError) Error() string {
	return "Invalid configuration"
}

func (p *GithubProvider) GetLatestRelease() model.Release {

	tagsQuery := p.mustQueryReleases()
	tags := tagsQuery.GetTags()

	// reverse order
	pattern := p.MustGetPattern()
	for i := len(tags) - 1; i >= 0; i-- {
		tag := tags[i]

		if p.GetVersionRegexp().MatchString(tag.getId()) {
			ref := tag.getCommitId()
			return model.Release{
				Project:        fmt.Sprintf("%s/%s", p.Owner, p.Repo),
				CurrentVersion: strings.Trim(tag.getMessage(), "\n"),
				Ref:            ref,
				Changelog:      p.getHistory(ref),
				VersionPattern: pattern,
			}
		}
	}

	return model.Release{
		Project:        fmt.Sprintf("%s/%s", p.Owner, p.Repo),
		CurrentVersion: model.FirstVersion,
		Ref:            "",
		Changelog:      p.getHistory(""),
		VersionPattern: pattern,
	}

}

func (p *GithubProvider) getHistory(fromRef string) []model.ReleaseItem {

	variables := p.defaultVariables()
	variables["branch"] = githubv4.String(p.mustGetBranch())
	variables["itemsCount"] = githubv4.Int(50)
	ts, _ := time.Parse(time.RFC3339, "1900-01-01T00:00:00Z")
	variables["since"] = githubv4.GitTimestamp{Time: ts}

	query := p.mustGetHistory(variables)
	nodes := getCommits(query)

	result := make([]model.ReleaseItem, 0)

	for _, n := range nodes {
		if n.Oid == fromRef {
			break
		}
		ri := model.NewReleaseItem(n.Author.Name, n.Author.Date, n.Message)
		result = append(result, ri)
	}

	return result

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

//MustGetPattern tries to fetch the config filee
func (p *GithubProvider) MustGetPattern() string {
	log.Debug("get pattern")
	if p.Pattern != "" {
		return p.Pattern
	}

	if p.config.Pattern != "" {
		p.Pattern = p.config.Pattern
		return p.Pattern
	}

	query := p.mustQueryConfigFile()

	if query.hasFile() {
		p.Pattern = query.mustGetConfig().Pattern
		log.WithField("pattern", p.Pattern).Debug("got pattern from github")
	} else {
		p.Pattern = model.DefaultPattern
		log.WithField("pattern", p.Pattern).Debug("got pattern from default")
	}
	return p.Pattern
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

func (p *GithubProvider) GetRelease(name string) (*model.Release, error) {

	_, _, _ = p.getReleaseBoundary(name)

	// get history from boundaries
	//p.mustGetHistory()

	//panic("implement me")
	return nil, nil
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

const DEFAULT_HUB_CONFIG = "~/.config/hub"

// readHubToken read token form hub config when available
// default location is ~/.config/hub
func ReadHubToken(f string) (string, error) {

	var v struct {
		Github []struct {
			Token string `yaml:"oauth_token"`
		} `yaml:"github.com,flow"`
	}

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
