package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/muchira007/jambo-green-go/initializers"
	"github.com/muchira007/jambo-green-go/models"
)

func init() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Connect to the database
	initializers.ConnectToDB()
}

func main() {
	log.Println("Starting migration")
	err := initializers.DB.AutoMigrate(&models.Product{})
	if err != nil {
		log.Fatalf("Failed to migrate database schema: %v", err)
	}
	log.Println("Migration completed")
}
