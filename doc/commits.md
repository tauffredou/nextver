# Commit messages

Nextver follows the *Conventional commit* specification : [https://www.conventionalcommits.org/](https://www.conventionalcommits.org/)

Examples:

commit message
```
feat(holidays): add towel to bag (#123)

A large towel can be used to lie on the sand AND to dry after swimming
```

would be analysed as follow:
```yaml
kind: feat
scope: holidays
title: add towel to bag (#123)
body: A large towel can be used to lie on the sand AND to dry after swimming

```

it would also be considered as a *MINOR* change.
