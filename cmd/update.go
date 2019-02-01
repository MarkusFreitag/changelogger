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
	"errors"
	"fmt"
	"os"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
)

var (
	force bool
)

func updateChecker() (*selfupdate.Release, error) {
	latest, found, err := selfupdate.DetectLatest("MarkusFreitag/changelogger")
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, errors.New("github.com/MarkusFreitag/changelogger does not have any releases")
	}

	current, err := semver.Parse(BuildVersion)
	if err != nil {
		return nil, err
	}

	if latest.Version.LTE(current) {
		return nil, nil
	}
	return latest, nil
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update to the latest version",
	Run: func(cmd *cobra.Command, args []string) {
		latest, err := updateChecker()
		handleError(err)

		if !force {
			var confirm bool
			prompt := &survey.Confirm{
				Message: fmt.Sprintf("%s => %s Continue?", BuildVersion, latest.Version.String()),
			}
			err = survey.AskOne(prompt, &confirm, nil)
			handleError(err)
			if !confirm {
				return
			}
		}

		exe, err := os.Executable()
		handleError(err)
		err = selfupdate.UpdateTo(latest.AssetURL, exe)
		handleError(err)
		fmt.Println("Successfully updated to version", latest.Version.String())
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.Flags().BoolVar(&force, "force", false, "update without confirmation")
}
