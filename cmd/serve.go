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
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dfirebaugh/bbook/pkg/parser"
	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve the .book directory with a local web server",
	Long:  "Starts a local web server to serve the .book directory and watches for changes in the ./src directory to rebuild and reload automatically.",
	Run: func(cmd *cobra.Command, args []string) {
		serveSite()
	},
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	clients   = make(map[*websocket.Conn]bool)
	broadcast = make(chan struct{})
)

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// logrus.Error(err)
		return
	}
	defer ws.Close()

	clients[ws] = true

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			// logrus.Errorf("error: %v", err)
			delete(clients, ws)
			break
		}
	}
}

func handleMessages() {
	for {
		<-broadcast
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, []byte("reload"))
			if err != nil {
				logrus.Errorf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func init() {
	rootCmd.AddCommand(serveCmd)
	logrus.SetLevel(logrus.TraceLevel)
}

func serveSite() {
	buildSite()

	siteURL := ensureTrailingSlash(conf.Output["html"].SiteURL)

	http.Handle(siteURL, http.StripPrefix(siteURL, http.FileServer(http.Dir(".book"))))

	http.HandleFunc("/ws", handleConnections)

	go handleMessages()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Error creating watcher: %v", err)
	}
	defer watcher.Close()

	err = watcher.Add("./src")
	if err != nil {
		log.Fatalf("Error adding watch on ./src: %v", err)
	}

	logrus.Println("Serving on http://localhost:5555" + conf.Output["html"].SiteURL)

	var debounceTimer *time.Timer
	debounceDuration := 500 * time.Millisecond

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove) != 0 {
					// log.Printf("File change detected: %s. Scheduling rebuild...", event.Name)
					if debounceTimer != nil {
						debounceTimer.Stop()
					}
					debounceTimer = time.AfterFunc(debounceDuration, func() {
						buildSite()

						select {
						case broadcast <- struct{}{}:
						default: // Avoid blocking if reload is already in progress
						}
					})
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Printf("Error watching filesystem: %v", err)
			}
		}
	}()

	log.Fatal(http.ListenAndServe(":5555", nil))
}

func readSummary() []parser.PageLink {
	file, err := os.Open("./src/SUMMARY.md")
	if err != nil {
		log.Fatalf("Error opening SUMMARY.md: %v", err)
	}
	defer file.Close()

	links, err := parser.ParseLinks(file, ".")
	if err != nil {
		log.Fatalf("Error parsing SUMMARY.md: %v", err)
	}

	return links
}

func ensureTrailingSlash(url string) string {
	if !strings.HasSuffix(url, "/") {
		return url + "/"
	}
	return url
}
