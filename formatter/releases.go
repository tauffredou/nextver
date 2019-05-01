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

func (f *ReleasesFormatter) Json() {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(f.releases)
}

func (f *ReleasesFormatter) Yaml() {
	encoder := yaml.NewEncoder(os.Stdout)
	_ = encoder.Encode(f.releases)
}

func (f *ReleasesFormatter) Console() {
	t := NewTable(os.Stdout, "ref", "release")
	for _, v := range f.releases {
		_ = t.AnalyseRow(v.Ref, v.CurrentVersion)
	}

	t.WriteHeaders()
	for _, v := range f.releases {
		t.WriteRow(v.Ref, v.CurrentVersion)
	}

}
