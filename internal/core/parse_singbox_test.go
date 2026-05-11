package core

import (
	"testing"

	"github.com/openlibrecommunity/olcvpn/internal/types"
)

func TestParseVLESS(t *testing.T) {
	tests := []struct {
		name    string
		uri     string
		wantErr bool
		check   func(*types.Profile) bool
	}{
		{
			name: "VLESS with Reality",
			uri:  "vless://uuid-here@example.com:443?security=reality&sni=example.com&pbk=pubkey123&sid=shortid&fp=chrome&flow=xtls-rprx-vision#TestServer",
			wantErr: false,
			check: func(p *types.Profile) bool {
				return p.SingBox != nil &&
					p.SingBox.Protocol == "vless" &&
					p.SingBox.Reality == true &&
					p.SingBox.RealityPubKey == "pubkey123" &&
					p.SingBox.RealityShortID == "shortid" &&
					p.Name == "TestServer"
			},
		},
		{
			name: "VLESS without Reality",
			uri:  "vless://uuid-here@example.com:443?security=tls&sni=example.com#TestServer",
			wantErr: false,
			check: func(p *types.Profile) bool {
				return p.SingBox != nil &&
					p.SingBox.Protocol == "vless" &&
					p.SingBox.TLS == true &&
					p.SingBox.Reality == false
			},
		},
		{
			name:    "VLESS missing UUID",
			uri:     "vless://@example.com:443",
			wantErr: true,
		},
		{
			name:    "VLESS Reality without pbk",
			uri:     "vless://uuid@example.com:443?security=reality&sni=example.com",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profile, err := ParseSingBoxURI(tt.uri)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseSingBoxURI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil && !tt.check(profile) {
				t.Errorf("ParseSingBoxURI() profile validation failed")
			}
		})
	}
}

func TestParseShadowsocks(t *testing.T) {
	tests := []struct {
		name    string
		uri     string
		wantErr bool
		check   func(*types.Profile) bool
	}{
		{
			name: "Shadowsocks basic",
			uri:  "ss://YWVzLTI1Ni1nY206cGFzc3dvcmQ=@example.com:8388#TestSS",
			wantErr: false,
			check: func(p *types.Profile) bool {
				return p.SingBox != nil &&
					p.SingBox.Protocol == "shadowsocks" &&
					p.SingBox.Address == "example.com" &&
					p.SingBox.Port == 8388
			},
		},
		{
			name:    "Shadowsocks missing port",
			uri:     "ss://YWVzLTI1Ni1nY206cGFzc3dvcmQ=@example.com",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profile, err := ParseSingBoxURI(tt.uri)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseSingBoxURI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil && !tt.check(profile) {
				t.Errorf("ParseSingBoxURI() profile validation failed")
			}
		})
	}
}

func TestParseTrojan(t *testing.T) {
	tests := []struct {
		name    string
		uri     string
		wantErr bool
		check   func(*types.Profile) bool
	}{
		{
			name: "Trojan basic",
			uri:  "trojan://password123@example.com:443?sni=example.com#TestTrojan",
			wantErr: false,
			check: func(p *types.Profile) bool {
				return p.SingBox != nil &&
					p.SingBox.Protocol == "trojan" &&
					p.SingBox.Password == "password123" &&
					p.SingBox.TLS == true &&
					p.SingBox.SNI == "example.com"
			},
		},
		{
			name:    "Trojan missing password",
			uri:     "trojan://@example.com:443",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profile, err := ParseSingBoxURI(tt.uri)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseSingBoxURI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil && !tt.check(profile) {
				t.Errorf("ParseSingBoxURI() profile validation failed")
			}
		})
	}
}

func TestParseVMess(t *testing.T) {
	// VMess использует base64 JSON, создаём тестовый
	// Для теста используем упрощённый вариант

	tests := []struct {
		name    string
		uri     string
		wantErr bool
	}{
		{
			name:    "VMess invalid base64",
			uri:     "vmess://invalid!!!",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseSingBoxURI(tt.uri)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseSingBoxURI() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseUnsupportedProtocol(t *testing.T) {
	_, err := ParseSingBoxURI("http://example.com")
	if err == nil {
		t.Error("Expected error for unsupported protocol")
	}
}
