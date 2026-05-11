package core

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/openlibrecommunity/olcvpn/internal/types"
)

func TestStorageEncryptDecrypt(t *testing.T) {
	tmpDir := t.TempDir()
	storage, err := NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	tests := []struct {
		name      string
		plaintext string
	}{
		{"Simple text", "hello world"},
		{"UUID", "550e8400-e29b-41d4-a716-446655440000"},
		{"Password", "MySecurePassword123!@#"},
		{"Empty string", ""},
		{"Long text", "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encrypted, err := storage.encryptSecret(tt.plaintext)
			if err != nil {
				t.Fatalf("encryptSecret() error = %v", err)
			}

			if tt.plaintext != "" && encrypted == tt.plaintext {
				t.Error("Encrypted text should not equal plaintext")
			}

			decrypted, err := storage.decryptSecret(encrypted)
			if err != nil {
				t.Fatalf("decryptSecret() error = %v", err)
			}

			if decrypted != tt.plaintext {
				t.Errorf("decryptSecret() = %v, want %v", decrypted, tt.plaintext)
			}
		})
	}
}

func TestStorageSaveLoad(t *testing.T) {
	tmpDir := t.TempDir()
	storage, err := NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	profiles := []*types.Profile{
		{
			ID:        "profile1",
			Name:      "Test VLESS",
			Engine:    types.EngineSingBox,
			CreatedAt: time.Now(),
			SingBox: &types.SingBoxProfile{
				Protocol: "vless",
				Address:  "example.com",
				Port:     443,
				UUID:     "secret-uuid-123",
				Password: "secret-password-456",
			},
		},
		{
			ID:        "profile2",
			Name:      "Test OlcRTC",
			Engine:    types.EngineOlcRTC,
			CreatedAt: time.Now(),
			OlcRTC: &types.OlcRTCProfile{
				Carrier:   "wbstream",
				Transport: "datachannel",
				RoomID:    "room123",
				Key:       "secret-key-789",
				ClientID:  "client1",
			},
		},
	}

	// Сохраняем
	if err := storage.SaveProfiles(profiles); err != nil {
		t.Fatalf("SaveProfiles() error = %v", err)
	}

	// Проверяем что файл создан
	profilesPath := filepath.Join(tmpDir, profilesFile)
	if _, err := os.Stat(profilesPath); os.IsNotExist(err) {
		t.Error("Profiles file was not created")
	}

	// Читаем файл напрямую и проверяем что секреты зашифрованы
	data, err := os.ReadFile(profilesPath)
	if err != nil {
		t.Fatalf("Failed to read profiles file: %v", err)
	}
	fileContent := string(data)
	if contains(fileContent, "secret-uuid-123") {
		t.Error("UUID should be encrypted in file")
	}
	if contains(fileContent, "secret-password-456") {
		t.Error("Password should be encrypted in file")
	}
	if contains(fileContent, "secret-key-789") {
		t.Error("Key should be encrypted in file")
	}

	// Загружаем
	loaded, err := storage.LoadProfiles()
	if err != nil {
		t.Fatalf("LoadProfiles() error = %v", err)
	}

	if len(loaded) != len(profiles) {
		t.Errorf("LoadProfiles() returned %d profiles, want %d", len(loaded), len(profiles))
	}

	// Проверяем что секреты расшифрованы правильно
	if loaded[0].SingBox.UUID != "secret-uuid-123" {
		t.Errorf("UUID not decrypted correctly: got %v", loaded[0].SingBox.UUID)
	}
	if loaded[0].SingBox.Password != "secret-password-456" {
		t.Errorf("Password not decrypted correctly: got %v", loaded[0].SingBox.Password)
	}
	if loaded[1].OlcRTC.Key != "secret-key-789" {
		t.Errorf("Key not decrypted correctly: got %v", loaded[1].OlcRTC.Key)
	}
}

func TestStorageLoadEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	storage, err := NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	// Загружаем из пустой директории
	profiles, err := storage.LoadProfiles()
	if err != nil {
		t.Fatalf("LoadProfiles() error = %v", err)
	}

	if len(profiles) != 0 {
		t.Errorf("LoadProfiles() returned %d profiles, want 0", len(profiles))
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
