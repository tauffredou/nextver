package formatter

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"os"
)

type ReleasesFormatter struct {
	releases []ReleaseDTO
}

func NewReleasesFormatter(releases []ReleaseDTO) *ReleasesFormatter {
	return &ReleasesFormatter{releases: releases}
}

func (c *ReleasesFormatter) Json() {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(c.releases)
}

func (rf *ReleasesFormatter) Yaml() {
	encoder := yaml.NewEncoder(os.Stdout)
	_ = encoder.Encode(rf.releases)
}

func (r *ReleasesFormatter) Console() {
	t := NewTable(os.Stdout, "ref", "release")
	for _, v := range r.releases {
		_ = t.AnalyseRow(v.Ref, v.CurrentVersion)
	}

	t.WriteHeaders()
	for _, v := range r.releases {
		t.WriteRow(v.Ref, v.CurrentVersion)
	}

}
