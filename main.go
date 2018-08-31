package main

import (
	"log"
	"os"

	"github.com/JonathonGore/api-check/cli"
)

func main() {
	app := cli.ConfigureCLI()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
