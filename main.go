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
		log.Fatal(`govvv: not enough arguments (example: "govvv build .")`)
	} else if args[1] != "build" {
		log.Fatal(`govvv: can be used only with subcommand "build"`)
	}

	cmd := exec.Command("go", args[1:]...)
	if err := cmd.Run(); err != nil {
		os.Exit(1) // TODO get exitcode from cmd (there must be a cross-platform pkg for this?)
	}
}
