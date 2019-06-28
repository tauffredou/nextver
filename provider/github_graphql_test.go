package provider

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shurcooL/githubv4"
	"github.com/stretchr/testify/assert"
)

func TestGithubProvider_getFirstTag(t *testing.T) {
	resp := mockResponseFile("../fixtures/github/latestReleases.response.json")

	p := &GithubProvider{client: mockGithubClient(resp)}

	q := p.mustQueryReleases()
	tag := q.GetTags()[0]
	assert.Equal(t, "v1.1.0", tag.getId())
	assert.Equal(t, "a3240571ac4bbe857a0cfad3b988942838e758d1", tag.getCommitId())

	//var q latestReleasesQuery
}

func TestGithubProvider_GetReleases(t *testing.T) {
	resp := mockResponse(`{
  "data": {
    "repository": {
      "refs": {
        "nodes": [
          {
            "name": "v1.0.1",
            "target": {
              "message": "v1.0.1\n",
              "target": {
                "oid": "1c23cc36d1383b82198af6ee04fe44b820b6a550"
              }
            }
          },
          {
            "name": "v1.1.0",
            "target": {
              "message": "v1.1.0\n",
              "target": {
                "oid": "a3240571ac4bbe857a0cfad3b988942838e758d1"
              }
            }
          }
        ]
      }
    }
  }
}`)

	p := &GithubProvider{client: mockGithubClient(resp)}

	_ = p.mustQueryReleases()
}

func TestGithubProvider_getReleaseBoundary_empty(t *testing.T) {
	resp := mockResponse(`{
  "data": {
    "repository": {
      "refs": {
        "nodes": []
      }
    }
  }
}`)

	p := &GithubProvider{client: mockGithubClient(resp)}

	first, last, err := p.getReleaseBoundary("v1.1.0")
	assert.NoError(t, err)
	assert.Equal(t, "", first)
	assert.Equal(t, "", last)
}

func TestGithubProvider_getReleaseBoundary_one(t *testing.T) {
	resp := mockResponse(`{
  "data": {
    "repository": {
      "refs": {
        "nodes": [
          {
            "name": "v1.1.0",
            "target": {
              "message": "v1.1.0\n",
              "target": {
                "oid": "a3240571ac4bbe857a0cfad3b988942838e758d1"
              }
            }
          }
        ]
      }
    }
  }
}`)

	p := &GithubProvider{client: mockGithubClient(resp)}

	first, last, err := p.getReleaseBoundary("v1.1.0")
	assert.NoError(t, err)
	assert.Equal(t, "a3240571ac4bbe857a0cfad3b988942838e758d1", first)
	assert.Equal(t, "", last)
}

func TestGithubProvider_getReleaseBoundary(t *testing.T) {
	resp := mockResponseFile("../fixtures/github/releases.response.json")

	p := &GithubProvider{client: mockGithubClient(resp)}

	first, last, err := p.getReleaseBoundary("v1.1.0")
	assert.NoError(t, err)
	assert.Equal(t, "a3240571ac4bbe857a0cfad3b988942838e758d1", first)
	assert.Equal(t, "1c23cc36d1383b82198af6ee04fe44b820b6a550", last)
}

func TestGithubProvider_MustGetPattern(t *testing.T) {
	resp := mockResponse(`{
  "data": {
    "repository": {
      "content": {
        "text": "version: \"beta\"\n\npattern: v-SEMVER\n"
      }
    }
  }
}`)
	p := &GithubProvider{
		client: mockGithubClient(resp),
		config: &GithubProviderConfig{
			Branch: "master",
		},
	}

	actual := p.MustGetPattern()
	expected := "v-SEMVER"
	assert.Equal(t, expected, actual)
}

func TestGithubProvider_MustGetPattern_definedInProvider(t *testing.T) {
	p := &GithubProvider{
		Pattern: "test2",
	}

	assert.Equal(t, "test2", p.MustGetPattern())
}
func TestGithubProvider_MustGetPattern_definedInConfiguration(t *testing.T) {

	p := &GithubProvider{
		config: &GithubProviderConfig{
			Pattern: "test3",
		},
	}

	assert.Equal(t, "test3", p.MustGetPattern())
}

func TestGithubProvider_historyQuery(t *testing.T) {
	resp := mockResponseFile("../fixtures/github/history.response.json")
	p := &GithubProvider{
		client: mockGithubClient(resp),
	}

	actual := p.mustGetHistory(nil)
	assert.Len(t, getCommits(actual), 5)
}

func mockResponseFile(f string) *http.ServeMux {
	content, err := ioutil.ReadFile(f)
	if err != nil {
		log.Fatal(err)
	}
	return mockResponse(string(content))
}

/* helpers */

func mockGithubClient(mux *http.ServeMux) *githubv4.Client {
	return githubv4.NewClient(&http.Client{Transport: localRoundTripper{handler: mux}})
}

func mockResponse(resp string) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/graphql", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		mustWrite(w, resp)
	})
	return mux
}

func mustWrite(w io.Writer, s string) {
	_, err := io.WriteString(w, s)
	if err != nil {
		panic(err)
	}
}

type localRoundTripper struct {
	handler http.Handler
}

func (l localRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	w := httptest.NewRecorder()
	l.handler.ServeHTTP(w, req)
	return w.Result(), nil
}
