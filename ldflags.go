package main

import (
	"bytes"
	"fmt"
	"strings"
)

// mkLdFlags will generate a string compatible to use in "go build --ldflags"
// with provided values.
func mkLdFlags(values map[string]string) string {
	var b bytes.Buffer
	var i int
	for k, v := range values {
		i++
		b.WriteString(fmt.Sprintf("%s=%q", k, v))
		if i != len(values) {
			b.WriteByte(' ')
		}
	}
	return b.String()
}

// addLdFlags appends the specified ldflags value to args. If a -ldflags= argument
// is present, it appends the given string to its value, if not, inserts it to
// the slice after the "build" argument and returns it.
func addLdFlags(args []string, ldflags string) []string {
	return nil
}

// findArgs looks for 'arg=...' values in args and returns its index
// or -1 if not found.
func findArg(args []string, arg string) int {
	key := arg + "="
	for i, v := range args {
		if strings.HasPrefix(v, key) {
			return i
		}
	}
	return -1
}

// normalize finds the -arg in the args and concats its value to the same
// argument. e.g. [-arg, foo] will be converted to [-arg="foo"] only once.
func normalizeArg(args []string, arg string) []string {
	for i, v := range args {
		if v == arg {
			// concat the next index
			if i == len(args)-1 { // flag has no value, return
				return args
			}
			val := fmt.Sprintf("%s=%q", arg, args[i+1])
			return append(append(args[:i], val), args[i+2:]...)
		}
	}
	return args
}
