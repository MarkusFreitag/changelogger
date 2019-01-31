package stringutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIndentLvl(t *testing.T) {
	require.Equal(t, 0, IndentLvl("- one liner"))
	require.Equal(t, 1, IndentLvl(" - one liner"))
	require.Equal(t, 2, IndentLvl("  - one liner"))
	require.Equal(t, 2, IndentLvl("  - multi\n    liner"))
}

func TestIncrIndent(t *testing.T) {
	require.Equal(t, "  - one liner", IncrIndent("- one liner", 2))
	require.Equal(t, "  - multi\n    liner", IncrIndent("- multi\n  liner", 2))
}

func TestDecrIndent(t *testing.T) {
	require.Equal(t, "- one liner", DecrIndent("  - one liner", 2))
	require.Equal(t, "- multi\n  liner", DecrIndent("  - multi\n    liner", 2))
}

func TestComment(t *testing.T) {
	require.Equal(t, "#one liner", Comment("one liner"))
	require.Equal(t, "#multi\n#liner", Comment("multi\nliner"))
}

func TestUncomment(t *testing.T) {
	require.Equal(t, "one liner", Uncomment("#one liner"))
	require.Equal(t, "multi\nliner", Uncomment("#multi\n#liner"))
}
