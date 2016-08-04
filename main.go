package main

import (
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
	} else if args[1] != "build" {
		log.Fatalf(`govvv: try "go %s" instead`, args[1])
	}

	args = args[1:]                   // strip the executable name
	args, err := addLdFlags(args, "") // add ldflags
	if err != nil {
		log.Fatal(err)
	}
	if err := execGoTool(args); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

// execGoTool invokes "go" with given arguments and passes the current
// process' standard streams.
func execGoTool(args []string) error {
	cmd := exec.Command("go", args...)
	cmd.Stdout, cmd.Stderr, cmd.Stdin = os.Stdout, os.Stderr, os.Stdin
	return cmd.Run()
}
