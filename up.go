package main

import (
	"strings"

	"errors"
	"github.com/urfave/cli"
	"math/rand"
	"os/exec"
)

// KeyNotFound is returned when no AWS access-key nor secret-key provided.
var KeyNotFound error = errors.New("please provide your AWS access key and secret key")

// UnknownRegion is returned when the provided region is not valid on AWS.
var UnknownRegion error = errors.New("there is no such region on AWS")

// UnSupportedInstanceType is returned when the provided instance is not
// supported in the selected region.
var UnSupportedInstanceType error = errors.New("instance type is not supported in the region")

// deployNode deploys node depending on the provider.
func deployNode(ctx *cli.Context, path string) error {
	provider := strings.ToLower(ctx.String("provider"))
	switch provider {
	case "aws":
		return deployToAWS(ctx, path)
	case "digital-ocean":
		return deployToDigitalOcean(ctx)
	default:
		return errors.New("unsupported service provider")
	}
}

// deployToAWS parses the AWS credentials and use terraform to deploy the node
// to AWS.
func deployToAWS(ctx *cli.Context, path string) error {
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
	config, err := NewConfig()
	if err != nil {
		return err
	}

	return nil
}

// deployToDigitalOcean parses the digital ocean credentials and use terraform
// to deploy the node to digital ocean.
func deployToDigitalOcean(ctx *cli.Context) error {
	panic("unimplemented")
}

// parseRegionAndInstance parses the region and the instance type from the
// cli parameters. It will randomly pick a region for the user if it's not
// specified. The default value for instance is `t2.small`.
func parseRegionAndInstance(ctx *cli.Context) (string, string, error) {
	region := strings.ToLower(ctx.String("region"))
	instance := strings.ToLower(ctx.String("instance"))

	// Parse the input region or pick one region randomly
	if region == "" {
		region = string(AllAwsRegions[rand.Intn(len(AllAwsRegions))])
	} else {
		if !stringInSlice(region, AllAwsRegions) {
			return "", "", UnknownRegion
		}
	}

	// Parse the input instance type or use the default one.
	if instance == "" {
		instance = "t2.small"
	} else {
		if region == EuWest3 && !stringInSlice(instance, AllAwsInstancesInEuWest3) {
			return "", "", UnSupportedInstanceType
		}
		if region == ApNorthEast1 && !stringInSlice(instance, AllAwsInstancesInApNortheast1) {
			return "", "", UnSupportedInstanceType
		}
		if !stringInSlice(instance, AllAwsInstances) {
			return "", "", UnSupportedInstanceType
		}
	}

	return region, instance, nil
}

// runTerraform initializes and applies terraform
func runTerraform() error {
	init := exec.Command("./terraform", "init")
	if err := init.Run(); err != nil {
		return err
	}
	if err := init.Wait(); err != nil {
		return err
	}

	apply := exec.Command("./terraform", "apply", "-auto-approve")
	if err := apply.Run(); err != nil {
		return err
	}
	return apply.Wait()
}
