package core

import (
	"testing"

	"github.com/openlibrecommunity/olcvpn/internal/types"
)

func TestParseOlcRTCURI(t *testing.T) {
	tests := []struct {
		name    string
		uri     string
		wantErr bool
		check   func(*testing.T, string, string, string, string, string, string)
	}{
		{
			name:    "Valid wbstream datachannel",
			uri:     "olcrtc://wbstream?datachannel@room123#abc123def456%client1$2x2",
			wantErr: false,
			check: func(t *testing.T, carrier, transport, roomID, key, clientID, mimo string) {
				if carrier != "wbstream" {
					t.Errorf("carrier = %s, want wbstream", carrier)
				}
				if transport != "datachannel" {
					t.Errorf("transport = %s, want datachannel", transport)
				}
				if roomID != "room123" {
					t.Errorf("roomID = %s, want room123", roomID)
				}
				if key != "abc123def456" {
					t.Errorf("key = %s, want abc123def456", key)
				}
				if clientID != "client1" {
					t.Errorf("clientID = %s, want client1", clientID)
				}
				if mimo != "2x2" {
					t.Errorf("mimo = %s, want 2x2", mimo)
				}
			},
		},
		{
			name:    "Valid jazz without MIMO",
			uri:     "olcrtc://jazz?datachannel@room456#key789%client2",
			wantErr: false,
			check: func(t *testing.T, carrier, transport, roomID, key, clientID, mimo string) {
				if carrier != "jazz" {
					t.Errorf("carrier = %s, want jazz", carrier)
				}
				if mimo != "" {
					t.Errorf("mimo = %s, want empty", mimo)
				}
			},
		},
		{
			name:    "Valid telemost vp8channel",
			uri:     "olcrtc://telemost?vp8channel@room789#keyabc%client3$4x4",
			wantErr: false,
			check: func(t *testing.T, carrier, transport, roomID, key, clientID, mimo string) {
				if carrier != "telemost" {
					t.Errorf("carrier = %s, want telemost", carrier)
				}
				if transport != "vp8channel" {
					t.Errorf("transport = %s, want vp8channel", transport)
				}
			},
		},
		{
			name:    "Missing @",
			uri:     "olcrtc://wbstream?datachannel",
			wantErr: true,
		},
		{
			name:    "Missing transport",
			uri:     "olcrtc://wbstream@room123#key%client",
			wantErr: true,
		},
		{
			name:    "Missing key",
			uri:     "olcrtc://wbstream?datachannel@room123%client",
			wantErr: true,
		},
		{
			name:    "Missing client ID",
			uri:     "olcrtc://wbstream?datachannel@room123#key",
			wantErr: true,
		},
		{
			name:    "Invalid carrier",
			uri:     "olcrtc://invalid?datachannel@room123#key%client",
			wantErr: true,
		},
		{
			name:    "Invalid transport",
			uri:     "olcrtc://wbstream?invalid@room123#key%client",
			wantErr: true,
		},
		{
			name:    "Not olcrtc URI",
			uri:     "http://example.com",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profile, err := ParseOlcRTCURI(tt.uri)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseOlcRTCURI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil && profile.OlcRTC != nil {
				tt.check(t,
					profile.OlcRTC.Carrier,
					profile.OlcRTC.Transport,
					profile.OlcRTC.RoomID,
					profile.OlcRTC.Key,
					profile.OlcRTC.ClientID,
					profile.OlcRTC.MIMO,
				)
			}
		})
	}
}

func TestBuildOlcRTCURI(t *testing.T) {
	tests := []struct {
		name     string
		carrier  string
		transport string
		roomID   string
		key      string
		clientID string
		mimo     string
		want     string
	}{
		{
			name:      "With MIMO",
			carrier:   "wbstream",
			transport: "datachannel",
			roomID:    "room123",
			key:       "key456",
			clientID:  "client1",
			mimo:      "2x2",
			want:      "olcrtc://wbstream?datachannel@room123#key456%client1$2x2",
		},
		{
			name:      "Without MIMO",
			carrier:   "jazz",
			transport: "datachannel",
			roomID:    "room789",
			key:       "keyabc",
			clientID:  "client2",
			mimo:      "",
			want:      "olcrtc://jazz?datachannel@room789#keyabc%client2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profile := &types.OlcRTCProfile{
				Carrier:   tt.carrier,
				Transport: tt.transport,
				RoomID:    tt.roomID,
				Key:       tt.key,
				ClientID:  tt.clientID,
				MIMO:      tt.mimo,
			}
			got := BuildOlcRTCURI(profile)
			if got != tt.want {
				t.Errorf("BuildOlcRTCURI() = %v, want %v", got, tt.want)
			}
		})
	}
}
