package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/monjofight/net_watch/pkg/database"
	"github.com/monjofight/net_watch/pkg/netflix"

	"github.com/joho/godotenv"
)

func init() {
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file")
		}
	} else if os.IsNotExist(err) {
		log.Println(".env file does not exist, skipping")
	} else {
		log.Println("Error checking for .env file:", err)
	}
}

func main() {
	db := database.InitDB()
	defer db.Close()

	title := getInput("Title: ")
	netflixSearch := netflix.NewSearchQuery(title)
	netflixSearch.SearchDuckDuckGo()

	for i, result := range netflixSearch.Results {
		fmt.Printf("%d. %s\n", i+1, result.Title)
	}

	optionStr := getInput("Select result number: ")
	optionNum, err := strconv.Atoi(optionStr)
	if err != nil || optionNum < 1 || optionNum > len(netflixSearch.Results) {
		fmt.Println("Invalid selection")
		return
	}

	selectedResult := netflixSearch.Results[optionNum-1]
	fmt.Println("ID:", selectedResult.ID)
	fmt.Println("You selected:")
	fmt.Println("Title:", selectedResult.Title)
	fmt.Println("Link:", selectedResult.Link)
	fmt.Println("Snippet:", selectedResult.Snippet)

	netflix := netflix.NewTitle(selectedResult.ID)
	err = netflix.Fetch()
	if err != nil {
		fmt.Println("Error fetching title:", err)
		return
	}

	database.SaveToDB(netflix, db)
}

func getInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')
	return text[:len(text)-1]
}
