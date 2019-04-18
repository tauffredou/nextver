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

	pattern  = kingpin.Flag("pattern", "Versionning pattern").Short('p').Default("vSEMVER").String()
	output   = kingpin.Flag("output", "Output format (console, json, yaml)").Short('o').Default("console").String()
	branch   = kingpin.Flag("branch", "Target branch (default branch if empty)").Short('b').String()
	logLevel = kingpin.Flag("log-level", "Log level").Default("info").String()

	color = kingpin.Flag("color", "Colorize output").Default("true").Bool()

	//get
	getCommand         = kingpin.Command("get", "")
	listReleaseCommand = githubCommand(getCommand, "releases", "List releases")
	changelogCommand   = githubCommand(getCommand, "changelog", "Get changelog")

	//create
	createCommand = kingpin.Command("create", "")

	createReleaseCommand = createCommand.Command("release", "Create release")
	releaseTemplate      = createCommand.Flag("template", "Template file").String()
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
	case "get releases":
		r := github().GetReleases()
		m := formatter.MapReleases(r)
		f = formatter.NewReleasesFormatter(m)

	case "create release":
		log.Warn("not implemented yet")

	case "get changelog":
		r := github().GetLatestRelease()
		dto := formatter.MapRelease(&r)
		f = formatter.NewChangelogFormatter(&dto, *color)
	}

	switch *output {
	case "console":
		f.Console()
		//formatter.ConsoleOutput(&r, *color)
	case "json":
		f.Json()
		//formatter.JsonOutput(&r)
	case "yaml":
		f.Yaml()
		//formatter.YamlOutput(&r)
	}

}
