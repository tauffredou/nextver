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
					Oid     string
					TagInfo TagNode `graphql:"... on Tag"`
				}
			} `graphql:"nodes"`
		} `graphql:"refs(refPrefix: \"refs/tags/\", last: 50, orderBy: {field: TAG_COMMIT_DATE, direction: ASC})"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

func (q *latestReleasesQuery) getTag(index int) (*TagNode, error) {
	return &q.Repository.Refs.Nodes[index].Target.TagInfo, nil
}

type TagNode struct {
	Message string
	Target  struct {
		Oid string
	}
}

func (t *TagNode) getCommitId() string {
	return t.Target.Oid
}

func (t *TagNode) getId() string {
	return strings.Trim(t.Message, "\n")
}

/*
tagsQuery retrieve the last 50 tags
graphql query:

query($owner: String!, $repo: String!){
  repository(owner: $owner, name: $repo) {
   refs(refPrefix: "refs/tags/", last: 50, orderBy: {field: TAG_COMMIT_DATE, direction: ASC}) {
     nodes {
       name
       target {
         oid
         ... on Tag {
           message
           target {
             oid
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
			TagNodes []struct {
				Name    string
				TagInfo TagNode `graphql:"target"`
			} `graphql:"nodes"`
		} `graphql:"refs(refPrefix: \"refs/tags/\", last: 50, orderBy: {field: TAG_COMMIT_DATE, direction: ASC})"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

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

query($branch: String!, $owner: String!, $repo: String!){
  repository(owner: $owner, name: $repo) {
    ref(qualifiedName: $branch ) {
      target {
        ... on Commit {
          history(first: 50) {
            nodes{
              message
              oid
              author{
                name
                email
              }
            }
            pageInfo {
              hasNextPage
            }
          }
        }
      }
    }
  }
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
					} `graphql:"history(first: $itemsCount)"`
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
