//go:build windows
// +build windows

package utils

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
)

// playWindowsMethod1 menggunakan Windows Media Player COM
func playWindowsMethod1(filePath string) *exec.Cmd {
	absPath, _ := filepath.Abs(filePath)
	script := fmt.Sprintf(`
$player = New-Object -ComObject WMPlayer.OCX
$player.settings.volume = 100
$player.URL = '%s'
$player.controls.play()

# Wait sampai audio selesai atau timeout
$timeout = 120
$elapsed = 0
while ($player.playState -ne 1 -and $elapsed -lt $timeout) {
    Start-Sleep -Milliseconds 500
    $elapsed++
}

$player.controls.stop()
$player.close()
`, absPath)

	return exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-NoProfile", "-Command", script)
}

// playWindowsMethod2 menggunakan SoundPlayer untuk WAV
func playWindowsMethod2(filePath string) *exec.Cmd {
	absPath, _ := filepath.Abs(filePath)
	script := fmt.Sprintf(`
Add-Type -AssemblyName System.Speech
$player = New-Object System.Media.SoundPlayer
$player.SoundLocation = '%s'
$player.PlaySync()
$player.Dispose()
`, absPath)

	return exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-NoProfile", "-Command", script)
}

// playWindowsMethod3 menggunakan presentationCore untuk MP3
func playWindowsMethod3(filePath string) *exec.Cmd {
	absPath, _ := filepath.Abs(filePath)
	script := fmt.Sprintf(`
Add-Type -AssemblyName presentationCore
$mediaPlayer = New-Object System.Windows.Media.MediaPlayer
$mediaPlayer.Open([System.Uri]::new('%s'))
$mediaPlayer.Play()

# Wait for the duration
Start-Sleep -Seconds 15

$mediaPlayer.Stop()
$mediaPlayer.Close()
`, absPath)

	return exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-NoProfile", "-Command", script)
}

// playWindowsMethod4 menggunakan VLC (jika terinstall)
func playWindowsMethod4(filePath string) *exec.Cmd {
	// Cek apakah VLC terinstall
	vlcPaths := []string{
		"C:\\Program Files\\VideoLAN\\VLC\\vlc.exe",
		"C:\\Program Files (x86)\\VideoLAN\\VLC\\vlc.exe",
	}

	for _, vlcPath := range vlcPaths {
		if _, err := exec.LookPath(vlcPath); err == nil {
			return exec.Command(vlcPath, "--play-and-exit", "--no-video", filePath)
		}
	}

	return nil
}

// playWindowsMethod5 menggunakan ffplay (jika terinstall)
func playWindowsMethod5(filePath string) *exec.Cmd {
	if _, err := exec.LookPath("ffplay"); err == nil {
		return exec.Command("ffplay", "-nodisp", "-autoexit", "-loglevel", "quiet", filePath)
	}
	return nil
}

// tryPlayWindows mencoba berbagai metode sampai ada yang berhasil
func (ap *AudioPlayer) playWindows(filePath string) *exec.Cmd {
	log.Println("🎵 Trying Windows audio playback methods...")

	// Method 1: Windows Media Player COM (Recommended)
	//log.Println("   Method 1: Windows Media Player COM")
	//cmd := playWindowsMethod1(filePath)

	// Uncomment untuk coba method lain jika method 1 gagal:

	// Method 3: PresentationCore MediaPlayer
	log.Println("   Method 3: PresentationCore MediaPlayer")
	cmd := playWindowsMethod3(filePath)

	// Method 4: VLC
	// cmd := playWindowsMethod4(filePath)
	// if cmd != nil {
	//     log.Println("   Method 4: VLC Player")
	//     return cmd
	// }

	// Method 5: FFPlay
	// cmd := playWindowsMethod5(filePath)
	// if cmd != nil {
	//     log.Println("   Method 5: FFPlay")
	//     return cmd
	// }

	return cmd
}
