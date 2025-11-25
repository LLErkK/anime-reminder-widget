package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// EnableAutoStart mengaktifkan auto-start aplikasi saat komputer menyala
func EnableAutoStart(appName string) error {
	switch runtime.GOOS {
	case "windows":
		return enableAutoStartWindows(appName)
	case "linux":
		return enableAutoStartLinux(appName)
	case "darwin":
		return enableAutoStartMacOS(appName)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

// DisableAutoStart menonaktifkan auto-start aplikasi
func DisableAutoStart(appName string) error {
	switch runtime.GOOS {
	case "windows":
		return disableAutoStartWindows(appName)
	case "linux":
		return disableAutoStartLinux(appName)
	case "darwin":
		return disableAutoStartMacOS(appName)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

// IsAutoStartEnabled mengecek apakah auto-start sudah aktif
func IsAutoStartEnabled(appName string) bool {
	switch runtime.GOOS {
	case "windows":
		return isAutoStartEnabledWindows(appName)
	case "linux":
		return isAutoStartEnabledLinux(appName)
	case "darwin":
		return isAutoStartEnabledMacOS(appName)
	default:
		return false
	}
}

// ===== WINDOWS =====
func enableAutoStartWindows(appName string) error {
	execPath, err := os.Executable()
	if err != nil {
		return err
	}

	// Dapatkan absolute path
	absPath, err := filepath.Abs(execPath)
	if err != nil {
		return err
	}

	// Buat .bat file untuk startup
	startupDir := filepath.Join(os.Getenv("APPDATA"), "Microsoft", "Windows", "Start Menu", "Programs", "Startup")
	batFile := filepath.Join(startupDir, appName+".bat")

	content := fmt.Sprintf(`@echo off
start "" "%s"`, absPath)

	err = os.WriteFile(batFile, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to create startup file: %v", err)
	}

	return nil
}

func disableAutoStartWindows(appName string) error {
	startupDir := filepath.Join(os.Getenv("APPDATA"), "Microsoft", "Windows", "Start Menu", "Programs", "Startup")
	batFile := filepath.Join(startupDir, appName+".bat")

	if _, err := os.Stat(batFile); err == nil {
		return os.Remove(batFile)
	}

	return nil
}

func isAutoStartEnabledWindows(appName string) bool {
	startupDir := filepath.Join(os.Getenv("APPDATA"), "Microsoft", "Windows", "Start Menu", "Programs", "Startup")
	batFile := filepath.Join(startupDir, appName+".bat")

	_, err := os.Stat(batFile)
	return err == nil
}

// ===== LINUX =====
func enableAutoStartLinux(appName string) error {
	execPath, err := os.Executable()
	if err != nil {
		return err
	}

	absPath, err := filepath.Abs(execPath)
	if err != nil {
		return err
	}

	// Buat .desktop file di ~/.config/autostart/
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	autostartDir := filepath.Join(homeDir, ".config", "autostart")
	if err := os.MkdirAll(autostartDir, 0755); err != nil {
		return err
	}

	desktopFile := filepath.Join(autostartDir, appName+".desktop")

	content := fmt.Sprintf(`[Desktop Entry]
Type=Application
Name=%s
Exec=%s
Hidden=false
NoDisplay=false
X-GNOME-Autostart-enabled=true`, appName, absPath)

	err = os.WriteFile(desktopFile, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to create autostart file: %v", err)
	}

	return nil
}

func disableAutoStartLinux(appName string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	desktopFile := filepath.Join(homeDir, ".config", "autostart", appName+".desktop")

	if _, err := os.Stat(desktopFile); err == nil {
		return os.Remove(desktopFile)
	}

	return nil
}

func isAutoStartEnabledLinux(appName string) bool {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	desktopFile := filepath.Join(homeDir, ".config", "autostart", appName+".desktop")
	_, err = os.Stat(desktopFile)
	return err == nil
}

// ===== macOS =====
func enableAutoStartMacOS(appName string) error {
	execPath, err := os.Executable()
	if err != nil {
		return err
	}

	absPath, err := filepath.Abs(execPath)
	if err != nil {
		return err
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// Buat .plist file di ~/Library/LaunchAgents/
	launchAgentsDir := filepath.Join(homeDir, "Library", "LaunchAgents")
	if err := os.MkdirAll(launchAgentsDir, 0755); err != nil {
		return err
	}

	plistFile := filepath.Join(launchAgentsDir, "com."+appName+".plist")

	content := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.%s</string>
    <key>ProgramArguments</key>
    <array>
        <string>%s</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <false/>
</dict>
</plist>`, appName, absPath)

	err = os.WriteFile(plistFile, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to create launch agent: %v", err)
	}

	return nil
}

func disableAutoStartMacOS(appName string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	plistFile := filepath.Join(homeDir, "Library", "LaunchAgents", "com."+appName+".plist")

	if _, err := os.Stat(plistFile); err == nil {
		return os.Remove(plistFile)
	}

	return nil
}

func isAutoStartEnabledMacOS(appName string) bool {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	plistFile := filepath.Join(homeDir, "Library", "LaunchAgents", "com."+appName+".plist")
	_, err = os.Stat(plistFile)
	return err == nil
}
