package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/tauffredou/nextver/formatter"
	"github.com/tauffredou/nextver/provider"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"path"
)

var (
	tokenFlag = kingpin.Flag("github-token", "Github token. Can be read form hub config file").Envar("GITHUB_TOKEN").String()

	repo     = kingpin.Flag("repo", "Repository").Default(".").Short('r').String()
	pattern  = kingpin.Flag("pattern", "Versionning pattern. Read from .nextver/config.yml by default").Short('p').String()
	output   = kingpin.Flag("output", "Output format (console, json, yaml)").Short('o').Default("console").String()
	branch   = kingpin.Flag("branch", "Target branch (default branch if empty)").Short('b').String()
	logLevel = kingpin.Flag("log-level", "Log level").Default("info").String()

	color = kingpin.Flag("color", "Colorize output").Default("true").Bool()

	//get
	getCommand       = kingpin.Command("get", "")
	_                = getCommand.Command("releases", "List releases")
	changelogCommand = getCommand.Command("changelog", "Get changelog")
	release          = changelogCommand.Flag("release", "Changelog for a specific release").Default("").String()
	_                = getCommand.Command("next-version", "Get next version")

	//create
	createCommand = kingpin.Command("create", "")

	_ = createCommand.Command("release", "Create release")
	_ = createCommand.Flag("template", "Template file").String()

	defaultHubConfig = path.Join(MustString(os.UserHomeDir()), ".config", "hub")
)

func MustString(s string, err error) string {
	if err != nil {
		log.Fatal(err)
	}
	return s
}

func githubToken() string {
	if *tokenFlag == "" {
		t, err := provider.ReadHubToken(defaultHubConfig)
		if err != nil {
			log.Fatal("required flag '--github-token'")
		}
		return t
	} else {
		return *tokenFlag
	}
}

func mustSetLoglevel(level string) {
	l, err := log.ParseLevel(level)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(l)
}

func main() {
	var err error
	parse := kingpin.Parse()

	mustSetLoglevel(*logLevel)
	log.SetOutput(os.Stderr)
	var f formatter.Formatter

	pf := provider.ProviderFactory{
		Pattern:     *pattern,
		TokenReader: githubToken,
	}

	prov, err := pf.CreateProvider(*repo)
	if err != nil {
		log.Fatal(err)
	}

	switch parse {
	case "get next-version":
		f = getNextVersion(prov)
		// Releases
	case "get releases":
		f = getReleases(prov)
	case "create release":
		log.Warn("not implemented yet")
	case "get changelog":
		f = getChangelog(prov)
	}

	if f != nil {
		switch *output {
		case "console":
			f.Console()
		case "json":
			f.Json()
		case "yaml":
			f.Yaml()
		}
	}

}

func getNextVersion(prov provider.Provider) formatter.Formatter {
	r, err := prov.GetRelease("")
	if err != nil {
		log.Fatal(err)
	}
	v, _ := r.NextVersion()
	return &formatter.SimpleFormatter{Key: "next-version", Value: v}
}

func getReleases(prov provider.Provider) formatter.Formatter {
	r, err := prov.GetReleases()
	if err != nil {
		log.Fatal(err)
	}
	m := formatter.MapReleases(r)
	return formatter.NewReleasesFormatter(m)
}

func getChangelog(prov provider.Provider) formatter.Formatter {
	r, err := prov.GetRelease(*release)
	checkErr(err)
	dto := formatter.MapRelease(r)
	return formatter.NewChangelogFormatter(&dto, *color)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
