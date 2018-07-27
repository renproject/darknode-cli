package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	// Create new cli application
	app := cli.NewApp()
	app.Name = "Darknode CLI"
	app.Usage = "A command-line tool for managing Darknodes."
	app.Version = "1.2.0"

	// Define sub-commands
	app.Commands = []cli.Command{
		{
			Name:  "up",
			Usage: "Deploy a new Darknode",
			Flags: []cli.Flag{
				NameFlag, TagsFlag, KeystoreFlag, PassphraseFlag, NetworkFlag, ConfigFlag,
				AwsFlag, AwsAccessKeyFlag, AwsSecretKeyFlag, AwsInstanceFlag, AwsRegionFlag, AwsElasticIpFlag,
				DoFlag, DoRegionFlag, DoSizeFlag, DoTokenFlag,
			},
			Action: func(c *cli.Context) error {
				return deployNode(c)
			},
		},
		{
			Name:    "destroy",
			Usage:   "Destroy one of your Darknode",
			Aliases: []string{"down"},
			Flags:   []cli.Flag{TagsFlag, ForceFlag},
			Action: func(c *cli.Context) error {
				return destroyNode(c)
			},
		},
		{
			Name:  "update",
			Usage: "Update your Darknodes to the latest software and configuration",
			Flags: []cli.Flag{TagsFlag, BranchFlag, UpdateConfigFlag},
			Action: func(c *cli.Context) error {
				return updateNode(c)
			},
		},
		{
			Name:  "ssh",
			Flags: []cli.Flag{},
			Usage: "SSH into one of your Darknode",
			Action: func(c *cli.Context) error {
				return sshNode(c)
			},
		},
		{
			Name:  "start",
			Flags: []cli.Flag{TagsFlag},
			Usage: "Start one of your Darknodes from a suspended state",
			Action: func(c *cli.Context) error {
				return startNode(c)
			},
		},
		{
			Name:  "stop",
			Flags: []cli.Flag{TagsFlag},
			Usage: "Stop one of your Darknodes by putting it into a suspended state",
			Action: func(c *cli.Context) error {
				return stopNode(c)
			},
		},
		{
			Name:  "list",
			Usage: "List all of your Darknodes",
			Flags: []cli.Flag{TagsFlag},
			Action: func(c *cli.Context) error {
				return listAllNodes(c)
			},
		},
		{
			Name:  "exec",
			Usage: "Exec scripts on nodes",
			Flags: []cli.Flag{TagsFlag, ScriptFlag},
			Action: func(c *cli.Context) error {
				return execScript(c)
			},
		},
	}

	// Show error message and display the help page for the app
	app.CommandNotFound = func(c *cli.Context, command string) {
		cli.ShowAppHelp(c)
		fmt.Fprintf(c.App.Writer, "%scommand %q not found%s.\n", RED, command, RESET)
	}

	// Start the app
	err := app.Run(os.Args)
	if err != nil {
		// Remove the timestamp for error message
		log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
		log.Fatal(err)
	}
}
