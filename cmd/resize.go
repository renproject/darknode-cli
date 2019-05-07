package main

import (
	"fmt"
	"io/ioutil"
	"path"
	"regexp"
	"strings"

	"github.com/urfave/cli"
)

func resize(ctx *cli.Context) error {
	name := ctx.Args().First()
	newSize := ctx.Args().Get(1)

	// Validate the name
	nodePath, err := validateDarknodeName(name)
	if err != nil {
		return err
	}
	if newSize == "" {
		return ErrInvalidInstanceSize
	}

	// Get main.tf file
	filePath := path.Join(nodePath, "main.tf")
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Check if it's aws or digital ocean
	if strings.Contains(string(data), `provider "aws"`) {
		return resizeAwsInstance(data, nodePath, filePath, newSize)
	} else if strings.Contains(string(data), `provider "digitalocean"`) {
		return resizeDoInstance(data, nodePath, filePath, newSize)
	} else {
		return ErrUnknownProvider
	}
}

func resizeAwsInstance(tfFile []byte, nodePath, tfPath, newSize string) error {
	reg, err := regexp.Compile(`variable "instance_type" \{\s+default = ".+"\s\}`)
	if err != nil {
		return err
	}
	replacement := fmt.Sprintf("variable \"instance_type\" {\n  default = \"%v\"\n}", newSize)
	newTF := reg.ReplaceAll(tfFile, []byte(replacement))
	if err := ioutil.WriteFile(tfPath, newTF, 0644); err != nil {
		return err
	}

	// Start running terraform
	fmt.Printf("\n%sResizing dark nodes ... %s\n", RESET, RESET)
	apply := fmt.Sprintf("cd %v && terraform apply -auto-approve -no-color", nodePath)
	err = run("bash", "-c", apply)
	if err != nil {
		// revert the `main.tf` file if fail to resize the droplet
		defer func() {
			if err := ioutil.WriteFile(tfPath, tfFile, 0644); err != nil {
				fmt.Println("fail to revert the change to `main.tf` file")
			}
		}()
		return err
	}

	// Update ip address to the multiAddress.out file
	update := fmt.Sprintf("cd %v && terraform output multiaddress > multiAddress.out", nodePath)
	return run("bash", "-c", update)
}

func resizeDoInstance(tfFile []byte, nodePath, tfPath, newSize string) error {
	// Replace with the new size in the `main.tf` file
	reg, err := regexp.Compile(`variable "size" \{\s+default = ".+"\s\}`)
	if err != nil {
		return err
	}
	replacement := fmt.Sprintf("variable \"size\" {\n\tdefault = \"%v\"\n}", newSize)
	newTF := reg.ReplaceAll(tfFile, []byte(replacement))
	if err := ioutil.WriteFile(tfPath, newTF, 0644); err != nil {
		return err
	}

	// Mark the droplet as tainted for recreating the droplet
	taint := fmt.Sprintf("cd %v && terraform taint digitalocean_droplet.darknode", nodePath)
	err = run("bash", "-c", taint)
	if err != nil {
		return err
	}

	// Start running terraform
	fmt.Printf("\n%sResizing dark nodes ... %s\n", RESET, RESET)
	apply := fmt.Sprintf("cd %v && terraform apply -auto-approve -no-color", nodePath)
	err = run("bash", "-c", apply)
	if err != nil {
		// revert the `main.tf` file if fail to resize the droplet
		defer func() {
			if err := ioutil.WriteFile(tfPath, tfFile, 0644); err != nil {
				fmt.Println("fail to revert the change to `main.tf` file")
			}
		}()
		return err
	}

	// Update ip address to the multiAddress.out file
	update := fmt.Sprintf("cd %v && terraform output multiaddress > multiAddress.out", nodePath)
	return run("bash", "-c", update)
}
