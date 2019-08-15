package main

import (
	"context"
	"github.com/republicprotocol/republic-go/cmd/darknode/config"
	"github.com/urfave/cli"
	"golang.org/x/crypto/ssh"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path"
	"strings"
	"text/template"
)

//TODO list from API call to stay up to date - https://www.googleapis.com/compute/v1/projects/{PROJECT}/zones
// TEMP regeneration: $gcloud compute zones list --format="value(name)" | awk '{ print "\""$0"\","}'
var gcpZones = []string{
	"us-east1-b",
	"us-east1-c",
	"us-east1-d",
	"us-east4-c",
	"us-east4-b",
	"us-east4-a",
	"us-central1-c",
	"us-central1-a",
	"us-central1-f",
	"us-central1-b",
	"us-west1-b",
	"us-west1-c",
	"us-west1-a",
	"europe-west4-a",
	"europe-west4-b",
	"europe-west4-c",
	"europe-west1-b",
	"europe-west1-d",
	"europe-west1-c",
	"europe-west3-c",
	"europe-west3-a",
	"europe-west3-b",
	"europe-west2-c",
	"europe-west2-b",
	"europe-west2-a",
	"asia-east1-b",
	"asia-east1-a",
	"asia-east1-c",
	"asia-southeast1-b",
	"asia-southeast1-a",
	"asia-southeast1-c",
	"asia-northeast1-b",
	"asia-northeast1-c",
	"asia-northeast1-a",
	"asia-south1-c",
	"asia-south1-b",
	"asia-south1-a",
	"australia-southeast1-b",
	"australia-southeast1-c",
	"australia-southeast1-a",
	"southamerica-east1-b",
	"southamerica-east1-c",
	"southamerica-east1-a",
	"asia-east2-a",
	"asia-east2-b",
	"asia-east2-c",
	"asia-northeast2-a",
	"asia-northeast2-b",
	"asia-northeast2-c",
	"europe-north1-a",
	"europe-north1-b",
	"europe-north1-c",
	"europe-west6-a",
	"europe-west6-b",
	"europe-west6-c",
	"northamerica-northeast1-a",
	"northamerica-northeast1-b",
	"northamerica-northeast1-c",
	"us-west2-a",
	"us-west2-b",
	"us-west2-c",
}

//TODO get options from API to stay up to date
// $gcloud compute machine-types list --format="value(name)" --filter="zone=europe-west1-b" | awk '{ print "\""$0"\","}'
var gcpMachineTypes = []string{
	"f1-micro",
	"g1-small",
	"m1-megamem-96",
	"m1-ultramem-160",
	"m1-ultramem-40",
	"m1-ultramem-80",
	"n1-highcpu-16",
	"n1-highcpu-2",
	"n1-highcpu-32",
	"n1-highcpu-4",
	"n1-highcpu-64",
	"n1-highcpu-8",
	"n1-highcpu-96",
	"n1-highmem-16",
	"n1-highmem-2",
	"n1-highmem-32",
	"n1-highmem-4",
	"n1-highmem-64",
	"n1-highmem-8",
	"n1-highmem-96",
	"n1-megamem-96",
	"n1-standard-1",
	"n1-standard-16",
	"n1-standard-2",
	"n1-standard-32",
	"n1-standard-4",
	"n1-standard-64",
	"n1-standard-8",
	"n1-standard-96",
	"n1-ultramem-160",
	"n1-ultramem-40",
	"n1-ultramem-80",
}

// awsTerraform contains all the fields needed to generate a terraform config file
// so that we can deploy the node on AWS.
type gcpTerraform struct {
	Name           string
	Zone           string
	Address        string
	MachineType    string
	SshPubKey      string
	SshPriKeyPath  string
	credentialPath string
	Port           string
	Path           string
	AllocationID   string
}

func deployToGCP(ctx *cli.Context) error {
	zone := strings.ToLower(ctx.String("gcp-zone"))
	machine_type := strings.ToLower(ctx.String("gcp-machine-type"))

	log.Println("zone: " + zone)
	if zone == "" {
		zone = gcpZones[rand.Intn(len(gcpZones))]
	}

	if machine_type == "" {
		machine_type = "n1-standard-1"
	}

	credentialPath, err := gcpCredentials(ctx)
	if err != nil {
		return err
	}

	network := ctx.String("network")
	network = strings.ToLower(network)
	if network != "testnet" && network != "mainnet" {
		return ErrUnknownNetwork
	}

	// Create node main.directory
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
	rsaKey := config.Keystore.RsaKey
	if err := WriteSshKey(rsaKey.PrivateKey, nodePath); err != nil {
		return err
	}
	pubKey, err := ssh.NewPublicKey(&rsaKey.PublicKey)
	if err != nil {
		return err
	}

	// Generate terraform config and start deploying
	if err := gcpTerraformConfig(ctx, config, pubKey, credentialPath, zone, machine_type); err != nil {
		return err
	}
	if err := runTerraform(nodePath); err != nil {
		return err
	}

	return outputURL(nodePath, name, network, pubKey.Marshal())
}

func gcpCredentials(ctx *cli.Context) (string, error) {
	jsonPath := ctx.String("gcp-credentials")
	//check if file exists
	data, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		log.Fatal(err)
	}

	googleCtx := context.Background()
	_, credErr := google.CredentialsFromJSON(googleCtx, data, "https://www.googleapis.com/auth/compute	")
	if credErr != nil {
		log.Fatal(credErr)
	}
	log.Println("valid credentials found on path " + jsonPath)
	return jsonPath, nil

}

func gcpTerraformConfig(ctx *cli.Context, config config.Config, key ssh.PublicKey, credentialPath, zone, machine_type string) error {
	name := ctx.String("name")
	nodePath := nodePath(name)

	tf := gcpTerraform{
		Name:           name,
		Zone:           zone,
		Address:        config.Address.String(),
		MachineType:    machine_type,
		SshPubKey:      strings.TrimSpace(StringfySshPubkey(key)),
		SshPriKeyPath:  path.Join(nodePath, "ssh_keypair"),
		credentialPath: credentialPath,
		Port:           config.Port,
		Path:           Directory,
		AllocationID:   ctx.String("aws-elastic-ip"),
	}

	templateFile := path.Join(Directory, "instance", "gcp", "gcp.tmpl")
	t := template.Must(template.New("aws.tmpl").Funcs(template.FuncMap{}).ParseFiles(templateFile))
	tfFile, err := os.Create(path.Join(nodePath, "main.tf"))
	if err != nil {
		return err
	}

	return t.Execute(tfFile, tf)
}
