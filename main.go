package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/github"
	log "github.com/sirupsen/logrus"

	"github.com/hostwithquantum/github-org-sync-action/utils"

	"github.com/hostwithquantum/github-org-sync-action/repo"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {
	org := os.Getenv("GITHUB_ORG")
	repositories := strings.Split(os.Getenv("GITHUB_REPOS"), " ")

	// this is the main copy which contains the files we want to sync
	skeletonRepository := os.Getenv("GITHUB_SKELETON")

	log.Info(fmt.Sprintf("Running for '%s' and syncing to: %s", org, repositories))
	log.Info(fmt.Sprintf("From: %s", skeletonRepository))

	// init CurrentUser (for auth)
	currentUser := repo.CurrentUser{
		Email: os.Getenv("GITHUB_EMAIL"),
		Name:  os.Getenv("GITHUB_USER"),
		Token: os.Getenv("GITHUB_ACCESS_TOKEN"),
	}

	githubClient := repo.NewGithub(currentUser, org)

	log.Info(fmt.Sprintf("Using: %s (of %s)", currentUser.Email, currentUser.Name))

	tmpDirectory := fmt.Sprintf("./tmp/%s", org)

	log.Info(fmt.Sprintf("Temp files will be created in: %s", tmpDirectory))

	// clone skeleton repository
	skeleton := repo.NewRepo(org, skeletonRepository, currentUser, tmpDirectory)

	skeletonGithub := utils.GithubLink(org, skeletonRepository)

	log.Info("Cloned skeleton")

	handler := repo.NewHandler(skeletonRepository, tmpDirectory)

	for _, repository := range repositories {

		target := repo.NewRepo(org, repository, currentUser, tmpDirectory)
		githubLink := utils.GithubLink(org, repository)
		log.Infof("Cloned: '%s'", githubLink)

		defaultBranch := target.GetDefaultBranch()
		log.Debugf("Default branch for '%s' is '%s'", githubLink, defaultBranch)

		target.CreateBranch("chore/update-workflows")

		handler.Sync(fmt.Sprintf("%s/%s", tmpDirectory, repository))

		if target.NeedsCommit() != true {
			log.Info(fmt.Sprintf("Repo '%s' is clean\n", githubLink))
			continue
		}

		commitMsg := fmt.Sprintf(
			"Chore: updates to .github/workflows\n\nSee %s@%s",
			skeletonGithub,
			skeleton.GetLastCommit(),
		)
		target.CommitAndPush(commitMsg, "chore/update-workflows")

		description := fmt.Sprintf(
			"Updating from %s: %s@%s",
			skeletonGithub,
			skeletonGithub,
			target.GetLastCommit())

		newPR := &github.NewPullRequest{
			Title:               github.String(fmt.Sprintf("[github-org-sync] Updating files via %s", skeletonGithub)),
			Head:                github.String("chore/update-workflows"),
			Base:                github.String(defaultBranch),
			Body:                github.String(description),
			MaintainerCanModify: github.Bool(true),
		}

		githubClient.CreatePullRequest(repository, newPR)
	}

	err := os.RemoveAll(tmpDirectory)
	utils.CheckIfError(err)
}
