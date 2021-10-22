package provider

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/renproject/darknode-cli/util"
	"github.com/urfave/cli"
)

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
	// Validate all input params
	if err := validateCommonParams(ctx); err != nil {
		return err
	}

	name := ctx.String("name")
	region, instance, err := p.validateRegionAndInstance(ctx)
	if err != nil {
		return err
	}

	// Get the latest darknode version
	latestVersion, err := util.LatestStableRelease()
	if err != nil {
		return err
	}

	// Initialize folder and files for the node
	if err := initNode(ctx); err != nil {
		return err
	}

	// Generate terraform config and start deploying
	if err := p.tfConfig(name, region, instance, latestVersion); err != nil {
		return err
	}
	if err := runTerraform(name); err != nil {
		return err
	}

	return outputURL(name)
}

func (p providerAws) validateRegionAndInstance(ctx *cli.Context) (string, string, error) {
	cred := credentials.NewStaticCredentials(p.accessKey, p.secretKey, "")
	region := strings.ToLower(strings.TrimSpace(ctx.String("aws-region")))
	instance := strings.ToLower(strings.TrimSpace(ctx.String("aws-instance")))

	// Get all available regions
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: cred,
	})
	service := ec2.New(sess)
	input := &ec2.DescribeRegionsInput{}
	result, err := service.DescribeRegions(input)
	if err != nil {
		return "", "", err
	}
	regions := make([]string, len(result.Regions))
	for i := range result.Regions {
		regions[i] = *result.Regions[i].RegionName
	}

	if region == "" {
		// Randomly select a region which has the given droplet size.
		indexes := rand.Perm(len(result.Regions))
		for _, index := range indexes {
			region = *result.Regions[index].RegionName
			if err := p.instanceTypesAvailability(cred, region, instance); err == nil {
				return region, instance, nil
			}
		}
		return "", "", fmt.Errorf("selected instance type [%v] is not available across all regions", instance)
	} else {
		err = p.instanceTypesAvailability(cred, region, instance)
		return region, instance, err
	}
}

func (p providerAws) instanceTypesAvailability(cred *credentials.Credentials, region, instance string) error {
	instanceSession, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: cred,
	})
	if err != nil {
		return err
	}
	service := ec2.New(instanceSession)
	instanceInput := &ec2.DescribeInstanceTypesInput{
		InstanceTypes: []*string{aws.String(instance)},
	}
	instanceResult, err := service.DescribeInstanceTypes(instanceInput)
	if err != nil {
		return err
	}
	for _, res := range instanceResult.InstanceTypes {
		if *res.InstanceType == instance {
			return nil
		}
	}
	return fmt.Errorf("instance not avaliable")
}
