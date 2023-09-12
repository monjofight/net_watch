package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleError(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}

	log.Println(err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	return true
}

func fetchFromDB(c *gin.Context, query string, args ...interface{}) (*sql.Rows, error) {
	db, _ := c.MustGet("db").(*sql.DB)
	rows, err := db.Query(query, args...)
	if handleError(c, err) {
		return nil, err
	}

	return rows, nil
}
