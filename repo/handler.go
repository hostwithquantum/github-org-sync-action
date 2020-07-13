package repo

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	gucci "github.com/noqcks/gucci/gucci"
	log "github.com/sirupsen/logrus"
)

// Handler ...
type Handler struct {
	Base      string
	Workflows []string
	Files     []string
}

// NewHandler ...
func NewHandler(repository string, baseDirectory string) Handler {
	baseDir := filepath.Join(baseDirectory, repository)

	workflows := getWorkflows(filepath.Join(baseDir, ".github/workflows"))
	files := getTemplates(baseDir)

	return Handler{
		Base:      repository,
		Workflows: workflows,
		Files:     files,
	}
}

// Sync ...
func (h Handler) Sync(target string, sources []string) {
	if ensureDirectory(target) == false {
		return
	}

	for _, file := range sources {
		var repoFile string

		if !isTemplate(file) {
			repoFile = filepath.Join(target, filepath.Base(file))

			err := copy(file, repoFile)
			CheckIfError(err)

			log.Debugf("Synced from '%s' to '%s'", file, repoFile)

			continue
		}

		log.Debugf("Working template: %s", file)

		repoFile = filepath.Join(
			target,
			strings.TrimSuffix(filepath.Base(file), filepath.Ext(file)),
		)

		log.Debugf("Generating: %s", repoFile)

		tmpl, err := gucci.LoadTemplateFile(file)
		CheckIfError(err)

		writer, err := os.Create(repoFile)
		CheckIfError(err)

		// GET most/all from env
		tmplVars := gucci.Env()
		tmplVars["GITHUB_REPOSITORY"] = filepath.Base(target)

		err = gucci.ExecuteTemplate(tmplVars, writer, tmpl)
		CheckIfError(err)

		log.Debugf("Created %s from template %s", repoFile, file)
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
	log.Debugf("Fetch workflows in: %s", path)
	return getFiles(path, ".yml")
}

func getTemplates(path string) []string {
	log.Debugf("Fetch templates from: %s", path)
	return getFiles(path, ".tpl")
}

func getFiles(path string, ext string) []string {
	var fetchedFiles []string

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) != ext {
			continue
		}

		fetchedFiles = append(fetchedFiles, filepath.Join(path, file.Name()))
	}

	return fetchedFiles
}

func isTemplate(path string) bool {
	if filepath.Ext(path) == ".tpl" {
		return true
	}

	return false
}
