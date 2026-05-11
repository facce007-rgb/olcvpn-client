package screens

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/openlibrecommunity/olcvpn/internal/core"
)

// SettingsScreen — экран настроек
type SettingsScreen struct {
	manager  *core.Manager
	showHome func()
}

// NewSettingsScreen создаёт новый экран настроек
func NewSettingsScreen(manager *core.Manager, showHome func()) *SettingsScreen {
	return &SettingsScreen{
		manager:  manager,
		showHome: showHome,
	}
}

// Content возвращает содержимое экрана
func (s *SettingsScreen) Content() fyne.CanvasObject {
	// SOCKS порт
	socksPortEntry := widget.NewEntry()
	socksPortEntry.SetText("2080")
	socksPortEntry.SetPlaceHolder("SOCKS5 port")

	// HTTP порт
	httpPortEntry := widget.NewEntry()
	httpPortEntry.SetText("2081")
	httpPortEntry.SetPlaceHolder("HTTP proxy port")

	// Kill Switch
	killSwitchCheck := widget.NewCheck("Enable Kill Switch", func(checked bool) {
		// TODO: реализовать Kill Switch
	})

	// Автозапуск
	autoStartCheck := widget.NewCheck("Start on system boot", func(checked bool) {
		// TODO: реализовать автозапуск
	})

	// Сохранить настройки
	saveBtn := widget.NewButton("Save Settings", func() {
		// TODO: сохранить настройки
	})

	settingsForm := container.NewVBox(
		widget.NewLabel("Network Settings"),
		widget.NewForm(
			widget.NewFormItem("SOCKS5 Port", socksPortEntry),
			widget.NewFormItem("HTTP Proxy Port", httpPortEntry),
		),
		widget.NewSeparator(),
		widget.NewLabel("General Settings"),
		killSwitchCheck,
		autoStartCheck,
		widget.NewSeparator(),
		saveBtn,
	)

	return container.NewBorder(
		widget.NewLabel("Settings"),
		nil,
		nil,
		nil,
		settingsForm,
	)
}
