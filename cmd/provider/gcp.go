package provider

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/renproject/darknode-cli/darknode"
	"github.com/renproject/darknode-cli/util"
	"github.com/urfave/cli"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/option"
)

var ErrInsufficientPermission = errors.New("insufficient permissions")

type GcpRegion struct {
	region string
	zones  []string
}

// AllGcpRegions
var AllGcpRegions = []GcpRegion{
	{"asia-east1", []string{"a", "b", "c"}},
	{"asia-east2", []string{"a", "b", "c"}},
	{"asia-northeast1", []string{"a", "b", "c"}},
	{"asia-northeast2", []string{"a", "b", "c"}},
	{"asia-south1", []string{"a", "b", "c"}},
	{"asia-southeast1", []string{"a", "b", "c"}},
	{"australia-southeast1", []string{"a", "b", "c"}},
	{"europe-north1", []string{"a", "b", "c"}},
	{"europe-west1", []string{"b", "c", "d"}},
	{"europe-west2", []string{"a", "b", "c"}},
	{"europe-west3", []string{"a", "b", "c"}},
	{"europe-west4", []string{"a", "b", "c"}},
	{"europe-west6", []string{"a", "b", "c"}},
	{"northamerica-northeast1", []string{"a", "b", "c"}},
	{"southamerica-east1", []string{"a", "b", "c"}},
	{"us-central1", []string{"a", "b", "c", "f"}},
	{"us-east1", []string{"b", "c", "d"}},
	{"us-east4", []string{"a", "b", "c"}},
	{"us-west1", []string{"a", "b", "c"}},
	{"us-west2", []string{"a", "b", "c"}},
}

type providerGcp struct {
	credFile string
}

func NewGcp(ctx *cli.Context) (Provider, error) {
	credFile := ctx.String("gcp-credentials")
	if _, err := os.Stat(credFile); err != nil {
		return nil, err
	}
	return providerGcp{
		credFile: credFile,
	}, nil
}

func (p providerGcp) Name() string {
	return NameGcp
}

func (p providerGcp) Deploy(ctx *cli.Context) error {
	name := ctx.String("name")
	tags := ctx.String("tags")

	latestVersion, err := util.LatestReleaseVersion()
	if err != nil {
		return err
	}
	projectID, err := p.projectID()
	if err != nil {
		return err
	}
	zone, machine, err := p.validateZoneAndMachine(ctx)
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
	if err := p.tfConfig(name, projectID, zone, machine, latestVersion); err != nil {
		return err
	}
	if err := runTerraform(name); err != nil {
		return err
	}

	return outputURL(name)
}

func (p providerGcp) projectID() (string, error) {
	data, err := ioutil.ReadFile(p.credFile)
	if err != nil {
		return "", err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	creds, error := google.CredentialsFromJSON(ctx, data, "https://www.googleapis.com/auth/cloud-platform")
	if error != nil {
		return "", err
	}
	service, err := cloudresourcemanager.NewService(ctx, option.WithCredentials(creds))
	if err != nil {
		return "", err
	}

	rb := &cloudresourcemanager.TestIamPermissionsRequest{
		Permissions: []string{"compute.instances.create", "compute.networks.create", "compute.firewalls.create"}, ForceSendFields: nil, NullFields: nil,
	}
	resp, err := service.Projects.TestIamPermissions(creds.ProjectID, rb).Context(ctx).Do()
	if err != nil {
		return "", err
	}
	if len(resp.Permissions) < 3 {
		return "", ErrInsufficientPermission
	}
	return creds.ProjectID, nil
}

func (p providerGcp) validateZoneAndMachine(ctx *cli.Context) (string, string, error) {
	zone := strings.ToLower(strings.TrimSpace(ctx.String("gcp-zone")))
	machine := strings.ToLower(strings.TrimSpace(ctx.String("gcp-machine")))

	// Select a random zone for user if they don't provide one.
	if zone == "" {
		region := AllGcpRegions[rand.Intn(len(AllGcpRegions))]
		return fmt.Sprintf("%v-%v", region.region, region.zones[rand.Intn(len(region.zones))]), machine, nil
	}

	// todo : validate the zone and machine type
	return zone, machine, nil
}
