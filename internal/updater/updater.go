package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"time"

	"github.com/openlibrecommunity/olcvpn/internal/types"
)

const (
	githubAPIURL = "https://api.github.com/repos/openlibrecommunity/olcvpn/releases/latest"
	currentVersion = "v0.1.0" // TODO: Генерировать из git tag при сборке
)

// Release представляет GitHub release
type Release struct {
	TagName     string  `json:"tag_name"`
	Name        string  `json:"name"`
	Body        string  `json:"body"`
	PublishedAt string  `json:"published_at"`
	Assets      []Asset `json:"assets"`
}

// Asset представляет файл в release
type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int64  `json:"size"`
}

// Updater проверяет обновления
type Updater struct {
	client *http.Client
}

// New создаёт новый Updater
func New() *Updater {
	return &Updater{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// CheckForUpdates проверяет наличие новой версии
func (u *Updater) CheckForUpdates() (*types.UpdateInfo, error) {
	req, err := http.NewRequest("GET", githubAPIURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := u.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch release info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var release Release
	if err := json.Unmarshal(body, &release); err != nil {
		return nil, fmt.Errorf("failed to parse release: %w", err)
	}

	// Проверяем есть ли новая версия
	if release.TagName == currentVersion {
		return &types.UpdateInfo{
			Available:      false,
			CurrentVersion: currentVersion,
			LatestVersion:  release.TagName,
		}, nil
	}

	// Находим подходящий asset для текущей платформы
	assetName := getAssetName()
	var downloadURL string
	var size int64

	for _, asset := range release.Assets {
		if asset.Name == assetName {
			downloadURL = asset.BrowserDownloadURL
			size = asset.Size
			break
		}
	}

	return &types.UpdateInfo{
		Available:      true,
		CurrentVersion: currentVersion,
		LatestVersion:  release.TagName,
		DownloadURL:    downloadURL,
		Size:           size,
		ReleaseNotes:   release.Body,
	}, nil
}

// getAssetName возвращает имя файла для текущей платформы
func getAssetName() string {
	switch runtime.GOOS {
	case "windows":
		return fmt.Sprintf("olcvpn-%s-%s.exe", runtime.GOOS, runtime.GOARCH)
	case "darwin":
		return fmt.Sprintf("olcvpn-%s-%s", runtime.GOOS, runtime.GOARCH)
	case "linux":
		return fmt.Sprintf("olcvpn-%s-%s", runtime.GOOS, runtime.GOARCH)
	default:
		return ""
	}
}

// DownloadUpdate скачивает обновление
func (u *Updater) DownloadUpdate(url string, dest string) error {
	resp, err := u.client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download update: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// TODO: Реализовать сохранение файла с прогрессом
	// TODO: Проверить подпись файла
	// TODO: Применить обновление

	return fmt.Errorf("not implemented yet")
}
