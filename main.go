package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func init() {
	log.SetFlags(0)
}

const (
	defaultPackage       = "main"
	flDryRun             = "-print"
	flDryRunPrintLdFlags = "-flags"
	flPackage            = "-pkg"
	flVersion            = "-version"
)

var (
	// govvvDirectives is mapping of govvv directives, which must be elided
	// when constructing the final go tool command, to a boolean which
	// indicates whether the directive takes an argument or not.
	govvvDirectives = map[string]bool{
		flDryRun:             false,
		flDryRunPrintLdFlags: false,
		flPackage:            true,
		flVersion:            true}
)

func main() {
	args := os.Args
	if len(args) < 2 {
		log.Println(`govvv: not enough arguments (try "govvv build .")`)
		log.Printf("version: %s", versionString())
		os.Exit(1)
	} else if args[1] != "build" && args[1] != "install" && args[1] != "list" && !isGovvvDirective(args[1]) {
		// do not wrap the entire 'go tool'
		// "list" is wrapped to be compatible with mitchellh/gox.
		log.Fatalf(`govvv: only works with "build", "install" and "list". try "go %s" instead`, args[1])
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("govvv: cannot get working directory: %v", err)
	}

	versionValues, err := GetFlags(wd, args)
	if err != nil {
		log.Fatalf("failed to collect values: %v", err)
	}

	ldflags, err := mkLdFlags(versionValues)
	if err != nil {
		log.Fatalf("failed to compile values: %v", err)
	}

	if _, ok := collectGovvvDirective(args, flDryRunPrintLdFlags); ok {
		fmt.Print(ldflags)
		return
	}

	args = args[1:] // rm executable name

	if args[0] == "build" || args[0] == "install" {
		args, err = addLdFlags(args, ldflags)
		if err != nil {
			log.Fatalf("failed to add ldflags to args: %v", err)
		}
	}

	if _, ok := collectGovvvDirective(args, flDryRun); ok {
		fmt.Println(goToolDryRunCmd(args))
		return
	}

	args = scrubGovvvDirectives(args)

	if err := execGoTool(args); err != nil {
		log.Fatalf("go tool: %v", err)
	}
}

// execGoTool invokes "go" with given arguments and passes the current
// process' standard streams.
func execGoTool(args []string) error {
	cmd := exec.Command("go", args...)
	cmd.Stdout, cmd.Stderr, cmd.Stdin = os.Stdout, os.Stderr, os.Stdin
	return cmd.Run()
}

// goToolDryRunCmd returns a POSIX shell-compatible command that would normally
// get executed. Not guaranteed to quote and escape the args very well.
func goToolDryRunCmd(args []string) string {
	var b bytes.Buffer
	b.WriteString("go")
	b.WriteRune(' ')
	printed := false
	for _, v := range scrubGovvvDirectives(args) {
		if printed {
			b.WriteString(" \\\n")
			b.WriteString("\t")
		}

		if strings.ContainsAny(v, " \"'\n\t") {
			v = strconv.QuoteToASCII(v)
		}
		b.WriteString(v)
		printed = true

	}
	return b.String()
}

// isGovvvDirective returns true if the arg is a govvv directive, and false
// otherwise.
func isGovvvDirective(arg string) bool {
	_, ok := govvvDirectives[arg]
	return ok
}

// scrubGovvvDirectives filters out govvv directs to return a clean set of args
// that can be passed to the go command.
func scrubGovvvDirectives(args []string) (filtered []string) {
	filtered = []string{}
	skipping := 0
	for _, arg := range args {
		if skipping > 0 {
			skipping--
			continue
		}
		if hasArgument, ok := govvvDirectives[arg]; ok {
			if hasArgument {
				skipping = 1
			}
		} else {
			filtered = append(filtered, arg)
		}
	}
	return filtered
}

// collectGovvvDirective searches the args array for a directive and a possible
// argument for that directive.  It returns that argument, or "", if the
// directive takes none, and an indication if the directive was actually found.
func collectGovvvDirective(args []string, directive string) (argument string, ok bool) {
	for i, arg := range args {
		if directive == arg {
			hasArgument, _ := govvvDirectives[arg]
			if !hasArgument {
				return "", true
			} else if i+1 < len(args) {
				return args[i+1], true
			} else {
				break
			}
		}
	}
	return "", false
}
