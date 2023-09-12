package function

import (
	"fmt"
	"net/http"

	"github.com/monjofight/net_watch/pkg/database"
	"github.com/monjofight/net_watch/pkg/netflix"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("Update", Update)
}

func Update(w http.ResponseWriter, r *http.Request) {
	db := database.InitDB()
	defer db.Close()

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
}
