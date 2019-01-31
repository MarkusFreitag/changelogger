// Copyright © 2019 Markus Freitag <fmarkus@mailbox.org>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/MarkusFreitag/changelogger/pkg/gitconfig"
	"github.com/MarkusFreitag/changelogger/pkg/parser"
	"github.com/Masterminds/semver"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
)

const VERSIONFORMATDEFAULT = "debian"

var (
	versionFormat string
	versionOnly   bool
)

var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Release related commands",
}

var releaseLastCmd = &cobra.Command{
	Use:   "last",
	Short: "Show last release",
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(changelogFile); os.IsNotExist(err) {
			handleError(fmt.Errorf("%s does not exist", changelogFile))
		}
		rels, err := parser.ReadFile(changelogFile)
		handleError(err)
		var lastRelease *parser.Release
		for _, rel := range rels {
			if rel.Released {
				lastRelease = rel
				break
			}
		}

		if lastRelease == nil {
			handleError(fmt.Errorf("no released version available"))
		}
		if versionOnly {
			fmt.Printf("v%s\n", lastRelease.Version.String())
			return
		}
		fmt.Println(lastRelease.Show())
	},
}

var releaseNewCmd = &cobra.Command{
	Use:   "new",
	Short: "Create new release",
	PreRun: func(cmd *cobra.Command, args []string) {
		if versionFormat == VERSIONFORMATDEFAULT {
			if val, ok := os.LookupEnv("CHANGELOGGER_VERSION_FORMAT"); ok {
				versionFormat = val
			}
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(changelogFile); os.IsNotExist(err) {
			handleError(fmt.Errorf("%s does not exist", changelogFile))
		}

		gitAuthor, err := gitconfig.GetGitAuthor()
		handleError(err)

		rels, err := parser.ReadFile(changelogFile)
		handleError(err)

		if len(rels) == 0 || rels[0].Released {
			handleError(fmt.Errorf("no unreleased version available"))
		}
		var lastVersion *semver.Version
		if len(rels) == 1 {
			lastVersion, err = semver.NewVersion("v0.0.0-1")
			handleError(err)
		} else {
			lastVersion = rels[1].Version
		}

		var bump string
		prompt := &survey.Select{
			Message: "Select the level of version bump:",
			Options: []string{"debian", "patch", "minor", "major"},
		}
		if versionFormat == "semver" {
			prompt.Options = []string{"patch", "minor", "major"}
		}
		err = survey.AskOne(prompt, &bump, nil)
		handleError(err)

		var newVersion semver.Version
		switch bump {
		case "debian":
			debian, err := strconv.Atoi(getFirstNumber(lastVersion.Prerelease()))
			handleError(err)
			newVersion, err = lastVersion.SetPrerelease(strconv.Itoa(debian + 1))
			handleError(err)
		case "patch":
			newVersion = lastVersion.IncPatch()
			if versionFormat == "debian" {
				newVersion = newVersion.IncPatch()
				newVersion, err = newVersion.SetPrerelease("1")
				handleError(err)
			}
		case "minor":
			newVersion = lastVersion.IncMinor()
			if versionFormat == "debian" {
				newVersion, err = newVersion.SetPrerelease("1")
				handleError(err)
			}
		case "major":
			newVersion = lastVersion.IncMajor()
			if versionFormat == "debian" {
				newVersion, err = newVersion.SetPrerelease("1")
				handleError(err)
			}
		}

		rels[0].Version = &newVersion
		rels[0].Date = time.Now()
		rels[0].By = *gitAuthor
		rels[0].GenerateHeader(bump)
		rels[0].GenerateFooter()

		blocks := make([]string, len(rels))
		for index, rel := range rels {
			blocks[index] = rel.Show()
		}

		file, err := os.Create(changelogFile)
		handleError(err)
		defer file.Close()

		_, err = io.Copy(file, strings.NewReader(strings.Join(blocks, "\n")))
		handleError(err)
		fmt.Printf("Released %s version %s (pre: %s)\n", bump, newVersion.String(), lastVersion.String())
	},
}

func getFirstNumber(s string) string {
	var buffer string
	for _, char := range s {
		if !unicode.IsDigit(char) {
			return buffer
		}
		buffer += string(char)
	}
	return buffer
}

func init() {
	rootCmd.AddCommand(releaseCmd)

	releaseCmd.AddCommand(releaseNewCmd)
	releaseNewCmd.Flags().StringVar(&versionFormat, "version-format", VERSIONFORMATDEFAULT, fmt.Sprintf("version format (default is %s)", VERSIONFORMATDEFAULT))
	releaseCmd.AddCommand(releaseLastCmd)
	releaseLastCmd.Flags().BoolVar(&versionOnly, "version-only", false, "show only version")
}
