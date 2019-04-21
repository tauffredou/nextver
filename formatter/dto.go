package formatter

import "time"

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
