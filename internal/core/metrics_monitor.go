package core

import (
	"context"
	"time"
)

// monitorMetrics периодически обновляет метрики подключения
func (m *Manager) monitorMetrics(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.updateMetrics()
		}
	}
}

// updateMetrics обновляет метрики
func (m *Manager) updateMetrics() {
	m.mu.RLock()
	profile := m.activeProfile
	m.mu.RUnlock()

	if profile == nil {
		return
	}

	// Обновляем latency если это sing-box профиль
	if profile.Engine == "singbox" && profile.SingBox != nil {
		latency, err := m.PingProfile(profile)
		if err == nil {
			// Сохраняем latency в миллисекундах
			// TODO: Добавить поле для хранения метрик в Manager
			m.sendLog("Latency: " + latency.String())
		}
	}

	// TODO: Получать BytesUp/BytesDown из sing-box stats API
	// sing-box предоставляет stats через experimental API
	// Требуется интеграция с libbox.StatsService
}
