package scheduler

import (
	"anime-reminder/database"
	"anime-reminder/models"
	"log"
	"time"
)

func Scheduler(stopCh <-chan bool) {
	// Check setiap 30 detik lebih efisien
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			checkAnimeSchedule()
		case <-stopCh:
			log.Println("Scheduler stopped")
			return
		}
	}
}

func checkAnimeSchedule() {
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

	for _, anime := range animes {
		// Bandingkan waktu: ambil jam dan menit dari anime.Time
		animeHour := anime.Time.Hour()
		animeMinute := anime.Time.Minute()

		currentHour := now.Hour()
		currentMinute := now.Minute()

		// Trigger jika waktu cocok (dalam rentang 1 menit)
		if animeHour == currentHour && animeMinute == currentMinute {
			triggerReminder(anime)
		}
	}
}

func triggerReminder(anime models.Anime) {
	log.Printf("🎬 Reminder: %s is airing now!", anime.Title)
	// TODO: Implementasi notifikasi atau play ringtone berdasarkan RingToneId
	// Contoh:
	// - Kirim notifikasi desktop
	// - Play audio file
	// - Kirim push notification
}
