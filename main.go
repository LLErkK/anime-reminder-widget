package main

import (
	"anime-reminder/database"
	"anime-reminder/scheduler"
	"anime-reminder/ui"
	"log"

	"fyne.io/fyne/v2/app"
)

func main() {
	// Create Fyne app instance (HANYA SATU KALI!)
	myApp := app.New()

	// Initialize database (GORM + SQLite)
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize upload directories
	if err := ui.InitUploadDirectories(); err != nil {
		log.Fatalf("Failed to initialize upload directories: %v", err)
	}

	// Start anime reminder scheduler in background
	stopCh := make(chan bool)
	go scheduler.Scheduler(stopCh)
	log.Println("✅ Anime reminder scheduler started")

	// Setup cleanup when Fyne app closes
	myApp.Lifecycle().SetOnStopped(func() {
		log.Println("🛑 Stopping scheduler...")
		stopCh <- true
		database.CloseDB()
		log.Println("Application closed gracefully")
	})

	// Create and show main window (pass app instance)
	mainWindow := ui.NewMainWindow(myApp)
	mainWindow.ShowAndRun() // Blocking call
}
