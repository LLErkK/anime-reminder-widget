package controllers

import (
	"anime-reminder/database"
	"anime-reminder/models"
	"anime-reminder/utils"
	"errors"
	"time"
)

type AnimeController struct{}

func (ac *AnimeController) Create(title, day, imagePath string, animeTime time.Time, ringToneId uint) (*models.Anime, error) {
	db := database.GetDB()
	validDay := false
	for _, d := range models.Days {
		if d == day {
			validDay = true
			break
		}
	}
	if !validDay {
		return nil, errors.New("invalid day")
	}

	anime := models.Anime{
		Title:      title,
		Day:        day,
		Time:       animeTime,
		ImagePath:  imagePath,
		RingToneId: ringToneId,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	result := db.Create(&anime)
	if result.Error != nil {
		return nil, result.Error
	}
	return &anime, nil
}

func (ac *AnimeController) GetAnimeById(id uint) (*models.Anime, error) {
	db := database.GetDB()
	var anime models.Anime
	result := db.First(&anime, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &anime, nil
}

func (ac *AnimeController) GetAnimeByTitle(title string) (*models.Anime, error) {
	db := database.GetDB()
	var anime models.Anime
	result := db.Where("title = ?", title).First(&anime)
	if result.Error != nil {
		return nil, result.Error
	}
	return &anime, nil
}

func (ac *AnimeController) GetAllAnimes() ([]models.Anime, error) {
	db := database.GetDB()
	var animes []models.Anime
	result := db.Find(&animes)
	if result.Error != nil {
		return nil, result.Error
	}
	return animes, nil
}

func (ac *AnimeController) UpdateAnime(title, day, imagePath string, animeTime time.Time, animeID, ringToneId uint) (*models.Anime, error) {
	db := database.GetDB()
	var anime models.Anime
	result := db.First(&anime, animeID)
	if result.Error != nil {
		return nil, result.Error
	}

	// Validasi day
	validDay := false
	for _, d := range models.Days {
		if d == day {
			validDay = true
			break
		}
	}
	if !validDay {
		return nil, errors.New("invalid day")
	}

	// HAPUS FILE LAMA jika ada file baru yang berbeda
	if imagePath != "" && imagePath != anime.ImagePath {
		err := utils.DeleteOldFileIfDifferent(anime.ImagePath, imagePath)
		if err != nil {
			// Log error tapi tetap lanjutkan update
			// (file mungkin sudah terhapus manual)
		}
	}

	// Update data
	animeInput := map[string]interface{}{
		"Title":      title,
		"Day":        day,
		"Time":       animeTime,
		"ImagePath":  imagePath,
		"RingToneId": ringToneId,
		"UpdatedAt":  time.Now(),
	}

	result = db.Model(&anime).Updates(animeInput)
	if result.Error != nil {
		return nil, result.Error
	}

	// Reload data untuk mendapatkan nilai terbaru
	db.First(&anime, animeID)
	return &anime, nil
}

func (ac *AnimeController) DeleteAnime(id uint) error {
	db := database.GetDB()

	// Ambil data anime dulu untuk mendapatkan path file
	var anime models.Anime
	result := db.First(&anime, id)
	if result.Error != nil {
		return result.Error
	}

	// Hapus file gambar jika ada
	if anime.ImagePath != "" {
		err := utils.DeleteOldFile(anime.ImagePath)
		if err != nil {
			// Log error tapi tetap lanjutkan delete
		}
	}

	// Delete dari database
	result = db.Delete(&models.Anime{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
