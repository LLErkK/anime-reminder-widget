package models

import "time"

type Anime struct {
	Id         uint   `gorm:"primary_key;auto_increment"`
	Title      string `gorm:"size:255"`
	Day        string `gorm:"size:50"`
	Time       time.Time
	ImagePath  string `gorm:"size:500"`
	RingToneId uint
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Gunakan kapitalisasi konsisten
var Days = []string{
	"Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu", "Minggu",
}
