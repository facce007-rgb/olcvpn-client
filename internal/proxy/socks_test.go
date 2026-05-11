package proxy

import (
	"testing"
)

func TestGenerateSOCKSCredentials(t *testing.T) {
	creds, err := GenerateSOCKSCredentials()
	if err != nil {
		t.Fatalf("GenerateSOCKSCredentials() error = %v", err)
	}

	if creds.Username == "" {
		t.Error("Username should not be empty")
	}
	if creds.Password == "" {
		t.Error("Password should not be empty")
	}

	// Username должен быть 32 символа (16 байт в hex)
	if len(creds.Username) != 32 {
		t.Errorf("Username length = %d, want 32", len(creds.Username))
	}

	// Password должен быть 64 символа (32 байта в hex)
	if len(creds.Password) != 64 {
		t.Errorf("Password length = %d, want 64", len(creds.Password))
	}

	// Проверяем что генерируются разные credentials
	creds2, err := GenerateSOCKSCredentials()
	if err != nil {
		t.Fatalf("GenerateSOCKSCredentials() error = %v", err)
	}

	if creds.Username == creds2.Username {
		t.Error("Generated credentials should be unique")
	}
	if creds.Password == creds2.Password {
		t.Error("Generated credentials should be unique")
	}
}

func TestGenerateCryptoRandom(t *testing.T) {
	tests := []struct {
		name   string
		length int
		want   int // expected hex string length
	}{
		{"16 bytes", 16, 32},
		{"32 bytes", 32, 64},
		{"8 bytes", 8, 16},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := generateCryptoRandom(tt.length)
			if err != nil {
				t.Fatalf("generateCryptoRandom() error = %v", err)
			}
			if len(result) != tt.want {
				t.Errorf("generateCryptoRandom() length = %d, want %d", len(result), tt.want)
			}

			// Проверяем что это валидный hex
			for _, c := range result {
				if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
					t.Errorf("Invalid hex character: %c", c)
				}
			}
		})
	}
}

// TestSOCKSAuthenticationRequired проверяет что SOCKS5 требует аутентификацию
// Этот тест должен падать если SOCKS5 запускается без auth (критическая уязвимость §8)
func TestSOCKSAuthenticationRequired(t *testing.T) {
	// Этот тест проверяет концептуально что credentials всегда генерируются
	// Реальная проверка подключения к SOCKS5 без auth должна быть в интеграционных тестах

	creds, err := LoadOrCreateSOCKSCredentials()
	if err != nil {
		// Может упасть если keyring недоступен в тестовой среде
		t.Skipf("Keyring not available in test environment: %v", err)
	}

	if creds.Username == "" || creds.Password == "" {
		t.Error("CRITICAL: SOCKS5 credentials are empty - authentication would be disabled!")
	}

	// Очищаем после теста
	DeleteSOCKSCredentials()
}

func TestLoadOrCreateSOCKSCredentials(t *testing.T) {
	// Очищаем перед тестом
	DeleteSOCKSCredentials()

	// Первый вызов должен создать новые credentials
	creds1, err := LoadOrCreateSOCKSCredentials()
	if err != nil {
		t.Skipf("Keyring not available in test environment: %v", err)
	}

	if creds1.Username == "" || creds1.Password == "" {
		t.Error("Credentials should not be empty")
	}

	// Второй вызов должен вернуть те же credentials
	creds2, err := LoadOrCreateSOCKSCredentials()
	if err != nil {
		t.Fatalf("LoadOrCreateSOCKSCredentials() error = %v", err)
	}

	if creds1.Username != creds2.Username {
		t.Error("Username should be the same on second call")
	}
	if creds1.Password != creds2.Password {
		t.Error("Password should be the same on second call")
	}

	// Очищаем после теста
	DeleteSOCKSCredentials()
}
