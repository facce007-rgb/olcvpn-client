package subscription

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/openlibrecommunity/olcvpn/internal/core"
	"github.com/openlibrecommunity/olcvpn/internal/types"
	"gopkg.in/yaml.v3"
)

// SubFormat определяет формат подписки
type SubFormat string

const (
	FormatBase64   SubFormat = "base64"   // V2Ray стандарт
	FormatClash    SubFormat = "clash"    // Clash YAML
	FormatSingBox  SubFormat = "singbox"  // sing-box JSON
	FormatOlcRTC   SubFormat = "olcrtc"   // olcrtc sub
	FormatUnknown  SubFormat = "unknown"
)

// Subscription представляет подписку
type Subscription struct {
	URL         string
	Format      SubFormat
	LastUpdated time.Time
	Profiles    []*types.Profile
}

// FetchAndParse загружает и парсит подписку
func FetchAndParse(url string) ([]*types.Profile, error) {
	// Загружаем данные
	data, err := fetchURL(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch subscription: %w", err)
	}

	// Определяем формат
	format := DetectFormat(data)

	// Парсим
	profiles, err := ParseSubscription(data, format)
	if err != nil {
		return nil, fmt.Errorf("failed to parse subscription: %w", err)
	}

	return profiles, nil
}

// fetchURL загружает данные по URL
func fetchURL(url string) ([]byte, error) {
	client := &http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: nil, // Используем системные сертификаты
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// DetectFormat определяет формат подписки
func DetectFormat(data []byte) SubFormat {
	// Пробуем JSON (sing-box или olcrtc)
	if json.Valid(data) {
		var obj map[string]interface{}
		if err := json.Unmarshal(data, &obj); err == nil {
			// sing-box имеет поля "inbounds", "outbounds"
			if _, hasInbounds := obj["inbounds"]; hasInbounds {
				return FormatSingBox
			}
			// olcrtc sub имеет поле "servers"
			if _, hasServers := obj["servers"]; hasServers {
				return FormatOlcRTC
			}
		}
	}

	// Пробуем YAML (Clash)
	var yamlObj map[string]interface{}
	if err := yaml.Unmarshal(data, &yamlObj); err == nil {
		if _, hasProxies := yamlObj["proxies"]; hasProxies {
			return FormatClash
		}
	}

	// Пробуем Base64 (V2Ray стандарт)
	decoded, err := base64.StdEncoding.DecodeString(string(data))
	if err == nil && len(decoded) > 0 {
		// Проверяем что это список URI
		lines := strings.Split(string(decoded), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			// Если хотя бы одна строка похожа на URI
			if strings.HasPrefix(line, "vless://") ||
				strings.HasPrefix(line, "vmess://") ||
				strings.HasPrefix(line, "ss://") ||
				strings.HasPrefix(line, "trojan://") ||
				strings.HasPrefix(line, "olcrtc://") {
				return FormatBase64
			}
			break
		}
	}

	return FormatUnknown
}

// ParseSubscription парсит подписку в зависимости от формата
func ParseSubscription(data []byte, format SubFormat) ([]*types.Profile, error) {
	switch format {
	case FormatBase64:
		return parseBase64(data)
	case FormatClash:
		return parseClash(data)
	case FormatSingBox:
		return parseSingBox(data)
	case FormatOlcRTC:
		return parseOlcRTC(data)
	default:
		return nil, fmt.Errorf("unknown subscription format")
	}
}

// parseBase64 парсит Base64 формат (V2Ray стандарт)
func parseBase64(data []byte) ([]*types.Profile, error) {
	// Декодируем Base64
	decoded, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	// Разбиваем на строки
	lines := strings.Split(string(decoded), "\n")
	profiles := make([]*types.Profile, 0)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Парсим URI
		profile, err := core.ParseSingBoxURI(line)
		if err != nil {
			// Пробуем olcrtc
			profile, err = core.ParseOlcRTCURI(line)
			if err != nil {
				// Пропускаем неподдерживаемые URI
				continue
			}
		}

		profiles = append(profiles, profile)
	}

	return profiles, nil
}

// parseClash парсит Clash YAML формат
func parseClash(data []byte) ([]*types.Profile, error) {
	var clash struct {
		Proxies []map[string]interface{} `yaml:"proxies"`
	}

	if err := yaml.Unmarshal(data, &clash); err != nil {
		return nil, fmt.Errorf("failed to parse clash yaml: %w", err)
	}

	profiles := make([]*types.Profile, 0)

	for _, proxy := range clash.Proxies {
		profile := clashProxyToProfile(proxy)
		if profile != nil {
			profiles = append(profiles, profile)
		}
	}

	return profiles, nil
}

// clashProxyToProfile конвертирует Clash proxy в Profile
func clashProxyToProfile(proxy map[string]interface{}) *types.Profile {
	proxyType, _ := proxy["type"].(string)
	name, _ := proxy["name"].(string)
	server, _ := proxy["server"].(string)
	port, _ := proxy["port"].(int)

	if name == "" || server == "" || port == 0 {
		return nil
	}

	profile := &types.Profile{
		Name:   name,
		Engine: types.EngineSingBox,
		SingBox: &types.SingBoxProfile{
			Address: server,
			Port:    port,
		},
	}

	switch proxyType {
	case "vless":
		profile.SingBox.Protocol = "vless"
		profile.SingBox.UUID, _ = proxy["uuid"].(string)
		if tls, ok := proxy["tls"].(bool); ok && tls {
			profile.SingBox.TLS = true
			profile.SingBox.SNI, _ = proxy["servername"].(string)
		}

	case "vmess":
		profile.SingBox.Protocol = "vmess"
		profile.SingBox.UUID, _ = proxy["uuid"].(string)
		if tls, _ := proxy["tls"].(bool); tls {
			profile.SingBox.TLS = true
		}

	case "trojan":
		profile.SingBox.Protocol = "trojan"
		profile.SingBox.Password, _ = proxy["password"].(string)
		profile.SingBox.TLS = true
		profile.SingBox.SNI, _ = proxy["sni"].(string)

	case "ss", "shadowsocks":
		profile.SingBox.Protocol = "shadowsocks"
		profile.SingBox.Password, _ = proxy["password"].(string)

	default:
		return nil
	}

	return profile
}

// parseSingBox парсит sing-box JSON формат
func parseSingBox(data []byte) ([]*types.Profile, error) {
	var config struct {
		Outbounds []map[string]interface{} `json:"outbounds"`
	}

	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse singbox json: %w", err)
	}

	profiles := make([]*types.Profile, 0)

	for _, outbound := range config.Outbounds {
		outboundType, _ := outbound["type"].(string)
		tag, _ := outbound["tag"].(string)

		// Пропускаем служебные outbound
		if outboundType == "direct" || outboundType == "block" || outboundType == "dns" {
			continue
		}

		profile := singboxOutboundToProfile(outbound, tag)
		if profile != nil {
			profiles = append(profiles, profile)
		}
	}

	return profiles, nil
}

// singboxOutboundToProfile конвертирует sing-box outbound в Profile
func singboxOutboundToProfile(outbound map[string]interface{}, tag string) *types.Profile {
	outboundType, _ := outbound["type"].(string)
	server, _ := outbound["server"].(string)
	port, _ := outbound["server_port"].(float64)

	if server == "" || port == 0 {
		return nil
	}

	profile := &types.Profile{
		Name:   tag,
		Engine: types.EngineSingBox,
		SingBox: &types.SingBoxProfile{
			Protocol: outboundType,
			Address:  server,
			Port:     int(port),
		},
	}

	// Извлекаем специфичные для протокола поля
	switch outboundType {
	case "vless", "vmess":
		profile.SingBox.UUID, _ = outbound["uuid"].(string)
	case "trojan", "shadowsocks":
		profile.SingBox.Password, _ = outbound["password"].(string)
	}

	return profile
}

// parseOlcRTC парсит olcrtc sub формат
func parseOlcRTC(data []byte) ([]*types.Profile, error) {
	var sub struct {
		Servers []struct {
			Name      string `json:"name"`
			Carrier   string `json:"carrier"`
			Transport string `json:"transport"`
			RoomID    string `json:"room_id"`
			Key       string `json:"key"`
			ClientID  string `json:"client_id"`
			MIMO      string `json:"mimo,omitempty"`
		} `json:"servers"`
	}

	if err := json.Unmarshal(data, &sub); err != nil {
		return nil, fmt.Errorf("failed to parse olcrtc json: %w", err)
	}

	profiles := make([]*types.Profile, 0)

	for _, server := range sub.Servers {
		profile := &types.Profile{
			Name:   server.Name,
			Engine: types.EngineOlcRTC,
			OlcRTC: &types.OlcRTCProfile{
				Carrier:   server.Carrier,
				Transport: server.Transport,
				RoomID:    server.RoomID,
				Key:       server.Key,
				ClientID:  server.ClientID,
				MIMO:      server.MIMO,
			},
		}
		profiles = append(profiles, profile)
	}

	return profiles, nil
}
