package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/tauffredou/nextver/formatter"
	"github.com/tauffredou/nextver/provider"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	token string
	owner string
	repo  string

	pattern      = kingpin.Flag("pattern", "Versionning pattern").Short('p').Default("vSEMVER").String()
	output       = kingpin.Flag("output", "Output format (console, json, yaml)").Short('o').Default("console").String()
	branch       = kingpin.Flag("branch", "Target branch (default branch if empty)").Short('b').String()
	logLevel     = kingpin.Flag("log-level", "Log level").Default("info").String()
	providerType = kingpin.Flag("provider", "provider").Default("local").String()

	color = kingpin.Flag("color", "Colorize output").Default("true").Bool()

	//get
	getCommand = kingpin.Command("get", "")
	_          = githubCommand(getCommand, "releases", "List releases")
	_          = githubCommand(getCommand, "changelog", "Get changelog")
	_          = githubCommand(getCommand, "next-version", "Get next version")

	//create
	createCommand = kingpin.Command("create", "")

	_ = createCommand.Command("release", "Create release")
	_ = createCommand.Flag("template", "Template file").String()
)

func githubCommand(command *kingpin.CmdClause, name string, help string) *kingpin.CmdClause {
	c := command.Command(name, help)
	c.Flag("github-token", "Github token").Envar("GITHUB_TOKEN").Required().StringVar(&token)
	c.Flag("github-owner", "Github owner").Required().StringVar(&owner)
	c.Flag("github-repo", "Github repo").Required().StringVar(&repo)
	return c
}

func github() *provider.GithubProvider {
	return provider.NewGithubProvider(owner, repo, token, *pattern, *branch)
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
		"token": token,
	}).Debug()

	log.WithField("command", parse).Debug("Action")

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
	r := github().GetLatestRelease()
	v, _ := r.NextVersion()
	return &formatter.SimpleFormatter{"next-version", v}
}

func getReleases() formatter.Formatter {
	r := github().GetReleases()
	m := formatter.MapReleases(r)
	return formatter.NewReleasesFormatter(m)
}

func getChangelog() formatter.Formatter {
	r := github().GetLatestRelease()
	dto := formatter.MapRelease(&r)
	return formatter.NewChangelogFormatter(&dto, *color)
}
