package utils

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"sync"
	"time"
)

// AudioPlayer struct untuk kontrol audio playback
type AudioPlayer struct {
	cmd       *exec.Cmd
	isPlaying bool
	mu        sync.Mutex
}

var globalPlayer = &AudioPlayer{}

// PlayAudio memutar file audio (blocking untuk satu instance)
func PlayAudio(filePath string, duration time.Duration) error {
	return globalPlayer.Play(filePath, duration)
}

// Play memutar audio file dengan durasi tertentu
func (ap *AudioPlayer) Play(filePath string, duration time.Duration) error {
	ap.mu.Lock()

	// Stop audio yang sedang playing
	if ap.isPlaying {
		log.Println("⚠️ Audio is already playing, stopping previous audio...")
		ap.stopInternal()
	}
	ap.mu.Unlock()

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = ap.playWindows(filePath)
	case "linux":
		cmd = ap.playLinux(filePath)
	case "darwin":
		cmd = ap.playMacOS(filePath)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	ap.mu.Lock()
	ap.cmd = cmd
	ap.isPlaying = true
	ap.mu.Unlock()

	// Start audio playback
	if err := cmd.Start(); err != nil {
		ap.mu.Lock()
		ap.isPlaying = false
		ap.mu.Unlock()
		return fmt.Errorf("failed to start audio playback: %v", err)
	}

	log.Printf("🔊 Audio playback started: %s", filePath)

	// Stop audio after duration
	if duration > 0 {
		go func() {
			time.Sleep(duration)
			ap.Stop()
		}()
	}

	// Wait for process to complete (non-blocking in goroutine)
	go func() {
		err := cmd.Wait()
		ap.mu.Lock()
		ap.isPlaying = false
		ap.mu.Unlock()

		if err != nil {
			// Ignore "killed" errors (normal when we stop manually)
			if err.Error() != "signal: killed" && err.Error() != "exit status 1" {
				log.Printf("⚠️ Audio playback ended with error: %v", err)
			} else {
				log.Println("✅ Audio playback completed")
			}
		} else {
			log.Println("✅ Audio playback completed")
		}
	}()

	return nil
}

// Stop menghentikan audio playback
func (ap *AudioPlayer) Stop() error {
	ap.mu.Lock()
	defer ap.mu.Unlock()
	return ap.stopInternal()
}

// stopInternal stops audio without locking (internal use only)
func (ap *AudioPlayer) stopInternal() error {
	if !ap.isPlaying || ap.cmd == nil || ap.cmd.Process == nil {
		return nil
	}

	err := ap.cmd.Process.Kill()
	ap.isPlaying = false

	if err != nil {
		log.Printf("⚠️ Failed to stop audio: %v", err)
	} else {
		log.Println("🛑 Audio playback stopped")
	}

	return err
}

// IsPlaying returns whether audio is currently playing
func (ap *AudioPlayer) IsPlaying() bool {
	ap.mu.Lock()
	defer ap.mu.Unlock()
	return ap.isPlaying
}

//// ===== WINDOWS =====
//func (ap *AudioPlayer) playWindows(filePath string) *exec.Cmd {
//	// Method 1: Menggunakan Windows Media Player COM object (paling reliable)
//	script := fmt.Sprintf(`
//$player = New-Object -ComObject WMPlayer.OCX
//$player.settings.volume = 100
//$player.URL = '%s'
//$player.controls.play()
//
//# Wait sampai audio selesai atau timeout
//$timeout = 60
//$elapsed = 0
//while ($player.playState -ne 1 -and $elapsed -lt $timeout) {
//    Start-Sleep -Milliseconds 500
//    $elapsed++
//}
//
//$player.controls.stop()
//$player.close()
//`, filePath)
//
//	return exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-NoProfile", "-Command", script)
//}

// ===== LINUX =====
func (ap *AudioPlayer) playLinux(filePath string) *exec.Cmd {
	// Coba beberapa audio player yang umum di Linux
	// Priority: paplay > aplay > ffplay > mpg123

	// Cek paplay (PulseAudio)
	if _, err := exec.LookPath("paplay"); err == nil {
		return exec.Command("paplay", filePath)
	}

	// Cek aplay (ALSA) - untuk .wav files
	if _, err := exec.LookPath("aplay"); err == nil {
		return exec.Command("aplay", filePath)
	}

	// Cek ffplay (FFmpeg)
	if _, err := exec.LookPath("ffplay"); err == nil {
		return exec.Command("ffplay", "-nodisp", "-autoexit", "-loglevel", "quiet", filePath)
	}

	// Cek mpg123
	if _, err := exec.LookPath("mpg123"); err == nil {
		return exec.Command("mpg123", "-q", filePath)
	}

	// Cek mplayer
	if _, err := exec.LookPath("mplayer"); err == nil {
		return exec.Command("mplayer", "-really-quiet", filePath)
	}

	// Fallback: cvlc (VLC command line)
	return exec.Command("cvlc", "--play-and-exit", "--quiet", filePath)
}

// ===== macOS =====
func (ap *AudioPlayer) playMacOS(filePath string) *exec.Cmd {
	// Menggunakan afplay (built-in di macOS)
	return exec.Command("afplay", filePath)
}

// PlayAudioAsync memutar audio secara asynchronous tanpa blocking
func PlayAudioAsync(filePath string, duration time.Duration) {
	go func() {
		err := PlayAudio(filePath, duration)
		if err != nil {
			log.Printf("❌ Failed to play audio: %v", err)
		} else {
			log.Printf("🔊 Playing audio: %s", filePath)
		}
	}()
}

// StopGlobalPlayer stops the global audio player
func StopGlobalPlayer() error {
	return globalPlayer.Stop()
}
