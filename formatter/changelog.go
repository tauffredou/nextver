package formatter

import (
	"encoding/json"
	"fmt"
	"github.com/Masterminds/sprig"
	"github.com/willf/pad"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"text/template"
)

type ChangelogFormatter struct {
	release  *ReleaseDTO
	colorize bool
	output   io.Writer
}

func NewChangelogFormatter(release *ReleaseDTO, colorize bool) *ChangelogFormatter {
	return &ChangelogFormatter{
		release:  release,
		colorize: colorize,
		output:   os.Stdout,
	}
}

func (c *ChangelogFormatter) Json() {
	encoder := json.NewEncoder(c.output)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(c.release)
}

func (c *ChangelogFormatter) Yaml() {
	encoder := yaml.NewEncoder(os.Stdout)
	_ = encoder.Encode(c.release)
}

func (c *ChangelogFormatter) Console() {
	r := c.release

	fmt.Printf("Current release version\t: %s\n", r.CurrentVersion)
	fmt.Printf("Next release version\t: %s\n", r.NextVersion)

	fmt.Println("\nChangelog:")

	if len(r.Changelog) == 0 {
		fmt.Println("No change since last release")
		return
	}

	t := NewTable(os.Stdout, "Date", "Author", "Kind", "Level", "Scope", "Title")
	if c.colorize {
		t.SetColorizer(func(row []string, index int) string {
			if index != 3 {
				return ""
			}
			switch row[3] {
			case "MAJOR":
				return "red"
			case "MINOR":
				return "blue"
			case "PATCH":
				return "yellow"
			}

			return ""
		})
	}

	for i := range r.Changelog {
		_ = t.AnalyseRow(consoleDateFormat,
			r.Changelog[i].Author,
			r.Changelog[i].Kind,
			r.Changelog[i].Level,
			r.Changelog[i].Scope,
			r.Changelog[i].Title)
	}

	t.WriteHeaders()
	for i := range r.Changelog {
		t.WriteRow(r.Changelog[i].Date.Format(consoleDateFormat),
			r.Changelog[i].Author,
			r.Changelog[i].Kind,
			r.Changelog[i].Level,
			r.Changelog[i].Scope,
			r.Changelog[i].Title)
	}

}

func (c *ChangelogFormatter) Template(text string) error {
	tpl, err := template.New("Release").
		Funcs(sprig.TxtFuncMap()).
		Funcs(FuncMap()).
		Parse(text)
	if err != nil {
		return err
	}

	return tpl.Execute(c.output, c.release)
}

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"padRight": pad.Right,
	}
}
