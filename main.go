package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

var (
	A string
	B string
)

func init() {
	log.SetFlags(0)
}

func main() {
	fmt.Println("A=", A, "B=", B)
	args := os.Args
	if len(args) < 2 {
		log.Fatal(`govvv: not enough arguments (try "govvv build .")`)
	} else if args[1] != "build" {
		log.Fatalf(`govvv: try "go %s" instead`, args[1])
	}

	cmd := exec.Command("go", args[1:]...)
	cmd.Stdout, cmd.Stderr, cmd.Stdin = os.Stdout, os.Stderr, os.Stdin

	if err := cmd.Run(); err != nil {
		os.Exit(1) // TODO get exitcode from cmd (there must be a cross-platform pkg for this?)
	}
}
