package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_mkldFlags_fails(t *testing.T) {
	_, err := mkLdFlags(map[string]string{"key space": "val"})
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "key contains whitespaces")

	_, err = mkLdFlags(map[string]string{"key": "val space"})
	require.Nil(t, err, "values can have spaces")
}

func Test_appendToFlag(t *testing.T) {
	v := "VALUE"
	cases := []struct{ in, out string }{
		{"-arg", "-arg=VALUE"},
		{"-arg=", "-arg=VALUE"},
		{"-arg=OLD", "-arg=OLD VALUE"},
		{"-arg=OLD ", "-arg=OLD VALUE"},
	}
	for _, c := range cases {
		out := appendToFlag(c.in, v)
		require.Equal(t, c.out, out, "input=%q", cases)
	}
}

func Test_addLdFlags(t *testing.T) {
	type testcase struct {
		in  []string
		out []string
	}
	validateCases := func(tc []testcase, ldflagsVal string) {
		for _, c := range tc {
			out, err := addLdFlags(c.in, ldflagsVal)
			require.Nil(t, err)
			require.Equal(t, c.out, out, "input args=%#v", c.in)
		}
	}

	{ // cannot find where to append ldflags
		_, err := addLdFlags([]string{"a", "b", "c"}, "NEW VALUE")
		require.NotNil(t, err)
		require.EqualError(t, err, "cannot locate where to append -ldflags")
	}
	{
		// modifies existing -ldflags
		val := "NEW"
		cases := []testcase{
			{
				[]string{"build", "-ldflags"},
				[]string{"build", "-ldflags=NEW"},
			},
			{
				[]string{"build", "-ldflags", "OLD"},
				[]string{"build", "-ldflags=OLD NEW"},
			},
			{
				[]string{"-v", "build", "-ldflags=OLD"},
				[]string{"-v", "build", "-ldflags=OLD NEW"},
			},
		}
		validateCases(cases, val)
	}
	{ // adds it after "build"
		val := "NEW VALUE"
		cases := []testcase{
			{
				[]string{"build"},
				[]string{"build", "-ldflags", val},
			},
			{
				[]string{"build", "-v"},
				[]string{"build", "-ldflags", val, "-v"},
			},
			{
				[]string{"-v", "build", "."},
				[]string{"-v", "build", "-ldflags", val, "."},
			},
			{
				[]string{"build", "-aflag", "-v", "."},
				[]string{"build", "-ldflags", val, "-aflag", "-v", "."},
			},
		}
		validateCases(cases, val)
	}
}

func Test_mkldFlags(t *testing.T) {
	{ // empty
		out, err := mkLdFlags(map[string]string{})
		require.Nil(t, err)
		require.Empty(t, out)
	}
	{ // normal input
		out, err := mkLdFlags(map[string]string{
			"key1": "val1",
			"key2": "val 2",
		})
		require.Nil(t, err)
		expected := []string{
			"-X key1=val1 -X key2='val 2'",
			"-X key2='val 2' -X key1=val1"}

		if out != expected[0] && out != expected[1] {
			t.Fatalf("output: %q, expected: either %q --or-- %q", out, expected[0], expected[1])
		}
	}
}

func Test_normalizeArg(t *testing.T) {
	cases := []struct {
		arg     string
		in, out []string
	}{
		{ // arg has no value
			"-foo",
			[]string{"a", "b", "-foo"},
			[]string{"a", "b", "-foo"},
		},
		{ // normalize at the end
			"-key1",
			[]string{"a", "b", "-key1", "value"},
			[]string{"a", "b", "-key1=value"},
		},
		{ // normalize at the beginning
			"-key2",
			[]string{"-key2", "value", "a", "b"},
			[]string{"-key2=value", "a", "b"},
		},
		{ // already in desired format
			"-key3",
			[]string{"a", "b", "-key3=value"},
			[]string{"a", "b", "-key3=value"},
		},
	}
	for _, c := range cases {
		out := normalizeArg(c.in, c.arg)
		require.EqualValues(t, c.out, out, "arg=%q input: %#v", c.arg, c.in)
	}
}

func Test_findArg(t *testing.T) {
	cases := []struct {
		in  []string
		key string
		out int
	}{
		{[]string{"foo", "bar", "quux"}, "none", -1},
		{[]string{"foo", "bar", "quux"}, "bar", 1},
		{[]string{"foo", "bar=val", "quux"}, "bar", 1},
		{[]string{"foo", "bar=val", "quux"}, "bar", 1},
		{[]string{"-arg1=bar", "-arg2"}, "-arg1", 0},
		{[]string{"-arg1=foo", "-arg2=foo"}, "-arg2", 1},
		{[]string{"-foo", "--bar"}, "-bar", -1},
	}
	for _, c := range cases {
		out := findArg(c.in, c.key)
		require.EqualValues(t, c.out, out, "key=%q args=%#v", c.key, c.in)
	}
}
