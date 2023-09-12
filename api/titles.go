package api

import (
	"database/sql"
	"net/http"

	"github.com/monjofight/net_watch/pkg/netflix"

	"github.com/gin-gonic/gin"
)

func getTitles(c *gin.Context) {
	query := `
		SELECT 
			t.id, t.name,
			COALESCE(SUM(CASE WHEN e.watched THEN 1 ELSE 0 END), 0) as watched_episodes,
			COUNT(e.id) as total_episodes
		FROM titles t
		LEFT JOIN episodes e ON t.id = e.title_id
		GROUP BY t.id
	`

	rows, err := fetchFromDB(c, query)
	if err != nil {
		return
	}
	defer rows.Close()

	var titles = make([]struct {
		netflix.Title
		AllWatched bool
	}, 0)

	for rows.Next() {
		var title netflix.Title
		var watchedEpisodes, totalEpisodes int
		err := rows.Scan(&title.ID, &title.Name, &watchedEpisodes, &totalEpisodes)
		if handleError(c, err) {
			return
		}
		allWatched := watchedEpisodes == totalEpisodes

		titles = append(titles, struct {
			netflix.Title
			AllWatched bool
		}{
			Title:      title,
			AllWatched: allWatched,
		})
	}

	c.JSON(http.StatusOK, titles)
}

func watchEpisodesOfTitle(c *gin.Context) {
	updateEpisodesOfTitleWatchedStatus(c, true)
}

func unwatchEpisodesOfTitle(c *gin.Context) {
	updateEpisodesOfTitleWatchedStatus(c, false)
}

func updateEpisodesOfTitleWatchedStatus(c *gin.Context, watched bool) {
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

	titleId := c.Param("titleId")

	stmt, err := db.Prepare("UPDATE episodes SET watched = $1 WHERE title_id = $2")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = stmt.Exec(watched, titleId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	statusMessage := "All episodes of the title are unwatched"
	if watched {
		statusMessage = "All episodes of the title are watched"
	}
	c.JSON(http.StatusOK, gin.H{"status": statusMessage})
}
