package api

import (
	"database/sql"
	"net/http"

	"github.com/monjofight/net_watch/pkg/netflix"

	"github.com/gin-gonic/gin"
)

func getSeasons(c *gin.Context) {
	titleId := c.Param("titleId")
	query := `
		SELECT 
			s.id, s.title_id, s.name,
			COALESCE(SUM(CASE WHEN e.watched THEN 1 ELSE 0 END), 0) as watched_episodes,
			COUNT(e.id) as total_episodes
		FROM seasons s
		LEFT JOIN episodes e ON s.id = e.season_id
		WHERE s.title_id = $1
		GROUP BY s.id
	`

	rows, err := fetchFromDB(c, query, titleId)
	if err != nil {
		return
	}
	defer rows.Close()

	var seasons []struct {
		netflix.Season
		AllWatched bool
	}

	for rows.Next() {
		var season netflix.Season
		var watchedEpisodes, totalEpisodes int
		err := rows.Scan(&season.ID, &season.TitleID, &season.Name, &watchedEpisodes, &totalEpisodes)
		if handleError(c, err) {
			return
		}

		seasons = append(seasons, struct {
			netflix.Season
			AllWatched bool
		}{
			Season:     season,
			AllWatched: watchedEpisodes == totalEpisodes,
		})
	}

	c.JSON(http.StatusOK, seasons)
}

func watchAllEpisodesOfSeason(c *gin.Context) {
	updateAllEpisodesOfSeasonWatchedStatus(c, true)
}

func unwatchAllEpisodesOfSeason(c *gin.Context) {
	updateAllEpisodesOfSeasonWatchedStatus(c, false)
}

func updateAllEpisodesOfSeasonWatchedStatus(c *gin.Context, watched bool) {
	dbValue, exists := c.Get("db")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database not initialized"})
		return
	}
	db, ok := dbValue.(*sql.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}

	seasonId := c.Param("seasonId")

	stmt, err := db.Prepare("UPDATE episodes SET watched = $1 WHERE season_id = $2")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = stmt.Exec(watched, seasonId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	statusMessage := "All episodes of the season marked as unwatched"
	if watched {
		statusMessage = "All episodes of the season marked as watched"
	}
	c.JSON(http.StatusOK, gin.H{"status": statusMessage})
}
