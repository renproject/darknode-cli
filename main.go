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
			Name:  "access-key",
			Value: "",
			Usage: "access key for your AWS account",
		},
		cli.StringFlag{
			Name:  "secret-key",
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
				path, err = mkdir(path)
				if err != nil {
					return err
				}

				return deployNode(c)
			},
		},
		{
			Name:  "destroy",
			Usage: "tear down the darkndoe and clean up everything",
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

// mkdir tries to create folder with the given name. If not success,
// it will try appending number appended to the name.
func mkdir(name string) (string, error) {
	// Try creating directory with the exact name we get
	if _, err := os.Stat(name); os.IsNotExist(err) {
		err = os.Mkdir(name, os.ModePerm)
		if err == nil {
			return name, nil
		}
	}

	// Try creating directory with number appended to the name.
	i := 1
	for {
		dirName := fmt.Sprintf("%v_%d", name, i)
		if _, err := os.Stat(dirName); os.IsNotExist(err) {
			err = os.Mkdir(dirName, os.ModePerm)
			if err != nil {
				return "", err
			}
			return dirName, nil
		}
		i++
	}
}
