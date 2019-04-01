package formatter

import (
	"encoding/json"
	"fmt"
	"github.com/jobteaser/github-release/provider"
	"gopkg.in/yaml.v2"
	"os"
)

func JsonOutput(release *provider.Release) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(mapRelease(release))
}

func YamlOutput(release *provider.Release) {

	_ = yaml.NewEncoder(os.Stdout).Encode(mapRelease(release))
}

func ConsoleOutput(release *provider.Release) {
	r := mapRelease(release)

	fmt.Printf("Current release version\t: %s\n", r.CurrentVersion)
	fmt.Printf("Next release version\t: %s\n", r.NextVersion)

	fmt.Println("\nCommit log:")

	if len(r.Changelog) == 0 {
		fmt.Println("No change since last release")
	} else {
		fmt.Printf("% 8s    %5s    %s\n", "Kind", "Level", "Message")
		fmt.Println("---------|---------|------------------------------")
		for i := range r.Changelog {
			ri := r.Changelog[i]
			fmt.Printf("% 8s |  %5s  | %s\n", ri.Kind, ri.Level, ri.Title)
		}
	}
}

type ReleaseDTO struct {
	Project        string           `json:"project"`
	Ref            string           `json:"ref,omitempty"`
	CurrentVersion string           `json:"current_version"`
	NextVersion    string           `json:"next_version"`
	Changelog      []ReleaseItemDTO `json:"changelog"`
	VersionPattern string           `json:"version_pattern"`
}

func mapRelease(release *provider.Release) ReleaseDTO {

	return ReleaseDTO{
		Project:        release.Project,
		Changelog:      mapReleaseItem(release.Changelog),
		NextVersion:    release.MustNextVersion(),
		CurrentVersion: release.CurrentVersion,
		Ref:            release.Ref,
		VersionPattern: release.VersionPattern,
	}
}

type ReleaseItemDTO struct {
	Kind   string `json:"kind,omitempty"`
	Scope  string `json:"scope,omitempty"`
	Title  string `json:"title"`
	Detail string `json:"detail,omitempty"`
	Level  string `json:"level"`
}

func mapReleaseItem(items []provider.ReleaseItem) []ReleaseItemDTO {
	res := make([]ReleaseItemDTO, len(items))
	for i := range items {
		item := items[i]
		res[i] = ReleaseItemDTO{
			Kind:   item.Kind,
			Level:  item.LevelName(),
			Title:  item.Title,
			Scope:  item.Scope,
			Detail: item.Detail,
		}
	}
	return res
}
