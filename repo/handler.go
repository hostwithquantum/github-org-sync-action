package repo

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

// Handler ...
type Handler struct {
	Base      string
	Workflows []string
}

// NewHandler ...
func NewHandler(repository string, baseDirectory string) Handler {
	workflows := getWorkflows(fmt.Sprintf("%s/%s/.github/workflows", baseDirectory, repository))
	return Handler{Base: repository, Workflows: workflows}
}

// Sync ...
func (h Handler) Sync(target string) {
	for _, file := range h.Workflows {
		repoFile := fmt.Sprintf(
			"%s/.github/workflows/%s", target, filepath.Base(file))

		err := copy(file, repoFile)
		CheckIfError(err)

		log.Debugf("Synced to '%s' to '%s'", file, repoFile)
	}

	log.Debugf("Synced files to '%s'", target)
}

func copy(src string, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func getWorkflows(path string) []string {
	var workflows []string

	log.Debugf("Fetch workflows in: %s", path)

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		workflows = append(workflows, fmt.Sprintf("%s/%s", path, file.Name()))
	}

	return workflows
}
