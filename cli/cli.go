package cli

import (
	"fmt"

	"github.com/JonathonGore/api-check/builder"
	"github.com/JonathonGore/api-check/config"
	"github.com/JonathonGore/api-check/parser"
	"github.com/JonathonGore/api-check/suite"
	"github.com/urfave/cli"
)

const (
	confFile         = ".ac.json"
	defaultVerbosity = true
)

// commandNotFound is executed when the user tries to execute an errorneous command.
func commandNotFound(c *cli.Context, cmd string) {
	fmt.Printf("%v has no commmand name '%v'\n", c.App.Name, cmd)
}

// generateAction defines the action that is run by invoking `api-check generate <name>`
func generateAction(c *cli.Context) error {
	if c.NArg() != 1 {
		return cli.NewExitError(fmt.Sprintf("%v %v requires exactly 1 argument", c.App.Name, c.Command.Name), 1)
	}

	filename, err := builder.CreateSkeletonFile(c.Args().First())
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	fmt.Printf("successfully created template file %v\n", filename)

	return nil
}

// verifyAction defines the action that is run by incoking `api-check verify <filename>`
func verifyAction(c *cli.Context) error {
	if c.NArg() == 0 {
		return cli.NewExitError(fmt.Sprintf("%v %v requires at least 1 argument", c.App.Name, c.Command.Name), 1)
	}

	conf, err := config.New(config.DefaultConfigFile)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("error parsing config file: %v", err), 1)
	}

	p := parser.New(conf)

	if _, err := p.Parse([]string(c.Args())); err != nil {
		return cli.NewExitError(fmt.Sprintf("%v", err), 1)
	}

	fmt.Println("validated successfuly")
	return nil
}

// runAction defines the action that is run by invoking `api-check run`
func runAction(c *cli.Context) error {
	suite.Verbose(defaultVerbosity)
	suite.RunStandalone()
	return nil
}

// buildCLICommands builds the list of available commands for the cli.
func buildCLICommands() []cli.Command {
	return []cli.Command{
		{
			Name:   "run",
			Usage:  "runs your test suites",
			Action: runAction,
		},
		{
			Name:   "verify",
			Usage:  "verify an api-check file",
			Action: verifyAction,
		},
		{
			Name:    "generate",
			Aliases: []string{"gen"},
			Usage:   "generate a skeleton api-check file",
			Action:  generateAction,
		},
	}
}

// ConfigureCLI configures a urfave cli app to enable api-check as an cli tool.
func ConfigureCLI() *cli.App {
	app := cli.NewApp()
	app.Name = "api-check"
	app.Usage = "automatically test your APIs"
	app.Version = "0.0.1"
	app.Commands = buildCLICommands()
	app.CommandNotFound = commandNotFound

	return app
}
