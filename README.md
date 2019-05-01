# 

This application calculates the next version from the git history.
It detects commit kinds using the ["conventional commits" specification](https://www.conventionalcommits.org/)

## Build
```
go mod download
go test ./...
go build
```

## Usage

```
usage: nextver [<flags>] <command> [<args> ...]

Flags:
      --help               Show context-sensitive help (also try --help-long and --help-man).
  -p, --pattern="vSEMVER"  Versionning pattern
  -o, --output="console"   Output format (console, json, yaml)
  -b, --branch=BRANCH      Target branch (default branch if empty)
      --log-level="info"   Log level
      --provider="local"   provider
      --color              Colorize output

Commands:
  help [<command>...]
    Show help.

  get releases --github-token=GITHUB-TOKEN --github-owner=GITHUB-OWNER --github-repo=GITHUB-REPO
    List releases

  get changelog --github-token=GITHUB-TOKEN --github-owner=GITHUB-OWNER --github-repo=GITHUB-REPO
    Get changelog

  get next-version --github-token=GITHUB-TOKEN --github-owner=GITHUB-OWNER --github-repo=GITHUB-REPO
    Get next version

```

### Authentication

```
export GITHUB_TOKEN=xxxxxxxxx
```

### Release pattern

The following keyworks are supported
- `SEMVER`: use semantic versionning (ex 1.0.5)
- `DATE`: use timestamping for rolling versionning. The opiniated format is YYYY-MM-DD-HHmmss

Those keywords can be used in any pattern. Some examples:
```
SEMVER        -> 1.0.5
vSEMVER       -> v1.0.5
rDATE         -> r2019-04-01-133742
DATE          -> 2019-04-01-133742
release-DATE  -> release-2019-04-01-133742
``` 

#### Default options
```
nextver --github-owner=tauffredou --github-repo=test-semver
```

Will output 
```
Current release version	: 0.0.0
Next release version	: v0.1.0

Commit log:
    Kind    Level    Message
---------|---------|------------------------------
```

#### Using date release
```
nextver --github-owner=tauffredou --github-repo=test-semver --pattern=myprefix-DATE
```
This will output
```
Current release version	: 0.0.0
Next release version	: myprefix-2019-04-01-112621

Commit log:
    Kind    Level    Message
...
```

