package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetFlags_error(t *testing.T) {
	repo := newRepo(t)
	defer os.RemoveAll(repo.dir)

	_, err := GetFlags(repo.dir, []string{})
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "failed to get commit")
}

func Test_date(t *testing.T) {
	v := date()
	require.Regexp(t, "^[0-9]{4}(-[0-9]{2}){2}T([0-9]{2}:){2}[0-9]{2}Z$", v)
}

func TestGetFlags(t *testing.T) {
	// prepare the repo
	repo := newRepo(t)
	defer os.RemoveAll(repo.dir)
	mkCommit(t, repo, "commit 1")
	mkCommit(t, repo, "commit 2")

	// read the flags
	fl, err := GetFlags(repo.dir, []string{})
	require.Nil(t, err)

	// validate the flags
	require.Regexp(t, "^[0-9]{4}(-[0-9]{2}){2}T([0-9]{2}:){2}[0-9]{2}Z$", fl["main.BuildDate"])
	require.Regexp(t, "^[0-9a-f]{4,15}$", fl["main.GitCommit"])
	require.Equal(t, "master", fl["main.GitBranch"])
	require.Equal(t, "clean", fl["main.GitState"])
	require.Equal(t, fl["main.GitCommit"], fl["main.GitSummary"])
}

func TestGetFlags_versionDefault(t *testing.T) {
	// prepare the repo
	repo := newRepo(t)
	defer os.RemoveAll(repo.dir)
	mkCommit(t, repo, "commit 1")

	// there is no main.Version flag
	fl, err := GetFlags(repo.dir, []string{})
	require.Nil(t, err)
	require.NotContains(t, fl, "main.Version")
}

func TestGetFlags_justVersionFlag(t *testing.T) {
	// prepare the repo
	repo := newRepo(t)
	defer os.RemoveAll(repo.dir)
	mkCommit(t, repo, "commit 1")

	// -version is specified and there is no VERSION file
	fl, err := GetFlags(repo.dir, []string{flVersion, "2.0.0-RC01"})
	require.Nil(t, err)
	require.Equal(t, "2.0.0-RC01", fl["main.Version"])
}

func TestGetFlags_versionFile(t *testing.T) {
	// prepare the repo
	repo := newRepo(t)
	defer os.RemoveAll(repo.dir)
	mkCommit(t, repo, "commit 1")

	// add version file and get the value back
	require.Nil(t, ioutil.WriteFile(filepath.Join(repo.dir, "VERSION"), []byte("2.0.0-beta\n"), 0600))
	fl, err := GetFlags(repo.dir, []string{})
	require.Nil(t, err)
	require.Equal(t, "2.0.0-beta", fl["main.Version"])
}

func TestGetFlags_versionFlagOverrides(t *testing.T) {
	// prepare the repo
	repo := newRepo(t)
	defer os.RemoveAll(repo.dir)
	mkCommit(t, repo, "commit 1")

	// add version file
	require.Nil(t, ioutil.WriteFile(filepath.Join(repo.dir, "VERSION"), []byte("2.0.0-beta\n"), 0600))

	// -version is specified and there is na VERSION file (flag takes precedence)
	fl, err := GetFlags(repo.dir, []string{flVersion, "2.0.0-RC01"})
	require.Nil(t, err)
	require.Equal(t, "2.0.0-RC01", fl["main.Version"])
}

func TestGetFlags_versionFileError(t *testing.T) {
	// prepare the repo
	repo := newRepo(t)
	defer os.RemoveAll(repo.dir)
	mkCommit(t, repo, "commit 1")

	// add version file and get the value back
	require.Nil(t, ioutil.WriteFile(filepath.Join(repo.dir, "VERSION"), []byte("2.0.0-beta\n"), 0000))
	fl, err := GetFlags(repo.dir, []string{})
	require.Nil(t, fl)
	require.Contains(t, err.Error(), "failed to read version file")
}

func TestGetFlags_pkgFlag(t *testing.T) {
	// prepare the repo
	repo := newRepo(t)
	defer os.RemoveAll(repo.dir)
	mkCommit(t, repo, "commit 1")
	mkCommit(t, repo, "commit 2")

	// read the flags for custom package
	pkg := "github.com/acct/coolproject/version"
	fl, err := GetFlags(repo.dir, []string{flPackage, pkg})
	require.Nil(t, err)

	// validate the flags
	require.Contains(t, fl, pkg+".BuildDate")
	require.Contains(t, fl, pkg+".GitCommit")
	require.Contains(t, fl, pkg+".GitBranch")
	require.Contains(t, fl, pkg+".GitState")
	require.Contains(t, fl, pkg+".GitSummary")
}

func Test_versionFromFile_notFound(t *testing.T) {
	dir := tmpDir(t)
	defer os.RemoveAll(dir)

	v, err := versionFromFile(dir)
	require.Equal(t, "", v)
	require.Nil(t, err)
}

func Test_versionFromFile_error(t *testing.T) {
	dir := tmpDir(t)
	defer os.RemoveAll(dir)

	require.Nil(t, ioutil.WriteFile(filepath.Join(dir, "VERSION"), nil, 0200)) // // no read perms

	_, err := versionFromFile(dir)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "failed to read version file")
}

func Test_versionFromFile(t *testing.T) {
	dir := tmpDir(t)
	defer os.RemoveAll(dir)

	require.Nil(t, ioutil.WriteFile(filepath.Join(dir, "VERSION"), []byte("\t 0.6.0.0 \n "), 0600)) // // no read perms

	v, err := versionFromFile(dir)
	require.Nil(t, err)
	require.Equal(t, "0.6.0.0", v)
}

// Test utilities

func tmpDir(t *testing.T) string {
	dir, err := ioutil.TempDir("", "")
	require.Nil(t, err, "failed to create test directory")
	return dir
}
