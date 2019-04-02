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
usage: nextver --github-token=GITHUB-TOKEN --github-owner=GITHUB-OWNER --github-repo=GITHUB-REPO [<flags>]

Flags:
      --help                     Show context-sensitive help (also try --help-long and --help-man).
      --github-token=GITHUB-TOKEN  
                                 Github token
      --github-owner=GITHUB-OWNER  
                                 Github owner
      --color                    Colorize output
      --github-repo=GITHUB-REPO  Github repo
  -p, --pattern="vSEMVER"        Versionning pattern
  -o, --output="console"         Output format (console, json, yaml)
  -b, --branch=BRANCH            Target branch (default branch if empty)


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

