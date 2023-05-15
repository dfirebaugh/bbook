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
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/dfirebaugh/docbook/pkg/parser"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "serve the .book dir in a local web server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		serveSite()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	logrus.SetLevel(logrus.TraceLevel)
}

func serveSite() {
	buildSite()

	http.Handle(conf.Output["html"].SiteURL, http.StripPrefix(conf.Output["html"].SiteURL, http.FileServer(http.Dir(".book"))))

	http.HandleFunc("/reload", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "<html><body><script>reloadPage();</script></body></html>")
	})

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	err = watcher.Add("./src")
	if err != nil {
		log.Fatal(err)
	}

	logrus.Println("serving on http://localhost:5555" + conf.Output["html"].SiteURL)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("Reloading site due to file change:", event.Name)

					buildSite()

					http.Get("http://localhost:5555/reload")
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("Error watching filesystem:", err)
			}
		}
	}()

	log.Fatal(http.ListenAndServe(":5555", nil))
}

func readSummary() []parser.PageLink {
	file, err := os.Open("./src/SUMMARY.md")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	links, err := parser.ParseLinks(file, ".")
	if err != nil {
		panic(err)
	}

	return links
}
