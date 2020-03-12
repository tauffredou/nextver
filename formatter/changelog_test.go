package formatter

import (
	"github.com/stretchr/testify/assert"
	"github.com/tauffredou/go-test/fake_clock"
	"github.com/tauffredou/nextver/model"
	"strings"
	"testing"
)

func TestChangelogFormatter_Template(t *testing.T) {

	clock := fake_clock.New()

	r := &ReleaseDTO{
		Project:        "test",
		Ref:            "test1234",
		CurrentVersion: "0.1.0",
		NextVersion:    "0.2.0",
		Changelog: []ReleaseItemDTO{
			{
				Kind:   "feat",
				Scope:  "mankind",
				Title:  "Save the cheerleader",
				Detail: "Save the world",
				Level:  model.ChangeLevelMinor,
				Author: "tauf",
				Date:   clock.Tick(),
			},
			{
				Kind:   "patch",
				Scope:  "auth",
				Title:  "fix security",
				Level:  model.ChangeLevelPatch,
				Author: "tauf",
				Date:   clock.Tick(),
			},
		},
		VersionPattern: "vSEMVER",
	}
	f := NewChangelogFormatter(r, false)

	tpl := `Changelog of {{.Project}}
version: {{.NextVersion}}
{{ "" }}
{{- if .HasChanges "MAJOR" }}
Breaking changes:
{{ end -}}
{{ range .ChangesByLevel "MAJOR" -}}
- {{.Level}} {{ .Title }}
{{ end -}}

{{- if .HasChanges "MINOR" }}
Features:
{{ end -}}
{{ range .ChangesByLevel "minor" -}}
- {{.Level}} {{ .Title }}
{{ end -}}

{{- if .HasChanges "Patch" }}
Security fixes:
{{ end -}}
{{ range .ChangesByLevel "patch" -}}
- {{.Level}} {{ .Title }}
{{ end -}}
`

	expected := `Changelog of test
version: 0.2.0

Features:
- MINOR Save the cheerleader

Security fixes:
- PATCH fix security
`

	sb := &strings.Builder{}
	f.output = sb
	err := f.Template(tpl)
	assert.NoError(t, err)
	assert.Equal(t, expected, sb.String())
}
