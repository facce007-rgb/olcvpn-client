package core

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/openlibrecommunity/olcvpn/internal/engine/olcrtc"
	"github.com/openlibrecommunity/olcvpn/internal/engine/singbox"
	"github.com/openlibrecommunity/olcvpn/internal/types"
)

// ConnectionStatus представляет статус подключения
type ConnectionStatus string

const (
	StatusDisconnected ConnectionStatus = "disconnected"
	StatusConnecting   ConnectionStatus = "connecting"
	StatusConnected    ConnectionStatus = "connected"
	StatusError        ConnectionStatus = "error"
)

// Manager — главный оркестратор VPN подключений
type Manager struct {
	mu              sync.RWMutex
	storage         *Storage
	profiles        []*types.Profile
	activeProfile   *types.Profile
	status          ConnectionStatus
	statusCh        chan string
	logCh           chan string
	singboxEngine   *singbox.Engine
	olcrtcEngine    *olcrtc.Engine
	cancelFunc      context.CancelFunc
	killSwitch      *KillSwitch
	socksPort       int
	httpPort        int
}

// NewManager создаёт новый Manager
func NewManager(dataDir string) (*Manager, error) {
	storage, err := NewStorage(dataDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage: %w", err)
	}

	profiles, err := storage.LoadProfiles()
	if err != nil {
		return nil, fmt.Errorf("failed to load profiles: %w", err)
	}

	killSwitch := NewKillSwitch()

	m := &Manager{
		storage:    storage,
		profiles:   profiles,
		status:     StatusDisconnected,
		statusCh:   make(chan string, 10),
		logCh:      make(chan string, 100),
		killSwitch: killSwitch,
		socksPort:  2080,
		httpPort:   2081,
	}

	return m, nil
}

// Connect подключается к профилю
func (m *Manager) Connect(profileID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.status == StatusConnected || m.status == StatusConnecting {
		return fmt.Errorf("already connected or connecting")
	}

	profile := m.findProfile(profileID)
	if profile == nil {
		return fmt.Errorf("profile not found: %s", profileID)
	}

	m.status = StatusConnecting
	m.activeProfile = profile
	m.sendLog(fmt.Sprintf("Connecting to %s...", profile.Name))

	ctx, cancel := context.WithCancel(context.Background())
	m.cancelFunc = cancel

	var err error
	switch profile.Engine {
	case types.EngineSingBox:
		err = m.connectSingBox(ctx, profile)
	case types.EngineOlcRTC:
		err = m.connectOlcRTC(ctx, profile)
	default:
		err = fmt.Errorf("unknown engine: %s", profile.Engine)
	}

	if err != nil {
		m.status = StatusError
		m.activeProfile = nil
		m.cancelFunc = nil
		return fmt.Errorf("failed to connect: %w", err)
	}

	m.status = StatusConnected
	m.sendLog(fmt.Sprintf("Connected to %s", profile.Name))

	// Запускаем мониторинг метрик
	go m.monitorMetrics(ctx)

	return nil
}

// connectSingBox подключается через sing-box
func (m *Manager) connectSingBox(ctx context.Context, profile *types.Profile) error {
	if profile.SingBox == nil {
		return fmt.Errorf("singbox profile is nil")
	}

	m.singboxEngine = singbox.NewEngine(profile.SingBox, m.statusCh, m.logCh)
	return m.singboxEngine.Start(ctx, m.socksPort, m.httpPort)
}

// connectOlcRTC подключается через olcrtc
func (m *Manager) connectOlcRTC(ctx context.Context, profile *types.Profile) error {
	if profile.OlcRTC == nil {
		return fmt.Errorf("olcrtc profile is nil")
	}

	m.olcrtcEngine = olcrtc.NewEngine(profile.OlcRTC, m.statusCh, m.logCh)
	return m.olcrtcEngine.Start(ctx)
}

// Disconnect отключается от текущего профиля
func (m *Manager) Disconnect() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.status == StatusDisconnected {
		return nil
	}

	m.sendLog("Disconnecting...")

	// Останавливаем движки
	if m.singboxEngine != nil {
		if err := m.singboxEngine.Stop(); err != nil {
			m.sendLog(fmt.Sprintf("Error stopping singbox: %v", err))
		}
		m.singboxEngine = nil
	}

	if m.olcrtcEngine != nil {
		if err := m.olcrtcEngine.Stop(); err != nil {
			m.sendLog(fmt.Sprintf("Error stopping olcrtc: %v", err))
		}
		m.olcrtcEngine = nil
	}

	// Отменяем контекст
	if m.cancelFunc != nil {
		m.cancelFunc()
		m.cancelFunc = nil
	}

	m.status = StatusDisconnected
	m.activeProfile = nil
	m.sendLog("Disconnected")

	return nil
}

// GetStatus возвращает текущий статус
func (m *Manager) GetStatus() ConnectionStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.status
}

// GetActiveProfile возвращает активный профиль
func (m *Manager) GetActiveProfile() *types.Profile {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.activeProfile
}

// ListProfiles возвращает список всех профилей
func (m *Manager) ListProfiles() []*types.Profile {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.profiles
}

// AddProfile добавляет новый профиль
func (m *Manager) AddProfile(profile *types.Profile) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if profile.ID == "" {
		profile.ID = uuid.New().String()
	}
	if profile.CreatedAt.IsZero() {
		profile.CreatedAt = time.Now()
	}

	m.profiles = append(m.profiles, profile)
	if err := m.storage.SaveProfiles(m.profiles); err != nil {
		return fmt.Errorf("failed to save profiles: %w", err)
	}

	m.sendLog(fmt.Sprintf("Profile added: %s", profile.Name))
	return nil
}

// DeleteProfile удаляет профиль
func (m *Manager) DeleteProfile(profileID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	index := -1
	for i, p := range m.profiles {
		if p.ID == profileID {
			index = i
			break
		}
	}

	if index == -1 {
		return fmt.Errorf("profile not found: %s", profileID)
	}

	// Нельзя удалить активный профиль
	if m.activeProfile != nil && m.activeProfile.ID == profileID {
		return fmt.Errorf("cannot delete active profile")
	}

	m.profiles = append(m.profiles[:index], m.profiles[index+1:]...)
	if err := m.storage.SaveProfiles(m.profiles); err != nil {
		return fmt.Errorf("failed to save profiles: %w", err)
	}

	m.sendLog(fmt.Sprintf("Profile deleted: %s", profileID))
	return nil
}

// ImportURI импортирует профиль из URI
func (m *Manager) ImportURI(uri string) (*types.Profile, error) {
	var profile *types.Profile
	var err error

	// Пробуем парсить как sing-box URI
	profile, err = ParseSingBoxURI(uri)
	if err == nil {
		profile.ID = uuid.New().String()
		profile.CreatedAt = time.Now()
		return profile, nil
	}

	// Пробуем парсить как olcrtc URI
	profile, err = ParseOlcRTCURI(uri)
	if err == nil {
		profile.ID = uuid.New().String()
		profile.CreatedAt = time.Now()
		return profile, nil
	}

	return nil, fmt.Errorf("failed to parse URI: unsupported format")
}

// StatusChannel возвращает канал для получения обновлений статуса
func (m *Manager) StatusChannel() <-chan string {
	return m.statusCh
}

// LogChannel возвращает канал для получения логов
func (m *Manager) LogChannel() <-chan string {
	return m.logCh
}

// findProfile находит профиль по ID
func (m *Manager) findProfile(id string) *types.Profile {
	for _, p := range m.profiles {
		if p.ID == id {
			return p
		}
	}
	return nil
}

// sendLog отправляет лог в канал (неблокирующая отправка)
func (m *Manager) sendLog(msg string) {
	select {
	case m.logCh <- msg:
	default:
	}
}

// Close закрывает Manager и освобождает ресурсы
func (m *Manager) Close() error {
	if err := m.Disconnect(); err != nil {
		return err
	}
	close(m.statusCh)
	close(m.logCh)
	return nil
}

// EnableKillSwitch включает Kill Switch
func (m *Manager) EnableKillSwitch() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.killSwitch == nil {
		return fmt.Errorf("kill switch not initialized")
	}

	m.killSwitch.Enable()
	m.sendLog("Kill Switch enabled")
	return nil
}

// DisableKillSwitch выключает Kill Switch
func (m *Manager) DisableKillSwitch() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.killSwitch == nil {
		return fmt.Errorf("kill switch not initialized")
	}

	m.killSwitch.Disable()
	m.sendLog("Kill Switch disabled")
	return nil
}

// IsKillSwitchEnabled возвращает состояние Kill Switch
func (m *Manager) IsKillSwitchEnabled() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.killSwitch == nil {
		return false
	}

	return m.killSwitch.IsEnabled()
}

// GenerateConfig генерирует sing-box JSON конфиг из профиля
func (m *Manager) GenerateConfig(profile *types.Profile) (string, error) {
	if profile == nil {
		return "", fmt.Errorf("profile is nil")
	}

	if profile.Engine != types.EngineSingBox {
		return "", fmt.Errorf("only singbox profiles supported")
	}

	if profile.SingBox == nil {
		return "", fmt.Errorf("singbox profile is nil")
	}

	// Получаем SOCKS credentials из storage
	creds, err := m.storage.GetSOCKSCredentials()
	if err != nil {
		return "", fmt.Errorf("failed to get SOCKS credentials: %w", err)
	}

	// Генерируем конфиг через singbox пакет
	configJSON, err := singbox.GenerateConfig(profile.SingBox, m.socksPort, m.httpPort, creds)
	if err != nil {
		return "", fmt.Errorf("failed to generate config: %w", err)
	}

	return string(configJSON), nil
}


