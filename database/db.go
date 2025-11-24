package database

import (
	"anime-reminder/models"
	"fmt"
	"gorm.io/driver/sqlite"
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open(sqlite.Open("anime_reminder.db"), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	migrate()
	return
}

func GetDB() *gorm.DB {
	return db
}

func migrate() {
	err := db.AutoMigrate(&models.Anime{}, &models.RingTone{})
	if err != nil {
		fmt.Println(err)
	}
}
