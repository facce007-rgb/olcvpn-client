package autostart

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"golang.org/x/sys/windows/registry"
)

// Autostart управляет автозапуском приложения
type Autostart struct {
	appName string
	appPath string
}

// New создаёт новый Autostart
func New(appName, appPath string) *Autostart {
	return &Autostart{
		appName: appName,
		appPath: appPath,
	}
}

// Enable включает автозапуск
func (a *Autostart) Enable() error {
	switch runtime.GOOS {
	case "windows":
		return a.enableWindows()
	case "linux":
		return a.enableLinux()
	case "darwin":
		return a.enableMacOS()
	default:
		return fmt.Errorf("autostart not supported on %s", runtime.GOOS)
	}
}

// Disable выключает автозапуск
func (a *Autostart) Disable() error {
	switch runtime.GOOS {
	case "windows":
		return a.disableWindows()
	case "linux":
		return a.disableLinux()
	case "darwin":
		return a.disableMacOS()
	default:
		return fmt.Errorf("autostart not supported on %s", runtime.GOOS)
	}
}

// IsEnabled проверяет включён ли автозапуск
func (a *Autostart) IsEnabled() (bool, error) {
	switch runtime.GOOS {
	case "windows":
		return a.isEnabledWindows()
	case "linux":
		return a.isEnabledLinux()
	case "darwin":
		return a.isEnabledMacOS()
	default:
		return false, fmt.Errorf("autostart not supported on %s", runtime.GOOS)
	}
}

// enableWindows включает автозапуск на Windows через Registry
func (a *Autostart) enableWindows() error {
	key, err := registry.OpenKey(registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\Run`,
		registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to open registry key: %w", err)
	}
	defer key.Close()

	if err := key.SetStringValue(a.appName, a.appPath); err != nil {
		return fmt.Errorf("failed to set registry value: %w", err)
	}

	return nil
}

// disableWindows выключает автозапуск на Windows
func (a *Autostart) disableWindows() error {
	key, err := registry.OpenKey(registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\Run`,
		registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to open registry key: %w", err)
	}
	defer key.Close()

	if err := key.DeleteValue(a.appName); err != nil && err != registry.ErrNotExist {
		return fmt.Errorf("failed to delete registry value: %w", err)
	}

	return nil
}

// isEnabledWindows проверяет автозапуск на Windows
func (a *Autostart) isEnabledWindows() (bool, error) {
	key, err := registry.OpenKey(registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\Run`,
		registry.QUERY_VALUE)
	if err != nil {
		return false, fmt.Errorf("failed to open registry key: %w", err)
	}
	defer key.Close()

	_, _, err = key.GetStringValue(a.appName)
	if err == registry.ErrNotExist {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to read registry value: %w", err)
	}

	return true, nil
}

// enableLinux включает автозапуск на Linux через systemd
func (a *Autostart) enableLinux() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	autostartDir := filepath.Join(homeDir, ".config", "autostart")
	if err := os.MkdirAll(autostartDir, 0755); err != nil {
		return fmt.Errorf("failed to create autostart directory: %w", err)
	}

	desktopFile := filepath.Join(autostartDir, a.appName+".desktop")
	content := fmt.Sprintf(`[Desktop Entry]
Type=Application
Name=%s
Exec=%s
Hidden=false
NoDisplay=false
X-GNOME-Autostart-enabled=true
`, a.appName, a.appPath)

	if err := os.WriteFile(desktopFile, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write desktop file: %w", err)
	}

	return nil
}

// disableLinux выключает автозапуск на Linux
func (a *Autostart) disableLinux() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	desktopFile := filepath.Join(homeDir, ".config", "autostart", a.appName+".desktop")
	if err := os.Remove(desktopFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove desktop file: %w", err)
	}

	return nil
}

// isEnabledLinux проверяет автозапуск на Linux
func (a *Autostart) isEnabledLinux() (bool, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false, fmt.Errorf("failed to get home directory: %w", err)
	}

	desktopFile := filepath.Join(homeDir, ".config", "autostart", a.appName+".desktop")
	_, err = os.Stat(desktopFile)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check desktop file: %w", err)
	}

	return true, nil
}

// enableMacOS включает автозапуск на macOS через LaunchAgent
func (a *Autostart) enableMacOS() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	launchAgentsDir := filepath.Join(homeDir, "Library", "LaunchAgents")
	if err := os.MkdirAll(launchAgentsDir, 0755); err != nil {
		return fmt.Errorf("failed to create LaunchAgents directory: %w", err)
	}

	plistFile := filepath.Join(launchAgentsDir, "com.olc.vpn.plist")
	content := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>Label</key>
	<string>com.olc.vpn</string>
	<key>ProgramArguments</key>
	<array>
		<string>%s</string>
	</array>
	<key>RunAtLoad</key>
	<true/>
	<key>KeepAlive</key>
	<false/>
</dict>
</plist>
`, a.appPath)

	if err := os.WriteFile(plistFile, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write plist file: %w", err)
	}

	return nil
}

// disableMacOS выключает автозапуск на macOS
func (a *Autostart) disableMacOS() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	plistFile := filepath.Join(homeDir, "Library", "LaunchAgents", "com.olc.vpn.plist")
	if err := os.Remove(plistFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove plist file: %w", err)
	}

	return nil
}

// isEnabledMacOS проверяет автозапуск на macOS
func (a *Autostart) isEnabledMacOS() (bool, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false, fmt.Errorf("failed to get home directory: %w", err)
	}

	plistFile := filepath.Join(homeDir, "Library", "LaunchAgents", "com.olc.vpn.plist")
	_, err = os.Stat(plistFile)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check plist file: %w", err)
	}

	return true, nil
}
