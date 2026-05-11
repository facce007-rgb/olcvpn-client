package types

import "time"

// EngineType определяет тип VPN движка
type EngineType string

const (
	EngineSingBox EngineType = "singbox"
	EngineOlcRTC  EngineType = "olcrtc"
)

// Profile представляет VPN профиль
type Profile struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Engine    EngineType      `json:"engine"`
	CreatedAt time.Time       `json:"created_at"`
	SingBox   *SingBoxProfile `json:"singbox,omitempty"`
	OlcRTC    *OlcRTCProfile  `json:"olcrtc,omitempty"`
}

// SingBoxProfile содержит настройки для sing-box движка
type SingBoxProfile struct {
	Protocol       string `json:"protocol"`        // vless, shadowsocks, trojan, vmess, tuic, hysteria2
	Address        string `json:"address"`
	Port           int    `json:"port"`
	UUID           string `json:"uuid,omitempty"`      // зашифрован в storage
	Password       string `json:"password,omitempty"`  // зашифрован в storage
	TLS            bool   `json:"tls"`
	SNI            string `json:"sni,omitempty"`
	Insecure       bool   `json:"insecure"`
	Reality        bool   `json:"reality"`
	RealityPubKey  string `json:"reality_pub_key,omitempty"`
	RealityShortID string `json:"reality_short_id,omitempty"`
	Transport      string `json:"transport,omitempty"` // ws, grpc, http, quic
	Path           string `json:"path,omitempty"`      // для ws/http
	RawConfig      string `json:"raw_config,omitempty"` // полный sing-box JSON если нужен
}

// OlcRTCProfile содержит настройки для olcrtc движка
type OlcRTCProfile struct {
	Carrier   string `json:"carrier"`    // wbstream | jazz | telemost
	Transport string `json:"transport"`  // datachannel | vp8channel | seichannel
	RoomID    string `json:"room_id"`
	Key       string `json:"key"`        // зашифрован, hex 32 байта
	ClientID  string `json:"client_id"`
	MIMO      string `json:"mimo,omitempty"` // опционально
}

// Metrics содержит метрики подключения
type Metrics struct {
	BytesUp   int64 `json:"bytes_up"`
	BytesDown int64 `json:"bytes_down"`
	LatencyMS int64 `json:"latency_ms"`
}
