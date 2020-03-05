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

package cmd // import "github.com/MarkusFreitag/changelogger/cmd"

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/MarkusFreitag/changelogger/pkg/editor"
	"github.com/MarkusFreitag/changelogger/pkg/gitconfig"
	"github.com/MarkusFreitag/changelogger/pkg/parser"
	"github.com/MarkusFreitag/changelogger/pkg/stringutil"
	"github.com/Masterminds/semver"
	"github.com/spf13/cobra"
)

var changelogFile string

func handleError(err error) {
	if err != nil {
		fmt.Printf("err: %s\n", err.Error())
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "changelogger",
	Short: "Create and update changelogs with ease",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if val, ok := os.LookupEnv("CHANGELOGGER_MIN_VERSION"); ok {
			minVersion, err := semver.NewVersion(val)
			handleError(err)
			bVersion, err := semver.NewVersion(BuildVersion)
			handleError(err)
			if bVersion.LessThan(minVersion) {
				fmt.Printf("Your current version: %s\n", bVersion.String())
				fmt.Printf("Required version: %s\n", minVersion.String())
				fmt.Println("Run `changelogger update`")
				os.Exit(1)
			}
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		rels := make(parser.Releases, 0)
		if _, err := os.Stat(changelogFile); err == nil {
			rels, err = parser.ReadFile(changelogFile)
			handleError(err)
		}

		gitAuthor, err := gitconfig.GetGitAuthor()
		handleError(err)

		if len(rels) == 0 {
			rels = append(rels, parser.NewRelease())
		} else if rels[0].Released {
			rels = parser.PrependRelease(rels, parser.NewRelease())
		}

		if _, ok := rels[0].Changes[gitAuthor.Name]; !ok {
			rels[0].Changes[gitAuthor.Name] = ""
		}

		existingChanges := rels[0].Changes.FormatByAuthor(gitAuthor.Name, 0)
		existingChanges = strings.TrimSpace(existingChanges)
		existingChanges = stringutil.Comment(existingChanges)
		userInput := existingChanges

		err = editor.Open(&userInput)
		handleError(err)
		userInput = strings.TrimSuffix(userInput, "\n")

		if userInput == existingChanges {
			fmt.Println("exit without writing")
			return
		}

		userInput = strings.TrimPrefix(userInput, existingChanges+"\n")
		userInput = stringutil.DecrIndent(userInput, stringutil.IndentLvl(userInput))

		if block := rels[0].Changes[gitAuthor.Name]; block != "" {
			rels[0].Changes[gitAuthor.Name] = strings.Join([]string{block, userInput}, "\n")
		} else {
			rels[0].Changes[gitAuthor.Name] = userInput
		}

		blocks := make([]string, len(rels))
		for index, rel := range rels {
			blocks[index] = rel.Show()
		}

		file, err := os.Create(changelogFile)
		handleError(err)
		defer file.Close()

		_, err = io.Copy(file, strings.NewReader(strings.Join(blocks, "\n")))
		handleError(err)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("err: %s\n", err.Error())
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&changelogFile, "file", "f", "CHANGELOG.md", "changelog file (default is CHANGELOG.md)")
}
