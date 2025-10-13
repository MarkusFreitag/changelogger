package gitconfig // import "github.com/MarkusFreitag/changelogger/pkg/gitconfig"

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"syscall"

	"github.com/MarkusFreitag/changelogger/pkg/parser"
)

func fetchConfigAttribute(key string) (string, error) {
	var stdout bytes.Buffer
	cmd := exec.Command("git", "config", "--get", "--null", key)
	cmd.Stdout = &stdout
	cmd.Stderr = io.Discard

	err := cmd.Run()
	if exitError, ok := err.(*exec.ExitError); ok {
		if waitStatus, ok := exitError.Sys().(syscall.WaitStatus); ok {
			if waitStatus.ExitStatus() == 1 {
				return "", fmt.Errorf("failed to get git config key %q", key)
			}
		}
		return "", err
	}

	return strings.TrimRight(stdout.String(), "\000"), nil
}

func GetGitAuthor() (*parser.Author, error) {
	var author parser.Author
	var err error

	author.Name, err = fetchConfigAttribute("user.name")
	if err != nil {
		return nil, err
	}

	author.Email, err = fetchConfigAttribute("user.email")
	if err != nil {
		return nil, err
	}

	if author.Name != "" && author.Email != "" {
		return &author, nil
	}

	return nil, errors.New("couldn't find an author in any config")
}
