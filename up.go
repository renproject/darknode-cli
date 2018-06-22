package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"strings"
	"os/exec"
	"os"
	"path"

	"github.com/republicprotocol/republic-go/cmd/darknode/config"
	"github.com/urfave/cli"
)

var terraformPath  =  path.Join(os.Getenv("HOME"), ".darknode/terraform")

// KeyNotFound is returned when no AWS access-key nor secret-key provided.
var KeyNotFound error = errors.New("please provide your AWS access key and secret key")

// deployNode deploys node depending on the provider.
func deployNode(ctx *cli.Context) error {
	provider := strings.ToLower(ctx.String("provider"))
	switch provider {
	case "aws":
		return deployToAWS(ctx)
	case "digital-ocean":
		return deployToDigitalOcean(ctx)
	default:
		return errors.New("unsupported service provider")
	}
}

// deployToAWS parses the AWS credentials and use terraform to deploy the node
// to AWS.
func deployToAWS(ctx *cli.Context) error {
	accessKey := ctx.String("access-key")
	secretKey := ctx.String("secret-key")
	if accessKey == "" || secretKey == "" {
		//TODO : Read FROM ~/aws/  FOLDER
		return KeyNotFound
	}

	region, instance, err := parseRegionAndInstance(ctx)
	if err != nil {
		return err
	}
	config, err := GetConfigOrGenerateNew()
	if err != nil {
		return err
	}
	pubKey, keyPath, err := NewSshKeyPair()
	if err != nil {
		return err
	}
	if err := generateTerraformConfig(config, accessKey, secretKey, region, instance, pubKey, keyPath); err != nil {
		return err
	}
	if err := runTerraform(); err != nil {
		return err
	}
	ip, err := getIp()
	if err != nil {
		return err
	}
	fmt.Printf("---------------------------------------------------------------\n")
	fmt.Printf("--- Congratulations! You darknode is deployed and running -----\n")
	fmt.Printf("--- Start regitering your node by going to the following URL --\n")
	fmt.Printf("--- https://darknode.republicprotocol.com/ip4/%v ----\n", ip)
	fmt.Printf("---------------------------------------------------------------\n")
	return nil
}

// runTerraform initializes and applies terraform
func runTerraform() error {

	init := exec.Command(terraformPath, "init")
	pipeToStd(init)
	if err := init.Start(); err != nil {
		return err
	}
	if err := init.Wait(); err != nil {
		return err
	}
	log.Println("Deploying dark nodes to AWS...")
	apply := exec.Command(terraformPath, "apply", "-auto-approve")
	pipeToStd(apply)
	if err := apply.Start(); err != nil {
		return err
	}
	return apply.Wait()
}

func generateTerraformConfig(config config.Config, accessKey, secretKey, region, instance, pubKey, keyPath string) error {
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
	mode := fmt.Sprintf(`
module "node-%v" {
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
    config = "./config.json"
    is_bootstrap = "false"
    port = "%v"
}`, config.Address, AMIs[region], region, avz, config.Address, instance, config.Port)

	terraformConfigPath :=  path.Join(os.Getenv("HOME"), ".darknode/main.tf")
	return ioutil.WriteFile(terraformConfigPath, []byte(terraformConfig+mode), 0600)
}

// deployToDigitalOcean parses the digital ocean credentials and use terraform
// to deploy the node to digital ocean.
func deployToDigitalOcean(ctx *cli.Context) error {
	panic("unimplemented")
}
