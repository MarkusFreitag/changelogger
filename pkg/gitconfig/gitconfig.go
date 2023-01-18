package gitconfig // import "github.com/MarkusFreitag/changelogger/pkg/gitconfig"

import (
	"errors"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/MarkusFreitag/changelogger/pkg/parser"
	"gopkg.in/ini.v1"
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

		config, err := ini.Load(file)
		if err != nil {
			return nil, err
		}

		var author parser.Author
		if section := config.Section("user"); section != nil {
			if key := section.Key("name"); key != nil {
				author.Name = key.String()
			}
			if key := section.Key("email"); key != nil {
				author.Email = key.String()
			}
		}

		if author.Name != "" && author.Email != "" {
			return &author, nil
		}
	}
	return nil, errors.New("couldn't find an author in any config")
}
