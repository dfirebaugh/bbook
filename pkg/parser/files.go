package parser

import (
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Page struct {
	Title   string
	Content string
	URL     string
	HTML    template.HTML
}

func ParseDir(srcDir string) []Page {
	var posts []Page

	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		if filepath.Ext(info.Name()) != ".md" {
			return nil
		}

		relPath, _ := filepath.Rel(srcDir, path)
		contentBytes, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}

		post := Page{
			Title:   strings.TrimSuffix(info.Name(), ".md"),
			Content: string(contentBytes),
			URL:     strings.TrimSuffix(relPath, filepath.Ext(relPath)) + ".html",
		}
		posts = append(posts, post)

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	return posts
}
