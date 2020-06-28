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

// Contains ...
func Contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

// GithubLink ...
func GithubLink(org string, name string) string {
	return fmt.Sprintf("%s/%s", org, name)
}
