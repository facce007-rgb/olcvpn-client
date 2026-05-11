package screens

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/openlibrecommunity/olcvpn/internal/core"
	"github.com/openlibrecommunity/olcvpn/internal/types"
)

// HomeScreen — главный экран в стиле Hiddify/v2RayTun
type HomeScreen struct {
	manager      *core.Manager
	showProfiles func()
	showLogs     func()

	// Основные элементы
	connectionCircle *canvas.Circle
	statusText       *widget.Label
	profileCard      *fyne.Container
	profileName      *widget.Label
	profileType      *widget.Label
	connectBtn       *widget.Button

	// Метрики
	delayIndicator *widget.Label
	uploadTotal    *widget.Label
	downloadTotal  *widget.Label

	// Цвета статуса (как в Hiddify)
	colorDisconnected color.Color
	colorConnecting   color.Color
	colorConnected    color.Color
	colorError        color.Color
}

// NewHomeScreen создаёт новый главный экран
func NewHomeScreen(manager *core.Manager, showProfiles func(), showLogs func()) *HomeScreen {
	s := &HomeScreen{
		manager:      manager,
		showProfiles: showProfiles,
		showLogs:     showLogs,

		// Цвета как в Hiddify
		colorDisconnected: color.NRGBA{R: 63, G: 81, B: 181, A: 255},  // Indigo
		colorConnecting:   color.NRGBA{R: 255, G: 193, B: 7, A: 255},  // Amber
		colorConnected:    color.NRGBA{R: 46, G: 125, B: 50, A: 255},  // Green
		colorError:        color.NRGBA{R: 211, G: 47, B: 47, A: 255},  // Red
	}

	// Большой круг подключения (как в Hiddify)
	s.connectionCircle = canvas.NewCircle(s.colorDisconnected)
	s.connectionCircle.StrokeWidth = 8
	s.connectionCircle.StrokeColor = s.colorDisconnected
	s.connectionCircle.FillColor = color.Transparent

	s.statusText = widget.NewLabel("Tap to Connect")
	s.statusText.Alignment = fyne.TextAlignCenter
	s.statusText.TextStyle = fyne.TextStyle{Bold: true}

	// Карточка профиля (как в Hiddify)
	s.profileName = widget.NewLabel("No Profile Selected")
	s.profileName.TextStyle = fyne.TextStyle{Bold: true}
	s.profileType = widget.NewLabel("")

	s.profileCard = container.NewVBox(
		s.profileName,
		s.profileType,
	)

	// Кнопка подключения
	s.connectBtn = widget.NewButton("CONNECT", s.onConnect)
	s.connectBtn.Importance = widget.HighImportance

	// Индикатор задержки (как в Hiddify)
	s.delayIndicator = widget.NewLabel("- ms")
	s.delayIndicator.Alignment = fyne.TextAlignCenter
	s.delayIndicator.TextStyle = fyne.TextStyle{Bold: true}

	// Метрики трафика
	s.uploadTotal = widget.NewLabel("↑ 0 B")
	s.downloadTotal = widget.NewLabel("↓ 0 B")

	return s
}

// Content возвращает содержимое экрана
func (s *HomeScreen) Content() fyne.CanvasObject {
	// Карточка профиля сверху (как в Hiddify)
	profileSection := container.NewVBox(
		container.NewPadded(
			container.New(
				layout.NewMaxLayout(),
				canvas.NewRectangle(color.NRGBA{R: 30, G: 30, B: 30, A: 255}),
				container.NewPadded(s.profileCard),
			),
		),
	)

	// Большой круг подключения в центре (как в Hiddify)
	connectionSection := container.NewVBox(
		layout.NewSpacer(),
		container.NewCenter(
			container.NewStack(
				// Большой круг (200x200)
				container.NewPadded(
					container.NewPadded(
						container.NewPadded(s.connectionCircle),
					),
				),
				// Текст статуса в центре круга
				container.NewVBox(
					layout.NewSpacer(),
					container.NewCenter(s.statusText),
					layout.NewSpacer(),
				),
			),
		),
		// Индикатор задержки под кругом
		container.NewCenter(
			container.NewVBox(
				widget.NewSeparator(),
				s.delayIndicator,
			),
		),
		layout.NewSpacer(),
	)

	// Футер с метриками (как в Hiddify)
	footerSection := container.NewVBox(
		widget.NewSeparator(),
		container.NewHBox(
			layout.NewSpacer(),
			s.uploadTotal,
			widget.NewLabel("  •  "),
			s.downloadTotal,
			layout.NewSpacer(),
		),
	)

	// Собираем всё вместе
	return container.NewBorder(
		profileSection,
		footerSection,
		nil,
		nil,
		connectionSection,
	)
}

// UpdateStatus обновляет статус подключения
func (s *HomeScreen) UpdateStatus(status string) {
	switch status {
	case "disconnected":
		s.statusText.SetText("Tap to Connect")
		s.connectionCircle.StrokeColor = s.colorDisconnected
		s.connectionCircle.FillColor = color.Transparent
		s.connectBtn.SetText("CONNECT")
		s.connectBtn.Enable()

	case "connecting":
		s.statusText.SetText("Connecting...")
		s.connectionCircle.StrokeColor = s.colorConnecting
		s.connectionCircle.FillColor = s.colorConnecting
		s.connectBtn.SetText("CANCEL")

	case "connected":
		s.statusText.SetText("Connected")
		s.connectionCircle.StrokeColor = s.colorConnected
		s.connectionCircle.FillColor = s.colorConnected
		s.connectBtn.SetText("DISCONNECT")
		s.connectBtn.Enable()

	case "error":
		s.statusText.SetText("Connection Failed")
		s.connectionCircle.StrokeColor = s.colorError
		s.connectionCircle.FillColor = color.Transparent
		s.connectBtn.SetText("RETRY")
		s.connectBtn.Enable()
	}

	// Обновляем информацию о профиле
	if profile := s.manager.GetActiveProfile(); profile != nil {
		s.profileName.SetText(profile.Name)
		s.profileType.SetText(fmt.Sprintf("%s • %s", profile.Engine, getProtocolName(profile)))
	} else {
		s.profileName.SetText("No Profile Selected")
		s.profileType.SetText("Tap + to add profile")
	}

	s.connectionCircle.Refresh()
}

// UpdateMetrics обновляет метрики
func (s *HomeScreen) UpdateMetrics() {
	metrics := s.manager.GetMetrics()

	s.uploadTotal.SetText(fmt.Sprintf("↑ %s", formatBytes(metrics.BytesUp)))
	s.downloadTotal.SetText(fmt.Sprintf("↓ %s", formatBytes(metrics.BytesDown)))

	if metrics.LatencyMS > 0 && metrics.LatencyMS < 65000 {
		s.delayIndicator.SetText(fmt.Sprintf("%d ms", metrics.LatencyMS))
	} else {
		s.delayIndicator.SetText("- ms")
	}
}

// onConnect обрабатывает нажатие кнопки подключения
func (s *HomeScreen) onConnect() {
	status := s.manager.GetStatus()

	if status == core.StatusConnected || status == core.StatusConnecting {
		// Отключаемся
		if err := s.manager.Disconnect(); err != nil {
			s.statusText.SetText(fmt.Sprintf("Error: %v", err))
		}
	} else {
		// Подключаемся
		profiles := s.manager.ListProfiles()
		if len(profiles) == 0 {
			s.showProfiles()
			return
		}

		// Подключаемся к первому профилю
		if err := s.manager.Connect(profiles[0].ID); err != nil {
			s.statusText.SetText(fmt.Sprintf("Error: %v", err))
		}
	}
}

// getProtocolName возвращает название протокола
func getProtocolName(profile *types.Profile) string {
	if profile.SingBox != nil {
		return profile.SingBox.Protocol
	}
	if profile.OlcRTC != nil {
		return fmt.Sprintf("olcrtc/%s", profile.OlcRTC.Carrier)
	}
	return "unknown"
}

// formatBytes форматирует байты в читаемый вид
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
