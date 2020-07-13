package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-github/github"
	log "github.com/sirupsen/logrus"

	"github.com/hostwithquantum/github-org-sync-action/repo"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

var gitHubActionVersion = "__REPLACED__"

var motd = `
 _______________________
< github-org-sync %s >
 -----------------------
`

func main() {
	org := os.Getenv("GITHUB_ORG")
	repositories := strings.Split(os.Getenv("GITHUB_REPOS"), " ")

	// this is the main copy which contains the files we want to sync
	skeletonRepository := os.Getenv("GITHUB_SKELETON")

	fmt.Printf(motd, gitHubActionVersion)

	log.Infof("Running for '%s' and syncing to: %s", org, repositories)
	log.Infof("From: %s", skeletonRepository)

	_, dryRun := os.LookupEnv("DRY_RUN")

	// init CurrentUser (for auth)
	currentUser := repo.CurrentUser{
		Email: os.Getenv("GITHUB_EMAIL"),
		Name:  os.Getenv("GITHUB_USER"),
		Token: os.Getenv("GITHUB_ACCESS_TOKEN"),
	}

	githubClient := repo.NewGithub(currentUser, org)

	log.Infof("Using: %s (of %s)", currentUser.Email, currentUser.Name)

	tmpDirectory := fmt.Sprintf("./tmp/%s", org)

	log.Infof("Temp files will be created in: %s", tmpDirectory)

	// clone skeleton repository
	skeleton := repo.NewRepo(org, skeletonRepository, currentUser, tmpDirectory)

	skeletonGithub := repo.GithubLink(org, skeletonRepository)

	log.Info("Cloned skeleton")

	handler := repo.NewHandler(skeletonRepository, tmpDirectory)

	for _, repository := range repositories {

		target := repo.NewRepo(org, repository, currentUser, tmpDirectory)
		githubLink := repo.GithubLink(org, repository)
		log.Infof("Cloned: '%s'", githubLink)

		defaultBranch := target.GetDefaultBranch()
		log.Debugf("Default branch for '%s' is '%s'", githubLink, defaultBranch)

		target.CreateBranch("chore/update-workflows")

		// sync .github/workflows
		var workflowDirTarget = filepath.Join(tmpDirectory, repository, ".github", "workflows")
		handler.Sync(
			workflowDirTarget,
			handler.Workflows,
		)

		handler.Sync(
			filepath.Join(tmpDirectory, repository),
			handler.Files,
		)

		if target.NeedsCommit() != true {
			log.Infof("Repo '%s' is clean\n", githubLink)
			continue
		}

		if dryRun {
			log.Info("[dry-run] Not committing, pushing or opening a PR.")
			continue
		}

		commitMsg := fmt.Sprintf(
			"Chore: automated updates to .github/workflows, templates, ...\n\nSee %s@%s",
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

	if !dryRun {
		err := os.RemoveAll(tmpDirectory)
		repo.CheckIfError(err)
	} else {
		log.Info("[dry-run] Leaving artifacts for inspection")
	}
}
