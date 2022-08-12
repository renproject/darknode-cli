package util

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/google/go-github/v44/github"
	"github.com/hashicorp/go-version"
	"github.com/renproject/darknode-cli/darknode"
	"golang.org/x/oauth2"
)

// GithubClient initialize the Github client. If an access token has been set
// as an environment, it will use it for oauth to avoid rate limiting.
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

// RateLimit checks if we get rate-limited by the Github API. It will return
// how many remaining requests we can make before getting rate-limited
func RateLimit(ctx context.Context, client *github.Client) (int, error) {
	rl, response, err := client.RateLimits(ctx)
	if err != nil {
		return 0, err
	}
	if response.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("cannot get github API rate limit info")
	}
	return rl.Core.Remaining, nil
}

// LatestRelease fetches the name of the latest Darknode release of given
// network.
func LatestRelease(network darknode.Network) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check the rate limit status of github api
	client := GithubClient(ctx)
	remaining, err := RateLimit(ctx, client)
	if err != nil {
		return "", err
	}
	if remaining < 10 {
		return "", fmt.Errorf("rate limited by github API, please set the env `GITHUB_TOKEN` with a personal access token")
	}

	// Construct the regex for release name
	reg := regexp.MustCompile("^v?[0-9]+\\.[0-9]+\\.[0-9]+$")
	latest, _ := version.NewVersion("0.0.0")

	opts := &github.ListOptions{
		Page:    0,
		PerPage: 100, // 100 maximum
	}
	for {
		releases, response, err := client.Repositories.ListReleases(ctx, "renproject", "darknode-release", opts)
		if err != nil {
			return "", err
		}

		// Verify the status code is 200.
		if err := VerifyStatusCode(response.Response, http.StatusOK); err != nil {
			return "", err
		}

		// Find the latest release tag for the given network
		for _, release := range releases {
			if !reg.MatchString(*release.TagName) {
				continue
			}
			ver, err := version.NewVersion(*release.TagName)
			if err != nil {
				return "", err
			}
			if ver.GreaterThan(latest) {
				latest = ver
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
