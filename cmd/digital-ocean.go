package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
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
	SF01   = "sfo1"
	SF02   = "sfo2"
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
	Size1GB   = "1gb"
	Size2GB   = "2gb"
	Size4GB   = "4gb"
	Size8GB   = "8gb"
	Size16GB  = "16gb"
	Size32GB  = "32gb"
	Size48GB  = "48gb"
	Size64GB  = "64gb"
)

var AllDoDropletSize = []string{
	Size512MB,
	Size1GB,
	Size2GB,
	Size4GB,
	Size8GB,
	Size16GB,
	Size32GB,
	Size48GB,
	Size64GB,
}

func parseDoRegionAndSize(ctx *cli.Context) (string, string, error) {
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
		// todo : output available droplet sizes.
		return "", "", ErrUnSupportedInstanceType
	}

	return region, size, nil
}

// deployToDo parses the digital ocean credentials and use terraform to
// deploy the node to digital ocean.
func deployToDo(ctx *cli.Context) error {
	token := ctx.String("do-token")

	// Parse DO related data.
	region, size, err := parseDoRegionAndSize(ctx)
	if err != nil {
		return err
	}

	// Create node directory
	name, err := createNodeDirectory(ctx)
	if err != nil {
		return err
	}
	nodeDir := nodeDirectory(name)

	// Generate config and ssh key for the node
	config, err := GetConfigOrGenerateNew(ctx, nodeDir)
	if err != nil {
		return err
	}
	_, err = NewSshKeyPair(nodeDir)
	if err != nil {
		return err
	}

	// Generate terraform config and start deploying
	if err := generateDoTFConfig(config, token, name, nodeDir, region, size); err != nil {
		return err
	}
	if err := runTerraform(nodeDir); err != nil {
		return err
	}

	return outputUrl(ctx, name, nodeDir)
}

// generateDoTFConfig generates the terraform config file for deploying to DO.
func generateDoTFConfig(config config.Config, token, name, nodeDir, region, size string) error {
	terraformConfig := fmt.Sprintf(`
variable "do_token" {
	default = "%v"
}

variable "name" {
	default = "%v"
}

variable "region" {
	default = "%v"
}

variable "size" {
	default = "%v"
}

variable "path" {
  default = "%v"
}

variable "id" {
  default = "%v"
}

variable "pub_key" {
  default = "%v/darknodes/%v/ssh_keypair.pub"
}

variable "pvt_key" {
  default = "%v/darknodes/%v/ssh_keypair"
}
	`, token, name, region, size, Directory, config.Address, Directory, name, Directory, name)

	err := ioutil.WriteFile(nodeDir+"/variables.tf", []byte(terraformConfig), 0644)
	if err != nil {
		return err
	}

	return copyFile(Directory+"/instance/digital-ocean/main.tf", nodeDir+"/main.tf")
}
