package formatter

import (
	"github.com/tauffredou/nextver/model"
)

const consoleDateFormat = "06/01/02 15:04"

type Formatter interface {
	Json()
	Yaml()
	Console()
}

func MapReleases(items []model.Release) []ReleaseDTO {
	res := make([]ReleaseDTO, len(items))
	for i := range items {
		res[i] = MapRelease(&items[i])
	}
	return res
}

func MapRelease(release *model.Release) ReleaseDTO {

	return ReleaseDTO{
		Project:        release.Project,
		Changelog:      mapReleaseItem(release.Changelog),
		NextVersion:    release.MustNextVersion(),
		CurrentVersion: release.CurrentVersion,
		Ref:            release.Ref,
		VersionPattern: release.VersionPattern,
	}
}

func mapReleaseItem(items []model.ReleaseItem) []ReleaseItemDTO {
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
