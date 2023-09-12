package netflix

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

func (ns *SearchQuery) SearchDuckDuckGo() error {
	if ns.Query == "" {
		return errors.New("search query cannot be empty")
	}

	query := url.QueryEscape(netflixSiteQuery + ns.Query)
	searchURL := duckDuckGoSearchURL + query

	fmt.Println("Searching:", searchURL)

	request, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	request.Header.Set("User-Agent", userAgent)
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return fmt.Errorf("failed to retrieve the webpage. Status: %d %s", response.StatusCode, http.StatusText(response.StatusCode))
	}

	return ns.parseSearchResults(response.Body)
}
