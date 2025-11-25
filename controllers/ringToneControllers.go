package controllers

import (
	"anime-reminder/database"
	"anime-reminder/models"
	"anime-reminder/utils"
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
		return nil, result.Error
	}
	return &ringTone, nil
}

func (rc *RingToneController) GetAllRingTone() ([]models.RingTone, error) {
	db := database.GetDB()
	var ringTones []models.RingTone
	result := db.Find(&ringTones)
	if result.Error != nil {
		fmt.Println(result.Error)
		return nil, result.Error
	}
	return ringTones, nil
}

func (rc *RingToneController) GetRingToneById(id uint) (*models.RingTone, error) {
	db := database.GetDB()
	var ringTone models.RingTone
	result := db.First(&ringTone, id)
	if result.Error != nil {
		fmt.Println(result.Error)
		return nil, result.Error
	}
	return &ringTone, nil
}

// UpdateRingTone - dengan penghapusan file lama
func (rc *RingToneController) UpdateRingTone(ringTone *models.RingTone) error {
	db := database.GetDB()

	// Ambil data lama untuk mendapatkan path file lama
	var oldRingTone models.RingTone
	result := db.First(&oldRingTone, ringTone.Id)
	if result.Error != nil {
		fmt.Println(result.Error)
		return result.Error
	}

	// HAPUS FILE LAMA jika ada file baru yang berbeda
	if ringTone.SongPath != "" && ringTone.SongPath != oldRingTone.SongPath {
		// Stop audio player jika sedang play file ini
		utils.StopGlobalPlayer()

		err := utils.DeleteOldFileIfDifferent(oldRingTone.SongPath, ringTone.SongPath)
		if err != nil {
			// Log error tapi tetap lanjutkan update
			fmt.Printf("Warning: Failed to delete old ringtone file: %v\n", err)
		}
	}

	// Update data
	result = db.Save(ringTone)
	if result.Error != nil {
		fmt.Println(result.Error)
		return result.Error
	}
	return nil
}

func (rc *RingToneController) DeleteRingTone(id uint) error {
	db := database.GetDB()

	// Ambil data ringtone dulu untuk mendapatkan path file
	var ringTone models.RingTone
	result := db.First(&ringTone, id)
	if result.Error != nil {
		fmt.Println(result.Error)
		return result.Error
	}

	// Stop audio player jika sedang play file ini
	utils.StopGlobalPlayer()

	// Hapus file audio jika ada
	if ringTone.SongPath != "" {
		err := utils.DeleteOldFile(ringTone.SongPath)
		if err != nil {
			// Log error tapi tetap lanjutkan delete
			fmt.Printf("Warning: Failed to delete ringtone file: %v\n", err)
		}
	}

	// Delete dari database
	result = db.Delete(&models.RingTone{}, id)
	if result.Error != nil {
		fmt.Println(result.Error)
		return result.Error
	}
	return nil
}
