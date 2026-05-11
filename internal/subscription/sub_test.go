package subscription

import (
	"encoding/base64"
	"testing"
)

func TestDetectFormat(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		expected SubFormat
	}{
		{
			name: "Base64 V2Ray",
			data: base64.StdEncoding.EncodeToString([]byte("vless://uuid@example.com:443?security=tls#Test\nvmess://...")),
			expected: FormatBase64,
		},
		{
			name: "Clash YAML",
			data: `proxies:
  - name: test
    type: vless
    server: example.com
    port: 443`,
			expected: FormatClash,
		},
		{
			name: "sing-box JSON",
			data: `{"inbounds": [], "outbounds": [{"type": "vless", "server": "example.com"}]}`,
			expected: FormatSingBox,
		},
		{
			name: "olcrtc JSON",
			data: `{"servers": [{"name": "test", "carrier": "wbstream"}]}`,
			expected: FormatOlcRTC,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			format := DetectFormat([]byte(tt.data))
			if format != tt.expected {
				t.Errorf("DetectFormat() = %v, want %v", format, tt.expected)
			}
		})
	}
}

func TestParseBase64(t *testing.T) {
	uris := "vless://uuid-here@example.com:443?security=tls&sni=example.com#Test1\n" +
		"ss://YWVzLTI1Ni1nY206cGFzc3dvcmQ=@example.com:8388#Test2"

	encoded := base64.StdEncoding.EncodeToString([]byte(uris))

	profiles, err := parseBase64([]byte(encoded))
	if err != nil {
		t.Fatalf("parseBase64() error = %v", err)
	}

	if len(profiles) != 2 {
		t.Errorf("parseBase64() returned %d profiles, want 2", len(profiles))
	}

	if profiles[0].Name != "Test1" {
		t.Errorf("First profile name = %s, want Test1", profiles[0].Name)
	}
}

func TestParseClash(t *testing.T) {
	yaml := `
proxies:
  - name: Test VLESS
    type: vless
    server: example.com
    port: 443
    uuid: test-uuid
    tls: true
    servername: example.com
  - name: Test Trojan
    type: trojan
    server: example.com
    port: 443
    password: test-pass
    sni: example.com
`

	profiles, err := parseClash([]byte(yaml))
	if err != nil {
		t.Fatalf("parseClash() error = %v", err)
	}

	if len(profiles) != 2 {
		t.Errorf("parseClash() returned %d profiles, want 2", len(profiles))
	}

	if profiles[0].SingBox.Protocol != "vless" {
		t.Errorf("First profile protocol = %s, want vless", profiles[0].SingBox.Protocol)
	}

	if profiles[1].SingBox.Protocol != "trojan" {
		t.Errorf("Second profile protocol = %s, want trojan", profiles[1].SingBox.Protocol)
	}
}

func TestParseSingBox(t *testing.T) {
	json := `{
		"outbounds": [
			{
				"type": "vless",
				"tag": "proxy1",
				"server": "example.com",
				"server_port": 443,
				"uuid": "test-uuid"
			},
			{
				"type": "direct",
				"tag": "direct"
			}
		]
	}`

	profiles, err := parseSingBox([]byte(json))
	if err != nil {
		t.Fatalf("parseSingBox() error = %v", err)
	}

	// direct должен быть пропущен
	if len(profiles) != 1 {
		t.Errorf("parseSingBox() returned %d profiles, want 1", len(profiles))
	}

	if profiles[0].Name != "proxy1" {
		t.Errorf("Profile name = %s, want proxy1", profiles[0].Name)
	}
}

func TestParseOlcRTC(t *testing.T) {
	json := `{
		"servers": [
			{
				"name": "Test Server",
				"carrier": "wbstream",
				"transport": "datachannel",
				"room_id": "room123",
				"key": "key123",
				"client_id": "client1"
			}
		]
	}`

	profiles, err := parseOlcRTC([]byte(json))
	if err != nil {
		t.Fatalf("parseOlcRTC() error = %v", err)
	}

	if len(profiles) != 1 {
		t.Errorf("parseOlcRTC() returned %d profiles, want 1", len(profiles))
	}

	if profiles[0].OlcRTC.Carrier != "wbstream" {
		t.Errorf("Carrier = %s, want wbstream", profiles[0].OlcRTC.Carrier)
	}
}
