package gitconfig_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/MarkusFreitag/changelogger/pkg/gitconfig"
)

func TestGetGitAuthor(t *testing.T) {
	var config = `[user]
  email = git@example.com
  name = Test Dummy`

	configPath := filepath.Join(t.TempDir(), "gitconfig")
	require.NoError(t, os.WriteFile(configPath, []byte(config), 0600))
	t.Setenv("GIT_CONFIG", configPath)

	author, err := gitconfig.GetGitAuthor()
	assert.NoError(t, err)
	assert.NotNil(t, author)
	assert.Equal(t, "git@example.com", author.Email)
	assert.Equal(t, "Test Dummy", author.Name)
}
