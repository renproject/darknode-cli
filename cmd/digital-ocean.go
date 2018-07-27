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

	SizeS1V1GB    = "s-1vcpu-1gb"
	SizeS1V2GB    = "s-1vcpu-2gb"
	SizeS1V3GB    = "s-1vcpu-3gb"
	SizeS2V2GB    = "s-2vcpu-2gb"
	SizeS3V1GB    = "s-3vcpu-1gb"
	SizeS2V4GB    = "s-2vcpu-4gb"
	SizeS4V8GB    = "s-4vcpu-8gb"
	SizeS6V16GB   = "s-6vcpu-16gb"
	SizeS8V32GB   = "s-8vcpu-32gb"
	SizeS12V48GB  = "s-12vcpu-48gb"
	SizeS16V64GB  = "s-16vcpu-64gb"
	SizeS20V96GB  = "s-20vcpu-96gb"
	SizeS24V128GB = "s-24vcpu-128gb"
	SizeS32V192GB = "s-32vcpu-192gb"

	SizeC1V2GB = "c-1vcpu-2gb"
	SizeC2     = "c-2"
	SizeC4     = "c-4"
	SizeC8     = "c-8"
	SizeC16    = "c-16"
	SizeC32    = "c-32"
	SizeC64    = "c-64"
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

	SizeS1V1GB,
	SizeS1V2GB,
	SizeS1V3GB,
	SizeS2V2GB,
	SizeS3V1GB,
	SizeS2V4GB,
	SizeS4V8GB,
	SizeS6V16GB,
	SizeS8V32GB,
	SizeS12V48GB,
	SizeS16V64GB,
	SizeS20V96GB,
	SizeS24V128GB,
	SizeS32V192GB,

	SizeC1V2GB,
	SizeC2,
	SizeC4,
	SizeC8,
	SizeC16,
	SizeC32,
	SizeC64,
}

func parseDoRegionAndSize(ctx *cli.Context) (string, string, error) {
	region := ctx.String("do-region")
	size := ctx.String("do-droplet")

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

	if token == "" {
		return ErrEmptyDoToken
	}
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
