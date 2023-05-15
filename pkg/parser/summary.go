package parser

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type PageLink struct {
	Header   string
	FilePath string
}

func ParseLinks(reader io.Reader, dirPath string) ([]PageLink, error) {
	var links []PageLink
	var re = regexp.MustCompile(`^- \[(.*)\]\((.*)\)$`)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindStringSubmatch(line)
		if len(matches) == 3 {
			link := PageLink{
				Header:   matches[1],
				FilePath: filepath.Join(dirPath, matches[2]), // Join with dirPath to preserve relative path
			}
			links = append(links, link)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return links, nil
}

func ParsePagesFromSummaryMD(filePath string) ([]Page, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var pages []Page
	linkPattern := regexp.MustCompile(`\[(.+)\]\((.+)\)`)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if matches := linkPattern.FindStringSubmatch(line); len(matches) == 3 {
			pages = append(pages, Page{
				Title: matches[1],
				URL:   matches[2],
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return pages, nil
}
