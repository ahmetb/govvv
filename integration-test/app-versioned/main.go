package main

import "fmt"

var (
	// These fields are populated by govvv
	BuildDate string
	GitCommit string
	GitBranch string
	GitState  string
)

func main() {
	fmt.Printf("BuildDate=%s\n", BuildDate)
	fmt.Printf("GitCommit=%s\n", GitCommit)
	fmt.Printf("GitBranch=%s\n", GitBranch)
	fmt.Printf("GitState=%s\n", GitState)
}
