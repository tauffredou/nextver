#Release pattern

The following keyworks are supported
- `SEMVER`: use semantic versioning (ex 1.0.5)
- `DATE`: use timestamping for rolling versioning. The opiniated format is YYYY-MM-DD-HHmmss

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
$ nextver get next-version --repo=github.com/tauffredou/test-semver
v1.0.0
```

Getting details

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

#### Using date release
```
$ nextver get next-version --repo=github.com/tauffredou/test-semver --pattern=myprefix-DATE
myprefix-2019-05-01-110159
```

The DATE pattern uses only time, ignoring the semantic versioning and the conventional commit convensions.

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

