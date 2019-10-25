package provider

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/republicprotocol/darknode-cli/darknode"
	"github.com/republicprotocol/darknode-cli/util"
	"github.com/urfave/cli"
)

// All available regions on AWS.
var AllAwsRegions = []string{
	"ap-northeast-1",
	"ap-northeast-2",
	"ap-south-1",
	"ap-southeast-1",
	"ap-southeast-2",
	"ca-central-1",
	"eu-central-1",
	"eu-north-1",
	"eu-west-1",
	"eu-west-2",
	"eu-west-3",
	"sa-east-1",
	"us-east-1",
	"us-east-2",
	"us-west-1",
	"us-west-2",
}

type providerAws struct {
	accessKey string
	secretKey string
}

func NewAws(ctx *cli.Context) (Provider, error) {
	accessKey := ctx.String("aws-access-key")
	secretKey := ctx.String("aws-secret-key")

	// Try reading the credential files if user does not provide credentials directly
	if accessKey == "" || secretKey == "" {
		cred := credentials.NewSharedCredentials("", ctx.String("aws-profile"))
		credValue, err := cred.Get()
		if err != nil {
			return nil, err
		}
		accessKey, secretKey = credValue.AccessKeyID, credValue.SecretAccessKey
		if accessKey == "" || secretKey == "" {
			return nil, err
		}
	}

	return providerAws{
		accessKey: accessKey,
		secretKey: secretKey,
	}, nil
}

func (p providerAws) Name() string {
	return NameAws
}

func (p providerAws) Deploy(ctx *cli.Context) error {
	name := ctx.String("name")
	tags := ctx.String("tags")

	region, instance, err := p.validateRegionAndInstance(ctx)
	if err != nil {
		return err
	}

	// Initialization
	network, err := darknode.NewNetwork(ctx.String("network"))
	if err != nil {
		return err
	}
	if err := initNode(name, tags, network); err != nil {
		return err
	}

	// Generate terraform config and start deploying
	if err := p.tfConfig(name, region, instance, ipfsUrl(network)); err != nil {
		return err
	}
	if err := runTerraform(name); err != nil {
		return err
	}

	return outputURL(name)
}

func (p providerAws) validateRegionAndInstance(ctx *cli.Context) (string, string, error) {
	// TODO : use aws api to validate the region and instance type
	region := strings.ToLower(strings.TrimSpace(ctx.String("aws-region")))
	instance := strings.ToLower(strings.TrimSpace(ctx.String("aws-instance")))

	if region == "" {
		region = AllAwsRegions[rand.Intn(len(AllAwsRegions))]
	}
	if !util.StringInSlice(region, AllAwsRegions) {
		return "", "", fmt.Errorf("aws region [%v] is not supported yet", region)
	}

	if instance == "" {
		return "", "", fmt.Errorf("instance type [%v] is not supported yet", instance)
	}

	return region, instance, nil
}
