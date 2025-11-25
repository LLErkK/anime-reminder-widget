package database

import (
	"anime-reminder/models"
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

// init() akan dipanggil otomatis saat package di-import
func init() {
	var err error
	db, err = gorm.Open(sqlite.Open("anime_reminder.db"), &gorm.Config{})
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	migrate()
}

// InitDB sebagai wrapper untuk kompatibilitas (dipanggil dari main.go)
func InitDB() error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}
	log.Println("✅ Database initialized successfully")
	return nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return db
}

// CloseDB closes the database connection
func CloseDB() error {
	if db == nil {
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	log.Println("🔒 Database connection closed")
	return sqlDB.Close()
}

// migrate runs auto-migration for models
func migrate() {
	err := db.AutoMigrate(&models.Anime{}, &models.RingTone{})
	if err != nil {
		fmt.Println("Migration error:", err)
	} else {
		log.Println("✅ Database migration completed")
	}
}
