package main

import (
	"log"
	"os"

	"fmt"
	"github.com/urfave/cli"
)

func main() {
	// Create new cli application
	app := cli.NewApp()

	// fixme: Define flags
	upFlags := []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Value: "darknode",
			Usage: "give your darknode a name",
		},
		cli.StringFlag{
			Name:  "provider",
			Value: "AWS",
			Usage: "cloud service provider you want to use for your darknode, default to AWS",
		},
		cli.StringFlag{
			Name:  "region",
			Value: "",
			Usage: "deployment region",
		},
		cli.StringFlag{
			Name:  "instance",
			Value: "",
			Usage: "instance type",
		},
		cli.StringFlag{
			Name:  "access_key",
			Value: "",
			Usage: "access key for your AWS account",
		},
		cli.StringFlag{
			Name:  "secret_key",
			Value: "",
			Usage: "secret key for your AWS account",
		},
	}

	// Define sub-commands
	app.Commands = []cli.Command{
		{
			Name:  "up",
			Usage: "deploying a new darknode",
			Flags: upFlags,
			Action: func(c *cli.Context) error {
				path := fmt.Sprintf("./%v", c.String("name"))
				if path == "./" {
					path = "./darknode"
				}
				var err error
				path, err = Mkdir(path)
				if err != nil {
					return err
				}

				return deployNode(c, path)
			},
		},
		{
			Name:  "destroy",
			Usage: "tear down the darknode and clean up everything",
			Action: func(c *cli.Context) error {
				panic("todo ")
			},
		},
		{
			Name:  "update",
			Usage: "update your darknode to the latest release",
			Action: func(c *cli.Context) error {
				panic("todo ")
			},
		},
		{
			Name:  "ssh",
			Usage: "ssh into your cloud service instance",
			Action: func(c *cli.Context) error {
				panic("todo ")
			},
		},
	}

	// Start the app
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
