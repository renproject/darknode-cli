package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path"
	"text/template"

	"github.com/republicprotocol/republic-go/cmd/darknode/config"
	"github.com/urfave/cli"
)

// Available regions on Digital Ocean.
var AllDoRegions = []string{
	"ams2",
	"ams3",
	"blr1",
	"fra1",
	"lon1",
	"nyc1",
	"nyc2",
	"nyc3",
	"sfo1",
	"sfo2",
	"sgp1",
	"tor1",
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

// Region is the json object returned by the digital-ocean API
type Region struct {
	Name      string   `json:"name"`
	Slug      string   `json:"slug"`
	Sizes     []string `json:"sizes"`
	Features  []string `json:"features"`
	Available bool     `json:"available"`
}

// doTerraform contains all the fields needed to generate a terraform config file
// so that we can deploy the node on Digital Ocean.
type doTerraform struct {
	Name    string
	Token   string
	Region  string
	Address string
	Size    string
	Path    string
	PubKey  string
	PvtKey  string
}

// deployToDo parses the digital ocean credentials and use terraform to
// deploy the node to digital ocean.
func deployToDo(ctx *cli.Context) error {
	token := ctx.String("do-token")
	if token == "" {
		return ErrEmptyDoToken
	}
	region, size, err := doRegionAndDroplet(ctx)
	if err != nil {
		return err
	}

	// Create node directory
	name := ctx.String("name")
	tags := ctx.String("tags")
	if err := mkdir(name, tags); err != nil {
		return err
	}
	nodePath := nodePath(name)

	// Generate config and ssh key for the node
	config, err := GetConfigOrGenerateNew(ctx, nodePath)
	if err != nil {
		return err
	}
	key, err := NewSshKeyPair(nodePath)
	if err != nil {
		return err
	}

	// Generate terraform config and start deploying
	if err := generateDoTFConfig(ctx, config, region, size); err != nil {
		return err
	}
	if err := runTerraform(nodePath); err != nil {
		return err
	}

	return outputURL(nodePath, name, key.Marshal())
}

func doRegionAndDroplet(ctx *cli.Context) (string, string, error) {
	region := ctx.String("do-region")
	droplet := ctx.String("do-droplet")

	regions, err := availableRegions(ctx)
	if err != nil {
		return "", "", err
	}

	// Parse the input region or pick one region randomly
	if region == "" {
		if len(regions) == 0 {
			return "", "", ErrNoAvailableRegion
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
			return "", "", ErrUnknownRegion
		}
		return chosenRegion.Slug, droplet, validateDroplet(droplet, chosenRegion.Slug, chosenRegion.Sizes)
	}
}

// validateDroplet validates whether the droplet is available in the region.
func validateDroplet(droplet, region string, droplets []string) error {
	if !StringInSlice(droplet, droplets) {
		fmt.Printf("[%v] is the selected droplet region.\n", region)
		fmt.Printf("Your account can only create below slugs in [%v]:\n", region)
		for i := range droplets {
			fmt.Println(droplets[i])
		}
		fmt.Println("You can find more details about these slugs from https://www.digitalocean.com/pricing")
		return ErrUnSupportedInstanceType
	}
	return nil
}

// availableRegions sends a GET request to the DO API to get all available
// regions and droplet sizes of the given DO token.
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

// generateDoTFConfig generates the terraform config file for deploying to DO.
func generateDoTFConfig(ctx *cli.Context, config config.Config, region, size string) error {
	name := ctx.String("name")
	token := ctx.String("do-token")
	nodePath := nodePath(name)

	tf := doTerraform{
		Name:    name,
		Token:   token,
		Region:  region,
		Address: config.Address.String(),
		Size:    size,
		Path:    Directory,
		PubKey:  path.Join(nodePath, "ssh_keypair.pub"),
		PvtKey:  path.Join(nodePath, "ssh_keypair"),
	}

	templateFile := path.Join(Directory, "instance", "do", "do.tmpl")
	t := template.Must(template.New("do.tmpl").Funcs(template.FuncMap{}).ParseFiles(templateFile))
	tfFile, err := os.Create(path.Join(nodePath, "main.tf"))
	if err != nil {
		return err
	}

	return t.Execute(tfFile, tf)
}
