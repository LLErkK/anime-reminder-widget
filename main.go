package main

import (
	"anime-reminder/database"
	"anime-reminder/scheduler"
	"anime-reminder/ui"
	"anime-reminder/utils"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
)

func main() {
	// Create Fyne app instance
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

	// Create main window
	mainWindow := ui.NewMainWindow(myApp)

	// Setup system tray (jika tersedia)
	if desk, ok := myApp.(desktop.App); ok {
		menu := setupSystemTray(myApp, mainWindow)
		desk.SetSystemTrayMenu(menu)
		log.Println("✅ System tray initialized")
	}

	// Intercept window close button - minimize to tray instead of exit
	mainWindow.SetCloseIntercept(func() {
		mainWindow.Hide()
		log.Println("🔽 Application minimized to system tray")
	})

	// Setup cleanup when app actually quits
	myApp.Lifecycle().SetOnStopped(func() {
		log.Println("🛑 Stopping scheduler...")
		stopCh <- true
		database.CloseDB()
		log.Println("Application closed gracefully")
	})

	// Show window and run
	mainWindow.ShowAndRun()
}

func setupSystemTray(myApp fyne.App, mainWindow *ui.MainWindow) *fyne.Menu {
	appName := "AnimeReminder"

	// Menu item untuk toggle auto-start
	autoStartItem := fyne.NewMenuItem("Enable Auto-Start", nil)
	updateAutoStartText := func() {
		if utils.IsAutoStartEnabled(appName) {
			autoStartItem.Label = "✓ Auto-Start Enabled"
		} else {
			autoStartItem.Label = "Enable Auto-Start"
		}
	}
	updateAutoStartText()

	autoStartItem.Action = func() {
		if utils.IsAutoStartEnabled(appName) {
			// Disable auto-start
			err := utils.DisableAutoStart(appName)
			if err != nil {
				log.Printf("❌ Failed to disable auto-start: %v", err)
			} else {
				log.Println("✅ Auto-start disabled")
			}
		} else {
			// Enable auto-start
			err := utils.EnableAutoStart(appName)
			if err != nil {
				log.Printf("❌ Failed to enable auto-start: %v", err)
			} else {
				log.Println("✅ Auto-start enabled")
			}
		}
		updateAutoStartText()
	}

	return fyne.NewMenu("Anime Reminder",
		fyne.NewMenuItem("Show", func() {
			mainWindow.Show()
			log.Println("🔼 Application restored from tray")
		}),
		fyne.NewMenuItem("Hide", func() {
			mainWindow.Hide()
			log.Println("🔽 Application hidden to tray")
		}),
		fyne.NewMenuItemSeparator(),
		autoStartItem,
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Quit", func() {
			myApp.Quit()
		}),
	)
}
