package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/republicprotocol/republic-go/cmd/darknode/config"
	"github.com/urfave/cli"
)

const (
	reset = "\x1b[0m"
	green = "\x1b[32;1m"
	red   = "\x1b[31;1m"
)

// ErrKeyNotFound is returned when no AWS access-key nor secret-key provided.
var ErrKeyNotFound = fmt.Errorf("%splease provide your AWS access key and secret key%s", red, reset)

// ErrNodeExist is returned when user tries to created a new node with name
// already exists.
var ErrNodeExist = fmt.Errorf("%snode with the name already exists%s", red, reset)

// ErrUnknownProvider is returned when user wants to deploy darknode to an
// unknown service provider
var ErrUnknownProvider = fmt.Errorf("%sunknown service provider%s", red, reset)

// ErrNilProvider is returned when the provider is nil.
var ErrNilProvider = fmt.Errorf("%sprovider cannot be nil%s", red, reset)

// deployNode deploys node depending on the provider.
func deployNode(ctx *cli.Context) error {
	provider := strings.ToLower(ctx.String("provider"))

	switch provider {
	case "":
		cli.ShowCommandHelp(ctx, "up")
		return ErrNilProvider
	case "aws":
		return deployToAWS(ctx)
	case "digital-ocean":
		return deployToDigitalOcean(ctx)
	default:
		return ErrUnknownProvider
	}
}

// deployToAWS parses the AWS credentials and use terraform to deploy the node
// to AWS.
func deployToAWS(ctx *cli.Context) error {
	accessKey := ctx.String("aws-access-key")
	secretKey := ctx.String("aws-secret-key")
	name := ctx.String("name")
	tags := ctx.String("tags")

	// Check input flags
	var nodeDirectory string
	if accessKey == "" || secretKey == "" {
		// Try reading the credentials from the default file.
		creds := credentials.NewSharedCredentials("", "default")
		credValue, err := creds.Get()
		if err != nil {
			return err
		}
		accessKey, secretKey = credValue.AccessKeyID, credValue.SecretAccessKey
		if accessKey == "" || secretKey == "" {
			return ErrKeyNotFound
		}
	}
	// Check darknode name and make directory for the node
	if name == "" {
		return ErrEmptyNodeName
	}
	if _, err := os.Stat(Directory + "/darknodes/" + name); !os.IsNotExist(err) {
		return ErrNodeExist
	}
	nodeDirectory = Directory + "/darknodes/" + name
	if err := os.Mkdir(nodeDirectory, 0777); err != nil {
		return err
	}

	// Store the tags
	if err := ioutil.WriteFile(nodeDirectory+"/tags.out", []byte(strings.ToLower(strings.TrimSpace(tags))), 0666); err != nil {
		return err
	}

	// Parse region and instance type
	region, instance, err := parseRegionAndInstance(ctx)
	if err != nil {
		return err
	}
	// Generate configs for the node
	config, err := GetConfigOrGenerateNew(ctx, nodeDirectory)
	if err != nil {
		return err
	}
	pubKey, err := NewSshKeyPair(nodeDirectory)
	if err != nil {
		return err
	}
	if err := generateTerraformConfig(config, accessKey, secretKey, region, instance, pubKey, nodeDirectory); err != nil {
		return err
	}
	if err := runTerraform(nodeDirectory); err != nil {
		if err := cleanUp(nodeDirectory); err != nil {
			return err
		}
		return err
	}

	ip, err := getIp(nodeDirectory)
	if err != nil {
		if err := cleanUp(nodeDirectory); err != nil {
			return err
		}
		return err
	}
	fmt.Printf("\n")
	fmt.Printf("%sCongratulations! Your Darknode is deployed and running%s.\n", green, reset)
	fmt.Printf("%sJoin the network by registering your Darknode at%s\n", green, reset)
	fmt.Printf("%shttps://darknode.republicprotocol.com/status/%v%s\n", green, ip, reset)
	fmt.Printf("\n")
	return nil
}

// runTerraform initializes and applies terraform
func runTerraform(nodeDirectory string) error {
	cmd := fmt.Sprintf("cd %v && terraform init", nodeDirectory)
	init := exec.Command("bash", "-c", cmd)
	pipeToStd(init)
	if err := init.Start(); err != nil {
		return err
	}
	if err := init.Wait(); err != nil {
		return err
	}
	fmt.Printf("%sDeploying dark nodes to AWS%s...\n", green, reset)
	cmd = fmt.Sprintf("cd %v && terraform apply -auto-approve", nodeDirectory)
	apply := exec.Command("bash", "-c", cmd)
	pipeToStd(apply)
	if err := apply.Start(); err != nil {
		return err
	}
	if err := apply.Wait(); err != nil {
		return err
	}
	return nil
}

func generateTerraformConfig(config config.Config, accessKey, secretKey, region, instance, pubKey, nodeDirectory string) error {
	terraformConfig := fmt.Sprintf(`
variable "access_key" {
	default = "%v"
}

variable "secret_key" {
	default = "%v"	
}

variable "ssh_public_key" {
	default = "%v"
}

variable "ssh_private_key_location" {
	default = "%v"
}
	`, accessKey, secretKey, strings.TrimSpace(pubKey), nodeDirectory+"/ssh_keypair")

	avz := region + AvailableZones[region][rand.Intn(len(AvailableZones[region]))]
	mode := fmt.Sprintf(`
module "node-%v" {
    source = "%v/instance"
    ami = "%v"
    region = "%v"
    avz = "%v"
    id = "%v"
    ec2_instance_type = "%v"
    ssh_public_key = "${var.ssh_public_key}"
    ssh_private_key_location = "${var.ssh_private_key_location}"
    access_key = "${var.access_key}"
    secret_key = "${var.secret_key}"
    config = "%v/config.json"
    is_bootstrap = "false"
    port = "%v"
    path = "%v"
}`, config.Address, Directory, AMIs[region], region, avz, config.Address, instance, nodeDirectory, config.Port, Directory)

	return ioutil.WriteFile(nodeDirectory+"/main.tf", []byte(terraformConfig+mode), 0600)
}

// deployToDigitalOcean parses the digital ocean credentials and use terraform
// to deploy the node to digital ocean.
func deployToDigitalOcean(ctx *cli.Context) error {
	panic("unimplemented")
}

// cleanUp removes the directory
func cleanUp(nodeDirectory string) error {
	cleanCmd := exec.Command("rm", "-rf", nodeDirectory)
	if err := cleanCmd.Start(); err != nil {
		return err
	}
	return cleanCmd.Wait()
}
