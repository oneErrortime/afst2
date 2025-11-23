package database

import (
	"log"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"library-go/config"
	"library-go/models"
)

var DB *gorm.DB

func InitDB() {
	var err error

	// Determine database driver based on URL
	dbURL := config.DatabaseURL
	if dbURL == "" {
		dbURL = "sqlite:library.db"
	}

	// For SQLite
	if strings.Contains(strings.ToLower(dbURL), "sqlite") || strings.HasSuffix(strings.ToLower(dbURL), ".db") {
		DB, err = gorm.Open(sqlite.Open("library.db"), &gorm.Config{})
		if err != nil {
			log.Fatal("Failed to connect to SQLite database:", err)
		}
	} else {
		// For PostgreSQL
		DB, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{})
		if err != nil {
			log.Fatal("Failed to connect to PostgreSQL database:", err)
		}
	}

	// Auto-migrate the schema
	err = DB.AutoMigrate(&models.User{}, &models.Book{}, &models.Reader{}, &models.Borrow{})
	if err != nil {
		log.Fatal("Failed to migrate database schema:", err)
	}

	log.Println("Database connected and migrated successfully")
}

func CloseDB() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Println("Error getting database instance:", err)
		return
	}
	sqlDB.Close()
}