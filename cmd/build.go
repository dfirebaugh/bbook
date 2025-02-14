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
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/dfirebaugh/bbook/templates"
	"github.com/dfirebaugh/bbook/web"

	"github.com/sirupsen/logrus"
	gohtml "golang.org/x/net/html"

	"github.com/dfirebaugh/bbook/pkg/config"
	"github.com/dfirebaugh/bbook/pkg/md"
	"github.com/dfirebaugh/bbook/pkg/parser"

	"github.com/spf13/cobra"
)

var conf *config.Config

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "build the book to the .book dir",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		buildSite()
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
	var err error
	conf, err = config.ReadConfig("book.toml")
	if err != nil {
		logrus.Error(err)
	}
}

func generateMDFiles() {
	for _, l := range readSummary() {
		writeFile(filepath.Join(conf.Book.Src, l.FilePath), []byte(fmt.Sprintf("\n# %s\n", l.Header)))
	}
}

func generateMDFilesIfNeeded() {
	for _, link := range readSummary() {
		filePath := filepath.Join(conf.Book.Src, link.FilePath)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			writeFile(filePath, []byte(fmt.Sprintf("\n# %s\n", link.Header)))
		}
	}
}

func copyAdditionalAssets() {
	for _, file := range conf.Output["html"].AdditionalJS {
		srcPath := filepath.Join(conf.Book.Src, "..", file)
		dstPath := filepath.Join(".book", file)
		if err := copyFile(srcPath, dstPath); err != nil {
			logrus.Errorf("Error copying JS file %s: %v", file, err)
		}
	}
	for _, file := range conf.Output["html"].AdditionalCSS {
		srcPath := filepath.Join(conf.Book.Src, "..", file)
		dstPath := filepath.Join(".book", file)
		if err := copyFile(srcPath, dstPath); err != nil {
			logrus.Errorf("Error copying CSS file %s: %v", file, err)
		}
	}
}

func buildSite() {
	logrus.Println("Building files to the `.book` directory...")

	generateMDFilesIfNeeded()

	if err := os.RemoveAll(".book"); err != nil {
		logrus.Errorf("Failed to clean up .book directory: %v", err)
		return
	}
	if err := os.Mkdir(".book", 0o777); err != nil {
		logrus.Errorf("Failed to create .book directory: %v", err)
		return
	}

	src := parser.ParseDir(conf.Book.Src)

	copyAdditionalAssets()

	if err := web.CopyStaticFiles(web.StaticFiles, "static", ".book/"); err != nil {
		logrus.Errorf("Failed to copy static files: %v", err)
		return
	}

	generatePages(src)

	logrus.Println("Build complete.")
}

func generatePages(pages []parser.Page) {
	tmpl, err := template.New("page").Parse(templates.PageTemplate)
	if err != nil {
		log.Fatal(err)
	}

	summaryLinks, err := parser.ParsePagesFromSummaryMD(filepath.Join(conf.Book.Src, "SUMMARY.md"))
	if err != nil {
		logrus.Error(err)
	}

	var links []parser.Page

	for _, l := range summaryLinks {
		links = append(links, parser.Page{
			URL:   l.URL,
			Title: l.Title,
		})
	}

	var wg sync.WaitGroup
	for i, page := range links {
		wg.Add(1)
		go func(page parser.Page, pages []parser.Page, index int) {
			if len(pages) == 1 {
				buildPage(page, tmpl, pages[0].URL, pages[0].URL)
				wg.Done()
				return
			}
			if index == 0 {
				buildPage(page, tmpl, pages[index+1].URL, pages[len(pages)-1].URL)
				// copy this to .book/index.html
				wg.Done()
				return
			}
			if index == len(pages)-1 {
				buildPage(page, tmpl, pages[0].URL, pages[index-1].URL)
				wg.Done()
				return
			}
			buildPage(page, tmpl, pages[index+1].URL, pages[index-1].URL)
			wg.Done()
		}(page, links, i)
	}

	wg.Wait()

	// copy the first link to index.html
	f, err := os.ReadFile(
		filepath.Join(".book", filepath.Base(mdLinkToHTMLLink(filepath.Join(conf.Book.Src, links[0].URL)))),
	)
	if err != nil {
		logrus.Error(err)
	}
	os.Remove(".book/index.html")
	writeFile(".book/index.html", f)
}

func buildPage(page parser.Page, tmpl *template.Template, nextPage string, previousPage string) {
	buildPath := filepath.Join(".book", mdLinkToHTMLLink(page.URL))
	buildDir := filepath.Dir(buildPath)
	go copyStaticFiles(conf.Book.Src, buildDir)

	err := os.MkdirAll(buildDir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(buildPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	contentBytes, err := os.ReadFile(filepath.Join(conf.Book.Src, page.URL))
	if err != nil {
		logrus.Error(err)
	}
	pageBytes := md.ToHTML([]byte(contentBytes))
	content, err := replaceMDWithHTMLInLinks(string(getBodyChildren(pageBytes)))
	if err != nil {
		logrus.Error(err)
	}
	navbar, err := replaceMDWithHTMLInLinks(string(buildNavLinks()))
	if err != nil {
		logrus.Error(err)
	}

	err = tmpl.Execute(f, struct {
		BookTitle       string
		SiteURL         string
		Title           string
		CurrentPage     string
		NextPage        string
		PreviousPage    string
		Theme           string
		Body            template.HTML
		NavLinks        template.HTML
		AdditionalJS    []string
		AdditionalCSS   []string
		HideSidebar     bool
		GithubRepoURL   string
		GithubRepoIcon  string
		EditURLTemplate string
	}{
		BookTitle:       conf.Book.Title,
		SiteURL:         conf.Output["html"].SiteURL,
		Title:           page.Title,
		CurrentPage:     page.URL,
		NextPage:        addSiteURL(mdLinkToHTMLLink(nextPage)),
		PreviousPage:    addSiteURL(mdLinkToHTMLLink(previousPage)),
		Body:            template.HTML(getBodyChildren([]byte(content))),
		NavLinks:        template.HTML(getBodyChildren([]byte(navbar))),
		Theme:           conf.Output["html"].DefaultTheme,
		AdditionalJS:    conf.Output["html"].AdditionalJS,
		AdditionalCSS:   conf.Output["html"].AdditionalCSS,
		HideSidebar:     conf.Output["html"].HideSidebar,
		GithubRepoURL:   conf.Output["html"].GithubRepoURL,
		GithubRepoIcon:  conf.Output["html"].GithubRepoIcon,
		EditURLTemplate: conf.Output["html"].EditURLTemplate,
	})
	if err != nil {
		log.Fatal(err)
	}
}

// writeFile will write file if the file doesn't exist
func writeFile(filePath string, content []byte) {
	err := os.MkdirAll(filepath.Dir(filePath), 0o755)
	if err != nil {
		logrus.Errorf("Failed to create directories for %s: %v", filePath, err)
		return
	}

	if existing, err := os.ReadFile(filePath); err == nil {
		if bytes.Equal(existing, content) {
			return
		}
	}

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		logrus.Errorf("Failed to open file %s: %v", filePath, err)
		return
	}
	defer f.Close()

	if _, err := f.Write(content); err != nil {
		logrus.Errorf("Failed to write to file %s: %v", filePath, err)
	}
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func copyStaticFiles(srcDir, dstDir string) error {
	err := filepath.Walk(srcDir, func(srcPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		relPath, _ := filepath.Rel(srcDir, srcPath)
		dstPath := filepath.Join(dstDir, relPath)

		err = os.MkdirAll(filepath.Dir(dstPath), os.ModePerm)
		if err != nil {
			return err
		}

		if filepath.Base(srcPath) == "SUMMARY.md" {
			return nil
		}

		if filepath.Ext(srcPath) == ".md" {
			return nil
		}

		return copyFile(srcPath, dstPath)
	})

	return err
}

func readNavLinks() []byte {
	summaryFile, err := os.Open("./src/SUMMARY.md")
	if err != nil {
		panic(err)
	}
	defer summaryFile.Close()
	b, err := io.ReadAll(summaryFile)
	if err != nil {
		logrus.Error(err)
	}
	return b
}

func buildNavLinks() []byte {
	return md.ToHTML(readNavLinks())
}

func mdLinkToHTMLLink(l string) string {
	result := strings.ReplaceAll(l, ".md", ".html")
	result = strings.TrimPrefix(result, ".")
	return result
}

func addSiteURL(l string) string {
	return strings.Join([]string{conf.Output["html"].SiteURL, l}, "")
}

func replaceMDWithHTMLInLinks(h string) (string, error) {
	doc, err := gohtml.Parse(strings.NewReader(h))
	if err != nil {
		return "", err
	}

	var f func(*gohtml.Node)
	f = func(n *gohtml.Node) {
		if n.Type == gohtml.ElementNode && n.Data == "a" {
			for i := range n.Attr {
				if n.Attr[i].Key == "href" {
					link := n.Attr[i].Val
					if link[0] == '/' || link[0] == '.' {
						n.Attr[i].Val = addSiteURL(mdLinkToHTMLLink(n.Attr[i].Val))
						continue
					}
					n.Attr[i].Val = mdLinkToHTMLLink(n.Attr[i].Val)
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	var buf bytes.Buffer
	w := io.Writer(&buf)
	err = gohtml.Render(w, doc)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func getBodyChildren(htmlContent []byte) []byte {
	doc, err := gohtml.Parse(bytes.NewReader(htmlContent))
	if err != nil {
		logrus.Error(err)
		return nil
	}

	var findBody func(*gohtml.Node) *gohtml.Node
	findBody = func(n *gohtml.Node) *gohtml.Node {
		if n.Type == gohtml.ElementNode && n.Data == "body" {
			return n
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if body := findBody(c); body != nil {
				return body
			}
		}
		return nil
	}

	body := findBody(doc)
	if body == nil {
		logrus.Error("no body element found")
		return nil
	}

	var buf bytes.Buffer
	w := io.Writer(&buf)
	for child := body.FirstChild; child != nil; child = child.NextSibling {
		err := gohtml.Render(w, child)
		if err != nil {
			logrus.Error(err)
			return nil
		}
	}

	return buf.Bytes()
}
