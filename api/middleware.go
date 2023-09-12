package api

import (
	"github.com/monjofight/net_watch/pkg/database"

	"github.com/gin-gonic/gin"
)

func databaseMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := database.InitDB()
		defer db.Close()
		c.Set("db", db)
		c.Next()
	}
}
