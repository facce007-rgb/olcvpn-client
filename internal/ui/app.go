package ui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"github.com/openlibrecommunity/olcvpn/internal/core"
	"github.com/openlibrecommunity/olcvpn/internal/ui/screens"
	olctheme "github.com/openlibrecommunity/olcvpn/internal/ui/theme"
)

// AppV2 — приложение в стиле v2RayTun
type AppV2 struct {
	app     fyne.App
	window  fyne.Window
	manager *core.Manager

	// Экраны
	homeScreen     *screens.HomeScreen
	profilesScreen *screens.ProfilesScreen
	logsScreen     *screens.LogsScreen
	settingsScreen *screens.SettingsScreen

	// Навигация
	tabs *container.AppTabs
}

// NewAppV2 создаёт новое приложение в стиле v2RayTun
func NewAppV2(app fyne.App, manager *core.Manager) *AppV2 {
	// Применяем тёмную тему
	app.Settings().SetTheme(&olctheme.OLCTheme{})

	window := app.NewWindow("OLC VPN")
	window.Resize(fyne.NewSize(400, 600))
	window.CenterOnScreen()

	a := &AppV2{
		app:     app,
		window:  window,
		manager: manager,
	}

	// Создаём экраны
	a.homeScreen = screens.NewHomeScreen(manager, a.showProfiles, a.showLogs)
	a.profilesScreen = screens.NewProfilesScreen(manager, a.showHome)
	a.profilesScreen.SetWindow(window)
	a.logsScreen = screens.NewLogsScreen(manager, a.showHome)
	a.settingsScreen = screens.NewSettingsScreen(manager, a.showHome)

	// Создаём табы (как в v2RayTun - нижняя навигация)
	a.tabs = container.NewAppTabs(
		container.NewTabItemWithIcon("Home", theme.HomeIcon(), a.homeScreen.Content()),
		container.NewTabItemWithIcon("Profiles", theme.ListIcon(), a.profilesScreen.Content()),
		container.NewTabItemWithIcon("Logs", theme.DocumentIcon(), a.logsScreen.Content()),
		container.NewTabItemWithIcon("Settings", theme.SettingsIcon(), a.settingsScreen.Content()),
	)
	a.tabs.SetTabLocation(container.TabLocationBottom)

	window.SetContent(a.tabs)

	// Запускаем мониторинг
	go a.monitorManager()

	return a
}

// Run запускает приложение
func (a *AppV2) Run() {
	a.window.ShowAndRun()
}

// showHome показывает главный экран
func (a *AppV2) showHome() {
	a.tabs.SelectIndex(0)
}

// showProfiles показывает экран профилей
func (a *AppV2) showProfiles() {
	a.tabs.SelectIndex(1)
}

// showLogs показывает экран логов
func (a *AppV2) showLogs() {
	a.tabs.SelectIndex(2)
}

// showSettings показывает экран настроек
func (a *AppV2) showSettings() {
	a.tabs.SelectIndex(3)
}

// monitorManager мониторит изменения в Manager
func (a *AppV2) monitorManager() {
	statusCh := a.manager.StatusChannel()
	logCh := a.manager.LogChannel()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case status, ok := <-statusCh:
			if !ok {
				return
			}
			// Обновляем UI при изменении статуса
			a.homeScreen.UpdateStatus(status)

		case log, ok := <-logCh:
			if !ok {
				return
			}
			// Добавляем лог
			timestamp := time.Now().Format("15:04:05")
			a.logsScreen.AddLog(fmt.Sprintf("[%s] %s", timestamp, log))

		case <-ticker.C:
			// Обновляем метрики каждую секунду
			a.homeScreen.UpdateMetrics()
		}
	}
}

// GetManager возвращает manager
func (a *AppV2) GetManager() *core.Manager {
	return a.manager
}
