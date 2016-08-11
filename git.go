package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type git struct {
	dir string
}

func (g git) exec(args ...string) (string, error) {
	var errOut bytes.Buffer
	c := exec.Command("git", args...)
	c.Dir = g.dir
	c.Stderr = &errOut
	out, err := c.Output()
	outStr := strings.TrimSpace(string(out))
	if err != nil {
		err = fmt.Errorf("git: error=%q stderr=%s", err, string(errOut.Bytes()))
	}
	return outStr, err
}

// Commit returns the short git commit hash.
func (g git) Commit() (string, error) {
	return g.exec("rev-parse", "--short", "HEAD")
}

// State returns the repository state indicating whether
// it is "clean" or "dirty".
func (g git) State() (string, error) {
	out, err := g.exec("status", "--porcelain")
	if err != nil {
		return "", err
	}
	if len(out) > 0 {
		return "dirty", nil
	}
	return "clean", nil
}

// Branch returns the branch name. If it is detached,
// or an error occurs, returns "HEAD".
func (g git) Branch() string {
	out, err := g.exec("symbolic-ref", "-q", "--short", "HEAD")
	if err != nil {
		// might be failed due to another reason, but assume it's
		// exit code 1 from `git symbolic-ref` in detached state.
		return "HEAD"
	}
	return out
}

// Summary returns the output of "git describe --tags --dirty --always".
func (g git) Summary() (string, error) {
	out, err := g.exec("describe", "--tags", "--dirty", "--always")
	if err != nil {
		return "", err
	}
	return out, err
}
