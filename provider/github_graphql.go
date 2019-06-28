// you can test queries using github explorer
// https://developer.github.com/v4/explorer/
package provider

import (
	"github.com/tauffredou/nextver/model"
	"gopkg.in/yaml.v2"
	"log"
	"strings"
	"time"
)

type PageInfo struct {
	HasNextPage bool
}

type latestReleasesQuery struct {
	Repository struct {
		Refs struct {
			Nodes []struct {
				Name   string
				Target struct {
					Oid           string
					TagInfo       TagInfo `graphql:"... on Tag"`
					CommittedDate string
				}
			} `graphql:"nodes"`
		} `graphql:"refs(refPrefix: \"refs/tags/\", last: 50, orderBy: {field: TAG_COMMIT_DATE, direction: ASC})"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

func (q *latestReleasesQuery) getTag(index int) (*TagInfo, error) {
	return &q.Repository.Refs.Nodes[index].Target.TagInfo, nil
}

/*
tagsQuery retrieve the last 50 tags
graphql query:

query ($owner: String!, $repo: String!) {
  repository(owner: $owner, name: $repo) {
    refs(refPrefix: "refs/tags/", last: 50, orderBy: {field: TAG_COMMIT_DATE, direction: ASC}) {
      nodes {
        name
        target {
          oid
          ... on Tag {
            oid
            message
            target {
              ... on Commit {
                oid
                committedDate
              }
            }
          }
        }
      }
    }
  }
}

*/
type tagsQuery struct {
	Repository struct {
		Refs struct {
			TagNodes []TagNode `graphql:"nodes"`
		} `graphql:"refs(refPrefix: \"refs/tags/\", last: 50, orderBy: {field: TAG_COMMIT_DATE, direction: ASC})"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

func (query *tagsQuery) GetTags() []TagNode { return query.Repository.Refs.TagNodes }

type TagNode struct {
	Name   string `graphql:"name"`
	Target struct {
		Oid     string  `graphql:"oid"`
		TagInfo TagInfo `graphql:"... on Tag"`
	} `graphql:"target"`
}

func (node *TagNode) getId() string       { return node.Target.TagInfo.getId() }
func (node *TagNode) getCommitId() string { return node.Target.TagInfo.getCommitId() }
func (node *TagNode) getMessage() string  { return node.Target.TagInfo.Message }

type TagInfo struct {
	Oid     string `graphql:"oid"`
	Message string `graphql:"message"`
	Target  struct {
		Commit struct {
			Oid           string `graphql:"oid"`
			CommittedDate string `graphql:"committedDate"`
		} `graphql:"... on Commit"`
	} `graphql:"target"`
}

func (t *TagInfo) getCommitId() string { return t.Target.Commit.Oid }
func (t *TagInfo) getId() string       { return strings.Trim(t.Message, "\n") }

/*
graphql query:

query($owner: String!, $repo: String!){
  repository(owner: $owner, name: $repo) {
    defaultBranchRef {
      name
    }
  }
}

*/
type defaultBranchQuery struct {
	Repository struct {
		DefaultBranchRef struct {
			Name string
		}
	} `graphql:"repository(owner: $owner, name: $name)"`
}

/*
graphql:

query (
  $release: String!,
  $owner: String!,
  $repo: String!,
  $since: String!,
) {
  repository(owner: $owner, name: $repo) {
    object(expression: $release) {
      ... on Commit {
        history(first: 20,since: $since) {
          pageInfo {
            endCursor
            startCursor
            hasNextPage
          }
          nodes {
            messageHeadline
            oid
            messageBody
            committedDate
          }
        }
      }
    }
  }
}

variables:
{
  "owner": "tauffredou",
  "repo": "test-semver",
  "release":"v1.1.0",
  "since": "2019-06-25T12:35:41Z"
}
*/
type historyQuery struct {
	Repository struct {
		Ref struct {
			Target struct {
				Commit struct {
					History struct {
						PageInfo PageInfo
						Nodes    []CommitNode
					} `graphql:"history(first: $itemsCount,since: $since)"`
				} `graphql:"... on Commit"`
			}
		} `graphql:"ref(qualifiedName: $branch)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

func getCommits(query *historyQuery) []CommitNode {
	return query.Repository.Ref.Target.Commit.History.Nodes
}

type CommitNode struct {
	Oid     string
	Message string
	Author  Author
}

type Author struct {
	Name  string
	Email string
	Date  time.Time
}

/*
Query config file content without checkout
graphql:

query ($owner: String!, $repo: String!, $file: String!) {
  repository(owner: $owner, name: $repo) {
    content: object(expression: $file) {
      ... on Blob {
        text
      }
    }
  }
}

*/
type configFileQuery struct {
	Repository struct {
		Content struct {
			Blob struct {
				Text string
			} `graphql:"... on Blob"`
		} `graphql:"content:object(expression: $file)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

func (query *configFileQuery) hasFile() bool {
	return query.Repository.Content.Blob.Text != ""
}

func (query *configFileQuery) getBytes() []byte {
	return []byte(query.Repository.Content.Blob.Text)
}

func (query *configFileQuery) mustGetConfig() *model.Config {
	var c model.Config
	err := yaml.Unmarshal(query.getBytes(), &c)
	if err != nil {
		log.Fatal(err)
	}
	return &c
}
