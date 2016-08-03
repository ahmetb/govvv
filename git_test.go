package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGitExists(t *testing.T) {
	_, err := exec.LookPath("git")
	require.Nil(t, err, "error: %+v", err)
}

func TestCommit_failsOnEmptyRepo(t *testing.T) {
	repo := newRepo(t)
	defer os.RemoveAll(repo.dir)

	_, err := repo.Commit()
	require.NotNil(t, err)
}

func TestCommit(t *testing.T) {
	repo := newRepo(t)
	defer os.RemoveAll(repo.dir)

	mkCommit(t, repo, "commit 1")
	c1, err := repo.Commit()
	require.Nil(t, err)
	require.NotEmpty(t, c1)

	mkCommit(t, repo, "commit 2")
	c2, err := repo.Commit()
	require.Nil(t, err)
	require.NotEmpty(t, c2)

	// commit hash changed
	require.NotEqual(t, c1, c2)
}

func TestState(t *testing.T) {
	repo := newRepo(t)
	defer os.RemoveAll(repo.dir)

	s1, err := repo.State()
	require.Nil(t, err)
	require.EqualValues(t, "clean", s1)

	f, err := ioutil.TempFile(repo.dir, "") // contaminate
	require.Nil(t, err, "failed to create test file")
	f.Close()

	s2, err := repo.State()
	require.Nil(t, err)
	require.EqualValues(t, "dirty", s2)

	require.Nil(t, os.Remove(f.Name()), "failed to rm test file")

	s3, err := repo.State()
	require.Nil(t, err)
	require.EqualValues(t, "clean", s3)
}

func TestBranch(t *testing.T) {
	repo := newRepo(t)
	defer os.RemoveAll(repo.dir)

	// default branch "master"
	mkCommit(t, repo, "commit 1")
	mkCommit(t, repo, "commit 2")
	require.EqualValues(t, "master", repo.Branch())

	// move into detached state: "HEAD"
	_, err := repo.exec("checkout", "HEAD~1")
	require.Nil(t, err)
	require.EqualValues(t, "HEAD", repo.Branch())

	// checkout into another branch
	_, err = repo.exec("checkout", "-b", "foo")
	require.Nil(t, err)
	require.EqualValues(t, "foo", repo.Branch())
}

// Test utilities

func newRepo(t *testing.T) git {
	dir, err := ioutil.TempDir("", "gitrepo")
	require.Nil(t, err, "failed to create test dir")

	repo := git{dir}
	_, err = repo.exec("init", "-q", dir)
	require.Nil(t, err, "failed to initialize git repo")
	return repo
}

func mkCommit(t *testing.T, repo git, msg string) {
	_, err := repo.exec("commit", "--allow-empty", "--message", msg)
	require.Nil(t, err, "failed to commit: %+v", err)
}
