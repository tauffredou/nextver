package formatter

import (
	"github.com/tauffredou/nextver/provider"
	"time"
)

const consoleDateFormat = "06/01/02 15:04"

type Formatter interface {
	Json()
	Yaml()
	Console()
}

type ReleaseDTO struct {
	Project        string           `json:"project"`
	Ref            string           `json:"ref,omitempty"`
	CurrentVersion string           `json:"current_version"`
	NextVersion    string           `json:"next_version"`
	Changelog      []ReleaseItemDTO `json:"changelog"`
	VersionPattern string           `json:"version_pattern"`
}

func MapReleases(items []provider.Release) []ReleaseDTO {
	res := make([]ReleaseDTO, len(items))
	for i := range items {
		res[i] = MapRelease(&items[i])
	}
	return res
}

func MapRelease(release *provider.Release) ReleaseDTO {

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
	Kind   string    `json:"kind,omitempty"`
	Scope  string    `json:"scope,omitempty"`
	Title  string    `json:"title"`
	Detail string    `json:"detail,omitempty"`
	Level  string    `json:"level"`
	Author string    `json:"author"`
	Date   time.Time `json:"date"`
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
			Date:   item.Date,
			Author: item.Author,
		}
	}
	return res
}
