package controllers

import (
	"anime-reminder/database"
	"anime-reminder/models"
	"fmt"
)

type RingToneController struct{}

func (rc *RingToneController) Create(name, songPath string) (*models.RingTone, error) {
	db := database.GetDB()

	ringTone := models.RingTone{
		Name:     name,
		SongPath: songPath,
	}

	result := db.Create(&ringTone)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return &ringTone, nil
}

func (rc *RingToneController) GetAllRingTone() ([]models.RingTone, error) {
	db := database.GetDB()
	var ringTones []models.RingTone
	result := db.Find(&ringTones)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return ringTones, nil
}

func (rc *RingToneController) GetRingToneById(id uint) (*models.RingTone, error) {
	db := database.GetDB()
	var ringTone models.RingTone
	result := db.First(&ringTone, id)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return &ringTone, nil
}

// ini langsung model bukan parameter
func (rc *RingToneController) UpdateRingTone(ringTone *models.RingTone) error {
	db := database.GetDB()
	result := db.Save(ringTone)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return nil
}
func (rc *RingToneController) DeleteRingTone(id uint) error {
	db := database.GetDB()
	result := db.Delete(&models.RingTone{}, id)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return nil
}
