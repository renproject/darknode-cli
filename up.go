package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os/exec"
	"strings"

	"encoding/json"
	"github.com/republicprotocol/republic-go/cmd/darknode/config"
	"github.com/urfave/cli"
)

// KeyNotFound is returned when no AWS access-key nor secret-key provided.
var KeyNotFound error = errors.New("please provide your AWS access key and secret key")

// deployNode deploys node depending on the provider.
func deployNode(ctx *cli.Context, path string) error {
	provider := strings.ToLower(ctx.String("provider"))
	switch provider {
	case "aws":
		return deployToAWS(ctx, path)
	case "digital-ocean":
		return deployToDigitalOcean(ctx)
	default:
		return errors.New("unsupported service provider")
	}
}

// deployToAWS parses the AWS credentials and use terraform to deploy the node
// to AWS.
func deployToAWS(ctx *cli.Context, path string) error {
	accessKey := ctx.String("access_key")
	secretKey := ctx.String("secret_key")
	if accessKey == "" || secretKey == "" {
		//TODO : Read FROM ~/aws/  FOLDER
		return KeyNotFound
	}

	region, instance, err := parseRegionAndInstance(ctx)
	if err != nil {
		return err
	}
	config, err := NewConfig()
	if err != nil {
		return err
	}
	configData, err := json.Marshal(config)
	if err != nil {
		return err
	}
	ioutil.WriteFile(fmt.Sprintf("%v/config.json", path), configData, 0600)

	pubKey, keyPath, err := NewSshKeyPair(path)
	if err != nil {
		return err
	}

	generateTerraformConfig(config, path, accessKey, secretKey, region, instance, pubKey, keyPath)

	return nil
}

// runTerraform initializes and applies terraform
func runTerraform() error {
	init := exec.Command("./terraform", "init")
	if err := init.Run(); err != nil {
		return err
	}
	if err := init.Wait(); err != nil {
		return err
	}

	apply := exec.Command("./terraform", "apply", "-auto-approve")
	if err := apply.Run(); err != nil {
		return err
	}
	return apply.Wait()
}

func generateTerraformConfig(config config.Config, path, accessKey, secretKey, region, instance, pubKey, keyPath string) error {
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

	`, accessKey, secretKey, strings.TrimSpace(pubKey), keyPath)

	avz := region + AvailableZones[region][rand.Intn(len(AvailableZones[region]))]
	//Fixme : does the bootstrap field important?
	mode := fmt.Sprintf(`
module "{node-%v}" {
    source = "./instance"
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
	bootstraps = []
    is_bootstrap = "false"
    port = "%v"
}`, config.Address, AMIs[region], region, avz, config.Address, instance, path, config.Port)

	return ioutil.WriteFile(fmt.Sprintf("%v/main.tf", path), []byte(terraformConfig+mode), 0600)
}

// deployToDigitalOcean parses the digital ocean credentials and use terraform
// to deploy the node to digital ocean.
func deployToDigitalOcean(ctx *cli.Context) error {
	panic("unimplemented")
}
