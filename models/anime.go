package models

import "time"

type Anime struct {
	Id         uint   `gorm:"primary_key;auto_increment"`
	Title      string `gorm:"size:255"`
	Day        string
	Time       time.Time
	ImagePath  string
	RingToneId uint
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

var Days = []string{
	"Senin", "selasa", "Rabu", "Kamis", "Jumat", "Sabtu", "Minggu",
}
