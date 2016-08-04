package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetValues_error(t *testing.T) {
	repo := newRepo(t)
	defer os.RemoveAll(repo.dir)

	_, err := GetFlags(repo.dir)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "failed to get commit")
}

func TestGetValues(t *testing.T) {
	// prepare the repo
	repo := newRepo(t)
	defer os.RemoveAll(repo.dir)
	mkCommit(t, repo, "commit 1")
	mkCommit(t, repo, "commit 2")

	// read the flags
	fl, err := GetFlags(repo.dir)
	require.Nil(t, err)

	// validate the flags
	require.Regexp(t, "^[0-9]{4}(-[0-9]{2}){2}T([0-9]{2}:){2}[0-9]{2}Z$", fl["main.BuildDate"])
	require.Regexp(t, "^[0-9a-f]{7}$", fl["main.GitCommit"])
	require.Equal(t, "master", fl["main.GitBranch"])
	require.Equal(t, "clean", fl["main.GitState"])
}
