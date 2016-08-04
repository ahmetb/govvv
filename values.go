package main

import (
	"fmt"
	"time"
)

// GetFlags collects data to be passed as ldflags.
func GetFlags(dir string) (map[string]string, error) {
	date := BuildDate()
	repo := git{dir}

	gitBranch := repo.Branch()
	gitCommit, err := repo.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to get commit: %v", err)
	}
	gitState, err := repo.State()
	if err != nil {
		return nil, fmt.Errorf("failed to get repository state: %v", err)
	}

	return map[string]string{
		"BuildDate": date,
		"GitCommit": gitCommit,
		"GitBranch": gitBranch,
		"GitState":  gitState,
	}, nil
}

// BuildDate returns the UTC date formatted in RFC 3339 layout.
func BuildDate() string {
	return time.Now().UTC().Format(time.RFC3339)
}
