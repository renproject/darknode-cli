package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/republicprotocol/republic-go/cmd/darknode/config"
	"github.com/urfave/cli"
)

// Available regions on Digital Ocean.
const (
	Global = "global"
	AMS2   = "ams2"
	AMS3   = "ams3"
	BLR1   = "blr1"
	FRA1   = "fra1"
	LON1   = "lon1"
	NYC1   = "nyc1"
	NYC2   = "nyc2"
	NYC3   = "nyc3"
	SF01   = "sf01"
	SF02   = "sf02"
	SGP1   = "sgp1"
	TOR1   = "tor1"
)

var AllDoRegions = []string{
	Global,
	AMS2,
	AMS3,
	BLR1,
	FRA1,
	LON1,
	NYC1,
	NYC2,
	NYC3,
	SF01,
	SF02,
	SGP1,
	TOR1,
}

// All available droplet size on digital ocean
const (
	Size512MB = "512mb"
	Size1GB = "1gb"
	Size2GB = "2gb"
	Size4GB = "4gb"
	Size8GB = "8gb"
	Size16GB = "16gb"
	Size32GB = "32gb"
	Size48GB = "48gb"
	Size64GB = "64gb"
)

var AllDoDropletSize = []string{
	Size512MB,
	Size1GB ,
	Size2GB ,
	Size4GB,
	Size8GB ,
	Size16GB,
	Size32GB ,
	Size48GB ,
	Size64GB ,
}

func doParseRegionAndInstance(ctx *cli.Context)( string ,string ,error ){
	region := ctx.String("do-region")
	size := ctx.String("do-size")

	// Parse the input region or pick one region randomly
	rand.Seed(time.Now().UTC().UnixNano())
	if region == "" {
		region = AllDoRegions[rand.Intn(len(AllDoRegions))]
	} else {
		if !StringInSlice(region, AllDoRegions) {
			return "", "", ErrUnknownRegion
		}
	}

	// Validate the droplet size
	if !StringInSlice(size, AllDoDropletSize) {
		// todo : oupput available droplet sizes.
		return "", "", ErrUnSupportedInstanceType
	}

	return region, size, nil
}


func generateTerraformConfigForDo(ctx *cli.Context, config config.Config, token string) error {

	terraformConfig := fmt.Sprintf(`
variable "do_token" {
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
	name := ctx.String("name")
	tags := ctx.String("tags")
	network := ctx.String("network")
	token := ctx.String("do-token")


	// Check digital ocean token
	if token == "" {
		return ErrEmptyDoToken
	}

	// Check region and droplet size
	region, size, err := doParseRegionAndInstance(ctx)
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
	if err := generateTerraformConfigForAws(ctx, config, accessKey, secretKey, region, instance, pubKey, nodeDirectory); err != nil {
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


}


