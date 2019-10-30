package provider

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/renproject/darknode-cli/darknode"
	"github.com/urfave/cli"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/option"
)

var ErrInsufficientPermission = errors.New("insufficient permissions")

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
	if err := p.tfConfig(name, projectID, zone, machine, ipfsUrl(network)); err != nil {
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
	log.Print("project id =", creds.ProjectID)
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
	log.Printf("number of permissions = %v", resp.Permissions)
	if len(resp.Permissions) < 3 {
		return "", ErrInsufficientPermission
	}
	return creds.ProjectID, nil
}

func (p providerGcp) validateZoneAndMachine(ctx *cli.Context) (string, string, error) {
	zone := strings.ToLower(strings.TrimSpace(ctx.String("gcp-zone")))
	machine := strings.ToLower(strings.TrimSpace(ctx.String("gcp-machine")))

	// todo : validate the zone and machine type
	return zone, machine, nil
}
