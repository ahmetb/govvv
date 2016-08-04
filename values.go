package main

import (
	"fmt"
	"time"
)

// GetFlags collects data to be passed as ldflags.
func GetFlags(dir string) (map[string]string, error) {
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
		"main.BuildDate": date(),
		"main.GitCommit": gitCommit,
		"main.GitBranch": gitBranch,
		"main.GitState":  gitState,
	}, nil
}

// date returns the UTC date formatted in RFC 3339 layout.
func date() string {
	return time.Now().UTC().Format(time.RFC3339)
}
