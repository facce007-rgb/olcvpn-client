package singbox

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/openlibrecommunity/olcvpn/internal/proxy"
	"github.com/openlibrecommunity/olcvpn/internal/types"
	"github.com/sagernet/sing-box"
	"github.com/sagernet/sing-box/option"
)

// Engine управляет sing-box процессом
type Engine struct {
	profile  *types.SingBoxProfile
	box      *box.Box
	cancel   context.CancelFunc
	statusCh chan<- string
	logCh    chan<- string
}

// NewEngine создаёт новый sing-box движок
func NewEngine(profile *types.SingBoxProfile, statusCh chan<- string, logCh chan<- string) *Engine {
	return &Engine{
		profile:  profile,
		statusCh: statusCh,
		logCh:    logCh,
	}
}

// Start запускает sing-box
func (e *Engine) Start(ctx context.Context, socksPort int, httpPort int) error {
	// Загружаем SOCKS credentials
	creds, err := proxy.LoadOrCreateSOCKSCredentials()
	if err != nil {
		return fmt.Errorf("failed to load SOCKS credentials: %w", err)
	}

	// Генерируем конфиг
	configData, err := GenerateConfig(e.profile, socksPort, httpPort, creds)
	if err != nil {
		return fmt.Errorf("failed to generate config: %w", err)
	}

	// Парсим конфиг в sing-box options
	options, err := parseConfig(configData)
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	// Создаём sing-box instance
	instance, err := box.New(box.Options{
		Context: ctx,
		Options: options,
	})
	if err != nil {
		return fmt.Errorf("failed to create sing-box instance: %w", err)
	}

	e.box = instance

	// Создаём контекст с отменой
	ctx, cancel := context.WithCancel(ctx)
	e.cancel = cancel

	// Запускаем sing-box в горутине
	go e.run(ctx)

	e.sendLog("sing-box engine started")
	e.sendStatus("connected")

	return nil
}

// run выполняет sing-box процесс
func (e *Engine) run(ctx context.Context) {
	e.sendLog("sing-box: starting box")

	// Запускаем sing-box
	if err := e.box.Start(); err != nil {
		e.sendLog(fmt.Sprintf("sing-box: start error: %v", err))
		e.sendStatus("error")
		return
	}

	e.sendLog("sing-box: running")

	// Ждём отмены контекста
	<-ctx.Done()

	e.sendLog("sing-box: stopping")

	// Останавливаем sing-box
	if err := e.box.Close(); err != nil {
		e.sendLog(fmt.Sprintf("sing-box: close error: %v", err))
	}

	e.sendLog("sing-box: stopped")
}

// Stop останавливает sing-box
func (e *Engine) Stop() error {
	if e.cancel != nil {
		e.cancel()
		e.cancel = nil
	}

	e.sendLog("sing-box engine stopped")
	e.sendStatus("disconnected")

	return nil
}

// parseConfig парсит JSON конфиг в sing-box options
func parseConfig(configData []byte) (option.Options, error) {
	var options option.Options
	if err := json.Unmarshal(configData, &options); err != nil {
		return option.Options{}, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return options, nil
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
