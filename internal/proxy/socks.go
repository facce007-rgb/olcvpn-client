package proxy

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/zalando/go-keyring"
)

const (
	serviceName = "olcvpn"
	socksUserKey = "socks5_username"
	socksPassKey = "socks5_password"
)

// SOCKSCredentials содержит учётные данные для SOCKS5
type SOCKSCredentials struct {
	Username string
	Password string
}

// GenerateSOCKSCredentials генерирует случайные credentials
func GenerateSOCKSCredentials() (SOCKSCredentials, error) {
	username, err := generateCryptoRandom(16)
	if err != nil {
		return SOCKSCredentials{}, fmt.Errorf("failed to generate username: %w", err)
	}
	password, err := generateCryptoRandom(32)
	if err != nil {
		return SOCKSCredentials{}, fmt.Errorf("failed to generate password: %w", err)
	}
	return SOCKSCredentials{
		Username: username,
		Password: password,
	}, nil
}

// generateCryptoRandom генерирует криптографически стойкую случайную строку
func generateCryptoRandom(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// LoadOrCreateSOCKSCredentials загружает или создаёт SOCKS5 credentials из keyring
func LoadOrCreateSOCKSCredentials() (SOCKSCredentials, error) {
	username, err := keyring.Get(serviceName, socksUserKey)
	if err != nil {
		// Credentials не найдены — создаём новые
		creds, err := GenerateSOCKSCredentials()
		if err != nil {
			return SOCKSCredentials{}, fmt.Errorf("failed to generate credentials: %w", err)
		}
		if err := keyring.Set(serviceName, socksUserKey, creds.Username); err != nil {
			return SOCKSCredentials{}, fmt.Errorf("failed to save username: %w", err)
		}
		if err := keyring.Set(serviceName, socksPassKey, creds.Password); err != nil {
			return SOCKSCredentials{}, fmt.Errorf("failed to save password: %w", err)
		}
		return creds, nil
	}

	password, err := keyring.Get(serviceName, socksPassKey)
	if err != nil {
		return SOCKSCredentials{}, fmt.Errorf("failed to get password from keyring: %w", err)
	}

	return SOCKSCredentials{
		Username: username,
		Password: password,
	}, nil
}

// DeleteSOCKSCredentials удаляет credentials из keyring
func DeleteSOCKSCredentials() error {
	if err := keyring.Delete(serviceName, socksUserKey); err != nil {
		return fmt.Errorf("failed to delete username: %w", err)
	}
	if err := keyring.Delete(serviceName, socksPassKey); err != nil {
		return fmt.Errorf("failed to delete password: %w", err)
	}
	return nil
}
