package api

import (
	"os"

	"github.com/gin-gonic/gin"
)

func RunServer() {
	r := gin.Default()
	r.Use(databaseMiddleware())

	defineRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}

func defineRoutes(r *gin.Engine) {
	r.GET("/titles", getTitles)
	r.POST("/titles/:titleId/watch", watchEpisodesOfTitle)
	r.POST("/titles/:titleId/unwatch", unwatchEpisodesOfTitle)

	r.GET("/titles/:titleId/seasons", getSeasons)
	r.POST("/titles/:titleId/seasons/:seasonId/watch", watchAllEpisodesOfSeason)
	r.POST("/titles/:titleId/seasons/:seasonId/unwatch", unwatchAllEpisodesOfSeason)

	r.GET("/titles/:titleId/seasons/:seasonId/episodes", getEpisodes)
	r.POST("/titles/:titleId/seasons/:seasonId/episodes/:episodeId/watch", watchEpisode)
	r.POST("/titles/:titleId/seasons/:seasonId/episodes/:episodeId/unwatch", unwatchEpisode)

	r.POST("/update", update)
}
