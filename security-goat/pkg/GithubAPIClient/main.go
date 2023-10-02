package GithubAPIClient

import (
	"strings"

	"github.com/google/go-github/v55/github"
)

func CreateGitHubClient(token string) *github.Client {
	// https://github.com/google/go-github/blob/master/example/tokenauth/main.go
	return github.NewClient(nil).WithAuthToken(token)
}

func CreateGitHubClientBasicAuth(username string, password string) *github.Client {
	tp := github.BasicAuthTransport{
		Username: strings.TrimSpace(username),
		Password: strings.TrimSpace(password),
	}

	return github.NewClient(tp.Client())
}
