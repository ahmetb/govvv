package main

import "fmt"

var (
	// Version is populated at compile time by govvv from ./VERSION
	Version string

	// GitCommit is populated at compile time by govvv.
	GitCommit string

	// GitState is populated at compile time by govvv.
	GitState string
)

func versionString() string {
	if Version == "" {
		return "N/A"
	}
	return fmt.Sprintf("%s@%s-%s", Version, GitCommit, GitState)
}
