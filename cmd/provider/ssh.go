package provider

import (
	"github.com/renproject/darknode-cli/darknode"
	"github.com/renproject/darknode-cli/util"
	"github.com/urfave/cli"
)

type providerSsh struct {
	user          string
	hostname      string
	priKeyPath    string
}

func NewSsh(ctx *cli.Context) (Provider, error) {
	user := ctx.String("ssh-user")
	hostname := ctx.String("ssh-hostname")
	priKeyPath := ctx.String("ssh-private-key")

	return providerSsh{
		user: user,
		hostname: hostname,
		priKeyPath: priKeyPath,
	}, nil
}

func (p providerSsh) Name() string {
	return NameSsh
}

func (p providerSsh) Deploy(ctx *cli.Context) error {
	name := ctx.String("name")
	tags := ctx.String("tags")

	latestVersion, err := util.LatestStableRelease()
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
	if err := p.tfConfig(name, latestVersion); err != nil {
		return err
	}
	if err := runTerraform(name); err != nil {
		return err
	}
	return outputURL(name)
}
