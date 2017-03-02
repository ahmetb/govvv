package main

import (
	"fmt"
	"github.com/ahmetalpbalkan/govvv/integration-test/app-different-package/mypkg"
)

func main() {
	fmt.Printf("Version=%s\n", mypkg.Version)
	fmt.Printf("BuildDate=%s\n", mypkg.BuildDate)
	fmt.Printf("GitCommit=%s\n", mypkg.GitCommit)
	fmt.Printf("GitBranch=%s\n", mypkg.GitBranch)
	fmt.Printf("GitState=%s\n", mypkg.GitState)
	fmt.Printf("GitSummary=%s\n", mypkg.GitSummary)
}
