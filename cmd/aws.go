package main

import (
	"math/rand"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/republicprotocol/republic-go/cmd/darknode/config"
	"github.com/urfave/cli"
	"golang.org/x/crypto/ssh"
)

// All available regions on AWS.
var AllAwsRegions = []string{
	"ap-northeast-1",
	"ap-northeast-2",
	// "ap-northeast-3", Terraform having issue support this provider
	"ap-south-1",
	"ap-southeast-1",
	"ap-southeast-2",
	"ca-central-1",
	"eu-central-1",
	"eu-west-1",
	"eu-west-2",
	"eu-west-3",
	"sa-east-1",
	"us-east-1",
	"us-east-2",
	"us-west-1",
	"us-west-2",
}

// awsDeployment parses the AWS credentials and use terraform to deploy the
// node to AWS.
func awsDeployment(ctx *cli.Context) error {
	region := strings.ToLower(ctx.String("aws-region"))
	instance := strings.ToLower(ctx.String("aws-instance"))
	if region == "" {
		region = AllAwsRegions[rand.Intn(len(AllAwsRegions))]
	}

	accessKey, secretKey, err := awsCredentials(ctx)
	if err != nil {
		return err
	}

	// Create node directory
	name := ctx.String("name")
	tags := ctx.String("tags")
	if err := mkdir(name, tags); err != nil {
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
	if err := awsTerraformConfig(ctx, config, key, accessKey, secretKey, nodeDir, region, instance); err != nil {
		return err
	}
	if err := runTerraform(nodeDir); err != nil {
		return err
	}

	return outputUrl(name, nodeDir)
}

// awsCredentials tries to get the AWS credentials from the user input
// or from the default aws credential file
func awsCredentials(ctx *cli.Context) (string, string, error) {
	profile := ctx.String("aws-profile")
	accessKey := ctx.String("aws-access-key")
	secretKey := ctx.String("aws-secret-key")

	// Try read the credential files if user does not provide them directly
	if accessKey == "" || secretKey == "" {
		creds := credentials.NewSharedCredentials("", profile)
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

// Terraform contains all the fields needed to generate a terraform config file
// so that we can deploy the node on AWS.
type Terraform struct {
	Name          string
	Source        string
	Region        string
	Address       string
	InstanceType  string
	SshPubKey     string
	SshPriKeyPath string
	AccessKey     string
	SecretKey     string
	Port          string
	Path          string
	AllocationID  string
}

// awsTerraformConfig generates the terraform config file for deploying to AWS.
func awsTerraformConfig(ctx *cli.Context, config config.Config, key ssh.PublicKey, accessKey, secretKey, nodeDir, region, instance string) error {
	tf := Terraform{
		Name:          ctx.String("name"),
		Region:        region,
		Address:       config.Address.String(),
		InstanceType:  instance,
		SshPubKey:     strings.TrimSpace(StringfySshPubkey(key)),
		SshPriKeyPath: path.Join(nodeDir, "ssh_keypair"),
		AccessKey:     accessKey,
		SecretKey:     secretKey,
		Port:          config.Port,
		Path:          Directory,
		AllocationID:  ctx.String("aws-elastic-ip"),
	}

	fmap := template.FuncMap{}
	templateFile := path.Join(Directory, "instance", "aws", "aws.tmpl")
	t := template.Must(template.New("aws.tmpl").Funcs(fmap).ParseFiles(templateFile))
	tfFile, err := os.Create(path.Join(nodeDir, "main.tf"))
	if err != nil {
		return err
	}

	return t.Execute(tfFile, tf)
}
