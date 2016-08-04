package main

import (
	"bytes"
	"fmt"
	"strings"
)

// mkLdFlags will generate a string compatible to use in "go build --ldflags"
// with provided values.
func mkLdFlags(values map[string]string) (string, error) {
	var b bytes.Buffer
	var i int
	for k, v := range values {
		if len(strings.Fields(k)) > 1 {
			return "", fmt.Errorf("cannot make ldflags for %q: key contains whitespaces", k)
		}
		if len(strings.Fields(v)) > 1 {
			return "", fmt.Errorf("cannot make ldflags for %q: value contains whitespaces", k)
		}

		i++
		b.WriteString(fmt.Sprintf("-X %s=%s", k, v))
		if i != len(values) {
			b.WriteByte(' ')
		}
	}
	return b.String(), nil
}

// addLdFlags appends the specified ldflags value to args right after the
// "build" argument. If a -ldflags argument is present, it returns error.
func addLdFlags(args []string, ldflags string) ([]string, error) {
	if findArg(args, "-ldflags") != -1 {
		return nil, fmt.Errorf("already have a ldflags flag")
	}
	buildIdx := findArg(args, "build")
	if buildIdx == -1 {
		return nil, fmt.Errorf("cannot locate where to append -ldflags")
	}

	// allocate a new slice to prevent modifying the old one
	newArgs := make([]string, buildIdx+1, len(args)+2)
	copy(newArgs, args[:buildIdx+1])
	newArgs = append(newArgs, "-ldflags", ldflags)
	newArgs = append(newArgs, args[buildIdx+1:]...)
	return newArgs, nil
}

// findArgs looks for 'arg' or 'arg=...' values in args and returns its index or
// -1 if not found.
func findArg(args []string, arg string) int {
	key := arg + "="
	for i, v := range args {
		if v == arg || strings.HasPrefix(v, key) {
			return i
		}
	}
	return -1
}

// normalize finds the -arg in the args and concats its value to the same
// argument. e.g. [-arg, foo] will be converted to [-arg="foo"] only once.
//
// TODO(ahmetb) this might be used in the future to modify existing -ldflags
// value perhaps.
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
