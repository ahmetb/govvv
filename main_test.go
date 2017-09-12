package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIsGovvvDirective(t *testing.T) {
	for directive := range govvvDirectives {
		require.True(t, isGovvvDirective(directive))
	}

	require.False(t, isGovvvDirective("-o"))
}

func TestScrubGovvvDirectives(t *testing.T) {
	require.Equal(t, []string{}, []string{}, "scrubGovvvDirectives should be fine with an empty arg array")

	require.Equal(t, []string{"build", "-o", "a.out"}, scrubGovvvDirectives([]string{"build", "-o", "a.out"}),
		"scrubGovvvDirectives should not touch normal go args")

	require.Equal(t, []string{}, scrubGovvvDirectives([]string{"-flags"}),
		"scrubGovvvDirectives should scrub -flags")

	require.Equal(t, []string{"build", "-o", "a.out"}, scrubGovvvDirectives([]string{"build", "-o", "a.out", "-print"}),
		"scrubGovvvDirectives should scrub -print")

	require.Equal(t, []string{"build", "-o", "a.out"},
		scrubGovvvDirectives([]string{"build", "-o", "a.out", "-pkg", "github.com/ahmetb/govvv"}),
		"scrubGovvvDirectives should scrub -pkg and its argument")
}

func TestCollectGovvvDirective(t *testing.T) {
	argument, ok := collectGovvvDirective([]string{}, flDryRun)
	require.False(t, ok)
	require.Equal(t, "", argument, "collectGovvvDirective should be fine with an empty arg array")

	argument, ok = collectGovvvDirective([]string{"build", "-o", "a.out", "-print"}, flDryRun)
	require.True(t, ok)
	require.Equal(t, "", argument, "collectGovvvDirective should find -print")

	argument, ok = collectGovvvDirective([]string{"-flags"}, flDryRunPrintLdFlags)
	require.True(t, ok)
	require.Equal(t, "", argument, "collectGovvvDirective should find -flags")

	argument, ok = collectGovvvDirective([]string{"build", "-o", "a.out", "-pkg", "github.com/ahmetb/govvv"}, flPackage)
	require.True(t, ok)
	require.Equal(t, "github.com/ahmetb/govvv", argument, "collectGovvvDirective should find -pkg")

	argument, ok = collectGovvvDirective([]string{"build", "-o", "a.out", "-pkg"}, flPackage)
	require.False(t, ok)
	require.Equal(t, "", argument, "collectGovvvDirective should catch missing argument for -pkg")
}
