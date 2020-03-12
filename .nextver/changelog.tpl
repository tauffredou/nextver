Changelog
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
