[![CircleCI](https://circleci.com/gh/tauffredou/nextver.svg?style=svg)](https://circleci.com/gh/tauffredou/nextver)

Nextver embrace best practices from both commit log and versioning. 
* leverages commit messages compliant with the ["conventional commits" specification](https://www.conventionalcommits.org/)
* calculates next version based on the git history

# Quick start
prerequisite: go

Install
```bash
go get github.com/tauffredou/nextver
```

Run 
```bash
nextver -r path/to/git/repo get changelog
```

# Documentation

- [configuration](doc/configuration.md) 
- Providers
  - [github](doc/providers/github.md)
  - [git](doc/providers/git.md)
- [versioning](doc/versioning.md) 
