package main

import (
	"fmt"

	"github.com/JonathonGore/api-check/config"
	"github.com/JonathonGore/api-check/parser"
	"github.com/JonathonGore/api-check/printer"
	"github.com/JonathonGore/api-check/runner"
)

func main() {
	fmt.Printf("Running go api-check\n\n")

	conf, err := config.New(".acconf.json")
	if err != nil {
		fmt.Printf("Unable to parse config object: %v", err)
		return
	}

	p := parser.New(conf)

	tests, err := p.Parse("api-check.json")
	if err != nil {
		fmt.Printf("Unable to parse tests: %v", err)
		return
	}

	reports := runner.RunTests(tests)

	printer.PrintReports(reports)
}
