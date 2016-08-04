package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func init() {
	log.SetFlags(0)
}

func main() {
	args := os.Args
	if len(args) < 2 {
		log.Fatal(`govvv: not enough arguments (try "govvv build .")`)
	} else if args[1] != "build" && args[1] != "install" {
		log.Fatalf(`govvv: only works with "build" and "install". try "go %s" instead`, args[1])
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("govvv: cannot get working directory: %v", err)
	}
	if err := process(wd, args[1:]); err != nil {
		log.Fatal(err)
	}
}

// process calls go tool with provided args (should not contain)
// the name of the current process in args[0])
func process(dir string, args []string) error {
	vals, err := GetFlags(dir)
	if err != nil {
		return fmt.Errorf("govvv: failed to collect values: %v", err)
	}
	ldflags, err := mkLdFlags(vals)
	if err != nil {
		return fmt.Errorf("govvv: failed to compile values: %v", err)
	}

	args, err = addLdFlags(args, ldflags)
	if err != nil {
		return fmt.Errorf("cannot add ldflags: %v", err)
	}
	if err := execGoTool(args); err != nil {
		return fmt.Errorf("go tool failed: %v", err)
	}
	return nil
}

// execGoTool invokes "go" with given arguments and passes the current
// process' standard streams.
func execGoTool(args []string) error {
	cmd := exec.Command("go", args...)
	cmd.Stdout, cmd.Stderr, cmd.Stdin = os.Stdout, os.Stderr, os.Stdin
	return cmd.Run()
}
