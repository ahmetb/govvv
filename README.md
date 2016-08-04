# govvv

The simple Go binary versioning tool that wraps the `go build` command. Stop
worrying about `-ldflags` and use `govvv`:

    go get github.com/ahmetalpbalkan/govvv
    govvv build .

## Build Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `main.GitCommit` | short commit hash of source tree | `0b5ed7a` |
| `main.GitBranch` | current branch name the code is built off | `master` |
| `main.GitState` | whether there are uncommitted changes | `clean` or `dirty` | 
| `main.BuildDate` | RFC3339 formatted UTC date | `2016-08-04T18:07:54Z` |
| `main.Version` | contents of `./VERSION` file, if exists | `2.0.0` |

## Installing govvv is easy

Download `govvv` via `go get` or from Releases tab (recommended).

    go get github.com/ahmetalpbalkan/govvv

## Using govvv is easy

Just add the Build Variables described above to your main package and build your
code with `govvv`.

```go
package main

import "fmt"

var (
	GitCommit, GitState string
)

func version() string { return GitCommit + "-" + GitState }

func main() {
	fmt.Printf("running build %s", version())
}
```

## govvv lets you specify custom `-ldflags`

Your existing `-ldflags` argument will still be preserved:

    govvv build -ldflags "-X main.BuildNumber=$buildnum" myapp

and the `-ldflags` constructed by govvv will be appended to your flag.


------

[![Build Status](https://travis-ci.org/ahmetalpbalkan/govvv.svg?branch=master)](https://travis-ci.org/ahmetalpbalkan/govvv)