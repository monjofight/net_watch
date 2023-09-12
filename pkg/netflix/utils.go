package netflix

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var (
	client              = &http.Client{}
	userAgent           = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36"
	duckDuckGoSearchURL = "https://html.duckduckgo.com/html/?q="
	netflixSiteQuery    = "site:https://www.netflix.com/jp "
)

func fetchURL(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func expandEscapeSequences(str string) string {
	re := regexp.MustCompile(`\\x([0-9A-Fa-f]{2})`)
	return re.ReplaceAllStringFunc(str, func(match string) string {
		code, _ := strconv.ParseUint(match[2:], 16, 8)
		return string(rune(code))
	})
}

func extractNetflixID(link string) string {
	parts := strings.Split(link, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""
}
