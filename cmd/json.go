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

	"github.com/MarkusFreitag/changelogger/pkg/parser"
	pj "github.com/hokaccha/go-prettyjson"
	"github.com/spf13/cobra"
)

var prettyPrint bool

var jsonCmd = &cobra.Command{
	Use:   "json",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(changelogFile); os.IsNotExist(err) {
			handleError(fmt.Errorf("%s does not exist", changelogFile))
		}
		rels, err := parser.ReadFile(changelogFile)
		handleError(err)
		f := pj.NewFormatter()
		if !prettyPrint {
			f.DisabledColor = true
			f.Indent = 0
			f.Newline = ""
		}
		s, _ := f.Marshal(rels)
		fmt.Println(string(s))
	},
}

func init() {
	rootCmd.AddCommand(jsonCmd)

	jsonCmd.Flags().BoolVar(&prettyPrint, "pretty-print", false, "indented and colorized output")
}
