package initializers

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		log.Fatal("DB_URL environment variable is not set")
	}

	log.Printf("DSN: %s", dsn) // Debugging line

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Successfully connected to the database")
}
