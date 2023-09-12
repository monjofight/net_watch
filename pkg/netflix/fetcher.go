package netflix

import (
	"encoding/json"
	"errors"
	"log"
	"regexp"
	"strconv"
)

func (title *Title) Fetch() error {
	url := buildNetflixURL(title.ID)
	log.Println("Fetching", url)

	body, err := fetchURL(url)
	if err != nil {
		return err
	}

	title.Name, err = extractTitleName(body)
	if err != nil {
		return err
	}

	jsonData, err := extractJSONContext(body)
	if err != nil {
		return err
	}

	return title.processSections(jsonData)
}

func buildNetflixURL(titleID int) string {
	return "https://www.netflix.com/jp/title/" + strconv.Itoa(titleID)
}

func extractTitleName(body []byte) (string, error) {
	correctedText := expandEscapeSequences(string(body))
	reTitle := regexp.MustCompile(`<title>(.*?)\s*\|\s*Netflix</title>`)
	titleMatch := reTitle.FindStringSubmatch(correctedText)
	if titleMatch != nil {
		return titleMatch[1], nil
	}
	return "", errors.New("title not found")
}

func extractJSONContext(body []byte) (map[string]interface{}, error) {
	correctedText := expandEscapeSequences(string(body))
	re := regexp.MustCompile(`netflix\.reactContext\s*=\s*({[\s\S]*?});`)
	contextMatch := re.FindStringSubmatch(correctedText)

	if contextMatch == nil {
		return nil, errors.New("pattern not found")
	}

	var jsonData map[string]interface{}
	err := json.Unmarshal([]byte(contextMatch[1]), &jsonData)
	return jsonData, err
}

func (title *Title) processSections(jsonData map[string]interface{}) error {
	sectionData, err := extractSectionData(jsonData)
	if err != nil {
		return err
	}

	for _, section := range sectionData {
		if isSeasonsAndEpisodesSection(section) {
			if err := title.processSeasons(section); err != nil {
				return err
			}
		}
	}
	return nil
}

func isSeasonsAndEpisodesSection(section interface{}) bool {
	secType, ok := section.(map[string]interface{})["type"].(string)
	return ok && secType == "seasonsAndEpisodes"
}

func (title *Title) processSeasons(section interface{}) error {
	seasonsData, err := extractSeasonsData(section.(map[string]interface{}))
	if err != nil {
		return err
	}

	for _, seasonData := range seasonsData {
		season, err := processSeason(seasonData, title)
		if err != nil {
			return err
		}
		title.Seasons = append(title.Seasons, season)
	}
	return nil
}

func processSeason(seasonData interface{}, title *Title) (Season, error) {
	var season Season
	season.ID = int(seasonData.(map[string]interface{})["seasonId"].(float64))
	season.Name = seasonData.(map[string]interface{})["seasonName"].(string)
	season.TitleID = title.ID

	episodesData, err := extractEpisodeData(seasonData.(map[string]interface{}))
	if err != nil {
		return season, err
	}

	for _, episodeData := range episodesData {
		episode := processEpisode(episodeData, &season, title)
		season.Episodes = append(season.Episodes, episode)
	}
	return season, nil
}

func processEpisode(episodeData interface{}, season *Season, title *Title) Episode {
	var episode Episode
	episode.ID = int(episodeData.(map[string]interface{})["episodeId"].(float64))
	episode.Name = episodeData.(map[string]interface{})["title"].(string)
	episode.Image = episodeData.(map[string]interface{})["artworkUrl"].(string)
	episode.SeasonID = season.ID
	episode.TitleID = title.ID
	return episode
}

func extractSectionData(jsonData map[string]interface{}) ([]interface{}, error) {
	data, ok := jsonData["models"].(map[string]interface{})
	if !ok {
		return nil, errors.New("models key missing or type mismatch")
	}

	nmTitleUI, ok := data["nmTitleUI"].(map[string]interface{})
	if !ok {
		return nil, errors.New("nmTitleUI key missing or type mismatch")
	}

	dataMap, ok := nmTitleUI["data"].(map[string]interface{})
	if !ok {
		return nil, errors.New("data key missing or type mismatch")
	}

	sectionData, ok := dataMap["sectionData"].([]interface{})
	if !ok {
		return nil, errors.New("sectionData key missing or type mismatch")
	}

	return sectionData, nil
}

func extractSeasonsData(section map[string]interface{}) ([]interface{}, error) {
	data, ok := section["data"].(map[string]interface{})
	if !ok {
		return nil, errors.New("data key missing or type mismatch")
	}

	seasonsData, ok := data["seasons"].([]interface{})
	if !ok {
		return nil, errors.New("seasons key missing or type mismatch")
	}

	return seasonsData, nil
}

func extractEpisodeData(seasonData map[string]interface{}) ([]interface{}, error) {
	episodesData, ok := seasonData["episodes"].([]interface{})
	if !ok {
		return nil, errors.New("episodes key missing or type mismatch")
	}

	return episodesData, nil
}
