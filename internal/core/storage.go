package core

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/openlibrecommunity/olcvpn/internal/types"
	"github.com/openlibrecommunity/olcvpn/internal/proxy"
	"github.com/zalando/go-keyring"
)

const (
	serviceName    = "olcvpn"
	masterKeyName  = "master_key"
	profilesFile   = "profiles.json"
)

// Storage управляет хранением профилей с шифрованием секретов
type Storage struct {
	dataDir   string
	masterKey []byte
}

// NewStorage создаёт новое хранилище
func NewStorage(dataDir string) (*Storage, error) {
	if err := os.MkdirAll(dataDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	s := &Storage{dataDir: dataDir}
	if err := s.loadOrCreateMasterKey(); err != nil {
		return nil, fmt.Errorf("failed to initialize master key: %w", err)
	}

	return s, nil
}

// loadOrCreateMasterKey загружает или создаёт мастер-ключ из системного keyring
func (s *Storage) loadOrCreateMasterKey() error {
	keyStr, err := keyring.Get(serviceName, masterKeyName)
	if err != nil {
		// Ключ не найден — создаём новый
		key := make([]byte, 32) // AES-256
		if _, err := rand.Read(key); err != nil {
			return fmt.Errorf("failed to generate master key: %w", err)
		}
		keyStr = base64.StdEncoding.EncodeToString(key)
		if err := keyring.Set(serviceName, masterKeyName, keyStr); err != nil {
			return fmt.Errorf("failed to save master key: %w", err)
		}
		s.masterKey = key
		return nil
	}

	key, err := base64.StdEncoding.DecodeString(keyStr)
	if err != nil {
		return fmt.Errorf("failed to decode master key: %w", err)
	}
	s.masterKey = key
	return nil
}

// encryptSecret шифрует секрет с помощью AES-256-GCM
func (s *Storage) encryptSecret(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	block, err := aes.NewCipher(s.masterKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// decryptSecret расшифровывает секрет
func (s *Storage) decryptSecret(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}

	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(s.masterKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertextBytes := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// SaveProfiles сохраняет профили с шифрованием секретов
func (s *Storage) SaveProfiles(profiles []*types.Profile) error {
	// Шифруем секреты перед сохранением
	encrypted := make([]*types.Profile, len(profiles))
	for i, p := range profiles {
		ep := *p // копия
		if ep.SingBox != nil {
			sb := *ep.SingBox
			var err error
			if sb.UUID != "" {
				sb.UUID, err = s.encryptSecret(sb.UUID)
				if err != nil {
					return fmt.Errorf("failed to encrypt UUID: %w", err)
				}
			}
			if sb.Password != "" {
				sb.Password, err = s.encryptSecret(sb.Password)
				if err != nil {
					return fmt.Errorf("failed to encrypt password: %w", err)
				}
			}
			ep.SingBox = &sb
		}
		if ep.OlcRTC != nil {
			oc := *ep.OlcRTC
			var err error
			if oc.Key != "" {
				oc.Key, err = s.encryptSecret(oc.Key)
				if err != nil {
					return fmt.Errorf("failed to encrypt key: %w", err)
				}
			}
			ep.OlcRTC = &oc
		}
		encrypted[i] = &ep
	}

	data, err := json.MarshalIndent(encrypted, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal profiles: %w", err)
	}

	path := filepath.Join(s.dataDir, profilesFile)
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write profiles: %w", err)
	}

	return nil
}

// LoadProfiles загружает профили с расшифровкой секретов
func (s *Storage) LoadProfiles() ([]*types.Profile, error) {
	path := filepath.Join(s.dataDir, profilesFile)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []*types.Profile{}, nil
		}
		return nil, fmt.Errorf("failed to read profiles: %w", err)
	}

	var encrypted []*types.Profile
	if err := json.Unmarshal(data, &encrypted); err != nil {
		return nil, fmt.Errorf("failed to unmarshal profiles: %w", err)
	}

	// Расшифровываем секреты
	profiles := make([]*types.Profile, len(encrypted))
	for i, ep := range encrypted {
		p := *ep // копия
		if p.SingBox != nil {
			sb := *p.SingBox
			var err error
			if sb.UUID != "" {
				sb.UUID, err = s.decryptSecret(sb.UUID)
				if err != nil {
					return nil, fmt.Errorf("failed to decrypt UUID: %w", err)
				}
			}
			if sb.Password != "" {
				sb.Password, err = s.decryptSecret(sb.Password)
				if err != nil {
					return nil, fmt.Errorf("failed to decrypt password: %w", err)
				}
			}
			p.SingBox = &sb
		}
		if p.OlcRTC != nil {
			oc := *p.OlcRTC
			var err error
			if oc.Key != "" {
				oc.Key, err = s.decryptSecret(oc.Key)
				if err != nil {
					return nil, fmt.Errorf("failed to decrypt key: %w", err)
				}
			}
			p.OlcRTC = &oc
		}
		profiles[i] = &p
	}

	return profiles, nil
}

// GetSOCKSCredentials загружает или создаёт SOCKS5 credentials
func (s *Storage) GetSOCKSCredentials() (proxy.SOCKSCredentials, error) {
	username, err := keyring.Get(serviceName, "socks5_username")
	if err != nil {
		// Credentials не найдены — создаём новые
		creds, err := proxy.GenerateSOCKSCredentials()
		if err != nil {
			return proxy.SOCKSCredentials{}, fmt.Errorf("failed to generate credentials: %w", err)
		}
		if err := keyring.Set(serviceName, "socks5_username", creds.Username); err != nil {
			return proxy.SOCKSCredentials{}, fmt.Errorf("failed to save SOCKS username: %w", err)
		}
		if err := keyring.Set(serviceName, "socks5_password", creds.Password); err != nil {
			return proxy.SOCKSCredentials{}, fmt.Errorf("failed to save SOCKS password: %w", err)
		}
		return creds, nil
	}

	password, err := keyring.Get(serviceName, "socks5_password")
	if err != nil {
		return proxy.SOCKSCredentials{}, fmt.Errorf("failed to get SOCKS password: %w", err)
	}

	return proxy.SOCKSCredentials{
		Username: username,
		Password: password,
	}, nil
}
