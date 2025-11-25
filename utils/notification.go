package utils

import (
	"fmt"
	"os/exec"
	"runtime"
)

// SendNotification mengirim notifikasi desktop
func SendNotification(title, message string) error {
	switch runtime.GOOS {
	case "windows":
		return sendNotificationWindows(title, message)
	case "linux":
		return sendNotificationLinux(title, message)
	case "darwin":
		return sendNotificationMacOS(title, message)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

// ===== WINDOWS =====
func sendNotificationWindows(title, message string) error {
	// Menggunakan PowerShell untuk notifikasi Windows 10/11
	script := fmt.Sprintf(`
[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] | Out-Null
[Windows.UI.Notifications.ToastNotification, Windows.UI.Notifications, ContentType = WindowsRuntime] | Out-Null
[Windows.Data.Xml.Dom.XmlDocument, Windows.Data.Xml.Dom.XmlDocument, ContentType = WindowsRuntime] | Out-Null

$APP_ID = 'AnimeReminder'

$template = @"
<toast>
    <visual>
        <binding template="ToastText02">
            <text id="1">%s</text>
            <text id="2">%s</text>
        </binding>
    </visual>
</toast>
"@

$xml = New-Object Windows.Data.Xml.Dom.XmlDocument
$xml.LoadXml($template)
$toast = New-Object Windows.UI.Notifications.ToastNotification $xml
[Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier($APP_ID).Show($toast)
`, title, message)

	cmd := exec.Command("powershell", "-Command", script)
	return cmd.Run()
}

// ===== LINUX =====
func sendNotificationLinux(title, message string) error {
	// Menggunakan notify-send (tersedia di kebanyakan distro Linux)
	cmd := exec.Command("notify-send", "-u", "normal", "-t", "5000", title, message)
	return cmd.Run()
}

// ===== macOS =====
func sendNotificationMacOS(title, message string) error {
	// Menggunakan osascript untuk notifikasi macOS
	script := fmt.Sprintf(`display notification "%s" with title "%s"`, message, title)
	cmd := exec.Command("osascript", "-e", script)
	return cmd.Run()
}
