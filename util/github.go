package util

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/google/go-github/v39/github"
	"github.com/hashicorp/go-version"
	"golang.org/x/oauth2"
)

// GithubClient initialize the github client. If an access token has been set as an environment,
// it will use it for oauth to avoid rate limiting.
func GithubClient(ctx context.Context) *github.Client {
	accessToken := os.Getenv("GITHUB_TOKEN")
	var client *http.Client
	if accessToken != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: accessToken},
		)
		client = oauth2.NewClient(ctx, ts)
	}

	return github.NewClient(client)
}

// LatestStableRelease checks the darknode release repo and return the version
// of the latest release.
func LatestStableRelease() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := github.NewClient(nil)
	opts := &github.ListOptions{
		PerPage: 50,
	}
	latest, err := version.NewVersion("0.0.0")
	if err != nil {
		return "", err
	}

	// Fetch all releases and find the latest stable release tag
	for {
		releases, response, err := client.Repositories.ListReleases(ctx, "renproject", "darknode-release", opts)
		if err != nil {
			return "", err
		}

		if response.StatusCode != http.StatusOK {
			return "", fmt.Errorf("cannot get latest darknode release from github, error code = %v", response.StatusCode)
		}

		verReg := "^v?[0-9]+\\.[0-9]+\\.[0-9]+$"
		for _, release := range releases {
			match, err := regexp.MatchString(verReg, *release.TagName)
			if err != nil {
				return "", err
			}
			if match {
				ver, err := version.NewVersion(*release.TagName)
				if err != nil {
					return "", err
				}
				if ver.GreaterThan(latest) {
					latest = ver
				}
			}
		}
		if response.NextPage == 0 {
			break
		}
		opts.Page = response.NextPage
	}
	if latest.String() == "0.0.0" {
		return "", errors.New("cannot find any stable release")
	}

	return latest.String(), nil
}
