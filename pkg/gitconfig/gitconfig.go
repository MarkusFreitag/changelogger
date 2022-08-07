package gitconfig // import "github.com/MarkusFreitag/changelogger/pkg/gitconfig"

import (
	"errors"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/MarkusFreitag/changelogger/pkg/parser"
	"github.com/muja/goconfig"
)

var files = []string{".git/config", "~/.gitconfig", "/etc/gitconfig", "~/.config/git/config"}

func GetGitAuthor() (*parser.Author, error) {
	user, err := user.Current()
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if strings.HasPrefix(file, "~") {
			file = strings.TrimPrefix(file, "~")
			file = filepath.Join(user.HomeDir, file)
		}
		if _, err := os.Stat(file); os.IsNotExist(err) {
			continue
		}
		bytes, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}
		config, _, err := goconfig.Parse(bytes)
		if err != nil {
			return nil, err
		}
		if config["user.name"] != "" && config["user.email"] != "" {
			return &parser.Author{Name: config["user.name"], Email: config["user.email"]}, nil
		}
	}
	return nil, errors.New("couldn't find an author in any config")
}
