package qrcode

import (
	"encoding/json"
	"fmt"
	"image"

	"github.com/openlibrecommunity/olcvpn/internal/types"
	"github.com/skip2/go-qrcode"
)

// Generator генерирует QR-коды для профилей
type Generator struct{}

// New создаёт новый Generator
func New() *Generator {
	return &Generator{}
}

// GenerateProfileQR генерирует QR-код для профиля
func (g *Generator) GenerateProfileQR(profile *types.Profile) (image.Image, error) {
	if profile == nil {
		return nil, fmt.Errorf("profile is nil")
	}

	// Сериализуем профиль в JSON
	data, err := json.Marshal(profile)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal profile: %w", err)
	}

	// Генерируем QR-код
	qr, err := qrcode.New(string(data), qrcode.Medium)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code: %w", err)
	}

	// Возвращаем как image.Image
	return qr.Image(256), nil
}

// GenerateURIQR генерирует QR-код для URI (vless://, olcrtc://, etc)
func (g *Generator) GenerateURIQR(uri string) (image.Image, error) {
	if uri == "" {
		return nil, fmt.Errorf("uri is empty")
	}

	qr, err := qrcode.New(uri, qrcode.Medium)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code: %w", err)
	}

	return qr.Image(256), nil
}

// GenerateSubscriptionQR генерирует QR-код для subscription URL
func (g *Generator) GenerateSubscriptionQR(url string) (image.Image, error) {
	if url == "" {
		return nil, fmt.Errorf("url is empty")
	}

	qr, err := qrcode.New(url, qrcode.Medium)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code: %w", err)
	}

	return qr.Image(256), nil
}

// ParseProfileQR парсит QR-код и возвращает профиль
func (g *Generator) ParseProfileQR(data string) (*types.Profile, error) {
	var profile types.Profile
	if err := json.Unmarshal([]byte(data), &profile); err != nil {
		return nil, fmt.Errorf("failed to parse profile: %w", err)
	}
	return &profile, nil
}
