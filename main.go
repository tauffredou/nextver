package main

import (
	"github.com/jobteaser/github-release/formatter"
	"github.com/jobteaser/github-release/provider"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	token   = kingpin.Flag("github-token", "Github token").Envar("GITHUB_TOKEN").Required().String()
	owner   = kingpin.Flag("github-owner", "Github owner").Required().String()
	color   = kingpin.Flag("color", "Colorize output").Default("true").Bool()
	repo    = kingpin.Flag("github-repo", "Github repo").Required().String()
	pattern = kingpin.Flag("pattern", "Versionning pattern").Short('p').Default("vSEMVER").String()
	output  = kingpin.Flag("output", "Output format (console, json, yaml)").Short('o').Default("console").String()
)

func main() {

	kingpin.Parse()

	p := provider.NewGithubProvider(*owner, *repo, *token, *pattern)
	r := p.GetLatestRelease()

	switch *output {
	case "console":
		formatter.ConsoleOutput(&r, *color)
	case "json":
		formatter.JsonOutput(&r)
	case "yaml":
		formatter.YamlOutput(&r)
	}

}
