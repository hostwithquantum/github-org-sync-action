package repo

import (
	"fmt"
	"os"
)

// CheckIfError ...
func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}

// GithubLink ...
func GithubLink(org string, name string) string {
	return fmt.Sprintf("%s/%s", org, name)
}
