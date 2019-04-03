package provider

import (
	"fmt"
	"github.com/shurcooL/githubv4"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"regexp"
	"strings"
)

type GithubProvider struct {
	Owner   string
	Repo    string
	client  *githubv4.Client
	Pattern string
	Branch  string
}

func NewGithubProvider(owner string, repo string, token string, pattern string, branch string) *GithubProvider {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	return &GithubProvider{
		Owner:   owner,
		Repo:    repo,
		client:  githubv4.NewClient(httpClient),
		Pattern: pattern,
		Branch:  branch,
	}
}

func (p *GithubProvider) GetLatestRelease() Release {

	var query struct {
		Repository struct {
			Refs struct {
				Edges []tagEdge
			} `graphql:"refs(refPrefix: \"refs/tags/\", last: 50, orderBy: {field: TAG_COMMIT_DATE, direction: ASC})"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	variables := map[string]interface{}{
		"owner": githubv4.String(p.Owner),
		"name":  githubv4.String(p.Repo),
	}

	err := p.client.Query(context.Background(), &query, variables)
	if err != nil {
		log.WithError(err).Fatal("cannot get last tags")
	}

	tags := query.Repository.Refs.Edges
	replacer := strings.NewReplacer(
		"SEMVER", SemverRegex,
		"DATE", DateRegexp,
	)

	re := regexp.MustCompile(replacer.Replace(p.Pattern))
	// reverse order
	for i := len(tags) - 1; i >= 0; i-- {
		if re.MatchString(tags[i].Node.Target.Tag.Message) {
			ref := tags[i].Node.Target.Tag.Target.Commit.Oid
			return Release{
				Project:        fmt.Sprintf("%s/%s", p.Owner, p.Repo),
				CurrentVersion: strings.Trim(tags[i].Node.Target.Tag.Message, "\n"),
				Ref:            ref,
				Changelog:      p.getHistory(ref),
				VersionPattern: p.Pattern,
			}
		}
	}

	return Release{
		Project:        fmt.Sprintf("%s/%s", p.Owner, p.Repo),
		CurrentVersion: FirstVersion,
		Ref:            "",
		Changelog:      p.getHistory(""),
		VersionPattern: p.Pattern,
	}

}

func (p *GithubProvider) getHistory(fromRef string) []ReleaseItem {

	variables := map[string]interface{}{
		"owner": githubv4.String(p.Owner),
		"name":  githubv4.String(p.Repo),
	}

	if p.Branch != "" {
		variables["branch"] = githubv4.String(p.Branch)
	} else {
		// Get default branch
		var query struct {
			Repository struct {
				DefaultBranchRef struct {
					Name string
				}
			} `graphql:"repository(owner: $owner, name: $name)"`
		}

		err := p.client.Query(context.Background(), &query, variables)
		if err != nil {
			log.Fatal(err)
		}
		variables["branch"] = githubv4.String(query.Repository.DefaultBranchRef.Name)
	}

	variables["itemsCount"] = githubv4.Int(50)
	var query struct {
		Repository struct {
			Ref struct {
				Target struct {
					Commit struct {
						History struct {
							PageInfo PageInfo
							Nodes    []CommitNode
						} `graphql:"history(first: $itemsCount)"`
					} `graphql:"... on Commit"`
				}
			} `graphql:"ref(qualifiedName: $branch)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	err := p.client.Query(context.Background(), &query, variables)
	if err != nil {
		log.Fatal(err)
	}
	nodes := query.Repository.Ref.Target.Commit.History.Nodes

	result := make([]ReleaseItem, 0)

	for i := range nodes {
		n := nodes[i]
		if n.Oid == fromRef {
			break
		}
		ri := NewReleaseItem(n.Author.Name, n.Author.Date, n.Message)
		result = append(result, ri)
	}

	return result

}
