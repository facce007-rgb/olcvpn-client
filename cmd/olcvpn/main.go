package main

import (
	_ "embed"
	"log"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/openlibrecommunity/olcvpn/internal/core"
	"github.com/openlibrecommunity/olcvpn/internal/ui"
)

//go:embed icon.png
var iconData []byte

func main() {
	// Создаём Fyne приложение
	fyneApp := app.NewWithID("com.olc.vpn")

	// Устанавливаем иконку
	if len(iconData) > 0 {
		icon := fyne.NewStaticResource("icon.png", iconData)
		fyneApp.SetIcon(icon)
	}

	// Определяем директорию данных
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get home directory: %v", err)
	}
	dataDir := filepath.Join(homeDir, ".olcvpn")

	log.Println("OLC VPN Client starting...")
	log.Printf("Data directory: %s", dataDir)

	// Создаём Manager
	manager, err := core.NewManager(dataDir)
	if err != nil {
		log.Fatalf("Failed to create manager: %v", err)
	}
	defer manager.Close()

	log.Printf("Loaded %d profiles", len(manager.ListProfiles()))

	// Создаём и запускаем UI в стиле v2RayTun
	vpnApp := ui.NewAppV2(fyneApp, manager)
	vpnApp.Run()
}
