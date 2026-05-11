package singbox

import (
	"encoding/json"
	"testing"

	"github.com/openlibrecommunity/olcvpn/internal/proxy"
	"github.com/openlibrecommunity/olcvpn/internal/types"
)

func TestGenerateConfig(t *testing.T) {
	creds := proxy.SOCKSCredentials{
		Username: "testuser",
		Password: "testpass",
	}

	tests := []struct {
		name    string
		profile *types.SingBoxProfile
		wantErr bool
	}{
		{
			name: "VLESS with Reality",
			profile: &types.SingBoxProfile{
				Protocol:       "vless",
				Address:        "example.com",
				Port:           443,
				UUID:           "test-uuid",
				TLS:            true,
				SNI:            "example.com",
				Reality:        true,
				RealityPubKey:  "pubkey123",
				RealityShortID: "shortid",
			},
			wantErr: false,
		},
		{
			name: "Shadowsocks",
			profile: &types.SingBoxProfile{
				Protocol: "shadowsocks",
				Address:  "example.com",
				Port:     8388,
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name:    "Nil profile",
			profile: nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configData, err := GenerateConfig(tt.profile, 2080, 2081, creds)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Проверяем что это валидный JSON
				var config map[string]interface{}
				if err := json.Unmarshal(configData, &config); err != nil {
					t.Errorf("Failed to unmarshal config: %v", err)
					return
				}

				// Проверяем наличие основных секций
				if _, ok := config["dns"]; !ok {
					t.Error("dns section missing")
				}
				if _, ok := config["inbounds"]; !ok {
					t.Error("inbounds section missing")
				}
				if _, ok := config["outbounds"]; !ok {
					t.Error("outbounds section missing")
				}
				if _, ok := config["route"]; !ok {
					t.Error("route section missing")
				}
			}
		})
	}
}

func TestGenerateDNSConfig(t *testing.T) {
	dns := generateDNSConfig()

	if dns["strategy"] != "prefer_ipv4" {
		t.Error("DNS strategy should be prefer_ipv4")
	}

	if _, ok := dns["servers"]; !ok {
		t.Error("DNS servers missing")
	}

	if _, ok := dns["rules"]; !ok {
		t.Error("DNS rules missing")
	}
}
