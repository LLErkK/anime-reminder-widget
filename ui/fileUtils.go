package ui

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

const (
	ImageUploadDir = "uploads/images"
	AudioUploadDir = "uploads/audio"
)

// InitUploadDirectories membuat direktori upload jika belum ada
func InitUploadDirectories() error {
	dirs := []string{ImageUploadDir, AudioUploadDir}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", dir, err)
		}
	}
	return nil
}

// CopyFile menyalin file dari source ke destination
func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Buat direktori tujuan jika belum ada
	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return err
	}

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return destFile.Sync()
}

// UploadImage meng-upload gambar dan mengembalikan path tujuan
func UploadImage(sourcePath string) (string, error) {
	if sourcePath == "" {
		return "", fmt.Errorf("source path is empty")
	}

	// Generate unique filename
	ext := filepath.Ext(sourcePath)
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("img_%d%s", timestamp, ext)
	destPath := filepath.Join(ImageUploadDir, filename)

	// Copy file
	if err := CopyFile(sourcePath, destPath); err != nil {
		return "", fmt.Errorf("failed to upload image: %v", err)
	}

	return destPath, nil
}

// UploadAudio meng-upload file audio dan mengembalikan path tujuan
func UploadAudio(sourcePath string) (string, error) {
	if sourcePath == "" {
		return "", fmt.Errorf("source path is empty")
	}

	// Generate unique filename
	ext := filepath.Ext(sourcePath)
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("audio_%d%s", timestamp, ext)
	destPath := filepath.Join(AudioUploadDir, filename)

	// Copy file
	if err := CopyFile(sourcePath, destPath); err != nil {
		return "", fmt.Errorf("failed to upload audio: %v", err)
	}

	return destPath, nil
}

// DeleteFile menghapus file
func DeleteFile(path string) error {
	if path == "" {
		return nil
	}

	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return nil // File sudah tidak ada
		}
		return err
	}

	return nil
}

// FileExists mengecek apakah file ada
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
