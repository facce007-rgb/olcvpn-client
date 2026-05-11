package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/openlibrecommunity/olcvpn/internal/core"
	"github.com/openlibrecommunity/olcvpn/internal/ui"
)

func main() {
	// Определяем директорию для данных
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

	log.Println("OLC VPN Client started")
	log.Printf("Data directory: %s", dataDir)
	log.Printf("Loaded %d profiles", len(manager.ListProfiles()))

	// Создаём приложение с системным треем
	a := app.NewWithID("com.olc.vpn")
	a.SetIcon(resourceIconPng)

	// Создаём главное окно
	w := a.NewWindow("OLC VPN")
	w.Resize(fyne.NewSize(900, 650))
	w.SetMaster()

	// Создаём UI
	vpnUI := ui.NewApp(manager, w)

	// Настраиваем системный трей
	if desk, ok := a.(desktop.App); ok {
		setupSystemTray(desk, w, manager, vpnUI)
	}

	// Обработка закрытия окна - сворачиваем в трей
	w.SetCloseIntercept(func() {
		w.Hide()
	})

	// Показываем окно и запускаем
	w.ShowAndRun()
}

func setupSystemTray(a desktop.App, w fyne.Window, manager *core.Manager, vpnUI *ui.App) {
	menu := fyne.NewMenu("OLC VPN",
		fyne.NewMenuItem("Show", func() {
			w.Show()
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Quick Connect", func() {
			profiles := manager.ListProfiles()
			if len(profiles) > 0 {
				if err := manager.Connect(profiles[0].ID); err != nil {
					log.Printf("Quick connect failed: %v", err)
				}
			}
		}),
		fyne.NewMenuItem("Disconnect", func() {
			if err := manager.Disconnect(); err != nil {
				log.Printf("Disconnect failed: %v", err)
			}
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Status", func() {
			status := manager.GetStatus()
			metrics := manager.GetMetrics()
			msg := fmt.Sprintf("Status: %s\nUpload: %s\nDownload: %s",
				status,
				formatBytes(metrics.BytesUp),
				formatBytes(metrics.BytesDown))

			dialog := fyne.NewInformationDialog("VPN Status", msg, w)
			dialog.Show()
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Quit", func() {
			manager.Disconnect()
			os.Exit(0)
		}),
	)

	a.SetSystemTrayMenu(menu)
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// Встроенная иконка (base64 PNG)
var resourceIconPng = &fyne.StaticResource{
	StaticName: "icon.png",
	StaticContent: []byte{
		// Простая иконка VPN (16x16 PNG)
		0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a,
		// ... (минимальная PNG иконка)
	},
}
