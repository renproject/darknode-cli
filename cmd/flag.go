package main

import "github.com/urfave/cli"

// General flags
var (
	NameFlag = cli.StringFlag{
		Name:  "name",
		Usage: "A unique human-readable `string` for identifying the Darknode",
	}
	TagsFlag = cli.StringFlag{
		Name:  "tags",
		Usage: "Multiple human-readable comma separated `strings` for identifying groups of Darknodes",
	}
	ScriptFlag = cli.StringFlag{
		Name:  "script",
		Usage: "path of the script file you want to run",
	}
	KeystoreFlag = cli.StringFlag{
		Name:  "keystore",
		Usage: "An optional keystore `file` that will be used for the Darknode",
	}
	PassphraseFlag = cli.StringFlag{
		Name:  "passphrase",
		Usage: "An optional `secret` for decrypting the keystore file",
	}
	ConfigFlag = cli.StringFlag{
		Name:  "config",
		Usage: "An optional configuration `file` for the Darknode",
	}
	NetworkFlag = cli.StringFlag{
		Name:  "network",
		Value: "mainnet",
		Usage: "Darkpool network of your node (default: mainnet)",
	}
	BranchFlag = cli.StringFlag{
		Name:  "branch, b",
		Usage: "Release `branch` used to update the software",
	}
	AddressFlag = cli.StringFlag{
		Name:  "address",
		Usage: "Ethereum address you want to withdraw the tokens to.",
	}
	UpdateConfigFlag = cli.BoolFlag{
		Name:  "config, c",
		Usage: "An optional configuration file used to update the configuration",
	}
	ForceFlag = cli.BoolFlag{
		Name:  "force, f",
		Usage: "Force destruction without interactive prompts",
	}
)

// AWS flags
var (
	AwsFlag = cli.BoolFlag{
		Name:  "aws",
		Usage: "AWS will be used to provision the Darknode",
	}
	AwsAccessKeyFlag = cli.StringFlag{
		Name:  "aws-access-key",
		Usage: "AWS access `key` for programmatic access",
	}
	AwsSecretKeyFlag = cli.StringFlag{
		Name:  "aws-secret-key",
		Usage: "AWS secret `key` for programmatic access",
	}
	AwsRegionFlag = cli.StringFlag{
		Name:  "aws-region",
		Usage: "An optional AWS region (default: random)",
	}
	AwsInstanceFlag = cli.StringFlag{
		Name:  "aws-instance",
		Value: "t3.micro",
		Usage: "An optional AWS EC2 instance type (default: t3.micro)",
	}
	AwsElasticIpFlag = cli.StringFlag{
		Name:  "aws-elastic-ip",
		Usage: "An optional allocation ID for an elastic IP address",
	}
	AwsProfileFlag = cli.StringFlag{
		Name:  "aws-profile",
		Value: "default",
		Usage: "Name of the profile containing the credentials",
	}
)

// Digital ocean flags
var (
	// Digital Ocean flags
	DoFlag = cli.BoolFlag{
		Name:  "do",
		Usage: "Digital Ocean will be used to provision the Darknode",
	}
	DoTokenFlag = cli.StringFlag{
		Name:  "do-token",
		Usage: "Digital ocean API token for programmatic access",
	}
	DoRegionFlag = cli.StringFlag{
		Name:  "do-region",
		Usage: "An optional digital ocean region (default: random)",
	}
	DoSizeFlag = cli.StringFlag{
		Name:  "do-droplet",
		Value: "s-1vcpu-1gb",
		Usage: "An optional digital ocean droplet size (default: s-1vcpu-1gb)",
	}
)

var (
	GcpFlag = cli.BoolFlag{
		Name:  "gcp",
		Usage: "Google Cloud Platform will be used to provision the Darknode",
	}
	GcpZoneFlag = cli.StringFlag{
		Name:  "gcp-zone",
		Usage: "An optional Google Cloud Zone (default: random)",
	}

	GcpCredFlag = cli.StringFlag{
		Name:  "gcp-credentials",
		Usage: "Service Account credential file (JSON) to be used",
	}
)
