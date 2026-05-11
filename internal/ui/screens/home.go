package screens

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/openlibrecommunity/olcvpn/internal/core"
)

// HomeScreen — главный экран в стиле v2RayTun
type HomeScreen struct {
	manager         *core.Manager
	showProfiles    func()
	showLogs        func()

	// Основные элементы
	statusCircle    *canvas.Circle
	statusText      *widget.Label
	profileName     *widget.Label
	connectBtn      *widget.Button

	// Метрики
	uploadSpeed     *widget.Label
	downloadSpeed   *widget.Label
	uploadTotal     *widget.Label
	downloadTotal   *widget.Label
	latency         *widget.Label
	duration        *widget.Label

	// Цвета статуса
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

		// Цвета как в v2RayTun
		colorDisconnected: color.NRGBA{R: 158, G: 158, B: 158, A: 255}, // Серый
		colorConnecting:   color.NRGBA{R: 255, G: 193, B: 7, A: 255},   // Жёлтый
		colorConnected:    color.NRGBA{R: 76, G: 175, B: 80, A: 255},   // Зелёный
		colorError:        color.NRGBA{R: 244, G: 67, B: 54, A: 255},   // Красный
	}

	// Создаём элементы
	s.statusCircle = canvas.NewCircle(s.colorDisconnected)
	s.statusCircle.StrokeWidth = 4
	s.statusCircle.StrokeColor = s.colorDisconnected
	s.statusCircle.FillColor = color.Transparent

	s.statusText = widget.NewLabel("Disconnected")
	s.statusText.Alignment = fyne.TextAlignCenter
	s.statusText.TextStyle = fyne.TextStyle{Bold: true}

	s.profileName = widget.NewLabel("Tap to select profile")
	s.profileName.Alignment = fyne.TextAlignCenter

	s.connectBtn = widget.NewButton("CONNECT", s.onConnect)
	s.connectBtn.Importance = widget.HighImportance

	// Метрики
	s.uploadSpeed = widget.NewLabel("0 B/s")
	s.downloadSpeed = widget.NewLabel("0 B/s")
	s.uploadTotal = widget.NewLabel("0 B")
	s.downloadTotal = widget.NewLabel("0 B")
	s.latency = widget.NewLabel("-")
	s.duration = widget.NewLabel("00:00:00")

	return s
}

// Content возвращает содержимое экрана
func (s *HomeScreen) Content() fyne.CanvasObject {
	// Верхняя часть - большой круг статуса (как в v2RayTun)
	statusContainer := container.NewVBox(
		layout.NewSpacer(),
		container.NewCenter(
			container.NewStack(
				container.NewPadded(s.statusCircle),
				container.NewVBox(
					layout.NewSpacer(),
					s.statusText,
					layout.NewSpacer(),
				),
			),
		),
		container.NewCenter(s.profileName),
		layout.NewSpacer(),
	)

	// Кнопка подключения
	connectContainer := container.NewCenter(
		container.NewPadded(s.connectBtn),
	)

	// Метрики в стиле v2RayTun - две колонки
	metricsGrid := container.New(
		layout.NewGridLayout(2),

		// Левая колонка
		s.createMetricCard("↑ Upload", s.uploadSpeed, s.uploadTotal),

		// Правая колонка
		s.createMetricCard("↓ Download", s.downloadSpeed, s.downloadTotal),
	)

	// Дополнительные метрики
	extraMetrics := container.New(
		layout.NewGridLayout(2),
		s.createSimpleMetric("Latency", s.latency),
		s.createSimpleMetric("Duration", s.duration),
	)

	// Нижние кнопки
	bottomButtons := container.NewHBox(
		layout.NewSpacer(),
		widget.NewButtonWithIcon("", theme.ContentAddIcon(), s.showProfiles),
		widget.NewButtonWithIcon("", theme.DocumentIcon(), s.showLogs),
		widget.NewButtonWithIcon("", theme.SettingsIcon(), func() {
			// TODO: показать настройки
		}),
		layout.NewSpacer(),
	)

	// Собираем всё вместе
	return container.NewBorder(
		nil,
		bottomButtons,
		nil,
		nil,
		container.NewVBox(
			statusContainer,
			connectContainer,
			widget.NewSeparator(),
			metricsGrid,
			extraMetrics,
		),
	)
}

// createMetricCard создаёт карточку метрики
func (s *HomeScreen) createMetricCard(title string, speed *widget.Label, total *widget.Label) fyne.CanvasObject {
	titleLabel := widget.NewLabel(title)
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}

	speed.TextStyle = fyne.TextStyle{Bold: true}
	total.Alignment = fyne.TextAlignCenter

	return container.NewVBox(
		titleLabel,
		speed,
		total,
	)
}

// createSimpleMetric создаёт простую метрику
func (s *HomeScreen) createSimpleMetric(title string, value *widget.Label) fyne.CanvasObject {
	titleLabel := widget.NewLabel(title)
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}

	return container.NewVBox(
		titleLabel,
		value,
	)
}

// UpdateStatus обновляет статус подключения
func (s *HomeScreen) UpdateStatus(status string) {
	switch status {
	case "disconnected":
		s.statusText.SetText("Disconnected")
		s.statusCircle.StrokeColor = s.colorDisconnected
		s.statusCircle.FillColor = color.Transparent
		s.connectBtn.SetText("CONNECT")
		s.connectBtn.Enable()

	case "connecting":
		s.statusText.SetText("Connecting...")
		s.statusCircle.StrokeColor = s.colorConnecting
		s.statusCircle.FillColor = s.colorConnecting
		s.connectBtn.SetText("CANCEL")

	case "connected":
		s.statusText.SetText("Connected")
		s.statusCircle.StrokeColor = s.colorConnected
		s.statusCircle.FillColor = s.colorConnected
		s.connectBtn.SetText("DISCONNECT")
		s.connectBtn.Enable()

	case "error":
		s.statusText.SetText("Error")
		s.statusCircle.StrokeColor = s.colorError
		s.statusCircle.FillColor = color.Transparent
		s.connectBtn.SetText("RETRY")
		s.connectBtn.Enable()
	}

	// Обновляем название профиля
	if profile := s.manager.GetActiveProfile(); profile != nil {
		s.profileName.SetText(profile.Name)
	} else {
		s.profileName.SetText("Tap to select profile")
	}

	s.statusCircle.Refresh()
}

// UpdateMetrics обновляет метрики
func (s *HomeScreen) UpdateMetrics() {
	metrics := s.manager.GetMetrics()

	// TODO: Вычислить скорость (нужно хранить предыдущие значения)
	s.uploadSpeed.SetText("0 B/s")
	s.downloadSpeed.SetText("0 B/s")

	s.uploadTotal.SetText(formatBytes(metrics.BytesUp))
	s.downloadTotal.SetText(formatBytes(metrics.BytesDown))

	if metrics.LatencyMS > 0 {
		s.latency.SetText(fmt.Sprintf("%d ms", metrics.LatencyMS))
	} else {
		s.latency.SetText("-")
	}

	// TODO: Вычислить длительность подключения
	s.duration.SetText("00:00:00")
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
