package netflix

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func (ns *SearchQuery) parseSearchResults(bodyReader io.Reader) error {
	document, err := goquery.NewDocumentFromReader(bodyReader)
	if err != nil {
		return fmt.Errorf("failed to parse document: %w", err)
	}

	ns.Results = []SearchResult{}
	document.Find("div.web-result").Each(func(index int, element *goquery.Selection) {
		result, err := extractSearchResult(element)
		if err == nil && result != nil {
			ns.Results = append(ns.Results, *result)
		}
	})
	return nil
}

func extractSearchResult(selection *goquery.Selection) (*SearchResult, error) {
	title := strings.TrimSpace(selection.Find("h2.result__title").Text())
	link := strings.TrimSpace(selection.Find("a.result__url").Text())
	snippet := selection.Find("a.result__snippet").Text()

	id, err := strconv.Atoi(extractNetflixID(link))
	if err != nil {
		return nil, fmt.Errorf("failed to extract Netflix ID: %w", err)
	}

	if title == "" || link == "" {
		return nil, nil
	}

	return &SearchResult{
		ID:      id,
		Title:   title,
		Link:    link,
		Snippet: snippet,
	}, nil
}
