package provider

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/renproject/darknode-cli/util"
	"github.com/urfave/cli"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

var ErrInsufficientPermission = errors.New("insufficient permissions")

type GcpRegion struct {
	region string
	zones  []string
}

type providerGcp struct {
	credFile string
}

func NewGcp(ctx *cli.Context) (Provider, error) {
	credFile, err := filepath.Abs(ctx.String("gcp-credentials"))
	if err != nil {
		return nil, err
	}
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
	// Validate all input params
	if err := validateCommonParams(ctx); err != nil {
		return err
	}

	// Validate name for GCP
	name := ctx.String("name")
	reg := "^[a-z]([-a-z0-9]{0,61}[a-z0-9])?$"
	match, err := regexp.MatchString(reg, name)
	if err != nil {
		return err
	}
	if !match {
		return errors.New("for google cloud, name must start with a lowercase letter followed by up to 62 lowercase letters, numbers, or hyphens, and cannot end with a hyphen")
	}

	latestVersion, err := util.LatestStableRelease()
	if err != nil {
		return err
	}
	projectID, err := p.projectID()
	if err != nil {
		return err
	}
	region, machine, err := p.validateRegionAndMachine(ctx)
	if err != nil {
		return err
	}

	// Initialize folder and files for the node
	if err := initNode(ctx); err != nil {
		return err
	}
	// Generate terraform config and start deploying
	if err := p.tfConfig(name, projectID, region, machine, latestVersion); err != nil {
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
	creds, err := google.CredentialsFromJSON(ctx, data, "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
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

func (p providerGcp) validateRegionAndMachine(ctx *cli.Context) (string, string, error) {
	region := strings.ToLower(strings.TrimSpace(ctx.String("gcp-region")))
	machine := strings.ToLower(strings.TrimSpace(ctx.String("gcp-machine")))

	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Retrieve project ID and initialize the computer service
	fileName, err := filepath.Abs(ctx.String("gcp-credentials"))
	if err != nil {
		return "", "", err
	}
	projectID, err := getProjectID(c, fileName)
	if err != nil {
		return "", "", err
	}
	service, err := compute.NewService(c, option.WithCredentialsFile(fileName))
	if err != nil {
		return "", "", err
	}

	// Select a random zone for user if they don't provide one.
	var zone string
	if region == "" {
		zone, err = randomZone(c, service, projectID)
		if err != nil {
			return "", "", err
		}
	} else {
		zone, err = validateZone(c, service, projectID, region)
		if err != nil {
			return "", "", err
		}
	}

	// Validate if the machine type is available in the zone.
	if err := validateMachineType(c, service, projectID, zone, machine); err != nil {
		return "", "", err
	}

	return zone, machine, nil
}

// Randomly pick a available zone
func randomZone(ctx context.Context, service *compute.Service, projectID string) (string, error) {
	// List all the available zones
	zones := []string{}
	req := service.Zones.List(projectID)
	if err := req.Pages(ctx, func(page *compute.ZoneList) error {
		for _, zone := range page.Items {
			zones = append(zones, zone.Name)
		}
		return nil
	}); err != nil {
		return "", err
	}
	if len(zones) == 0 {
		return "", errors.New("cannot find any available zone under your project")
	}

	return zones[rand.Intn(len(zones))], nil
}

func validateZone(ctx context.Context, service *compute.Service, projectID, input string) (string, error) {
	reg := regexp.MustCompile("^(?P<region>[a-z]+-[a-z]+[1-9])(-(?P<zone>[a-g]))?$")
	if !reg.MatchString(input) {
		return "", errors.New("invalid region name")
	}

	values := util.CaptureGroups("^(?P<region>[a-z]+-[a-z]+[1-9])(-(?P<zone>[a-g]))?$", input)
	zone := values["zone"]

	// If user doesn't provide a zone in the region, we only need to validate the
	// region and randomly select a zone in the region.
	if zone == "" {
		availableRegion, err := service.Regions.Get(projectID, input).Context(ctx).Do()
		if err != nil {
			return "", err
		}

		zone = path.Base(availableRegion.Zones[rand.Intn(len(availableRegion.Zones))])
		return zone, nil
	}

	// If user gives both region and zone, we only need to validate the zone.
	_, err := service.Zones.Get(projectID, input).Context(ctx).Do()
	return input, err
}

func validateMachineType(ctx context.Context, service *compute.Service, projectID, zone, machine string) error {
	// List all the available zones
	res, err := service.MachineTypes.Get(projectID, zone, machine).Context(ctx).Do()
	if err != nil {
		return err
	}
	if res.Name != machine {
		return fmt.Errorf("%v type is not available from your project", machine)
	}
	return nil
}

func getProjectID(ctx context.Context, fileName string) (string, error) {
	// Parse the credential file
	credFile, err := filepath.Abs(fileName)
	if err != nil {
		return "", err
	}
	jsonFile, err := os.Open(credFile)
	if err != nil {
		return "", err
	}
	defer jsonFile.Close()
	data, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return "", err
	}

	// Get the project ID
	credentials, err := google.CredentialsFromJSON(ctx, data, compute.ComputeScope)
	if err != nil {
		return "", err
	}
	return credentials.ProjectID, nil
}
