package repo

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
	"github.com/hostwithquantum/github-org-sync-action/utils"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

// Github ...
type Github struct {
	user    CurrentUser
	org     string
	context context.Context
	client  *github.Client
}

// NewGithub ...
func NewGithub(user CurrentUser, org string) Github {
	// init context
	ctx := context.Background()

	o := Github{
		user:    user,
		org:     org,
		context: ctx,
		client:  createAuthorizedClient(ctx, user.Token),
	}
	return o
}

// CreatePullRequest ...
func (g Github) CreatePullRequest(repository string, pr *github.NewPullRequest) {
	pullRequest, response, err := g.client.PullRequests.Create(
		g.context,
		g.org,
		repository,
		pr,
	)
	if response.StatusCode == 422 {
		log.Info("Pull request already exists. Calling it a win!")
		return
	}

	// handle other errors
	utils.CheckIfError(err)

	log.Info(fmt.Sprintf("PR created: %s", pullRequest.GetHTMLURL()))
}

func createAuthorizedClient(context context.Context, token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: token,
		},
	)
	tc := oauth2.NewClient(context, ts)
	client := github.NewClient(tc)
	return client
}
