package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/republicprotocol/republic-go/cmd/darknode/config"
	"github.com/urfave/cli"
	"golang.org/x/crypto/ssh"
)

// Available regions on AWS.
const (
	ApNorthEast1 = "ap-northeast-1"
	ApNorthEast2 = "ap-northeast-2"
	ApSouth1     = "ap-south-1"
	ApSouthEast1 = "ap-southeast-1"
	ApSouthEast2 = "ap-southeast-2"
	CaCentral1   = "ca-central-1"
	EuCentral1   = "eu-central-1"
	EuWest1      = "eu-west-1"
	EuWest2      = "eu-west-2"
	EuWest3      = "eu-west-3"
	SaEast1      = "sa-east-1"
	UsEast1      = "us-east-1"
	UsEast2      = "us-east-2"
	UsWest1      = "us-west-1"
	UsWest2      = "us-west-2"
)

var AllAwsRegions = []string{
	ApNorthEast1,
	ApNorthEast2,
	ApSouth1,
	ApSouthEast1,
	ApSouthEast2,
	CaCentral1,
	EuCentral1,
	EuWest1,
	EuWest2,
	EuWest3,
	SaEast1,
	UsEast1,
	UsEast2,
	UsWest1,
	UsWest2,
}

// AMIs maps the region to the AMI id.
var AMIs = map[string]string{
	ApNorthEast1: "ami-82c928fd",
	ApNorthEast2: "ami-d0cf66be",
	ApSouth1:     "ami-1118397e",
	ApSouthEast1: "ami-b64866ca",
	ApSouthEast2: "ami-fbb66399",
	CaCentral1:   "ami-e3189987",
	EuCentral1:   "ami-331d3bd8",
	EuWest1:      "ami-0b91bd72",
	EuWest2:      "ami-cc6d8eab",
	EuWest3:      "ami-e7cf7e9a",
	SaEast1:      "ami-e8da8984",
	UsEast1:      "ami-7ad76705",
	UsEast2:      "ami-f3211396",
	UsWest1:      "ami-ef415d8f",
	UsWest2:      "ami-22741f5a",
}

// AvailableZones maps the region to its available zones.
var AvailableZones = map[string][]string{
	ApNorthEast1: {"a", "c", "d"},
	ApNorthEast2: {"a", "c"},
	ApSouth1:     {"a", "b"},
	ApSouthEast1: {"a", "b", "c"},
	ApSouthEast2: {"a", "b", "c"},
	CaCentral1:   {"a", "b"},
	EuCentral1:   {"a", "b", "c"},
	EuWest1:      {"a", "b", "c"},
	EuWest2:      {"a", "b", "c"},
	EuWest3:      {"a", "b", "c"},
	SaEast1:      {"a", "c"},
	UsEast1:      {"a", "b", "c", "d", "e", "f"},
	UsEast2:      {"a", "b", "c"},
	UsWest1:      {"b", "c"},
	UsWest2:      {"a", "b", "c"},
}

// Available instance types on AWS.
const (
	T2Nano    = "t2.nano"
	T2Micro   = "t2.micro"
	T2Small   = "t2.small"
	T2Medium  = "t2.medium"
	T2Large   = "t2.large"
	T2XLarge  = "t2.xlarge"
	T2XXLarge = "t2.xxlarge"

	M4Large    = "m4.large"
	M4XLarge   = "m4.xlarge"
	M42XLarge  = "m4.2xlarge"
	M44XLarge  = "m4.4xlarge"
	M410XLarge = "m4.10xlarge"
	M416XLarge = "m4.16xlarge"

	M5Large    = "m5.large"
	M5XLarge   = "m5.xlarge"
	M52XLarge  = "m5.2xlarge"
	M54XLarge  = "m5.4xlarge"
	M512XLarge = "m5.12xlarge"
	M524XLarge = "m5.24xlarge"
)

// AllAwsInstances contains all instance types available on AWS
var AllAwsInstances = []string{
	T2Nano,
	T2Micro,
	T2Small,
	T2Medium,
	T2Large,
	T2XLarge,
	T2XXLarge,
	M4Large,
	M4XLarge,
	M42XLarge,
	M44XLarge,
	M410XLarge,
	M416XLarge,
	M5Large,
	M5XLarge,
	M52XLarge,
	M54XLarge,
	M512XLarge,
	M524XLarge,
}

// AllAwsInstances contains all instance types available in eu-west-3 region
var AllAwsInstancesInEuWest3 = []string{
	T2Nano,
	T2Micro,
	T2Small,
	T2Medium,
	T2Large,
	T2XLarge,
	T2XXLarge,
	M5Large,
	M5XLarge,
	M52XLarge,
	M54XLarge,
	M512XLarge,
	M524XLarge,
}

// AllAwsInstances contains all instance types available in ap-northeast-1
// region
var AllAwsInstancesInApNortheast1 = []string{
	T2Nano,
	T2Micro,
	T2Small,
	T2Medium,
	T2Large,
	T2XLarge,
	T2XXLarge,
	M4Large,
	M4XLarge,
	M42XLarge,
	M44XLarge,
	M410XLarge,
	M416XLarge,
}

// deployToAws parses the AWS credentials and use terraform to deploy the node
// to AWS.
func deployToAws(ctx *cli.Context) error {
	// Parse AWS related data.
	accessKey, secretKey, err := parseAwsCredentials(ctx)
	if err != nil {
		return err
	}
	region, instance, err := parseAwsRegionAndInstance(ctx)
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
	key, err := NewSshKeyPair(nodeDir)
	if err != nil {
		return err
	}

	// Generate terraform config and start deploying
	if err := generateAwsTFConfig(ctx, config, key, accessKey, secretKey, nodeDir, region, instance); err != nil {
		return err
	}
	if err := runTerraform(nodeDir); err != nil {
		return err
	}

	return outputUrl(ctx, name, nodeDir)
}

// parseAwsRegionAndInstance parses the region and the instance type from the
// cli parameters. It will randomly pick a region for the user if it's not
// specified. The default value for instance is `t2.small`.
func parseAwsRegionAndInstance(ctx *cli.Context) (string, string, error) {
	region := strings.ToLower(ctx.String("aws-region"))
	instance := strings.ToLower(ctx.String("aws-instance"))

	// Parse the input region or pick one region randomly
	rand.Seed(time.Now().UTC().UnixNano())
	if region == "" {
		region = AllAwsRegions[rand.Intn(len(AllAwsRegions))]
	} else {
		if !StringInSlice(region, AllAwsRegions) {
			return "", "", ErrUnknownRegion
		}
	}

	// Parse the input instance type or use the default one.
	if region == EuWest3 && !StringInSlice(instance, AllAwsInstancesInEuWest3) {
		return "", "", ErrUnSupportedInstanceType
	}
	if region == ApNorthEast1 && !StringInSlice(instance, AllAwsInstancesInApNortheast1) {
		return "", "", ErrUnSupportedInstanceType
	}
	if !StringInSlice(instance, AllAwsInstances) {
		return "", "", ErrUnSupportedInstanceType
	}

	return region, instance, nil
}

// parseAwsCredentials tries to get the AWS credentials from the user input
// or from the default aws credential file
func parseAwsCredentials(ctx *cli.Context) (string, string, error) {
	accessKey := ctx.String("aws-access-key")
	secretKey := ctx.String("aws-secret-key")

	// Try getting AWS credentials from the input or the default file.
	if accessKey == "" || secretKey == "" {
		creds := credentials.NewSharedCredentials("", "default")
		credValue, err := creds.Get()
		if err != nil {
			return "", "", err
		}
		accessKey, secretKey = credValue.AccessKeyID, credValue.SecretAccessKey
		if accessKey == "" || secretKey == "" {
			return "", "", ErrKeyNotFound
		}
	}

	return accessKey, secretKey, nil
}

// generateAwsTFConfig generates the terraform config file for deploying to AWS.
func generateAwsTFConfig(ctx *cli.Context, config config.Config, key ssh.PublicKey, accessKey, secretKey, nodeDir, region, instance string) error {
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
	`, accessKey, secretKey, strings.TrimSpace(StringfySshPubkey(key)), nodeDir+"/ssh_keypair")

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
}`, config.Address, Directory, tfFolder, AMIs[region], region, avz, config.Address, instance, nodeDir, config.Port, Directory, allocationConfig)

	return ioutil.WriteFile(nodeDir+"/main.tf", []byte(terraformConfig+mode), 0644)
}
