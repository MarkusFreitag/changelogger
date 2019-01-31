package stringutil // import "github.com/MarkusFreitag/changelogger/pkg/stringutil"

import "strings"

func IndentLvl(block string) int {
	if !strings.HasPrefix(block, " ") {
		return 0
	}
	lvl := 10000
	for _, line := range strings.Split(block, "\n") {
		for pos, char := range line {
			if char != ' ' && pos < lvl {
				lvl = pos
			}
		}
	}
	return lvl
}

func IncrIndent(block string, count int) string {
	lines := strings.Split(block, "\n")
	for index, line := range lines {
		lines[index] = strings.Repeat(" ", count) + line
	}
	return strings.Join(lines, "\n")
}

func DecrIndent(block string, count int) string {
	lines := strings.Split(block, "\n")
	for index, line := range lines {
		lines[index] = strings.TrimPrefix(line, strings.Repeat(" ", count))
	}
	return strings.Join(lines, "\n")
}

func Comment(block string) string {
	lines := strings.Split(block, "\n")
	for index, line := range lines {
		lines[index] = "#" + line
	}
	return strings.Join(lines, "\n")
}

func Uncomment(block string) string {
	lines := strings.Split(block, "\n")
	for index, line := range lines {
		lines[index] = strings.TrimPrefix(line, "#")
	}
	return strings.Join(lines, "\n")
}

func SplitBlock(block, splitter string) (string, string) {
	parts := strings.SplitN(block, splitter, 2)
	if len(parts) > 1 {
		return parts[0], strings.TrimPrefix(block, parts[0]+"\n")
	}
	return block, ""
}
