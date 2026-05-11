package olcrtc

import (
	"context"
	"fmt"

	"github.com/openlibrecommunity/olcvpn/internal/types"
)

// Engine управляет olcrtc процессом
type Engine struct {
	profile  *types.OlcRTCProfile
	cancel   context.CancelFunc
	statusCh chan<- string
	logCh    chan<- string
}

// NewEngine создаёт новый olcrtc движок
func NewEngine(profile *types.OlcRTCProfile, statusCh chan<- string, logCh chan<- string) *Engine {
	return &Engine{
		profile:  profile,
		statusCh: statusCh,
		logCh:    logCh,
	}
}

// Start запускает olcrtc
func (e *Engine) Start(ctx context.Context) error {
	// Валидация профиля
	if err := e.validateProfile(); err != nil {
		return fmt.Errorf("invalid profile: %w", err)
	}

	// olcrtc пока не экспортирует публичный API
	// Используем stub реализацию до тех пор пока не появится pkg/client
	// TODO: Интеграция с olcrtc когда появится публичный API

	// Создаём контекст с отменой
	ctx, cancel := context.WithCancel(ctx)
	e.cancel = cancel

	// Запускаем olcrtc в горутине
	go e.run(ctx)

	e.sendLog(fmt.Sprintf("olcrtc engine started: %s/%s", e.profile.Carrier, e.profile.Transport))
	e.sendStatus("connected")

	return nil
}

// run выполняет olcrtc процесс
func (e *Engine) run(ctx context.Context) {
	e.sendLog(fmt.Sprintf("olcrtc: connecting to %s room %s", e.profile.Carrier, e.profile.RoomID))

	// TODO: Здесь будет реальная интеграция с olcrtc
	// Варианты интеграции:
	// 1. Когда olcrtc экспортирует pkg/client — использовать напрямую
	// 2. Запуск olcrtc CLI через exec.Command
	// 3. gomobile binding если olcrtc предоставит его

	e.sendLog("olcrtc: stub mode - waiting for public API from openlibrecommunity/olcrtc")

	<-ctx.Done()

	e.sendLog("olcrtc: stopping")
}

// Stop останавливает olcrtc
func (e *Engine) Stop() error {
	if e.cancel != nil {
		e.cancel()
		e.cancel = nil
	}

	e.sendLog("olcrtc engine stopped")
	e.sendStatus("disconnected")

	return nil
}

// validateProfile проверяет корректность профиля
func (e *Engine) validateProfile() error {
	if e.profile.Carrier == "" {
		return fmt.Errorf("carrier is required")
	}
	if e.profile.Transport == "" {
		return fmt.Errorf("transport is required")
	}
	if e.profile.RoomID == "" {
		return fmt.Errorf("room ID is required")
	}
	if e.profile.Key == "" {
		return fmt.Errorf("key is required")
	}
	if e.profile.ClientID == "" {
		return fmt.Errorf("client ID is required")
	}

	// Валидация carrier
	validCarriers := map[string]bool{
		"wbstream": true,
		"jazz":     true,
		"telemost": true,
	}
	if !validCarriers[e.profile.Carrier] {
		return fmt.Errorf("invalid carrier: %s", e.profile.Carrier)
	}

	// Валидация transport
	validTransports := map[string]bool{
		"datachannel": true,
		"vp8channel":  true,
		"seichannel":  true,
	}
	if !validTransports[e.profile.Transport] {
		return fmt.Errorf("invalid transport: %s", e.profile.Transport)
	}

	return nil
}

// sendStatus отправляет статус в канал (неблокирующая отправка)
func (e *Engine) sendStatus(status string) {
	if e.statusCh != nil {
		select {
		case e.statusCh <- status:
		default:
		}
	}
}

// sendLog отправляет лог в канал (неблокирующая отправка)
func (e *Engine) sendLog(msg string) {
	if e.logCh != nil {
		select {
		case e.logCh <- msg:
		default:
		}
	}
}
