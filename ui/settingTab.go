package ui

import (
	"anime-reminder/utils"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

const AppName = "AnimeReminder"

func (mw *MainWindow) createSettingsTab() fyne.CanvasObject {
	// Auto-start toggle
	autoStartCheck := widget.NewCheck("Run at system startup", nil)
	autoStartCheck.SetChecked(utils.IsAutoStartEnabled(AppName))

	autoStartCheck.OnChanged = func(checked bool) {
		var err error
		if checked {
			err = utils.EnableAutoStart(AppName)
			if err != nil {
				dialog.ShowError(fmt.Errorf("failed to enable auto-start: %v", err), mw.window)
				autoStartCheck.SetChecked(false)
				return
			}
			dialog.ShowInformation("Success", "Auto-start enabled! The app will run when your computer starts.", mw.window)
		} else {
			err = utils.DisableAutoStart(AppName)
			if err != nil {
				dialog.ShowError(fmt.Errorf("failed to disable auto-start: %v", err), mw.window)
				autoStartCheck.SetChecked(true)
				return
			}
			dialog.ShowInformation("Success", "Auto-start disabled.", mw.window)
		}
	}

	// Test Notification Button
	testNotifBtn := widget.NewButton("Test Notification", func() {
		err := utils.SendNotification("🎬 Test Notification", "This is a test notification from Anime Reminder!")
		if err != nil {
			dialog.ShowError(fmt.Errorf("notification test failed: %v", err), mw.window)
		} else {
			dialog.ShowInformation("Success", "Notification sent! Check your system tray.", mw.window)
		}
	})

	// Test Audio Button (test dengan ringtone yang ada)
	testAudioBtn := widget.NewButton("Test Audio Player", func() {
		ringTones, err := mw.ringToneController.GetAllRingTone()
		if err != nil || len(ringTones) == 0 {
			dialog.ShowError(fmt.Errorf("no ringtones found. Please add a ringtone first"), mw.window)
			return
		}

		// Ambil ringtone pertama untuk testing
		firstRingTone := ringTones[0]

		// Cek apakah file exists
		if _, err := os.Stat(firstRingTone.SongPath); os.IsNotExist(err) {
			dialog.ShowError(fmt.Errorf("audio file not found: %s", firstRingTone.SongPath), mw.window)
			return
		}

		utils.PlayAudioAsync(firstRingTone.SongPath, 10*time.Second)

		// Get absolute path untuk ditampilkan
		absPath, _ := filepath.Abs(firstRingTone.SongPath)

		dialog.ShowInformation("Playing Audio",
			fmt.Sprintf("Playing: %s\nDuration: 10 seconds\n\nFile: %s\n\n⚠️ If you don't hear anything:\n1. Check volume is not muted\n2. Check file format (MP3 recommended)\n3. Try playing file manually in Windows Media Player",
				firstRingTone.Name, absPath),
			mw.window)
	})

	// Stop Audio Button
	stopAudioBtn := widget.NewButton("Stop Audio", func() {
		err := utils.StopGlobalPlayer()
		if err != nil {
			dialog.ShowError(fmt.Errorf("failed to stop audio: %v", err), mw.window)
		} else {
			dialog.ShowInformation("Success", "Audio playback stopped", mw.window)
		}
	})
	stopAudioBtn.Importance = widget.WarningImportance

	// Info tentang aplikasi
	appInfo := widget.NewLabel(
		"Anime Reminder v1.0\n\n" +
			"This application helps you track anime schedules and sends reminders.\n\n" +
			"When minimized, the app runs in the system tray and continues monitoring your anime schedule.",
	)
	appInfo.Wrapping = fyne.TextWrapWord

	// Tombol untuk membersihkan database (opsional)
	clearDataBtn := widget.NewButton("Clear All Data", func() {
		dialog.ShowConfirm(
			"Clear All Data",
			"Are you sure you want to delete all anime and ringtone data? This cannot be undone!",
			func(confirmed bool) {
				if confirmed {
					// TODO: Implementasi clear all data
					dialog.ShowInformation("Info", "Feature coming soon!", mw.window)
				}
			},
			mw.window,
		)
	})
	clearDataBtn.Importance = widget.DangerImportance

	// Layout settings
	content := container.NewVBox(
		widget.NewCard("Startup Settings", "", container.NewVBox(
			autoStartCheck,
			widget.NewLabel("Enable this to automatically run the app when your computer starts."),
		)),
		widget.NewSeparator(),
		widget.NewCard("Testing", "", container.NewVBox(
			widget.NewLabel("Test reminder features:"),
			container.NewGridWithColumns(3, testNotifBtn, testAudioBtn, stopAudioBtn),
			widget.NewLabel("⚠️ Make sure you have added at least one ringtone before testing audio."),
		)),
		widget.NewSeparator(),
		widget.NewCard("About", "", appInfo),
		widget.NewSeparator(),
		widget.NewCard("Danger Zone", "", container.NewVBox(
			widget.NewLabel("⚠️ Warning: This action cannot be undone"),
			clearDataBtn,
		)),
	)

	return container.NewVScroll(content)
}
