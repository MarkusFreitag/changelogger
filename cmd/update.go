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
	"os"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
)

const REPO = "MarkusFreitag/changelogger"

var (
	force bool
)

func updateChecker() (*selfupdate.Release, error) {
	latest, found, err := selfupdate.DetectLatest(REPO)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, fmt.Errorf("%s does not have any releases", REPO)
	}

	current, err := semver.ParseTolerant(BuildVersion)
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
	RunE: func(cmd *cobra.Command, args []string) error {
		latest, err := updateChecker()
		if err != nil {
			return err
		}
		if latest == nil {
			fmt.Println("Already up to date")
			return nil
		}

		if !force {
			var confirm bool
			prompt := &survey.Confirm{
				Message: fmt.Sprintf("%s => %s Continue?", BuildVersion, latest.Version.String()),
			}
			err = survey.AskOne(prompt, &confirm, nil)
			if err != nil {
				return err
			}
			if !confirm {
				return nil
			}
		}

		exe, err := os.Executable()
		if err != nil {
			return err
		}
		err = selfupdate.UpdateTo(latest.AssetURL, exe)
		if err != nil {
			return err
		}
		fmt.Println("Successfully updated to version", latest.Version.String())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.Flags().BoolVar(&force, "force", false, "update without confirmation")
}
