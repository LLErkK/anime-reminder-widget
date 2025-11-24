package models

import "time"

type RingTone struct {
	Id        uint
	Name      string
	Song_path string
	CreatedAt time.Time
	UpdatedAt time.Time
}
