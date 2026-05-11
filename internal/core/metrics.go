package core

import (
	"context"
	"fmt"
	"net"
	"sync/atomic"
	"time"

	"github.com/openlibrecommunity/olcvpn/internal/types"
)

// metricsCollector собирает метрики в фоне
type metricsCollector struct {
	bytesUp   atomic.Int64
	bytesDown atomic.Int64
	latencyMS atomic.Int64
	cancel    context.CancelFunc
}

// GetMetrics возвращает текущие метрики
func (m *Manager) GetMetrics() types.Metrics {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.singboxEngine == nil {
		return types.Metrics{
			BytesUp:   0,
			BytesDown: 0,
			LatencyMS: 0,
		}
	}

	// TODO: Интегрировать с sing-box stats API когда будет доступен
	// Пока возвращаем заглушку
	return types.Metrics{
		BytesUp:   0,
		BytesDown: 0,
		LatencyMS: 0,
	}
}

// PingProfile измеряет latency до сервера профиля
func (m *Manager) PingProfile(profile *types.Profile) (time.Duration, error) {
	if profile == nil {
		return 0, fmt.Errorf("profile is nil")
	}

	var addr string
	var port int

	switch profile.Engine {
	case types.EngineSingBox:
		if profile.SingBox == nil {
			return 0, fmt.Errorf("singbox profile is nil")
		}
		addr = profile.SingBox.Address
		port = profile.SingBox.Port
	case types.EngineOlcRTC:
		// olcrtc не имеет прямого адреса для пинга
		return 0, fmt.Errorf("ping not supported for olcrtc profiles")
	default:
		return 0, fmt.Errorf("unknown engine: %s", profile.Engine)
	}

	start := time.Now()
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", addr, port), 5*time.Second)
	if err != nil {
		return 0, fmt.Errorf("failed to connect: %w", err)
	}
	defer conn.Close()

	latency := time.Since(start)
	return latency, nil
}
