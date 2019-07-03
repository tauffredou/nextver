package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/tauffredou/nextver/formatter"
	"github.com/tauffredou/nextver/model"
	"github.com/tauffredou/nextver/provider"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"path"
)

var (
	token string

	repo         = kingpin.Flag("repo", "Repository").Short('r').String()
	pattern      = kingpin.Flag("pattern", "Versionning pattern. Read from .nextver/config.yml by default").Short('p').String()
	output       = kingpin.Flag("output", "Output format (console, json, yaml)").Short('o').Default("console").String()
	branch       = kingpin.Flag("branch", "Target branch (default branch if empty)").Short('b').String()
	logLevel     = kingpin.Flag("log-level", "Log level").Default("info").String()
	providerType = kingpin.Flag("provider", "provider").Default("local").String()

	color = kingpin.Flag("color", "Colorize output").Default("true").Bool()

	//get
	getCommand       = kingpin.Command("get", "")
	_                = githubCommand(getCommand, "releases", "List releases")
	changelogCommand = githubCommand(getCommand, "changelog", "Get changelog")
	release          = changelogCommand.Flag("release", "Changelog for a specific release").String()
	_                = githubCommand(getCommand, "next-version", "Get next version")

	//create
	createCommand = kingpin.Command("create", "")

	_ = createCommand.Command("release", "Create release")
	_ = createCommand.Flag("template", "Template file").String()

	defaultHubConfig = path.Join(MustString(os.UserHomeDir()), ".config", "hub")
	githubProvider   *provider.GithubProvider
)

func MustString(s string, err error) string {
	if err != nil {
		log.Fatal(err)
	}
	return s
}

func githubCommand(command *kingpin.CmdClause, name string, help string) *kingpin.CmdClause {
	c := command.Command(name, help)
	c.Flag("github-token", "Github token. Can be read form hub config file").Envar("GITHUB_TOKEN").StringVar(&token)

	return c
}

// github lazy load github provider
func github() *provider.GithubProvider {
	if githubProvider == nil {

		if token == "" {
			t, err := provider.ReadHubToken(defaultHubConfig)
			if err != nil {
				log.Fatalf("required flag '--%s'", "github-token")
			}
			token = t
		}

		r, err := provider.ParseRepo(*repo)
		if err != nil {
			log.Fatalf("Cannot parse repository")
		}
		githubRepo := r.(provider.GithubRepository)
		githubProvider, _ = provider.NewGithubProvider(githubRepo.Owner, githubRepo.Repo, token, &provider.GithubProviderConfig{
			Pattern: *pattern,
			Branch:  *branch,
		})
	}

	return githubProvider

}

func MustSetLoglevel(level string) {
	l, err := log.ParseLevel(level)
	if err != nil {
		log.Panic(err)
	}
	log.SetLevel(l)
}

func main() {
	parse := kingpin.Parse()

	MustSetLoglevel(*logLevel)
	log.WithFields(log.Fields{
		"pattern": *pattern,
	}).Debug(parse)

	var f formatter.Formatter

	switch parse {
	case "get next-version":
		f = getNextVersion()
		// Releases
	case "get releases":
		f = getReleases()
	case "create release":
		log.Warn("not implemented yet")
	case "get changelog":
		f = getChangelog()
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

func getNextVersion() formatter.Formatter {
	r := github().GetNextRelease()
	v, _ := r.NextVersion()
	return &formatter.SimpleFormatter{Key: "next-version", Value: v}
}

func getReleases() formatter.Formatter {
	r, err := github().GetReleases()
	if err != nil {
		log.Fatal(err)
	}
	m := formatter.MapReleases(r)
	return formatter.NewReleasesFormatter(m)
}

func getChangelog() formatter.Formatter {
	var r *model.Release
	if *release == "" {
		r = github().GetNextRelease()
	} else {
		var err error
		r, err = github().GetRelease(*release)
		checkErr(err)
	}
	dto := formatter.MapRelease(r)
	return formatter.NewChangelogFormatter(&dto, *color)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
