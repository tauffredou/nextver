package formatter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mitchellh/colorstring"

	pad "github.com/willf/pad/utf8"
	"gopkg.in/yaml.v2"
	"os"
)

type ChangelogFormatter struct {
	release  *ReleaseDTO
	colorize bool
}

func NewChangelogFormatter(release *ReleaseDTO, colorize bool) *ChangelogFormatter {
	return &ChangelogFormatter{
		release:  release,
		colorize: colorize,
	}
}

func (c *ChangelogFormatter) Json() {
	encoder := json.NewEncoder(os.Stdout)
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

	var cols = []int{
		len("Date"),
		len("Author"),
		len("Kind"),
		len("Level"),
		len("Scope"),
		len("Title"),
	}
	for i := range r.Changelog {
		row := []string{
			consoleDateFormat,
			r.Changelog[i].Author,
			r.Changelog[i].Kind,
			r.Changelog[i].Level,
			r.Changelog[i].Scope,
			r.Changelog[i].Title,
		}

		for j := 0; j < 6; j++ {
			l := len(row[j])
			if l > cols[j] {
				cols[j] = l
			}
		}
	}

	var buffer bytes.Buffer

	buffer.WriteString(" ")
	buffer.WriteString(pad.Right("date", cols[0], " "))
	buffer.WriteString(" | ")
	buffer.WriteString(pad.Right("author", cols[1], " "))
	buffer.WriteString(" | ")
	buffer.WriteString(pad.Right("kind", cols[2], " "))
	buffer.WriteString(" | ")
	buffer.WriteString(pad.Right("Level", cols[3], " "))
	buffer.WriteString(" | ")
	buffer.WriteString(pad.Right("Scope", cols[4], " "))
	buffer.WriteString(" | ")
	buffer.WriteString(pad.Right("Message", cols[5], " "))

	fmt.Println(buffer.String())

	buffer = bytes.Buffer{}

	buffer.WriteString(" ")
	buffer.WriteString(pad.Right("", cols[0], "-"))
	buffer.WriteString(" | ")
	buffer.WriteString(pad.Right("", cols[1], "-"))
	buffer.WriteString(" | ")
	buffer.WriteString(pad.Right("", cols[2], "-"))
	buffer.WriteString(" | ")
	buffer.WriteString(pad.Right("", cols[3], "-"))
	buffer.WriteString(" | ")
	buffer.WriteString(pad.Right("", cols[4], "-"))
	buffer.WriteString(" | ")
	buffer.WriteString(pad.Right("", cols[5], "-"))

	fmt.Println(buffer.String())

	for i := range r.Changelog {
		ri := r.Changelog[i]
		var buffer bytes.Buffer

		if c.colorize {
			switch ri.Level {
			case "MAJOR":
				buffer.WriteString("[red]")
			case "MINOR":
				buffer.WriteString("[blue]")
			case "PATCH":
				buffer.WriteString("[yellow]")
			}
		}

		buffer.WriteString(" ")
		buffer.WriteString(pad.Right(ri.Date.Format(consoleDateFormat), cols[0], " "))
		buffer.WriteString(" | ")
		buffer.WriteString(pad.Right(ri.Author, cols[1], " "))
		buffer.WriteString(" | ")
		buffer.WriteString(pad.Right(ri.Kind, cols[2], " "))
		buffer.WriteString(" | ")
		buffer.WriteString(pad.Right(ri.Level, cols[3], " "))
		buffer.WriteString(" | ")
		buffer.WriteString(pad.Right(ri.Scope, cols[4], " "))
		buffer.WriteString(" | ")
		buffer.WriteString(pad.Right(ri.Title, cols[5], " "))

		fmt.Println(colorstring.Color(buffer.String()))
	}
}
