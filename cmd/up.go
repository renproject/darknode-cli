package main

import (
	"encoding/json"
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

var Providers = []string{"aws", "digitalocean"}

// deployNode deploys node depending on the provider.
func deployNode(ctx *cli.Context) error {
	aws := ctx.Bool("aws")
	digitalOcean := ctx.Bool("digitalocean")

	// Make sure only one provider is provided
	counter, provider := 0, ""
	for i, j := range []bool{aws, digitalOcean} {
		if j {
			counter++
			provider = Providers[i]
		}
	}
	if counter == 0 {
		return ErrNilProvider
	} else if counter > 1 {
		return ErrMultipleProviders
	}

	switch provider {
	case "aws":
		return deployToAWS(ctx)
	case "digitalocean":
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
	network := ctx.String("network")
	name := ctx.String("name")
	tags := ctx.String("tags")

	// Try getting AWS credentials from the input or the default file.
	if accessKey == "" || secretKey == "" {
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

	// Parse region and instance type
	region, instance, err := parseRegionAndInstance(ctx)
	if err != nil {
		return err
	}
	// Generate configs for the node
	config, err := GetConfigOrGenerateNew(ctx)
	if err != nil {
		return err
	}

	// Check darknode name and make directory for the node
	if name == "" {
		return ErrEmptyNodeName
	}
	if _, err := os.Stat(Directory + "/darknodes/" + name); !os.IsNotExist(err) {
		return ErrNodeExist
	}
	nodeDirectory := Directory + "/darknodes/" + name
	if err := os.Mkdir(nodeDirectory, 0777); err != nil {
		return err
	}
	// Store the tags
	if err := ioutil.WriteFile(nodeDirectory+"/tags.out", []byte(strings.TrimSpace(tags)), 0666); err != nil {
		return err
	}
	// Write the config to file
	configData, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(nodeDirectory+"/config.json", configData, 0600); err != nil {
		return err
	}
	// Generate new ssk key pair
	pubKey, err := NewSshKeyPair(nodeDirectory)
	if err != nil {
		if err := cleanUp(nodeDirectory); err != nil {
			return err
		}
		return err
	}
	if err := generateTerraformConfig(ctx, config, accessKey, secretKey, region, instance, pubKey, nodeDirectory); err != nil {
		if err := cleanUp(nodeDirectory); err != nil {
			return err
		}
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

	// Update node to different branch according to the network.
	switch network {
	case "testnet":
	case "falcon":
		err = updateSingleNode(name, "develop", false)
	case "nightly":
		err = updateSingleNode(name, "nightly", false)
	}

	fmt.Printf("\n")
	fmt.Printf("%sCongratulations! Your Darknode is deployed and running%s.\n", GREEN, RESET)
	fmt.Printf("%sJoin the network by registering your Darknode at%s\n", GREEN, RESET)
	fmt.Printf("%shttps://darknode.republicprotocol.com/status/%v%s\n", GREEN, ip, RESET)
	fmt.Printf("\n")
	return err
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

	fmt.Printf("%sDeploying dark nodes to AWS%s...\n", GREEN, RESET)

	cmd = fmt.Sprintf("cd %v && terraform apply -auto-approve", nodeDirectory)
	apply := exec.Command("bash", "-c", cmd)
	pipeToStd(apply)
	if err := apply.Start(); err != nil {
		return err
	}
	return apply.Wait()
}

func generateTerraformConfig(ctx *cli.Context, config config.Config, accessKey, secretKey, region, instance, pubKey, nodeDirectory string) error {
	allocationID := ctx.String("aws-elastic-ip")

	allocationConfig, tfFolder := "", "std"
	if allocationID != "" {
		allocationConfig = fmt.Sprintf(`allocation_id = "%v"`, allocationID)
		tfFolder = "eip"
	}

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
    source = "%v/instance/%v"
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
    port = "%v"
    path = "%v"
    %v
}`, config.Address, Directory, tfFolder, AMIs[region], region, avz, config.Address, instance, nodeDirectory, config.Port, Directory, allocationConfig)

	return ioutil.WriteFile(nodeDirectory+"/main.tf", []byte(terraformConfig+mode), 0600)
}

// deployToDigitalOcean parses the digital ocean credentials and use terraform
// to deploy the node to digital ocean.
func deployToDigitalOcean(ctx *cli.Context) error {
	panic("unimplemented")
}
