package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

const versionFile = "VERSION"

// GetFlags collects data to be passed as ldflags.
func GetFlags(dir, pkg string) (map[string]string, error) {
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
	gitSummary, err := repo.Summary()
	if err != nil {
		return nil, fmt.Errorf("failed to get repository summary: %v", err)
	}

	v := map[string]string{
		pkg + ".BuildDate":  date(),
		pkg + ".GitCommit":  gitCommit,
		pkg + ".GitBranch":  gitBranch,
		pkg + ".GitState":   gitState,
		pkg + ".GitSummary": gitSummary,
	}

	return v, nil
}

// date returns the UTC date formatted in RFC 3339 layout.
func date() string {
	return time.Now().UTC().Format(time.RFC3339)
}

// versionFromFile looks for a file named VERSION in dir if it exists and
// returns its contents by trimming the whitespace around it. If the file
// does not exist, it does not return any errors
func versionFromFile(dir string) (string, error) {
	fp := filepath.Join(dir, versionFile)
	b, err := ioutil.ReadFile(fp)
	if os.IsNotExist(err) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("failed to read version file %s: %v", fp, err)
	}
	return string(bytes.TrimSpace(b)), nil
}
