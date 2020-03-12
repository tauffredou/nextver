package formatter

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"text/template"
)

type ReleasesFormatter struct {
	releases []ReleaseDTO
	output   io.Writer
}

func (f *ReleasesFormatter) Template(text string) error {
	tpl, err := template.New("Release").Parse(text)
	if err != nil {
		return err
	}

	return tpl.Execute(f.output, f.releases)
}

func NewReleasesFormatter(releases []ReleaseDTO) *ReleasesFormatter {
	return &ReleasesFormatter{
		releases: releases,
		output:   os.Stdout,
	}
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
