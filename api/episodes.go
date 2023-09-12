package api

import (
	"database/sql"
	"net/http"

	"github.com/monjofight/net_watch/pkg/netflix"

	"github.com/gin-gonic/gin"
)

func getEpisodes(c *gin.Context) {
	seasonId := c.Param("seasonId")
	query := `
		SELECT id, title_id, season_id, name, image, watched 
		FROM episodes 
		WHERE season_id = $1
	`

	rows, err := fetchFromDB(c, query, seasonId)
	if err != nil {
		return
	}
	defer rows.Close()

	var episodes []netflix.Episode
	for rows.Next() {
		var episode netflix.Episode
		err := rows.Scan(&episode.ID, &episode.TitleID, &episode.SeasonID, &episode.Name, &episode.Image, &episode.Watched)
		if handleError(c, err) {
			return
		}
		episodes = append(episodes, episode)
	}

	c.JSON(http.StatusOK, episodes)
}

func watchEpisode(c *gin.Context) {
	updateEpisodeWatchedStatus(c, true)
}

func unwatchEpisode(c *gin.Context) {
	updateEpisodeWatchedStatus(c, false)
}

func updateEpisodeWatchedStatus(c *gin.Context, watched bool) {
	db, _ := c.MustGet("db").(*sql.DB)

	episodeId := c.Param("episodeId")

	stmt, err := db.Prepare("UPDATE episodes SET watched = $1 WHERE id = $2")
	if handleError(c, err) {
		return
	}

	_, err = stmt.Exec(watched, episodeId)
	if handleError(c, err) {
		return
	}

	statusMessage := "Episode marked as unwatched"
	if watched {
		statusMessage = "Episode marked as watched"
	}
	c.JSON(http.StatusOK, gin.H{"status": statusMessage})
}
