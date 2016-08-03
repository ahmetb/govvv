package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_mkldFlags(t *testing.T) {
	cases := []struct {
		in  map[string]string
		out string
	}{
		{map[string]string{}, ""},
		{map[string]string{"a": "b"}, `a="b"`},
		{map[string]string{"main.Version": "b"}, `main.Version="b"`},
		{map[string]string{"foo": "bar quux"}, `foo="bar quux"`},
	}
	for _, c := range cases {
		out := mkLdFlags(c.in)
		require.EqualValues(t, c.out, out, "input: %#v", c.in)
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
			[]string{"a", "b", `-key1="value"`},
		},
		{ // normalize at the beginning
			"-key2",
			[]string{"-key2", "value", "a", "b"},
			[]string{`-key2="value"`, "a", "b"},
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
		{[]string{"foo", "bar=", "quux"}, "bar", 1},
		{[]string{"-arg1", "-arg2", "-arg3"}, "-arg2", -1},
		{[]string{"-arg1=bar", "-arg2"}, "-arg1", 0},
		{[]string{"-arg1=foo", "-arg2=foo"}, "-arg2", 1},
		{[]string{"-arg1", "-arg2=foo", "-arg3"}, "-arg4", -1},
	}
	for _, c := range cases {
		out := findArg(c.in, c.key)
		require.EqualValues(t, c.out, out, "key=%q args=%#v", c.key, c.in)
	}
}
