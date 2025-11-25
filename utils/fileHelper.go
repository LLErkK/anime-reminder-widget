package utils

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

// DeleteOldFile menghapus file lama jika ada
func DeleteOldFile(oldFilePath string) error {
	if oldFilePath == "" {
		return nil // Tidak ada file lama
	}

	// Check if file exists
	if _, err := os.Stat(oldFilePath); os.IsNotExist(err) {
		log.Printf("⚠️ Old file not found: %s (skipping)", oldFilePath)
		return nil // File tidak ada, skip
	}

	// Tunggu sebentar untuk memastikan file tidak sedang digunakan
	time.Sleep(500 * time.Millisecond)

	// Try to delete with retry mechanism
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		err := os.Remove(oldFilePath)
		if err == nil {
			log.Printf("🗑️ Old file deleted: %s", oldFilePath)
			return nil
		}

		// If file is being used, wait and retry
		if i < maxRetries-1 {
			log.Printf("⚠️ File is busy, retrying... (%d/%d)", i+1, maxRetries)
			time.Sleep(time.Second * time.Duration(i+1))
		} else {
			log.Printf("❌ Failed to delete old file %s after %d retries: %v", oldFilePath, maxRetries, err)
			return err
		}
	}

	return nil
}

// DeleteOldFileIfDifferent menghapus file lama hanya jika berbeda dengan file baru
func DeleteOldFileIfDifferent(oldFilePath, newFilePath string) error {
	if oldFilePath == "" || newFilePath == "" {
		return nil
	}

	// Jika path sama, jangan hapus
	if oldFilePath == newFilePath {
		return nil
	}

	// Convert to absolute path untuk perbandingan
	oldAbs, _ := filepath.Abs(oldFilePath)
	newAbs, _ := filepath.Abs(newFilePath)

	if oldAbs == newAbs {
		return nil
	}

	return DeleteOldFile(oldFilePath)
}

// CleanupOrphanedFiles membersihkan file yang tidak terpakai di direktori
func CleanupOrphanedFiles(directory string, usedFiles []string) error {
	// Buat map untuk cek file yang masih digunakan
	usedMap := make(map[string]bool)
	for _, file := range usedFiles {
		absPath, _ := filepath.Abs(file)
		usedMap[absPath] = true
	}

	// Scan direktori
	files, err := os.ReadDir(directory)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fullPath := filepath.Join(directory, file.Name())
		absPath, _ := filepath.Abs(fullPath)

		// Jika file tidak ada dalam daftar used, hapus
		if !usedMap[absPath] {
			err := os.Remove(fullPath)
			if err != nil {
				log.Printf("❌ Failed to delete orphaned file %s: %v", fullPath, err)
			} else {
				log.Printf("🗑️ Orphaned file deleted: %s", fullPath)
			}
		}
	}

	return nil
}
