package repo

import (
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

// Repo ...
type Repo struct {
	org    string
	name   string
	user   CurrentUser
	target string
	repo   *git.Repository
}

// NewRepo ...
func NewRepo(org string, name string, user CurrentUser, tmpDirectory string) Repo {

	cloneDirectory := fmt.Sprintf("%s/%s", tmpDirectory, name)

	repo := clone(org, name, cloneDirectory, user.Token)

	o := Repo{org, name, user, cloneDirectory, repo}
	return o
}

// CommitAndPush ...
func (r Repo) CommitAndPush(message string, branch string) {
	var hashes []plumbing.Hash

	branchRef := plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch))

	worktree, err := r.repo.Worktree()
	CheckIfError(err)

	status, err := worktree.Status()
	CheckIfError(err)

	log.Debugf("Current status is: %v", status.IsClean())
	if status.IsClean() {
		log.Info("Skipping commit: no changes")
		return
	}

	for file := range status {
		hash, err := worktree.Add(file)
		CheckIfError(err)

		log.Debugf("Staged file: %s", file)

		hashes = append(hashes, hash)
	}

	commit, err := worktree.Commit(message, &git.CommitOptions{
		All: true,
		Author: &object.Signature{
			Name:  r.user.Name,
			Email: r.user.Email,
			When:  time.Now(),
		},
	})
	CheckIfError(err)
	log.Debugf("Created commit: %s", commit.String())

	err = r.repo.Push(&git.PushOptions{
		Auth:       authMethod(r.user.Token),
		Force:      true,
		RemoteName: "origin",
		RefSpecs: []config.RefSpec{
			config.RefSpec(branchRef + ":" + branchRef),
		},
	})
	CheckIfError(err)
	log.Debugf("Pushed branch: %s!", branchRef)
}

// CreateBranch ...
func (r Repo) CreateBranch(branch string) {
	worktree, err := r.repo.Worktree()
	CheckIfError(err)

	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branch),
		Create: true,
	})
	CheckIfError(err)

	log.Debugf("Created new branch and switched workspace")
}

// GetDefaultBranch ...
func (r Repo) GetDefaultBranch() string {
	ref, err := r.repo.Head()
	CheckIfError(err)

	return ref.Name().Short()
}

// GetLastCommit ...
func (r Repo) GetLastCommit() plumbing.Hash {
	headRef, err := r.repo.Head()
	CheckIfError(err)
	log.Debugf("Last commit in '%s' is: %s", GithubLink(r.org, r.name), headRef.Hash())
	return headRef.Hash()
}

// NeedsCommit ...
func (r Repo) NeedsCommit() bool {
	worktree, err := r.repo.Worktree()
	CheckIfError(err)

	status, err := worktree.Status()
	CheckIfError(err)

	log.Debugf("Repo clean?: %v", status.IsClean())

	if status.IsClean() {
		return false
	}
	return true
}

func authMethod(token string) *http.BasicAuth {
	return &http.BasicAuth{
		Username: "x-oauth-basic",
		Password: token,
	}
}

func clone(org string, name string, target string, token string) *git.Repository {
	repo, err := git.PlainClone(target, false, &git.CloneOptions{
		Auth:     authMethod(token),
		URL:      createGithubURL(org, name),
		Progress: os.Stdout,
	})
	CheckIfError(err)

	log.Debugf("Cloned repository '%s' to '%s'", GithubLink(org, name), target)
	return repo
}

func createGithubURL(org string, name string) string {
	return fmt.Sprintf("https://github.com/%s", GithubLink(org, name))
}
