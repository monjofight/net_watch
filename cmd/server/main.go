package main

import (
	"log"
	"os"

	api "github.com/monjofight/net_watch/api"

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
	api.RunServer()
}
