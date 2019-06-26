package provider

import (
	"fmt"
	"github.com/shurcooL/githubv4"
	log "github.com/sirupsen/logrus"
	"github.com/tauffredou/nextver/model"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"regexp"
	"strings"
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

	var query latestReleasesQuery

	err := p.queryLatestRelease(&query)
	if err != nil {
		log.WithError(err).Fatal("cannot get last tags")
	}

	tags := query.Repository.Refs.Nodes

	// reverse order
	pattern := p.MustGetPattern()
	for i := len(tags) - 1; i >= 0; i-- {
		tag, _ := query.getTag(i)

		if p.GetVersionRegexp().MatchString(tag.getId()) {
			ref := tag.getCommitId()
			return model.Release{
				Project:        fmt.Sprintf("%s/%s", p.Owner, p.Repo),
				CurrentVersion: strings.Trim(tag.Message, "\n"),
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

func (p *GithubProvider) queryLatestRelease(query *latestReleasesQuery) error {
	variables := map[string]interface{}{
		"owner": githubv4.String(p.Owner),
		"name":  githubv4.String(p.Repo),
	}

	return p.client.Query(context.Background(), query, variables)
}

func (p *GithubProvider) getHistory(fromRef string) []model.ReleaseItem {

	variables := p.defaultVariables()
	variables["branch"] = githubv4.String(p.mustGetBranch())
	variables["itemsCount"] = githubv4.Int(50)

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
		if p.tagFilter(v.TagInfo) {
			tag := p.tagMapper(v.TagInfo, nil)
			r = append(r, tag)
		}
	}

	return r
}

func (p *GithubProvider) tagFilter(v TagNode) bool {
	return p.GetVersionRegexp().MatchString(v.getId())
}

func (p *GithubProvider) tagMapper(tag TagNode, changeLog []model.ReleaseItem) model.Release {
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

func (p *GithubProvider) defaultVariables() map[string]interface{} {
	return map[string]interface{}{
		"owner": githubv4.String(p.Owner),
		"name":  githubv4.String(p.Repo),
	}
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
