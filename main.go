package main

import (
	"anime-reminder/ui"
	"log"
)

func main() {
	// Initialize upload directories
	if err := ui.InitUploadDirectories(); err != nil {
		log.Fatalf("Failed to initialize upload directories: %v", err)
	}

	// Create and show main window
	mainWindow := ui.NewMainWindow()
	mainWindow.Show()
}
