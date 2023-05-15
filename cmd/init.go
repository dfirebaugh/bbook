/*
Copyright © 2023 Dustin Firebaugh<dafirebaugh@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"errors"
	"os"

	"github.com/dfirebaugh/docbook/pkg/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [bookname]",
	Short: "initialize a book",
	Long:  ``,
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		initBook(args[0])
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func initBook(bookName string) {
	os.Mkdir("src", 0777)
	writeConfig(bookName)
	writeInitialSummaryMD()
}

func writeInitialSummaryMD() {
	if _, err := os.Stat("./src/SUMMARY.md"); errors.Is(err, os.ErrNotExist) {
		logrus.Trace("creating summary file")
		f, _ := os.OpenFile("./src/SUMMARY.md", os.O_WRONLY|os.O_CREATE, 0666)
		defer f.Close()

		f.Write([]byte("\n# Summary\n- [Intro](./intro.md)"))
	}
}

func writeConfig(bookName string) {
	logrus.Trace("writing initial config file")

	output := make(map[string]config.Theme)
	output["html"] = config.Theme{
		DefaultTheme: "simple",
		SiteURL:      "/book/",
	}

	config.WriteConfig(&config.Config{
		Book: config.BookConfig{
			Authors: []string{},
			Src:     "src",
			Title:   bookName,
		},
		Output: output,
	}, "book.toml")
}
