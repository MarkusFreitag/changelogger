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
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/MarkusFreitag/changelogger/pkg/parser"
	"github.com/MarkusFreitag/changelogger/pkg/stringutil"
	"github.com/spf13/cobra"
)

var (
	changelogPath   string
	templatePath    string
	binaryName      string
	temp            *template.Template
	releases        parser.Releases
	out             = bytes.NewBufferString("")
	defaultTemplate = `{{.BinaryName}} ({{.Version}}) UNSTABLE; urgency=medium

{{.Text}}

 -- {{.AuthorName}} <{{.AuthorMail}}>  {{.Date}}
`
)

func loadTemplate() (*template.Template, error) {
	if _, err := os.Stat(templatePath); err == nil {
		customTemplate, err := ioutil.ReadFile(templatePath)
		if err != nil {
			return nil, err
		}
		return template.New("changelog").Parse(string(customTemplate))
	}
	return template.New("changelog").Parse(defaultTemplate)
}

type templateData struct {
	BinaryName string
	Version    string
	Text       string
	AuthorName string
	AuthorMail string
	Date       string
}

var debianCmd = &cobra.Command{
	Use:   "debian",
	Short: "Generate debian changelog file",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if _, err := os.Stat(changelogPath); !os.IsNotExist(err) && !force {
			return fmt.Errorf("`%s` already exist, use --force to overwrite it", changelogPath)
		}
		if value, ok := os.LookupEnv("DEBIAN_NAME"); ok {
			binaryName = value
		} else {
			return errors.New("`DEBIAN_NAME` env var is not set")
		}

		var err error
		temp, err = loadTemplate()
		if err != nil {
			return err
		}

		releases, err = parser.ReadFile(changelogFile)
		return err
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		if _, err := os.Stat(templatePath); err == nil {
			return os.Remove(templatePath)
		}
		return ioutil.WriteFile(changelogPath, out.Bytes(), 0644)
	},
}

var debianDummyCmd = &cobra.Command{
	Use:   "dummy",
	Short: "Generate dummy debian changelog",
	RunE: func(cmd *cobra.Command, args []string) error {
		var rel *parser.Release
		for _, r := range releases {
			if r.Released {
				rel = r
				break
			}
		}
		if rel == nil {
			return errors.New("no released version found")
		}

		data := templateData{
			BinaryName: binaryName,
			Version:    rel.Version.Original(),
			AuthorName: rel.By.Name,
			AuthorMail: rel.By.Email,
			Text:       "* See CHANGELOG.md",
			Date:       rel.Date.Format("Mon, 02 Jan 2006 15:04:05 -0700"),
		}

		return temp.ExecuteTemplate(out, "changelog", data)
	},
}

var debianFullCmd = &cobra.Command{
	Use:   "full",
	Short: "Generate full debian changelog",
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, rel := range releases {
			if !rel.Released {
				continue
			}

			var text string
			if value, ok := rel.Changes[rel.By.Name]; ok && len(rel.Changes) == 1 {
				text = value
			} else {
				for _, author := range rel.Changes.SortedAuthors() {
					text += fmt.Sprintf("[ %s ]\n", author)
					text += stringutil.IncrIndent(rel.Changes[author], 2)
					text += "\n"
				}
			}
			text = strings.TrimSuffix(text, "\n")
			text = stringutil.IncrIndent(text, 2)

			data := templateData{
				BinaryName: binaryName,
				Version:    rel.Version.Original(),
				AuthorName: rel.By.Name,
				AuthorMail: rel.By.Email,
				Text:       text,
				Date:       rel.Date.Format("Mon, 02 Jan 2006 15:04:05 -0700"),
			}

			err := temp.ExecuteTemplate(out, "changelog", data)
			if err != nil {
				return err
			}

			out.WriteString("\n")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(debianCmd)
	debianCmd.PersistentFlags().StringVar(&changelogPath, "output-file", "debian/changelog", "")
	debianCmd.PersistentFlags().StringVar(&templatePath, "template-file", "debian/changelog.template", "")
	debianCmd.PersistentFlags().BoolVar(&force, "force", false, "Overwrite an existing file")

	debianCmd.AddCommand(debianDummyCmd)
	debianCmd.AddCommand(debianFullCmd)
}
