package provider

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func cleanRepo(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		log.Fatal(err)
	}
}

func cloneRepo(url string, directory string) {
	var err error

	_, err = git.PlainClone(directory, false, &git.CloneOptions{
		URL:               url,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func TestGitProvider_GetReleases(t *testing.T) {
	path := fmt.Sprintf(filepath.Join(os.TempDir(), "nextver-%d"), time.Now().UnixNano())
	cloneRepo("https://github.com/tauffredou/test-semver.git", path)
	defer cleanRepo(path)

	p := NewGitProvider(path, "vSEMVER")
	actual := p.GetReleases()
	assert.Equal(t, 2, len(actual))
	assert.Equal(t, "v1.1.0", actual[0].CurrentVersion)
	assert.Equal(t, "v1.0.1", actual[1].CurrentVersion)

}

func TestGitProvider_tagFilter(t *testing.T) {

	tests := []struct {
		tag  string
		want bool
	}{
		{"refs/tags/someTag", false},
		{"refs/tags/v1.0.0", true},
	}

	p := NewGitProvider("", "vSEMVER")

	for _, tt := range tests {
		t.Run(tt.tag, func(t *testing.T) {
			reference := plumbing.NewReferenceFromStrings(tt.tag, "a")
			assert.Equal(t, p.tagFilter(reference), tt.want)
		})
	}
}
