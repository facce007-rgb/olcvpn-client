package screens

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/openlibrecommunity/olcvpn/internal/core"
)

// LogsScreen — экран логов
type LogsScreen struct {
	manager  *core.Manager
	showHome func()
	logsList *widget.List
	logs     []string
}

// NewLogsScreen создаёт новый экран логов
func NewLogsScreen(manager *core.Manager, showHome func()) *LogsScreen {
	s := &LogsScreen{
		manager:  manager,
		showHome: showHome,
		logs:     make([]string, 0),
	}

	s.logsList = widget.NewList(
		func() int {
			return len(s.logs)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Log entry")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id < len(s.logs) {
				obj.(*widget.Label).SetText(s.logs[id])
			}
		},
	)

	return s
}

// Content возвращает содержимое экрана
func (s *LogsScreen) Content() fyne.CanvasObject {
	clearBtn := widget.NewButton("Clear Logs", s.onClearLogs)

	return container.NewBorder(
		widget.NewLabel("Logs"),
		clearBtn,
		nil,
		nil,
		s.logsList,
	)
}

// AddLog добавляет новую запись в лог
func (s *LogsScreen) AddLog(log string) {
	s.logs = append(s.logs, log)

	// Ограничиваем количество логов
	if len(s.logs) > 1000 {
		s.logs = s.logs[len(s.logs)-1000:]
	}

	s.logsList.Refresh()

	// Автоскролл к последнему элементу
	s.logsList.ScrollToBottom()
}

// onClearLogs очищает логи
func (s *LogsScreen) onClearLogs() {
	s.logs = make([]string, 0)
	s.logsList.Refresh()
}
