package main

import (
	"log"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2/app"
	"github.com/openlibrecommunity/olcvpn/internal/core"
	"github.com/openlibrecommunity/olcvpn/internal/ui"
)

func main() {
	// Создаём Fyne приложение
	fyneApp := app.NewWithID("com.olc.vpn")

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
