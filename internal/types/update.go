package types

// UpdateInfo содержит информацию об обновлении
type UpdateInfo struct {
	Available      bool   `json:"available"`
	CurrentVersion string `json:"current_version"`
	LatestVersion  string `json:"latest_version"`
	DownloadURL    string `json:"download_url,omitempty"`
	Size           int64  `json:"size,omitempty"`
	ReleaseNotes   string `json:"release_notes,omitempty"`
}
