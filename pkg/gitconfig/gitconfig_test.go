package gitconfig

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetGitAuthor(t *testing.T) {
	files = []string{}
	author, err := GetGitAuthor()
	require.Equal(t, "couldn't find an author in any config", err.Error())
	require.Nil(t, author)

	files = []string{"testdata/gitconfig"}
	author, err = GetGitAuthor()
	require.Nil(t, err)
	require.Equal(t, "Test Dummy", author.Name)
	require.Equal(t, "git@example.com", author.Email)
}
