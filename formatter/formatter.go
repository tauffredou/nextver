package formatter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jobteaser/github-release/provider"
	"github.com/mitchellh/colorstring"
	"github.com/willf/pad"
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

func ConsoleOutput(release *provider.Release, colorize bool) {
	r := mapRelease(release)

	fmt.Printf("Current release version\t: %s\n", r.CurrentVersion)
	fmt.Printf("Next release version\t: %s\n", r.NextVersion)

	fmt.Println("\nCommit log:")

	if len(r.Changelog) == 0 {
		fmt.Println("No change since last release")
		return
	}

	var cols = []int{
		len("Kind"),
		len("Level"),
		len("Scope"),
		len("Title"),
	}
	for i := range r.Changelog {
		row := []string{
			r.Changelog[i].Kind,
			r.Changelog[i].Level,
			r.Changelog[i].Scope,
			r.Changelog[i].Title,
		}

		for j := 0; j < 4; j++ {
			l := len(row[j])
			if l > cols[j] {
				cols[j] = l
			}
		}
	}

	var buffer bytes.Buffer

	buffer.WriteString(" ")
	buffer.WriteString(pad.Right("kind", cols[0], " "))
	buffer.WriteString(" | ")
	buffer.WriteString(pad.Right("Level", cols[1], " "))
	buffer.WriteString(" | ")
	buffer.WriteString(pad.Right("Scope", cols[2], " "))
	buffer.WriteString(" | ")
	buffer.WriteString(pad.Right("Message", cols[3], " "))

	fmt.Println(buffer.String())

	buffer = bytes.Buffer{}

	buffer.WriteString(" ")
	buffer.WriteString(pad.Right("", cols[0], "-"))
	buffer.WriteString(" | ")
	buffer.WriteString(pad.Right("", cols[1], "-"))
	buffer.WriteString(" | ")
	buffer.WriteString(pad.Right("", cols[2], "-"))
	buffer.WriteString(" | ")
	buffer.WriteString(pad.Right("", cols[3], "-"))

	fmt.Println(buffer.String())

	for i := range r.Changelog {
		ri := r.Changelog[i]
		var buffer bytes.Buffer

		if colorize {
			switch ri.Level {
			case "MAJOR":
				buffer.WriteString("[red]")
			case "MINOR":
				buffer.WriteString("[blue]")
			case "PATCH":
				buffer.WriteString("[yellow]")
			}
		}

		buffer.WriteString(" ")
		buffer.WriteString(pad.Right(ri.Kind, cols[0], " "))
		buffer.WriteString(" | ")
		buffer.WriteString(pad.Right(ri.Level, cols[1], " "))
		buffer.WriteString(" | ")
		buffer.WriteString(pad.Right(ri.Scope, cols[2], " "))
		buffer.WriteString(" | ")
		buffer.WriteString(pad.Right(ri.Title, cols[3], " "))

		fmt.Println(colorstring.Color(buffer.String()))
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
