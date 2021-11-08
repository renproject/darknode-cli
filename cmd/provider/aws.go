package provider

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/fatih/color"
	"github.com/renproject/darknode-cli/darknode"
	"github.com/renproject/darknode-cli/util"
	"github.com/urfave/cli"
)

// All available regions on AWS.
var AllAwsRegions = []string{
	"ap-northeast-1",
	"ap-northeast-2",
	"ap-south-1",
	"ap-southeast-1",
	"ap-southeast-2",
	"ca-central-1",
	"eu-central-1",
	"eu-north-1",
	"eu-west-1",
	"eu-west-2",
	"eu-west-3",
	"sa-east-1",
	"us-east-1",
	"us-east-2",
	"us-west-1",
	"us-west-2",
}

type providerAws struct {
	accessKey string
	secretKey string
}

func NewAws(ctx *cli.Context) (Provider, error) {
	accessKey := ctx.String("aws-access-key")
	secretKey := ctx.String("aws-secret-key")

	// Try reading the credential files if user does not provide credentials directly
	if accessKey == "" || secretKey == "" {
		cred := credentials.NewSharedCredentials("", ctx.String("aws-profile"))
		credValue, err := cred.Get()
		if err != nil {
			return nil, err
		}
		accessKey, secretKey = credValue.AccessKeyID, credValue.SecretAccessKey
		if accessKey == "" || secretKey == "" {
			return nil, err
		}
	}

	return providerAws{
		accessKey: accessKey,
		secretKey: secretKey,
	}, nil
}

func (p providerAws) Name() string {
	return NameAws
}

func (p providerAws) Deploy(ctx *cli.Context) error {
	name := ctx.String("name")
	tags := ctx.String("tags")
	config := ctx.String("config")

	latestVersion, err := util.LatestStableRelease()
	if err != nil {
		return err
	}
	region, instance, err := p.validateRegionAndInstance(ctx)
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
	if err := p.tfConfig(name, region, instance, latestVersion); err != nil {
		return err
	}
	if err := runTerraform(name); err != nil {
		return err
	}

	return outputURL(name)
}

func (p providerAws) DeployMultiple(ctx *cli.Context) error {
	prefix := ctx.String("name")
	tags := ctx.String("tags")
	n := ctx.Int("n")

	if n <= 0 {
		return fmt.Errorf("invalid n, must be possitive number")
	}

	regions, err := p.AvailableRegions()
	if err != nil {
		return err
	}

	ch := make(chan awsTerraform)
	wg := new(sync.WaitGroup)

	// Start 20 background workers
	for i := 0; i < 5; i++ {
		go func() {
			for {
				func() {
					config := <-ch
					name := config.Name
					defer wg.Done()

					// Initialization
					network, err := darknode.NewNetwork(ctx.String("network"))
					if err != nil {
						color.Red("invalid network, %v", err)
						return
					}
					if err := initNode(name, tags, network, ""); err != nil {
						color.Red("failed to initialize darknode, %v", err)
						return
					}

					t, err := template.New("aws").Parse(awsTemplate)
					if err != nil {
						color.Red("failed to initialize aws template, %v", err)
						return
					}
					tfFile, err := os.Create(filepath.Join(util.NodePath(name), "main.tf"))
					if err != nil {
						color.Red("failed to create terraform config, %v", err)
						return
					}
					if err := t.Execute(tfFile, config); err != nil {
						color.Red("failed to execute terraform, %v", err)
						return
					}

					if err := runTerraformSilent(name); err != nil {
						color.Red("failed to create terraform config, %v", err)
						return
					}
				}()
			}
		}()
	}

	latestVersion, err := util.LatestStableRelease()
	if err != nil {
		return err
	}

	// Found the starting index
	startIndex := 1
	for ; startIndex <= 1000; startIndex ++ {
		name := fmt.Sprintf("%v-%v", prefix, startIndex)
		if err := util.ValidateNodeName(name); err != nil {
			break
		}
	}
	if startIndex == 999 {
		return fmt.Errorf("try using a different prefix for your darknodes")
	}

	for i := startIndex; i < startIndex + n ; i++ {
		name := fmt.Sprintf("%v-%v", prefix, i)
		tf := awsTerraform{
			Name:          name,
			Region:        regions[i%len(regions)],
			InstanceType:  "t3.micro",
			ConfigPath:    fmt.Sprintf("~/.darknode/darknodes/%v/config.json", name),
			PubKeyPath:    fmt.Sprintf("~/.darknode/darknodes/%v/ssh_keypair.pub", name),
			PriKeyPath:    fmt.Sprintf("~/.darknode/darknodes/%v/ssh_keypair", name),
			AccessKey:     p.accessKey,
			SecretKey:     p.secretKey,
			ServiceFile:   darknodeService,
			LatestVersion: latestVersion,
		}
		ch <- tf
		wg.Add(1)
		color.Green("Deploying %v", name)
	}

	wg.Wait()
	return nil
}

func (p providerAws) validateRegionAndInstance(ctx *cli.Context) (string, string, error) {
	// TODO : use aws api to validate the region and instance type
	region := strings.ToLower(strings.TrimSpace(ctx.String("aws-region")))
	instance := strings.ToLower(strings.TrimSpace(ctx.String("aws-instance")))

	if region == "" {
		region = AllAwsRegions[rand.Intn(len(AllAwsRegions))]
	}
	if !util.StringInSlice(region, AllAwsRegions) {
		return "", "", fmt.Errorf("aws region [%v] is not supported yet", region)
	}

	if instance == "" {
		return "", "", fmt.Errorf("instance type [%v] is not supported yet", instance)
	}

	return region, instance, nil
}

// AvailableRegions returns all available regions
func (p providerAws) AvailableRegions() ([]string, error) {
	cred := credentials.NewStaticCredentials(p.accessKey, p.secretKey, "")
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: cred,
	})
	service := ec2.New(sess)
	input := &ec2.DescribeRegionsInput{}
	result, err := service.DescribeRegions(input)
	if err != nil {
		return nil, err
	}
	regions := make([]string, 0, len(result.Regions))
	for i := range result.Regions {
		regions = append(regions, *result.Regions[i].RegionName)
	}
	return regions, nil
}
