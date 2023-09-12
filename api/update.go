package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/monjofight/net_watch/pkg/database"
	"github.com/monjofight/net_watch/pkg/netflix"

	"github.com/gin-gonic/gin"
)

func update(c *gin.Context) {
	db, _ := c.MustGet("db").(*sql.DB)

	titles := database.GetTitles(db)
	if len(titles) == 0 {
		fmt.Println("No titles found in the database.")
		return
	}

	for _, title := range titles {
		fmt.Printf("Updating episodes for: %s\n", title.Name)
		netflixInstance := netflix.NewTitle(title.ID)
		netflixInstance.Fetch()
		database.SaveToDB(netflixInstance, db)
	}

	fmt.Println("All titles updated!")
	c.JSON(http.StatusOK, gin.H{"status": "All titles updated!"})
}
