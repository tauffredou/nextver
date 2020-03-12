package formatter

import (
	"strings"
	"time"
)

type ChangeLevel string

type ReleaseItemDTO struct {
	Kind   string    `json:"kind,omitempty"`
	Scope  string    `json:"scope,omitempty"`
	Title  string    `json:"title"`
	Detail string    `json:"detail,omitempty"`
	Level  string    `json:"level"`
	Author string    `json:"author"`
	Date   time.Time `json:"date"`
}

type ReleaseDTO struct {
	Project        string           `json:"project"`
	Ref            string           `json:"ref,omitempty"`
	CurrentVersion string           `json:"current_version"`
	NextVersion    string           `json:"next_version"`
	Changelog      []ReleaseItemDTO `json:"changelog"`
	VersionPattern string           `json:"version_pattern"`
}

func (r *ReleaseDTO) HasChanges(level string) bool {
	return len(r.ChangesByLevel(level)) > 0
}

func (r *ReleaseDTO) ChangesByLevel(level string) []ReleaseItemDTO {
	level = strings.ToUpper(level)
	res := []ReleaseItemDTO{}
	for i := range r.Changelog {
		if r.Changelog[i].Level == level {
			res = append(res, r.Changelog[i])
		}
	}
	return res
}
