package models

import "time"

type RingTone struct {
	Id        uint `gorm:"primary_key;auto_increment"`
	Name      string
	SongPath  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
