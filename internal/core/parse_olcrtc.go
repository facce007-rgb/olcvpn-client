package core

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/openlibrecommunity/olcvpn/internal/types"
)

// ParseOlcRTCURI парсит olcrtc:// URI
// Формат: olcrtc://<Carrier>?<Transport>@<RoomID>#<Key>%<ClientID>$<MIMO>
// Пример: olcrtc://wbstream?datachannel@room123#abc123%client1$2x2
func ParseOlcRTCURI(uri string) (*types.Profile, error) {
	if !strings.HasPrefix(uri, "olcrtc://") {
		return nil, fmt.Errorf("not an olcrtc URI")
	}

	// Убираем схему
	rest := strings.TrimPrefix(uri, "olcrtc://")

	// Парсим компоненты
	// Carrier?Transport@RoomID#Key%ClientID$MIMO

	// Разделяем по @
	parts := strings.SplitN(rest, "@", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid olcrtc URI: missing @")
	}

	carrierTransport := parts[0]
	roomKeyClientMIMO := parts[1]

	// Carrier и Transport
	ctParts := strings.SplitN(carrierTransport, "?", 2)
	if len(ctParts) != 2 {
		return nil, fmt.Errorf("invalid olcrtc URI: missing transport")
	}
	carrier := ctParts[0]
	transport := ctParts[1]

	// RoomID, Key, ClientID, MIMO
	// RoomID#Key%ClientID$MIMO
	var roomID, key, clientID, mimo string

	// Разделяем по #
	roomKeyParts := strings.SplitN(roomKeyClientMIMO, "#", 2)
	if len(roomKeyParts) != 2 {
		return nil, fmt.Errorf("invalid olcrtc URI: missing key")
	}
	roomID = roomKeyParts[0]
	keyClientMIMO := roomKeyParts[1]

	// Разделяем по %
	keyClientParts := strings.SplitN(keyClientMIMO, "%", 2)
	if len(keyClientParts) != 2 {
		return nil, fmt.Errorf("invalid olcrtc URI: missing client ID")
	}
	key = keyClientParts[0]
	clientMIMO := keyClientParts[1]

	// Разделяем по $
	clientMIMOParts := strings.SplitN(clientMIMO, "$", 2)
	clientID = clientMIMOParts[0]
	if len(clientMIMOParts) == 2 {
		mimo = clientMIMOParts[1]
	}

	// Валидация
	if carrier == "" {
		return nil, fmt.Errorf("missing carrier")
	}
	if transport == "" {
		return nil, fmt.Errorf("missing transport")
	}
	if roomID == "" {
		return nil, fmt.Errorf("missing room ID")
	}
	if key == "" {
		return nil, fmt.Errorf("missing key")
	}
	if clientID == "" {
		return nil, fmt.Errorf("missing client ID")
	}

	// Валидация carrier
	validCarriers := map[string]bool{
		"wbstream": true,
		"jazz":     true,
		"telemost": true,
	}
	if !validCarriers[carrier] {
		return nil, fmt.Errorf("invalid carrier: %s (must be wbstream, jazz, or telemost)", carrier)
	}

	// Валидация transport
	validTransports := map[string]bool{
		"datachannel": true,
		"vp8channel":  true,
		"seichannel":  true,
	}
	if !validTransports[transport] {
		return nil, fmt.Errorf("invalid transport: %s", transport)
	}

	name := fmt.Sprintf("%s/%s (%s)", carrier, transport, roomID)

	profile := &types.Profile{
		Name:   name,
		Engine: types.EngineOlcRTC,
		OlcRTC: &types.OlcRTCProfile{
			Carrier:   carrier,
			Transport: transport,
			RoomID:    roomID,
			Key:       key,
			ClientID:  clientID,
			MIMO:      mimo,
		},
	}

	return profile, nil
}

// BuildOlcRTCURI создаёт olcrtc:// URI из профиля
func BuildOlcRTCURI(profile *types.OlcRTCProfile) string {
	uri := fmt.Sprintf("olcrtc://%s?%s@%s#%s%%%s",
		url.PathEscape(profile.Carrier),
		url.PathEscape(profile.Transport),
		url.PathEscape(profile.RoomID),
		url.PathEscape(profile.Key),
		url.PathEscape(profile.ClientID),
	)
	if profile.MIMO != "" {
		uri += "$" + url.PathEscape(profile.MIMO)
	}
	return uri
}
