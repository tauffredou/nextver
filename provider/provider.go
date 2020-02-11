package provider

import (
	"github.com/tauffredou/nextver/model"
	"regexp"
	"strings"
)

type Provider interface {
	GetReleases() ([]model.Release, error)
	GetRelease(name string) (*model.Release, error)
}

func GetVersionRegexp(pattern string) *regexp.Regexp {
	replacer := strings.NewReplacer(
		"SEMVER", model.SemverRegex,
		"DATE", model.DateRegexp,
	)
	return regexp.MustCompile("^" + replacer.Replace(pattern) + "$")
}
