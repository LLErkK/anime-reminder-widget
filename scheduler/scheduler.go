package scheduler

import (
	"anime-reminder/controllers"
	"anime-reminder/database"
	"anime-reminder/models"
	"anime-reminder/utils"
	"fmt"
	"log"
	"time"
)

// Scheduler runs the anime reminder checker
func Scheduler(stopCh <-chan bool) {
	// Check setiap 30 detik lebih efisien
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	// Map untuk tracking anime yang sudah di-trigger hari ini
	triggeredToday := make(map[uint]time.Time)

	for {
		select {
		case <-ticker.C:
			checkAnimeSchedule(triggeredToday)
		case <-stopCh:
			log.Println("Scheduler stopped")
			return
		}
	}
}

func checkAnimeSchedule(triggeredToday map[uint]time.Time) {
	now := time.Now()
	today := now.Weekday().String()

	dayMap := map[string]string{
		"Monday":    "Senin",
		"Tuesday":   "Selasa",
		"Wednesday": "Rabu",
		"Thursday":  "Kamis",
		"Friday":    "Jumat",
		"Saturday":  "Sabtu",
		"Sunday":    "Minggu",
	}

	todayID := dayMap[today]

	db := database.GetDB()
	var animes []models.Anime

	// Ambil anime yang jadwalnya hari ini
	result := db.Where("day = ?", todayID).Find(&animes)
	if result.Error != nil {
		log.Printf("Error fetching anime schedule: %v", result.Error)
		return
	}

	// Cleanup triggered map jika sudah ganti hari
	currentDate := now.Format("2006-01-02")
	for animeID, lastTriggered := range triggeredToday {
		if lastTriggered.Format("2006-01-02") != currentDate {
			delete(triggeredToday, animeID)
		}
	}

	for _, anime := range animes {
		// Cek apakah anime ini sudah di-trigger hari ini
		if lastTriggered, exists := triggeredToday[anime.Id]; exists {
			// Jika sudah di-trigger dalam 1 jam terakhir, skip
			if time.Since(lastTriggered) < 1*time.Hour {
				continue
			}
		}

		// Bandingkan waktu: ambil jam dan menit dari anime.Time
		animeHour := anime.Time.Hour()
		animeMinute := anime.Time.Minute()

		currentHour := now.Hour()
		currentMinute := now.Minute()

		// Trigger jika waktu cocok (dalam rentang 1 menit)
		if animeHour == currentHour && animeMinute == currentMinute {
			triggerReminder(anime)
			triggeredToday[anime.Id] = now
		}
	}
}

func triggerReminder(anime models.Anime) {
	log.Printf("🎬 Reminder: %s is airing now!", anime.Title)

	// 1. Kirim notifikasi desktop
	title := "🎬 Anime Reminder"
	message := fmt.Sprintf("%s is airing now!\n%s at %s",
		anime.Title,
		anime.Day,
		anime.Time.Format("15:04"))

	err := utils.SendNotification(title, message)
	if err != nil {
		log.Printf("⚠️ Failed to send notification: %v", err)
	} else {
		log.Printf("✅ Notification sent for: %s", anime.Title)
	}

	// 2. Play ringtone jika ada
	if anime.RingToneId > 0 {
		ringToneController := &controllers.RingToneController{}
		ringTone, err := ringToneController.GetRingToneById(anime.RingToneId)

		if err != nil {
			log.Printf("⚠️ Failed to get ringtone: %v", err)
		} else if ringTone.SongPath != "" {
			// Play audio selama 15 detik
			duration := 30 * time.Second
			utils.PlayAudioAsync(ringTone.SongPath, duration)
			log.Printf("🔊 Playing ringtone: %s", ringTone.Name)
		}
	}

	// 3. Log reminder ke file (opsional)
	logReminder(anime)
}

func logReminder(anime models.Anime) {
	// Bisa disimpan ke file log atau database untuk history
	log.Printf("📝 Reminder logged: [%s] %s at %s",
		time.Now().Format("2006-01-02 15:04:05"),
		anime.Title,
		anime.Time.Format("15:04"))
}
