package model

import (
	"regexp"
	"strings"
	"time"
)

type ReleaseItem struct {
	ID     string
	Kind   string
	Scope  string
	Title  string
	Detail string
	Level  byte
	Author string
	Date   time.Time
}

func (ri *ReleaseItem) LevelName() string {
	switch ri.Level {
	case MAJOR:
		return "MAJOR"
	case MINOR:
		return "MINOR"
	case PATCH:
		return "PATCH"
	}
	return ""
}

func NewReleaseItem(author string, date time.Time, message string) ReleaseItem {

	ri := ReleaseItem{
		Author: author,
		Date:   date,
	}

	// Read first line
	index := strings.Index(message, "\n")

	var fl string
	if index == -1 {
		fl = message
		ri.Detail = ""
	} else {
		fl = message[0:index]
		ri.Detail = strings.Trim(message[index+1:], "\n ")
	}

	re := regexp.MustCompile(ConventionalCommitRegexp)

	lower := strings.ToLower(fl)
	if re.MatchString(lower) {
		data := re.FindStringSubmatch(lower)
		ri.Kind = strings.ToLower(data[1])
		ri.Scope = data[3]
		ri.Title = data[4]
	} else {
		ri.Title = strings.Trim(fl, "\n ")
	}

	switch {
	case strings.Contains(message, "BREAKING CHANGE"):
		ri.Level = MAJOR
	case ri.Kind == "feat":
		ri.Level = MINOR
	case ri.Kind == "fix":
		ri.Level = PATCH
	default:
		ri.Level = UNDEFINED
	}

	return ri
}
