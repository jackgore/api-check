package main

import (
	"fmt"

	"github.com/JonathonGore/api-check/parser"
	"github.com/JonathonGore/api-check/runner"
)

func main() {
	fmt.Printf("Running go api-check\n")

	tests, err := parser.Parse("api-check.json")
	if err != nil {
		fmt.Printf("Unable to parse tests: %v", err)
		return
	}

	report := runner.RunTests(tests)

	fmt.Printf("Run report: %+v\n", report)
}
