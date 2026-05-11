package proxy

import (
	"context"
	"fmt"
	"net"
	"runtime"
)

// TUNConfig — конфигурация TUN интерфейса
type TUNConfig struct {
	Name        string
	Address     string
	MTU         int
	AutoRoute   bool
	StrictRoute bool
}

// DefaultTUNConfig возвращает конфигурацию по умолчанию
func DefaultTUNConfig() TUNConfig {
	return TUNConfig{
		Name:        "olcvpn0",
		Address:     "10.0.0.1/30",
		MTU:         1500,
		AutoRoute:   true,
		StrictRoute: true,
	}
}

// IsTUNSupported проверяет поддержку TUN на платформе
func IsTUNSupported() bool {
	switch runtime.GOOS {
	case "windows", "linux", "darwin":
		return true
	case "android", "ios":
		// На мобильных платформах TUN управляется системой
		return false
	default:
		return false
	}
}

// GetTUNInterfaceIP возвращает IP адрес TUN интерфейса
func GetTUNInterfaceIP(name string) (net.IP, error) {
	iface, err := net.InterfaceByName(name)
	if err != nil {
		return nil, fmt.Errorf("interface not found: %w", err)
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return nil, fmt.Errorf("failed to get addresses: %w", err)
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok {
			if ipnet.IP.To4() != nil {
				return ipnet.IP, nil
			}
		}
	}

	return nil, fmt.Errorf("no IPv4 address found")
}

// TUNInterface — обёртка для TUN интерфейса
type TUNInterface struct {
	config TUNConfig
	ctx    context.Context
	cancel context.CancelFunc
}

// NewTUNInterface создаёт новый TUN интерфейс
func NewTUNInterface(config TUNConfig) *TUNInterface {
	return &TUNInterface{
		config: config,
	}
}

// Start запускает TUN интерфейс
func (t *TUNInterface) Start() error {
	if !IsTUNSupported() {
		return fmt.Errorf("TUN not supported on %s", runtime.GOOS)
	}

	t.ctx, t.cancel = context.WithCancel(context.Background())

	// TUN интерфейс создаётся и управляется sing-box
	// Эта структура используется только для хранения конфигурации
	return nil
}

// Stop останавливает TUN интерфейс
func (t *TUNInterface) Stop() error {
	if t.cancel != nil {
		t.cancel()
	}
	return nil
}

// GetConfig возвращает конфигурацию
func (t *TUNInterface) GetConfig() TUNConfig {
	return t.config
}
