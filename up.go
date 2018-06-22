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

	"github.com/republicprotocol/republic-go/cmd/darknode/config"
	"github.com/urfave/cli"
)

// ErrKeyNotFound is returned when no AWS access-key nor secret-key provided.
var ErrKeyNotFound = errors.New("please provide your AWS access key and secret key")

// ErrNodeExist is returned when trying to created a new node with name already
// exists.
var ErrNodeExist = errors.New("node with the name already exists")

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
	name := ctx.String("name")

	// Check input flags
	var nodeDirectory string
	if accessKey == "" || secretKey == "" {
		//TODO : Read FROM ~/aws/  FOLDER
		return ErrKeyNotFound
	}
	if name == ""{
		for i := 1; ;i ++ {
			if _, err := os.Stat(Directory+ fmt.Sprintf("/darknodes/darknode%d",i)); os.IsNotExist(err) {
				nodeDirectory = Directory+ fmt.Sprintf("/darknodes/darknode%d",i)
				break
			}
		}
	} else {
		if _, err := os.Stat(Directory+ "/darknodes/"+ name); !os.IsNotExist(err) {
			return ErrNodeExist
		}
		nodeDirectory = Directory+ "/darknodes/"+  name
	}
	os.Mkdir(nodeDirectory, 0777)

	// Parse region and instance type
	region, instance, err := parseRegionAndInstance(ctx)
	if err != nil {
		return err
	}

	// Generate configs for the node
	config, err := GetConfigOrGenerateNew(nodeDirectory)
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
		return err
	}
	ip, err := getIp(nodeDirectory)
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
func runTerraform(nodeDirectory string ) error {

	cmd := fmt.Sprintf("cd %v && terraform init",nodeDirectory)
	init := exec.Command( "bash", "-c", cmd)
	pipeToStd(init)
	if err := init.Start(); err != nil {
		return err
	}
	if err := init.Wait(); err != nil {
		return err
	}
	log.Println("Deploying dark nodes to AWS...")
	cmd = fmt.Sprintf("cd %v && terraform apply -auto-approve",nodeDirectory)
	apply := exec.Command( "bash", "-c", cmd)
	pipeToStd(apply)
	if err := apply.Start(); err != nil {
		return err
	}
	if err:=  apply.Wait();err !=nil {
		return err
	}
	return  exec.Command("cd","-").Run()
}

func generateTerraformConfig(config config.Config, accessKey, secretKey, region, instance, pubKey , nodeDirectory string) error {
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
	`, accessKey, secretKey, strings.TrimSpace(pubKey), nodeDirectory +  "/ssh_keypair" )

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

	return ioutil.WriteFile( nodeDirectory  + "/main.tf", []byte(terraformConfig+mode), 0600)
}

// deployToDigitalOcean parses the digital ocean credentials and use terraform
// to deploy the node to digital ocean.
func deployToDigitalOcean(ctx *cli.Context) error {
	panic("unimplemented")
}
