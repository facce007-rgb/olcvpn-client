package singbox

import (
	"encoding/json"
	"fmt"

	"github.com/openlibrecommunity/olcvpn/internal/types"
	"github.com/openlibrecommunity/olcvpn/internal/proxy"
)

// GenerateConfig создаёт sing-box JSON конфиг из профиля
func GenerateConfig(profile *types.SingBoxProfile, socksPort int, httpPort int, creds proxy.SOCKSCredentials) ([]byte, error) {
	if profile == nil {
		return nil, fmt.Errorf("profile is nil")
	}

	config := map[string]any{
		"log": map[string]any{
			"level": "info",
		},
		"dns":       generateDNSConfig(),
		"inbounds":  generateInbounds(socksPort, httpPort, creds),
		"outbounds": generateOutbounds(profile),
		"route":     generateRoute(),
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config: %w", err)
	}

	return data, nil
}

// generateDNSConfig создаёт DNS конфигурацию согласно §7
func generateDNSConfig() map[string]any {
	return map[string]any{
		"servers": []map[string]any{
			{
				"tag":     "remote",
				"address": "tls://1.1.1.1",
				"detour":  "proxy",
			},
			{
				"tag":     "local",
				"address": "223.5.5.5",
				"detour":  "direct",
			},
		},
		"rules": []map[string]any{
			{
				"outbound": "any",
				"server":   "local",
			},
			{
				"clash_mode": "global",
				"server":     "remote",
			},
		},
		"strategy": "prefer_ipv4",
	}
}

// generateInbounds создаёт inbound конфигурацию с обязательной аутентификацией
func generateInbounds(socksPort int, httpPort int, creds proxy.SOCKSCredentials) []map[string]any {
	inbounds := []map[string]any{
		{
			"type":        "socks",
			"tag":         "socks-in",
			"listen":      "127.0.0.1",
			"listen_port": socksPort,
			"users": []map[string]string{
				{
					"username": creds.Username,
					"password": creds.Password,
				},
			},
			// UDP отключён при включённой auth (§8 Фикс 2)
			"udp_over_tcp": false,
		},
		{
			"type":        "http",
			"tag":         "http-in",
			"listen":      "127.0.0.1",
			"listen_port": httpPort,
			"users": []map[string]string{
				{
					"username": creds.Username,
					"password": creds.Password,
				},
			},
		},
	}

	return inbounds
}

// generateOutbounds создаёт outbound конфигурацию
func generateOutbounds(profile *types.SingBoxProfile) []map[string]any {
	outbounds := []map[string]any{
		generateProxyOutbound(profile),
		{
			"type": "direct",
			"tag":  "direct",
		},
		{
			"type": "block",
			"tag":  "block",
		},
	}

	return outbounds
}

// generateProxyOutbound создаёт proxy outbound в зависимости от протокола
func generateProxyOutbound(profile *types.SingBoxProfile) map[string]any {
	outbound := map[string]any{
		"type":        profile.Protocol,
		"tag":         "proxy",
		"server":      profile.Address,
		"server_port": profile.Port,
	}

	switch profile.Protocol {
	case "vless":
		outbound["uuid"] = profile.UUID
		if profile.Reality {
			outbound["flow"] = "xtls-rprx-vision"
		}
		if profile.TLS {
			outbound["tls"] = generateTLSConfig(profile)
		}
		if profile.Transport != "" {
			outbound["transport"] = generateTransportConfig(profile)
		}

	case "vmess":
		outbound["uuid"] = profile.UUID
		outbound["security"] = "auto"
		if profile.TLS {
			outbound["tls"] = generateTLSConfig(profile)
		}
		if profile.Transport != "" {
			outbound["transport"] = generateTransportConfig(profile)
		}

	case "trojan":
		outbound["password"] = profile.Password
		if profile.TLS {
			outbound["tls"] = generateTLSConfig(profile)
		}
		if profile.Transport != "" {
			outbound["transport"] = generateTransportConfig(profile)
		}

	case "shadowsocks":
		outbound["method"] = "aes-256-gcm" // дефолт, можно расширить
		outbound["password"] = profile.Password

	case "tuic":
		outbound["uuid"] = profile.UUID
		outbound["password"] = profile.Password
		outbound["congestion_control"] = "bbr"
		if profile.TLS {
			outbound["tls"] = generateTLSConfig(profile)
		}

	case "hysteria2":
		outbound["password"] = profile.Password
		if profile.TLS {
			outbound["tls"] = generateTLSConfig(profile)
		}
	}

	return outbound
}

// generateTLSConfig создаёт TLS конфигурацию
func generateTLSConfig(profile *types.SingBoxProfile) map[string]any {
	tls := map[string]any{
		"enabled":     true,
		"server_name": profile.SNI,
		"insecure":    profile.Insecure,
		"utls": map[string]any{
			"enabled":     true,
			"fingerprint": "chrome", // всегда chrome, не random (§7)
		},
	}

	if profile.Reality {
		tls["reality"] = map[string]any{
			"enabled":    true,
			"public_key": profile.RealityPubKey,
			"short_id":   profile.RealityShortID,
		}
	}

	return tls
}

// generateTransportConfig создаёт transport конфигурацию
func generateTransportConfig(profile *types.SingBoxProfile) map[string]any {
	transport := map[string]any{
		"type": profile.Transport,
	}

	switch profile.Transport {
	case "ws":
		if profile.Path != "" {
			transport["path"] = profile.Path
		}
	case "http":
		if profile.Path != "" {
			transport["path"] = profile.Path
		}
	case "grpc":
		if profile.Path != "" {
			transport["service_name"] = profile.Path
		}
	}

	return transport
}

// generateRoute создаёт routing конфигурацию
func generateRoute() map[string]any {
	return map[string]any{
		"rules": []map[string]any{
			{
				"protocol": []string{"dns"},
				"outbound": "dns-out",
			},
		},
		"final": "proxy",
	}
}
