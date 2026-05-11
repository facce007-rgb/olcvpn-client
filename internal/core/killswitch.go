package core

import (
	"fmt"
)

// KillSwitch — блокировка трафика при отключении VPN
type KillSwitch struct {
	enabled bool
}

// NewKillSwitch создаёт новый Kill Switch
func NewKillSwitch() *KillSwitch {
	return &KillSwitch{
		enabled: false,
	}
}

// Enable включает Kill Switch
func (ks *KillSwitch) Enable() {
	ks.enabled = true
}

// Disable выключает Kill Switch
func (ks *KillSwitch) Disable() {
	ks.enabled = false
}

// IsEnabled возвращает состояние Kill Switch
func (ks *KillSwitch) IsEnabled() bool {
	return ks.enabled
}

// GetStatus возвращает статус Kill Switch
func (ks *KillSwitch) GetStatus() string {
	if ks.enabled {
		return "enabled"
	}
	return "disabled"
}

// Toggle переключает состояние Kill Switch
func (ks *KillSwitch) Toggle() bool {
	ks.enabled = !ks.enabled
	return ks.enabled
}

// Set устанавливает состояние Kill Switch
func (ks *KillSwitch) Set(enabled bool) {
	ks.enabled = enabled
}

// String возвращает строковое представление
func (ks *KillSwitch) String() string {
	return fmt.Sprintf("KillSwitch{enabled: %v}", ks.enabled)
}
