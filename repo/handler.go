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
	dotGithub := filepath.Join(target, ".github")
	if ensureDirectory(dotGithub) == false {
		return
	}

	workflows := filepath.Join(dotGithub, "workflows")
	if ensureDirectory(workflows) == false {
		return
	}

	for _, file := range h.Workflows {
		repoFile := filepath.Join(workflows, filepath.Base(file))

		err := copy(file, repoFile)
		CheckIfError(err)

		log.Debugf("Synced to '%s' to '%s'", file, repoFile)
	}

	log.Debugf("Synced files to '%s'", target)
}

func ensureDirectory(directory string) bool {
	stat, err := os.Stat(directory)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(directory, os.ModePerm)
		}

		// handle all errors
		CheckIfError(err)
		return true
	}

	if stat.IsDir() {
		return true
	}

	log.Errorf("Path '%s' exists, but is not a directory!", directory)
	return false
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
