package formatter

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"os"
	"text/template"
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
	tmpl, err := template.New("releases").Parse(`
Releases:
{{ range . }}
{{ .CurrentVersion }}
{{- end }}
`)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, r.releases)
	if err != nil {
		panic(err)
	}

}
