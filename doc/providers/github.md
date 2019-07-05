# Github provider 

Github provider uses the Github API to fetch data. You don't have to clone the repository locally.

It requires a **github token** to query the API

```
$ nextver get changelog --repo=github.com/tauffredou/test-semver 
Current release version : 0.0.0
Next release version    : v1.0.0

Changelog:
 Date           │ Author           │ Kind    │ Level │ Scope │ Title        
 ━━━━━━━━━━━━━━━┿━━━━━━━━━━━━━━━━━━┿━━━━━━━━━┿━━━━━━━┿━━━━━━━┿━━━━━━━━━━━━━━
 19/03/30 12:44 │ Thomas Auffredou │ chore   │ MAJOR │ test2 │ some change  
 19/03/30 12:29 │ Thomas Auffredou │ feat    │ MINOR │ test  │ super feature
 19/03/30 11:53 │ Thomas Auffredou │ initial │       │       │ commit       

```

## Authentication

Create a Github access token with at least the following scopes:

![github_scope](../images/github_scopes.png)


Then use this token using the option
```
nextver --github-token=xxxxxxxxx ...
```
or exporting the environment variable
```
export GITHUB_TOKEN=xxxxxxxxx
```

Nextver will also use the [hub](https://github.com/github/hub) configuration file if it exists (default location: ~/.config/hub). 

Resolution order is as follow: *parameter* > *environment variable* > *configuration file*

Parameter will take priority 
