package provider

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/digitalocean/godo"
	"github.com/renproject/darknode-cli/util"
	"github.com/urfave/cli"
)

type providerDo struct {
	token string
}

func NewDo(ctx *cli.Context) (Provider, error) {
	token := ctx.String("do-token")

	return providerDo{
		token: token,
	}, nil
}

func (p providerDo) Name() string {
	return NameDo
}

func (p providerDo) Deploy(ctx *cli.Context) error {
	// Validate all input params
	if err := validateCommonParams(ctx); err != nil {
		return err
	}

	name := ctx.String("name")
	region, droplet, err := p.validateRegionAndDroplet(ctx)
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
	if err := p.tfConfig(name, region.Slug, droplet, latestVersion); err != nil {
		return err
	}
	if err := runTerraform(name); err != nil {
		return err
	}
	return outputURL(name)
}

func (p providerDo) validateRegionAndDroplet(ctx *cli.Context) (godo.Region, string, error) {
	region := strings.ToLower(strings.TrimSpace(ctx.String("do-region")))
	droplet := strings.ToLower(strings.TrimSpace(ctx.String("do-droplet")))
	c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Fetch all available regions
	client := godo.NewFromToken(p.token)
	regions, response, err := client.Regions.List(c, nil)
	if err != nil {
		return godo.Region{}, "", err
	}
	if err := util.VerifyStatusCode(response.Response, http.StatusOK); err != nil {
		return godo.Region{}, "", err
	}
	if len(regions) == 0 {
		return godo.Region{}, "", fmt.Errorf("account has no available region")
	}

	// Validate the given region and droplet type. Will use a random region
	// if not specified.
	if region == "" {
		// Randomly select a region which has the given droplet size.
		indexes := rand.Perm(len(regions))
		for _, index := range indexes {
			if util.StringInSlice(droplet, regions[index].Sizes) {
				if regions[index].Available {
					return regions[index], droplet, nil
				}
			}
		}
		return godo.Region{}, "", fmt.Errorf("selected droplet [%v] not available across all regions", droplet)
	} else {
		for _, r := range regions {
			if r.Slug == region {
				if util.StringInSlice(droplet, r.Sizes) {
					return r, droplet, nil
				}
				return godo.Region{}, "", fmt.Errorf("selected droplet [%v] not available in region %v", droplet, region)
			}
		}
		return godo.Region{}, "", fmt.Errorf("region [%v] is not avaliable", region)
	}
}
