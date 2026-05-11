package screens

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/openlibrecommunity/olcvpn/internal/core"
	"github.com/openlibrecommunity/olcvpn/internal/types"
)

// ProfilesScreen — экран профилей в стиле v2RayTun
type ProfilesScreen struct {
	manager      *core.Manager
	showHome     func()
	window       fyne.Window
	profilesList *widget.List
	profiles     []*types.Profile
}

// NewProfilesScreen создаёт новый экран профилей
func NewProfilesScreen(manager *core.Manager, showHome func()) *ProfilesScreen {
	s := &ProfilesScreen{
		manager:  manager,
		showHome: showHome,
		profiles: manager.ListProfiles(),
	}

	s.profilesList = widget.NewList(
		func() int {
			return len(s.profiles)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewIcon(theme.ComputerIcon()),
				widget.NewLabel("Profile Name"),
				widget.NewLabel("vless://..."),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			profile := s.profiles[id]
			container := obj.(*fyne.Container)

			// Иконка в зависимости от типа
			icon := container.Objects[0].(*widget.Icon)
			if profile.Engine == types.EngineOlcRTC {
				icon.SetResource(theme.MediaVideoIcon())
			} else {
				icon.SetResource(theme.ComputerIcon())
			}

			// Название
			nameLabel := container.Objects[1].(*widget.Label)
			nameLabel.SetText(profile.Name)

			// Адрес
			addrLabel := container.Objects[2].(*widget.Label)
			if profile.SingBox != nil {
				addrLabel.SetText(fmt.Sprintf("%s://%s:%d",
					profile.SingBox.Protocol,
					profile.SingBox.Address,
					profile.SingBox.Port))
			} else if profile.OlcRTC != nil {
				addrLabel.SetText(fmt.Sprintf("olcrtc://%s/%s",
					profile.OlcRTC.Carrier,
					profile.OlcRTC.Transport))
			}
		},
	)

	s.profilesList.OnSelected = func(id widget.ListItemID) {
		s.onProfileSelected(id)
	}

	return s
}

// Content возвращает содержимое экрана
func (s *ProfilesScreen) Content() fyne.CanvasObject {
	// Заголовок
	header := container.NewHBox(
		widget.NewButtonWithIcon("", theme.NavigateBackIcon(), s.showHome),
		widget.NewLabel("Profiles"),
		widget.NewLabel(""), // Spacer
	)

	// Кнопки действий
	addBtn := widget.NewButtonWithIcon("Add", theme.ContentAddIcon(), s.onAddProfile)
	importBtn := widget.NewButtonWithIcon("Import", theme.DownloadIcon(), s.onImportProfile)
	scanBtn := widget.NewButtonWithIcon("Scan QR", theme.SearchIcon(), s.onScanQR)

	actions := container.NewHBox(
		addBtn,
		importBtn,
		scanBtn,
	)

	// Список профилей
	return container.NewBorder(
		container.NewVBox(header, actions),
		nil,
		nil,
		nil,
		s.profilesList,
	)
}

// onProfileSelected обрабатывает выбор профиля
func (s *ProfilesScreen) onProfileSelected(id widget.ListItemID) {
	profile := s.profiles[id]

	// Показываем диалог с действиями
	connectBtn := widget.NewButton("Connect", func() {
		if err := s.manager.Connect(profile.ID); err != nil {
			dialog.ShowError(err, s.window)
		} else {
			s.showHome()
		}
	})

	editBtn := widget.NewButton("Edit", func() {
		// TODO: показать диалог редактирования
	})

	deleteBtn := widget.NewButton("Delete", func() {
		dialog.ShowConfirm("Delete Profile",
			fmt.Sprintf("Delete profile '%s'?", profile.Name),
			func(ok bool) {
				if ok {
					if err := s.manager.DeleteProfile(profile.ID); err != nil {
						dialog.ShowError(err, s.window)
					} else {
						s.refreshProfiles()
					}
				}
			}, s.window)
	})

	content := container.NewVBox(
		widget.NewLabel(profile.Name),
		connectBtn,
		editBtn,
		deleteBtn,
	)

	dialog.ShowCustom("Profile Actions", "Close", content, s.window)
}

// onAddProfile обрабатывает добавление профиля
func (s *ProfilesScreen) onAddProfile() {
	// TODO: показать диалог добавления профиля
	dialog.ShowInformation("Add Profile", "Manual profile creation coming soon", s.window)
}

// onImportProfile обрабатывает импорт профиля
func (s *ProfilesScreen) onImportProfile() {
	uriEntry := widget.NewEntry()
	uriEntry.SetPlaceHolder("vless://... or olcrtc://...")

	dialog.ShowForm("Import Profile", "Import", "Cancel",
		[]*widget.FormItem{
			widget.NewFormItem("URI", uriEntry),
		},
		func(ok bool) {
			if !ok {
				return
			}

			profile, err := s.manager.ImportURI(uriEntry.Text)
			if err != nil {
				dialog.ShowError(err, s.window)
				return
			}

			if err := s.manager.AddProfile(profile); err != nil {
				dialog.ShowError(err, s.window)
				return
			}

			s.refreshProfiles()
			dialog.ShowInformation("Success",
				fmt.Sprintf("Profile '%s' imported", profile.Name),
				s.window)
		}, s.window)
}

// onScanQR обрабатывает сканирование QR-кода
func (s *ProfilesScreen) onScanQR() {
	// TODO: реализовать сканирование QR через веб-камеру
	dialog.ShowInformation("Scan QR", "QR code scanning coming soon", s.window)
}

// refreshProfiles обновляет список профилей
func (s *ProfilesScreen) refreshProfiles() {
	s.profiles = s.manager.ListProfiles()
	s.profilesList.Refresh()
}

// SetWindow устанавливает окно для диалогов
func (s *ProfilesScreen) SetWindow(window fyne.Window) {
	s.window = window
}
