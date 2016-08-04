package main

import "fmt"

var (
	// GitCommit is provided by govvv
	GitCommit string

	// MyVariable is provided by user via -ldflags
	MyVariable string
)

func main() {
	fmt.Printf("MyVariable=%s\n", MyVariable)
	fmt.Printf("GitCommit=%s\n", GitCommit)
}
