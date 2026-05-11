package core

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/openlibrecommunity/olcvpn/internal/types"
)

// ParseSingBoxURI парсит URI протоколов sing-box (vless://, ss://, trojan://, vmess://)
func ParseSingBoxURI(uri string) (*types.Profile, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("invalid URI: %w", err)
	}

	switch u.Scheme {
	case "vless":
		return parseVLESS(u)
	case "ss":
		return parseShadowsocks(u)
	case "trojan":
		return parseTrojan(u)
	case "vmess":
		return parseVMess(uri)
	default:
		return nil, fmt.Errorf("unsupported protocol: %s", u.Scheme)
	}
}

// parseVLESS парсит vless:// URI
func parseVLESS(u *url.URL) (*types.Profile, error) {
	uuid := u.User.Username()
	if uuid == "" {
		return nil, fmt.Errorf("missing UUID")
	}

	host := u.Hostname()
	portStr := u.Port()
	if portStr == "" {
		portStr = "443"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid port: %w", err)
	}

	query := u.Query()
	name := u.Fragment
	if name == "" {
		name = fmt.Sprintf("%s:%d", host, port)
	}

	profile := &types.Profile{
		Name:   name,
		Engine: types.EngineSingBox,
		SingBox: &types.SingBoxProfile{
			Protocol: "vless",
			Address:  host,
			Port:     port,
			UUID:     uuid,
		},
	}

	// TLS/Reality
	security := query.Get("security")
	if security == "tls" || security == "reality" {
		profile.SingBox.TLS = true
		profile.SingBox.SNI = query.Get("sni")
		if security == "reality" {
			profile.SingBox.Reality = true
			profile.SingBox.RealityPubKey = query.Get("pbk")
			profile.SingBox.RealityShortID = query.Get("sid")
			if profile.SingBox.RealityPubKey == "" {
				return nil, fmt.Errorf("reality requires pbk parameter")
			}
		}
	}

	// Transport
	profile.SingBox.Transport = query.Get("type")
	profile.SingBox.Path = query.Get("path")

	// Insecure
	if query.Get("allowInsecure") == "1" {
		profile.SingBox.Insecure = true
	}

	return profile, nil
}

// parseShadowsocks парсит ss:// URI
func parseShadowsocks(u *url.URL) (*types.Profile, error) {
	// ss://base64(method:password)@host:port#name
	var password string

	if u.User != nil {
		password, _ = u.User.Password()
	} else {
		// Попытка декодировать из userinfo
		userinfo := strings.Split(u.Host, "@")[0]
		decoded, err := base64.RawURLEncoding.DecodeString(userinfo)
		if err != nil {
			decoded, err = base64.StdEncoding.DecodeString(userinfo)
			if err != nil {
				return nil, fmt.Errorf("failed to decode userinfo: %w", err)
			}
		}
		parts := strings.SplitN(string(decoded), ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid ss userinfo format")
		}
		password = parts[1]
	}

	host := u.Hostname()
	portStr := u.Port()
	if portStr == "" {
		return nil, fmt.Errorf("missing port")
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid port: %w", err)
	}

	name := u.Fragment
	if name == "" {
		name = fmt.Sprintf("%s:%d", host, port)
	}

	return &types.Profile{
		Name:   name,
		Engine: types.EngineSingBox,
		SingBox: &types.SingBoxProfile{
			Protocol: "shadowsocks",
			Address:  host,
			Port:     port,
			Password: password,
		},
	}, nil
}

// parseTrojan парсит trojan:// URI
func parseTrojan(u *url.URL) (*types.Profile, error) {
	password := u.User.Username()
	if password == "" {
		return nil, fmt.Errorf("missing password")
	}

	host := u.Hostname()
	portStr := u.Port()
	if portStr == "" {
		portStr = "443"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid port: %w", err)
	}

	query := u.Query()
	name := u.Fragment
	if name == "" {
		name = fmt.Sprintf("%s:%d", host, port)
	}

	profile := &types.Profile{
		Name:   name,
		Engine: types.EngineSingBox,
		SingBox: &types.SingBoxProfile{
			Protocol: "trojan",
			Address:  host,
			Port:     port,
			Password: password,
			TLS:      true,
		},
	}

	// SNI
	if sni := query.Get("sni"); sni != "" {
		profile.SingBox.SNI = sni
	}

	// Transport
	profile.SingBox.Transport = query.Get("type")
	profile.SingBox.Path = query.Get("path")

	return profile, nil
}

// parseVMess парсит vmess:// URI
func parseVMess(uri string) (*types.Profile, error) {
	// vmess://base64(json)
	if !strings.HasPrefix(uri, "vmess://") {
		return nil, fmt.Errorf("not a vmess URI")
	}

	encoded := strings.TrimPrefix(uri, "vmess://")
	decoded, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		decoded, err = base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			return nil, fmt.Errorf("failed to decode vmess: %w", err)
		}
	}

	var vmess struct {
		Add  string `json:"add"`
		Port any    `json:"port"`
		ID   string `json:"id"`
		Net  string `json:"net"`
		Type string `json:"type"`
		Host string `json:"host"`
		Path string `json:"path"`
		TLS  string `json:"tls"`
		SNI  string `json:"sni"`
		PS   string `json:"ps"`
	}

	if err := json.Unmarshal(decoded, &vmess); err != nil {
		return nil, fmt.Errorf("failed to unmarshal vmess: %w", err)
	}

	port := 0
	switch v := vmess.Port.(type) {
	case string:
		port, _ = strconv.Atoi(v)
	case float64:
		port = int(v)
	case int:
		port = v
	}

	name := vmess.PS
	if name == "" {
		name = fmt.Sprintf("%s:%d", vmess.Add, port)
	}

	profile := &types.Profile{
		Name:   name,
		Engine: types.EngineSingBox,
		SingBox: &types.SingBoxProfile{
			Protocol:  "vmess",
			Address:   vmess.Add,
			Port:      port,
			UUID:      vmess.ID,
			Transport: vmess.Net,
			Path:      vmess.Path,
			TLS:       vmess.TLS == "tls",
			SNI:       vmess.SNI,
		},
	}

	return profile, nil
}
