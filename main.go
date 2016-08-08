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

func main() {
	args := os.Args
	if len(args) < 2 {
		log.Println(`govvv: not enough arguments (try "govvv build .")`)
		log.Printf("version: %s", versionString())
		os.Exit(1)
	} else if args[1] != "build" && args[1] != "install" && args[1] != "list" {
		// do not wrap the entire 'go tool'
		// "list" is wrapped to be compatible with mitchellh/gox.
		log.Fatalf(`govvv: only works with "build", "install" and "list". try "go %s" instead`, args[1])
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("govvv: cannot get working directory: %v", err)
	}
	args, err = prepArgs(wd, args[1:])
	if err != nil {
		log.Fatalf("govvv: failed to prepare args: %v", err)
	}
	if findArg(args, flDryRun) != -1 {
		fmt.Println(goToolDryRunCmd(args))
	} else {
		if err := execGoTool(args); err != nil {
			log.Fatalf("go tool: %v", err)
		}
	}

}

// prepArgs prepares the arguments with correct ldflags if
// args[0] is "build" or "install".
func prepArgs(dir string, args []string) ([]string, error) {
	if len(args) == 0 || (args[0] != "build" && args[0] != "install") {
		return args, nil
	}
	vals, err := GetFlags(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to collect values: %v", err)
	}
	ldflags, err := mkLdFlags(vals)
	if err != nil {
		return nil, fmt.Errorf("failed to compile values: %v", err)
	}

	args, err = addLdFlags(args, ldflags)
	if err != nil {
		return nil, fmt.Errorf("cannot add ldflags: %v", err)
	}
	return args, nil
}

const (
	flDryRun = "-print"
)

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
	for _, v := range args {
		if v == flDryRun {
			continue
		}
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
