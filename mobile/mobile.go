package mobile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/openlibrecommunity/olcvpn/internal/core"
	"github.com/openlibrecommunity/olcvpn/internal/types"
	"github.com/sagernet/sing-box/experimental/libbox"
)

// VPNCore — главный API для мобильных платформ
type VPNCore struct {
	manager         *core.Manager
	statusCallback  StatusCallback
	keychainStorage KeychainStorage
	packetFlow      PacketTunnelFlow
	platformBox     *libbox.BoxService
}

// StatusCallback — интерфейс для получения обновлений статуса
type StatusCallback interface {
	OnStatusChanged(status, message string)
}

// KeychainStorage — интерфейс для iOS Keychain (только iOS)
type KeychainStorage interface {
	Set(key, value string) error
	Get(key string) (string, error)
	Delete(key string) error
}

// PacketTunnelFlow — интерфейс для iOS PacketTunnelProvider (только iOS)
type PacketTunnelFlow interface {
	ReadPacket() []byte
	WritePacket(data []byte) bool
}

// NewVPNCore создаёт новый VPNCore
func NewVPNCore() *VPNCore {
	return &VPNCore{}
}

// Initialize инициализирует VPNCore с директорией данных
func (v *VPNCore) Initialize(dataDir string) error {
	if dataDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		dataDir = filepath.Join(homeDir, ".olcvpn")
	}

	manager, err := core.NewManager(dataDir)
	if err != nil {
		return fmt.Errorf("failed to create manager: %w", err)
	}

	v.manager = manager
	go v.monitorStatus()

	return nil
}

// Connect подключается к профилю (для Android и desktop)
func (v *VPNCore) Connect(profileJSON string) error {
	if v.manager == nil {
		return fmt.Errorf("VPNCore not initialized")
	}

	var profile types.Profile
	if err := json.Unmarshal([]byte(profileJSON), &profile); err != nil {
		return fmt.Errorf("failed to parse profile: %w", err)
	}

	return v.manager.Connect(profile.ID)
}

// ConnectIOS подключается к профилю на iOS через platform.Interface
func (v *VPNCore) ConnectIOS(profileJSON string) error {
	if v.manager == nil {
		return fmt.Errorf("VPNCore not initialized")
	}

	if v.packetFlow == nil {
		return fmt.Errorf("PacketTunnelFlow not set")
	}

	var profile types.Profile
	if err := json.Unmarshal([]byte(profileJSON), &profile); err != nil {
		return fmt.Errorf("failed to parse profile: %w", err)
	}

	// Создаём sing-box конфиг
	configJSON, err := v.manager.GenerateConfig(&profile)
	if err != nil {
		return fmt.Errorf("failed to generate config: %w", err)
	}

	// Создаём libbox.BoxService для iOS
	// libbox.BoxService использует platform.Interface для работы с пакетами
	service, err := libbox.NewService(configJSON, &iosCommandServer{
		flow: v.packetFlow,
	})
	if err != nil {
		return fmt.Errorf("failed to create box service: %w", err)
	}

	v.platformBox = service

	// Запускаем
	if err := service.Start(); err != nil {
		return fmt.Errorf("failed to start box service: %w", err)
	}

	return nil
}

// StartWithTunFd запускает VPN с TUN file descriptor (только Android)
func (v *VPNCore) StartWithTunFd(fd int) error {
	if v.manager == nil {
		return fmt.Errorf("VPNCore not initialized")
	}

	// Получаем активный профиль
	profile := v.manager.GetActiveProfile()
	if profile == nil {
		return fmt.Errorf("no active profile")
	}

	// Генерируем конфиг
	configJSON, err := v.manager.GenerateConfig(profile)
	if err != nil {
		return fmt.Errorf("failed to generate config: %w", err)
	}

	// Создаём libbox.BoxService для Android с TUN fd
	service, err := libbox.NewService(configJSON, &androidCommandServer{
		tunFd: fd,
	})
	if err != nil {
		return fmt.Errorf("failed to create box service: %w", err)
	}

	v.platformBox = service

	// Запускаем
	if err := service.Start(); err != nil {
		return fmt.Errorf("failed to start box service: %w", err)
	}

	return nil
}

// Disconnect отключается от VPN
func (v *VPNCore) Disconnect() error {
	if v.manager == nil {
		return fmt.Errorf("VPNCore not initialized")
	}

	// Останавливаем platformBox если запущен
	if v.platformBox != nil {
		if err := v.platformBox.Close(); err != nil {
			return fmt.Errorf("failed to close box service: %w", err)
		}
		v.platformBox = nil
	}

	return v.manager.Disconnect()
}

// GetStatus возвращает текущий статус
func (v *VPNCore) GetStatus() string {
	if v.manager == nil {
		return "disconnected"
	}
	return string(v.manager.GetStatus())
}

// GetBytesUp возвращает количество отправленных байт
func (v *VPNCore) GetBytesUp() int64 {
	if v.manager == nil {
		return 0
	}
	metrics := v.manager.GetMetrics()
	return metrics.BytesUp
}

// GetBytesDown возвращает количество полученных байт
func (v *VPNCore) GetBytesDown() int64 {
	if v.manager == nil {
		return 0
	}
	metrics := v.manager.GetMetrics()
	return metrics.BytesDown
}

// GetLatencyMS возвращает latency в миллисекундах
func (v *VPNCore) GetLatencyMS() int64 {
	if v.manager == nil {
		return 0
	}
	metrics := v.manager.GetMetrics()
	return metrics.LatencyMS
}

// ImportURI импортирует профиль из URI и возвращает JSON
func (v *VPNCore) ImportURI(uri string) (string, error) {
	if v.manager == nil {
		return "", fmt.Errorf("VPNCore not initialized")
	}

	profile, err := v.manager.ImportURI(uri)
	if err != nil {
		return "", err
	}

	if err := v.manager.AddProfile(profile); err != nil {
		return "", err
	}

	data, err := json.Marshal(profile)
	if err != nil {
		return "", fmt.Errorf("failed to marshal profile: %w", err)
	}

	return string(data), nil
}

// ListProfiles возвращает JSON массив профилей
func (v *VPNCore) ListProfiles() string {
	if v.manager == nil {
		return "[]"
	}

	profiles := v.manager.ListProfiles()
	data, err := json.Marshal(profiles)
	if err != nil {
		return "[]"
	}

	return string(data)
}

// DeleteProfile удаляет профиль по ID
func (v *VPNCore) DeleteProfile(id string) error {
	if v.manager == nil {
		return fmt.Errorf("VPNCore not initialized")
	}
	return v.manager.DeleteProfile(id)
}

// SetStatusCallback устанавливает callback для обновлений статуса
func (v *VPNCore) SetStatusCallback(cb StatusCallback) {
	v.statusCallback = cb
}

// SetKeychainStorage устанавливает iOS Keychain storage (только iOS)
func (v *VPNCore) SetKeychainStorage(storage KeychainStorage) {
	v.keychainStorage = storage
}

// SetPacketFlow устанавливает PacketTunnelFlow (только iOS)
func (v *VPNCore) SetPacketFlow(flow PacketTunnelFlow) {
	v.packetFlow = flow
}

// GetSupportedTransports возвращает JSON массив поддерживаемых транспортов для carrier
func (v *VPNCore) GetSupportedTransports(carrier string) string {
	if runtime.GOOS == "ios" {
		return `["datachannel"]`
	}

	switch carrier {
	case "telemost":
		return `["vp8channel", "seichannel", "datachannel"]`
	case "wbstream", "jazz":
		return `["datachannel"]`
	default:
		return `["datachannel"]`
	}
}

// monitorStatus мониторит изменения статуса и вызывает callback
func (v *VPNCore) monitorStatus() {
	if v.manager == nil {
		return
	}

	statusCh := v.manager.StatusChannel()
	logCh := v.manager.LogChannel()

	for {
		select {
		case status, ok := <-statusCh:
			if !ok {
				return
			}
			if v.statusCallback != nil {
				v.statusCallback.OnStatusChanged(status, "")
			}
		case log, ok := <-logCh:
			if !ok {
				return
			}
			if v.statusCallback != nil {
				v.statusCallback.OnStatusChanged("log", log)
			}
		}
	}
}

// iosCommandServer реализует libbox.CommandServer для iOS
type iosCommandServer struct {
	flow PacketTunnelFlow
}

func (s *iosCommandServer) ServiceReload() error {
	return nil
}

func (s *iosCommandServer) GetSystemProxyStatus() *libbox.SystemProxyStatus {
	return &libbox.SystemProxyStatus{Available: false}
}

func (s *iosCommandServer) SetSystemProxyEnabled(enabled bool) error {
	return fmt.Errorf("not supported on iOS")
}

func (s *iosCommandServer) AutoDetectInterfaceControl(fd int32) error {
	// На iOS не требуется - система автоматически исключает Extension из туннеля
	return nil
}

func (s *iosCommandServer) ClearDNSCache() {
	// DNS cache clearing не требуется на iOS
}

func (s *iosCommandServer) CloseDefaultInterfaceMonitor(monitor libbox.InterfaceUpdateListener) error {
	// Interface monitoring не используется
	return nil
}

func (s *iosCommandServer) FindConnectionOwner(ipProto int32, srcIP string, srcPort int32, destIP string, destPort int32) (int32, error) {
	// Connection owner tracking не используется на iOS
	return -1, nil
}

func (s *iosCommandServer) GetInterfaces() (libbox.NetworkInterfaceIterator, error) {
	// Interface enumeration не требуется - iOS управляет этим автоматически
	return nil, fmt.Errorf("not implemented")
}

func (s *iosCommandServer) IncludeAllNetworks() bool {
	return false
}

func (s *iosCommandServer) OpenDefaultInterfaceMonitor(listener libbox.InterfaceUpdateListener) (libbox.InterfaceUpdateListener, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *iosCommandServer) PackageNameByUid(uid int32) (string, error) {
	return "", fmt.Errorf("not implemented")
}

func (s *iosCommandServer) ReadWIFIState() *libbox.WIFIState {
	return &libbox.WIFIState{}
}

func (s *iosCommandServer) UnderNetworkExtension() bool {
	return true // Мы работаем в Network Extension на iOS
}

func (s *iosCommandServer) UsePlatformAutoDetectInterfaceControl() bool {
	return true
}

func (s *iosCommandServer) UsePlatformDefaultInterfaceMonitor() bool {
	return false
}

func (s *iosCommandServer) UsePlatformInterfaceGetter() bool {
	return false
}

func (s *iosCommandServer) OpenTun(options libbox.TunOptions) (int32, error) {
	// На iOS TUN управляется через PacketTunnelProvider, не через fd
	return -1, fmt.Errorf("OpenTun not supported on iOS - use PacketTunnelProvider")
}

func (s *iosCommandServer) SendNotification(notification *libbox.Notification) error {
	// Notifications не используются
	return nil
}

func (s *iosCommandServer) StartDefaultInterfaceMonitor(listener libbox.InterfaceUpdateListener) error {
	// Interface monitoring не используется
	return nil
}

func (s *iosCommandServer) UIDByPackageName(packageName string) (int32, error) {
	// Package UID lookup не используется на iOS
	return -1, fmt.Errorf("not implemented")
}

func (s *iosCommandServer) UseProcFS() bool {
	return false // iOS не использует procfs
}

func (s *iosCommandServer) WriteLog(message string) {
	// Логирование можно добавить позже
}

// androidCommandServer реализует libbox.CommandServer для Android
type androidCommandServer struct {
	tunFd int
}

func (s *androidCommandServer) ServiceReload() error {
	return nil
}

func (s *androidCommandServer) GetSystemProxyStatus() *libbox.SystemProxyStatus {
	return &libbox.SystemProxyStatus{Available: false}
}

func (s *androidCommandServer) SetSystemProxyEnabled(enabled bool) error {
	return fmt.Errorf("not supported on Android")
}

func (s *androidCommandServer) AutoDetectInterfaceControl(fd int32) error {
	// На Android VpnService.protect() вызывается автоматически через libbox
	return nil
}

func (s *androidCommandServer) ClearDNSCache() {
	// DNS cache clearing не требуется на Android
}

func (s *androidCommandServer) CloseDefaultInterfaceMonitor(monitor libbox.InterfaceUpdateListener) error {
	// Interface monitoring не используется
	return nil
}

func (s *androidCommandServer) FindConnectionOwner(ipProto int32, srcIP string, srcPort int32, destIP string, destPort int32) (int32, error) {
	// Connection owner tracking не используется на Android
	return -1, nil
}

func (s *androidCommandServer) GetInterfaces() (libbox.NetworkInterfaceIterator, error) {
	// Interface enumeration не требуется - Android VpnService управляет этим
	return nil, fmt.Errorf("not implemented")
}

func (s *androidCommandServer) IncludeAllNetworks() bool {
	return false
}

func (s *androidCommandServer) OpenDefaultInterfaceMonitor(listener libbox.InterfaceUpdateListener) (libbox.InterfaceUpdateListener, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *androidCommandServer) PackageNameByUid(uid int32) (string, error) {
	return "", fmt.Errorf("not implemented")
}

func (s *androidCommandServer) ReadWIFIState() *libbox.WIFIState {
	return &libbox.WIFIState{}
}

func (s *androidCommandServer) UnderNetworkExtension() bool {
	return false // На Android это не Network Extension
}

func (s *androidCommandServer) UsePlatformAutoDetectInterfaceControl() bool {
	return true
}

func (s *androidCommandServer) UsePlatformDefaultInterfaceMonitor() bool {
	return false
}

func (s *androidCommandServer) UsePlatformInterfaceGetter() bool {
	return false
}

func (s *androidCommandServer) OpenTun(options libbox.TunOptions) (int32, error) {
	// На Android TUN fd передаётся через StartWithTunFd
	return int32(s.tunFd), nil
}

func (s *androidCommandServer) SendNotification(notification *libbox.Notification) error {
	// Notifications не используются
	return nil
}

func (s *androidCommandServer) StartDefaultInterfaceMonitor(listener libbox.InterfaceUpdateListener) error {
	// Interface monitoring не используется
	return nil
}

func (s *androidCommandServer) UIDByPackageName(packageName string) (int32, error) {
	// Package UID lookup не используется
	return -1, fmt.Errorf("not implemented")
}

func (s *androidCommandServer) UseProcFS() bool {
	return false // Android не использует procfs в VPN контексте
}

func (s *androidCommandServer) WriteLog(message string) {
	// Логирование можно добавить позже
}
