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

	// Создаём Manager
	manager, err := core.NewManager(dataDir)
	if err != nil {
		log.Fatalf("Failed to create manager: %v", err)
	}
	defer manager.Close()

	// Создаём и запускаем UI в стиле v2RayTun
	app := ui.NewAppV2(fyneApp, manager)
	app.Run()
}
