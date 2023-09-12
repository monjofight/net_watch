package main

import (
	"log"
	"os"

	_ "github.com/monjofight/net_watch"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
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
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}
	if err := funcframework.Start(port); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}
}
