<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [](#)
  - [Build](#build)
  - [Usage](#usage)
    - [Authentication](#authentication)
    - [Release pattern](#release-pattern)
      - [Default options](#default-options)
      - [Using date release](#using-date-release)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

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

Create a Github access token with at least the following scopes:

![github_scope](doc/images/github_scopes.png)


Then use this token using the option
```
nextver --github-token=xxxxxxxxx ...
```
or exporting the environment variable
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
SEMVER (vSEMVER) is the default release pattern
```
$ nextver get next-version --github-owner=tauffredou --github-repo=test-semver
v1.0.0
```

Getting details

```
$ nextver get changelog --github-owner=tauffredou --github-repo=test-semver 
Current release version : 0.0.0
Next release version    : v1.0.0

Changelog:
 Date           │ Author           │ Kind    │ Level │ Scope │ Title        
 ━━━━━━━━━━━━━━━┿━━━━━━━━━━━━━━━━━━┿━━━━━━━━━┿━━━━━━━┿━━━━━━━┿━━━━━━━━━━━━━━
 19/03/30 12:44 │ Thomas Auffredou │ chore   │ MAJOR │ test2 │ some change  
 19/03/30 12:29 │ Thomas Auffredou │ feat    │ MINOR │ test  │ super feature
 19/03/30 11:53 │ Thomas Auffredou │ initial │       │       │ commit       

```

#### Using date release
```
$ nextver get next-version --github-owner=tauffredou --github-repo=test-semver --pattern=myprefix-DATE
myprefix-2019-05-01-110159
```

The DATE pattern uses only time, ignoring the semantic versionning and the conventional commit convensions.

```
Current release version : 0.0.0
Next release version    : myprefix-2019-05-01-111206

Changelog:
 Date           │ Author           │ Kind    │ Level │ Scope │ Title        
 ━━━━━━━━━━━━━━━┿━━━━━━━━━━━━━━━━━━┿━━━━━━━━━┿━━━━━━━┿━━━━━━━┿━━━━━━━━━━━━━━
 19/03/30 12:44 │ Thomas Auffredou │ chore   │ MAJOR │ test2 │ some change  
 19/03/30 12:29 │ Thomas Auffredou │ feat    │ MINOR │ test  │ super feature
 19/03/30 11:53 │ Thomas Auffredou │ initial │       │       │ commit  
```

