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
			v = fmt.Sprintf("'%s'", v) // surround it with single quotes
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
// "build" or "install" arguments. If a -ldflags argument is already present, it
// normalizes the argument (converts [-ldflags, val] into [-ldflags=val]) and
// appends the given ldflags value.
func addLdFlags(args []string, ldflags string) ([]string, error) {
	if ldIdx := findArg(args, "-ldflags"); ldIdx != -1 { // -ldflag exists, normalize and append
		args = normalizeArg(args, "-ldflags")
		args[ldIdx] = appendToFlag(args[ldIdx], ldflags)
		return args, nil
	}

	// -ldflags argument does not exist in args.
	// find where to insert the new argument (after "build" or "install")
	insertIdx := findArg(args, "build")
	if insertIdx == -1 {
		insertIdx = findArg(args, "install")
	}
	if insertIdx == -1 {
		return nil, fmt.Errorf("cannot locate where to append -ldflags")
	}

	// allocate a new slice to prevent modifying the old one
	newArgs := make([]string, insertIdx+1, len(args)+2)
	copy(newArgs, args[:insertIdx+1])
	newArgs = append(newArgs, "-ldflags", ldflags)
	newArgs = append(newArgs, args[insertIdx+1:]...)
	return newArgs, nil
}

// appendToFlag appends val to -arg or -arg=... format. If the flag is missing a
// value, it adds a "=" to the flag before appending the value. If a value
// already exists, inserts a space character before appending the value.
func appendToFlag(arg, val string) string {
	if !strings.ContainsRune(arg, '=') {
		arg = arg + "="
	}
	if arg[len(arg)-1] != '=' && arg[len(arg)-1] != ' ' {
		arg += " "
	}
	arg += val
	return arg
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
func normalizeArg(args []string, arg string) []string {
	idx := -1
	for i, v := range args {
		if v == arg {
			idx = i
			break
		}
	}
	if idx == -1 || idx == len(args)-1 { // not found OR -arg has no succeding element
		return args
	}
	newArg := fmt.Sprintf("%s=%s", args[idx], args[idx+1]) // merge values
	args[idx] = newArg                                     // modify the arg
	return append(args[:idx+1], args[idx+2:]...)
}
