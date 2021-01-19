package provider

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"

	"github.com/renproject/darknode-cli/darknode"
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
	name := ctx.String("name")
	tags := ctx.String("tags")
	config := ctx.String("config")

	latestVersion, err := util.LatestStableRelease()
	if err != nil {
		return err
	}
	region, droplet, err := validateRegionAndDroplet(ctx)
	if err != nil {
		return err
	}

	// Initialization
	network, err := darknode.NewNetwork(ctx.String("network"))
	if err != nil {
		return err
	}
	if err := initNode(name, tags, network, config); err != nil {
		return err
	}

	// Generate terraform config and start deploying
	if err := p.tfConfig(name, region, droplet, latestVersion); err != nil {
		return err
	}
	if err := runTerraform(name); err != nil {
		return err
	}
	return outputURL(name)
}

func validateRegionAndDroplet(ctx *cli.Context) (string, string, error) {
	region := ctx.String("do-region")
	droplet := ctx.String("do-droplet")

	regions, err := availableRegions(ctx)
	if err != nil {
		return "", "", err
	}

	// Parse the input region or pick one region randomly
	if region == "" {
		if len(regions) == 0 {
			return "", "", ErrRegionNotAvailable
		}
		randomRegion := regions[rand.Intn(len(regions))]
		return randomRegion.Slug, droplet, validateDroplet(droplet, randomRegion.Slug, randomRegion.Sizes)
	} else {
		var chosenRegion Region
		for i := range regions {
			if region == regions[i].Slug {
				chosenRegion = regions[i]
				break
			}
		}
		if chosenRegion.Name == "" {
			return "", "", ErrRegionNotAvailable
		}
		return chosenRegion.Slug, droplet, validateDroplet(droplet, chosenRegion.Slug, chosenRegion.Sizes)
	}
}

// validateDroplet validates whether the droplet is available in the region.
func validateDroplet(droplet, region string, droplets []string) error {
	if !util.StringInSlice(droplet, droplets) {
		fmt.Printf("[%v] is the selected droplet region.\n", region)
		fmt.Printf("Your account can only create below slugs in [%v]:\n", region)
		for i := range droplets {
			fmt.Println(droplets[i])
		}
		fmt.Println("You can find more details about these slugs from https://www.digitalocean.com/pricing")
		return ErrInstanceTypeNotAvailable
	}
	return nil
}

// Region is the json object returned by the digital-ocean API
type Region struct {
	Name      string   `json:"name"`
	Slug      string   `json:"slug"`
	Sizes     []string `json:"sizes"`
	Features  []string `json:"features"`
	Available bool     `json:"available"`
}

// availableRegions sends a GET request to Digital Ocean API to get all available regions and droplet sizes of the given
// Digital Ocean token.
func availableRegions(ctx *cli.Context) ([]Region, error) {
	token := ctx.String("do-token")

	url := "https://api.digitalocean.com/v2/regions"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// Check the response status code
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(data))
	}

	// Unmarshal the response
	regions := struct {
		Regions []Region `json:"regions"`
	}{}
	err = json.Unmarshal(data, &regions)
	if err != nil {
		return nil, err
	}
	availableRegions := make([]Region, 0)
	for _, i := range regions.Regions {
		if i.Available {
			availableRegions = append(availableRegions, i)
		}
	}
	return availableRegions, nil
}
