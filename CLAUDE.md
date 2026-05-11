# CLAUDE.md — OLC VPN Client

> Этот файл читается Claude Code в начале каждой сессии.
> Единственный источник истины об архитектуре, соглашениях и решениях проекта.
> **Прочитай полностью перед любыми изменениями в коде.**

---

## 1. Что такое этот проект

**OLC VPN Client** — кросс-платформенный VPN-клиент на Go, аналог V2RayTun/NekoBox/Happ.
Объединяет два движка в одном GUI:

- **sing-box** (`github.com/sagernet/sing-box v1.11.0`) — VLESS+Reality, TLS, TUIC, Hysteria2, Shadowsocks, VMess, Trojan
- **olcrtc** (`github.com/openlibrecommunity/olcrtc`) — WebRTC-туннель через российские белые сервисы (WB Stream, Jazz, Яндекс Telemost)

Пользователь получает один клиент для всего: обычные протоколы + WebRTC-туннель как резерв когда DPI блокирует всё остальное.

### Целевые платформы

| Платформа | GUI | Приоритет |
|-----------|-----|-----------|
| Android | gomobile AAR + Kotlin/Compose | Высший |
| Windows | Fyne v2 | Высший |
| Linux | Fyne v2 | Средний |
| macOS | Fyne v2 | Средний |
| iOS | gomobile xcframework + SwiftUI | Средний |

---

## 2. Быстрый старт — команды

```bash
# Запуск для разработки (desktop)
go run ./cmd/olcvpn/

# Тесты
go test ./...

# Сборка всех desktop платформ
mage cross

# Сборка Android AAR
mage android

# Сборка iOS xcframework
mage ios

# Линтер
golangci-lint run ./...
```

Зависимости: `go install github.com/magefile/mage@latest`
iOS/Android: `go install golang.org/x/mobile/cmd/gomobile@latest && gomobile init`

---

## 3. Структура репозитория

```
olcvpn/
├── cmd/olcvpn/main.go            # точка входа desktop
├── internal/
│   ├── core/
│   │   ├── manager.go            # главный оркестратор
│   │   ├── profile.go            # типы Profile, SingBoxProfile, OlcRTCProfile
│   │   ├── storage.go            # JSON-хранилище + шифрование секретов
│   │   └── metrics.go            # трафик up/down, latency
│   ├── engine/
│   │   ├── singbox/
│   │   │   ├── engine.go         # Start/Stop sing-box
│   │   │   ├── config.go         # генерация sing-box JSON из Profile
│   │   │   └── parse.go          # парсер vless://, ss://, trojan://, vmess://
│   │   └── olcrtc/
│   │       ├── engine.go         # Start/Stop olcrtc
│   │       └── parse.go          # парсер olcrtc:// URI
│   ├── proxy/
│   │   ├── tun.go                # TUN интерфейс
│   │   └── socks.go              # SOCKS5 с аутентификацией (обязательно, см. §8)
│   ├── subscription/
│   │   └── sub.go                # загрузка sub-ссылок (4 формата)
│   └── ui/
│       ├── app.go                # Fyne приложение
│       ├── screens/
│       │   ├── home.go
│       │   ├── profiles.go
│       │   ├── add_profile.go
│       │   ├── logs.go
│       │   └── settings.go
│       └── theme/theme.go        # тёмная тема
├── mobile/
│   ├── mobile.go                 # gomobile API — общий для Android и iOS
│   └── bridge.go                 # вспомогательные типы и интерфейсы
├── android/app/                  # Android Studio проект
│   └── src/main/java/.../
│       ├── OlcVpnService.kt
│       └── CoreBridge.kt
├── ios/
│   ├── OlcVPN/                   # SwiftUI приложение (main target)
│   │   ├── ContentView.swift
│   │   ├── ProfilesView.swift
│   │   ├── VPNManager.swift      # NETunnelProviderManager
│   │   └── KeychainBridge.swift  # iOS Keychain реализация
│   └── OlcVPNExtension/          # Network Extension (отдельный target!)
│       └── PacketTunnelProvider.swift
├── magefile.go
├── go.mod
└── CLAUDE.md
```

Не создавай файлы за пределами этой структуры без явного обсуждения.

---

## 4. Ключевые зависимости

| Пакет | Версия | Назначение |
|-------|--------|------------|
| `fyne.io/fyne/v2` | v2.5.0 | Desktop GUI |
| `github.com/sagernet/sing-box` | v1.11.0 | VPN движок |
| `github.com/openlibrecommunity/olcrtc` | latest | WebRTC туннель |
| `github.com/zalando/go-keyring` | v0.2.3 | Keyring (desktop/Android) |
| `github.com/google/uuid` | v1.6.0 | UUID профилей |
| `golang.org/x/mobile` | latest | gomobile |
| `gopkg.in/yaml.v3` | v3.0.1 | Clash YAML |

Не добавляй новые зависимости без явного обсуждения.

---

## 5. Архитектурные принципы

```
GUI / Swift / Compose  →  Manager  →  Engine (singbox | olcrtc)  →  Network
```

- GUI никогда не обращается к engine напрямую — только через `Manager`.
- Engine не знает о GUI — только каналы `statusCh` и `logCh`.
- `Manager` — единственный владелец жизненного цикла движков.
- Публичные методы `Manager` всегда возвращают `error`. Игнорировать ошибки запрещено.
- Логируй через `logCh`, не через `fmt.Println`.
- `Manager` защищён `sync.RWMutex`: читающие методы — `RLock`, пишущие — `Lock`.
- Движки запускаются в горутинах. Stop всегда через `context.CancelFunc`.
- Blocking-операции никогда не вызываются в UI-горутине.

---

## 6. Модель данных профиля

```go
// internal/core/profile.go

type EngineType string

const (
    EngineSingBox EngineType = "singbox"
    EngineOlcRTC  EngineType = "olcrtc"
)

type Profile struct {
    ID        string          `json:"id"`
    Name      string          `json:"name"`
    Engine    EngineType      `json:"engine"`
    CreatedAt time.Time       `json:"created_at"`
    SingBox   *SingBoxProfile `json:"singbox,omitempty"`
    OlcRTC    *OlcRTCProfile  `json:"olcrtc,omitempty"`
}

type SingBoxProfile struct {
    Protocol       string `json:"protocol"`
    Address        string `json:"address"`
    Port           int    `json:"port"`
    UUID           string `json:"uuid,omitempty"`      // хранится зашифрованным
    Password       string `json:"password,omitempty"`  // хранится зашифрованным
    TLS            bool   `json:"tls"`
    SNI            string `json:"sni,omitempty"`
    Insecure       bool   `json:"insecure"`
    Reality        bool   `json:"reality"`
    RealityPubKey  string `json:"reality_pub_key,omitempty"`
    RealityShortID string `json:"reality_short_id,omitempty"`
    Transport      string `json:"transport,omitempty"`
    Path           string `json:"path,omitempty"`
    RawConfig      string `json:"raw_config,omitempty"`
}

type OlcRTCProfile struct {
    Carrier   string `json:"carrier"`    // wbstream | jazz | telemost
    Transport string `json:"transport"`  // datachannel | vp8channel | seichannel
    RoomID    string `json:"room_id"`
    Key       string `json:"key"`        // зашифрован, hex 32 байта
    ClientID  string `json:"client_id"`
    MIMO      string `json:"mimo,omitempty"`
}
```

---

## 7. Генерация sing-box конфига

VLESS + Reality — обязательные поля:
```json
{
  "type": "vless",
  "tag": "proxy",
  "server": "<address>",
  "server_port": 443,
  "uuid": "<из keyring>",
  "flow": "xtls-rprx-vision",
  "tls": {
    "enabled": true,
    "server_name": "<sni>",
    "utls": { "enabled": true, "fingerprint": "chrome" },
    "reality": {
      "enabled": true,
      "public_key": "<pubkey>",
      "short_id": "<shortid>"
    }
  }
}
```

Всегда `"fingerprint": "chrome"`, не `"random"`.

DNS — всегда включать в конфиг:
```json
{
  "dns": {
    "servers": [
      { "tag": "remote", "address": "tls://1.1.1.1", "detour": "proxy" },
      { "tag": "local",  "address": "223.5.5.5",     "detour": "direct" }
    ],
    "rules": [
      { "outbound": "any", "server": "local" },
      { "clash_mode": "global", "server": "remote" }
    ],
    "strategy": "prefer_ipv4"
  }
}
```

Outbound `"direct"` и `"block"` обязательны — sing-box требует их для route.

---

## 8. КРИТИЧЕСКАЯ УЯЗВИМОСТЬ — SOCKS5 без аутентификации

> Опубликовано runetfreedom, Хабр, 7 апреля 2026.
> PoC: https://github.com/runetfreedom/per-app-split-bypass-poc
> v2RayTun пофиксил (v5.22.71+). Happ пофиксил.
> OLC VPN Client должен быть безопасным с первого дня.

Android не изолирует loopback между приложениями. Шпионское ПО может подключиться
к открытому SOCKS5 напрямую, минуя VpnService, и узнать реальный IP VPN-сервера.

### Фикс 1: SOCKS5 всегда с аутентификацией

```go
// internal/proxy/socks.go

type SOCKSCredentials struct {
    Username string
    Password string
}

func GenerateSOCKSCredentials() SOCKSCredentials {
    return SOCKSCredentials{
        Username: generateCryptoRandom(16),
        Password: generateCryptoRandom(32),
    }
}

// В sing-box конфиге ВСЕГДА:
// { "type": "socks", "listen": "127.0.0.1", "listen_port": 2080,
//   "users": [{ "username": "...", "password": "..." }] }
```

### Фикс 2: UDP отключить при включённой auth

SOCKS5 не аутентифицирует UDP-ассоциации:
```go
if settings.SOCKSAuthEnabled {
    inbound["udp_over_tcp"] = false // принудительно TCP-only
}
```

### Фикс 3: sing-box API только с токеном

```go
// Если ClashAPI нужен — обязательно с secret.
// Если secret пустой — не добавляй секцию experimental вообще.
"clash_api": {
    "external_controller": "127.0.0.1:9090",
    "secret": "<случайный токен из keyring>"
}
```

### Фикс 4: HTTP прокси тоже с auth

```json
{ "type": "http", "listen": "127.0.0.1", "listen_port": 2081,
  "users": [{ "username": "...", "password": "..." }] }
```

### Хранение credentials

```go
func LoadOrCreateSOCKSCredentials() (SOCKSCredentials, error) {
    username, err := keyring.Get("olcvpn", "socks5_username")
    if err != nil {
        creds := GenerateSOCKSCredentials()
        keyring.Set("olcvpn", "socks5_username", creds.Username)
        keyring.Set("olcvpn", "socks5_password", creds.Password)
        return creds, nil
    }
    password, _ := keyring.Get("olcvpn", "socks5_password")
    return SOCKSCredentials{username, password}, nil
}
```

Чеклист:
- [ ] SOCKS5 запускается с username/password (случайные, из keyring)
- [ ] UDP отключён при включённой auth
- [ ] HTTP прокси тоже с auth
- [ ] sing-box API с secret или не включается
- [ ] Credentials не хранятся в plaintext JSON

---

## 9. Хранение секретов профиля

UUID, Password, Key никогда не хранятся в открытом виде в JSON файлах:

```go
// AES-256-GCM, мастер-ключ из системного keyring
// Windows: DPAPI
// Linux/macOS: github.com/zalando/go-keyring
// Android: Android Keystore через JNI
// iOS: iOS Keychain через gomobile bridge (см. §11.6)

func encryptSecret(plaintext string, masterKey []byte) (string, error)
func decryptSecret(ciphertext string, masterKey []byte) (string, error)
```

---

## 10. olcrtc — детальный анализ совместимости с iOS

> Прочитай этот раздел перед реализацией olcrtc на iOS.

olcrtc — чистый Go, использует pion/webrtc (pure Go WebRTC). Уже имеет `mage mobile`
для Android gomobile. Для iOS собирается через `gomobile bind -target ios`.

### Факт: routing loop НЕ проблема на iOS

На iOS сокеты процесса Network Extension автоматически исключены из туннеля.
olcrtc подключается к Telemost/Jazz/WB Stream через реальную сеть, а не через TUN.
Это встроено в iOS. Ничего дополнительного делать не нужно.
(На Android нужно явно вызывать `VpnService.protect(socket)` — на iOS это лишнее.)

### Анализ транспортов по платформам

| Транспорт | Носитель | Android | Desktop | iOS | Причина для iOS |
|-----------|----------|---------|---------|-----|-----------------|
| `datachannel` | wbstream, jazz | OK | OK | OK | Чистые WebRTC data channels, нет видео |
| `vp8channel` | telemost | OK | OK | Рискованно | Нет аппаратного VP8 на iOS, software encoder CPU-тяжёлый |
| `seichannel` | telemost | OK | OK | Не рекомендуется | SEI injection в H.264, H.264 из Go на iOS сложно |
| `videochannel` | telemost | OK | OK | Проблематично | Видео qrcode real-time нестабильно в Extension |

### Почему vp8channel опасен на iOS

Network Extension работает в отдельном процессе с жёстким лимитом памяти.
На старых устройствах это 50 МБ, на новых до 200 МБ.
sing-box уже занимает значительную часть этого бюджета.
VP8 software encoder (libvpx) + WebRTC stack поверх — высокий риск OOM.
При OOM kill iOS убивает Extension молча: VPN просто падает без объяснений.

### Решение для iOS: только datachannel

Поддерживай на iOS только `datachannel` транспорт. Это правильный компромисс:
- `wbstream + datachannel` и `jazz + datachannel` работают полноценно.
- Telemost + видео-транспорты — показывай пользователю понятное объяснение.

```go
// mobile/mobile.go

func (v *VPNCore) GetSupportedTransports(carrier string) string {
    if isIOS() {
        // На iOS только datachannel — безопасно по памяти
        return `["datachannel"]`
    }
    switch carrier {
    case "telemost":
        return `["vp8channel", "seichannel", "datachannel"]`
    default:
        return `["datachannel"]`
    }
}
```

### Итог по olcrtc на iOS

olcrtc работает на iOS при использовании datachannel транспорта. wbstream и jazz
через datachannel — полноценный туннель. Видео-транспорты не поддерживаем на iOS
из-за рисков OOM в Network Extension.

---

## 11. iOS — реализация

### 11.1 Два таргета в Xcode — обязательно

```
OlcVPN (main app)  ←→  OlcVPNExtension (Network Extension)
```

Go-ядро работает в Extension. Основное приложение управляет Extension
через `NETunnelProviderManager`. Это два отдельных процесса.

### 11.2 gomobile для iOS

```go
// magefile.go
func IOS() error {
    return sh.Run("gomobile", "bind",
        "-target", "ios",
        "-o", "ios/Frameworks/VPNCore.xcframework",
        "./mobile/")
}
```

API в `mobile/mobile.go` тот же что для Android. Меняется только нативная обёртка.

### 11.3 PacketTunnelProvider.swift

```swift
class PacketTunnelProvider: NEPacketTunnelProvider {
    var core: VPNCore?

    override func startTunnel(
        options: [String: NSObject]?,
        completionHandler: @escaping (Error?) -> Void
    ) {
        let config = protocolConfiguration as! NETunnelProviderProtocol
        let profileJSON = config.providerConfiguration?["profile"] as? String ?? ""

        let settings = NEPacketTunnelNetworkSettings(tunnelRemoteAddress: "127.0.0.1")
        settings.ipv4Settings = NEIPv4Settings(
            addresses: ["10.0.0.1"],
            subnetMasks: ["255.255.255.0"]
        )
        settings.ipv4Settings?.includedRoutes = [NEIPv4Route.default()]
        settings.dnsSettings = NEDNSSettings(servers: ["1.1.1.1"])
        settings.mtu = 1500

        setTunnelNetworkSettings(settings) { error in
            guard error == nil else { completionHandler(error); return }
            self.core = VPNCore()
            // Используй ConnectIOS — без TUN fd, через platform.Interface
            let err = self.core?.connectIOS(profileJSON)
            completionHandler(err)
        }
    }

    override func stopTunnel(
        with reason: NEProviderStopReason,
        completionHandler: @escaping () -> Void
    ) {
        core?.disconnect()
        completionHandler()
    }
}
```

### 11.4 ConnectIOS в mobile.go

На iOS sing-box использует `platform.Interface` вместо TUN fd.
Именно так устроен Hiddify на iOS — изучи их реализацию как референс.

```go
// mobile/mobile.go

// ConnectIOS — для iOS вместо Connect(profileJSON).
// sing-box box.NewWithOptions + WithPlatformInterface.
func (v *VPNCore) ConnectIOS(profileJSON string) error

// StartWithTunFd — для Android.
func (v *VPNCore) StartWithTunFd(fd int) error

// PacketTunnelFlow — интерфейс gomobile bridge для чтения/записи пакетов.
type PacketTunnelFlow interface {
    ReadPacket() []byte
    WritePacket(data []byte) bool
}

func (v *VPNCore) SetPacketFlow(flow PacketTunnelFlow)
```

### 11.5 Обязательный entitlement

В `OlcVPNExtension/OlcVPNExtension.entitlements`:
```xml
<key>com.apple.developer.networking.networkextension</key>
<array>
    <string>packet-tunnel-provider</string>
</array>
```

Этот entitlement нужно запросить у Apple через developer portal. Без него Extension не запустится.

### 11.6 Keychain вместо keyring на iOS

`github.com/zalando/go-keyring` не работает на iOS. Используй iOS Keychain через bridge:

```go
// mobile/bridge.go

// KeychainStorage реализуется в Swift, передаётся в Go через gomobile.
type KeychainStorage interface {
    Set(key, value string) error
    Get(key string) (string, error)
    Delete(key string) error
}

func (v *VPNCore) SetKeychainStorage(storage KeychainStorage)
```

```swift
// ios/OlcVPN/KeychainBridge.swift

class KeychainBridge: NSObject, MobileKeychainStorage {
    func set(_ key: String?, value: String?) throws {
        // SecItemAdd / SecItemUpdate
    }
    func get(_ key: String?) throws -> String {
        // SecItemCopyMatching
        return ""
    }
    func delete(_ key: String?) throws {
        // SecItemDelete
    }
}
```

### 11.7 Распространение на iOS в России

App Store в России удаляет VPN-приложения по требованию РКН. Варианты:
- TestFlight — до 10 000 пользователей, без проверки РКН, бесплатно
- US/EU аккаунт разработчика — разместить под нероссийским аккаунтом
- AltStore/Sideloading — пользователь устанавливает сам

---

## 12. Android — реализация

```bash
# magefile.go
func Android() error {
    return sh.Run("gomobile", "bind",
        "-target", "android",
        "-androidapi", "21",
        "-o", "android/app/libs/vpncore.aar",
        "./mobile/")
}
```

```kotlin
// OlcVpnService.kt
class OlcVpnService : VpnService() {
    override fun onStartCommand(intent: Intent?, flags: Int, startId: Int): Int {
        val builder = Builder()
            .addAddress("10.0.0.1", 30)
            .addRoute("0.0.0.0", 0)
            .addDnsServer("1.1.1.1")
            .setSession("OLC VPN")
            .setMtu(1500)
        val vpnInterface = builder.establish()
        val fd = vpnInterface?.detachFd() ?: return START_NOT_STICKY
        coreBridge.startWithTunFd(fd) // передаём fd в Go
        return START_STICKY
    }
}
```

---

## 13. gomobile API — полный список

```go
// mobile/mobile.go — компилируется для Android и iOS

type VPNCore struct { manager *core.Manager }

func NewVPNCore() *VPNCore
func (v *VPNCore) Connect(profileJSON string) error         // desktop/Android
func (v *VPNCore) ConnectIOS(profileJSON string) error      // iOS: через platform.Interface
func (v *VPNCore) StartWithTunFd(fd int) error              // Android: с TUN fd
func (v *VPNCore) Disconnect() error
func (v *VPNCore) GetStatus() string                        // disconnected|connecting|connected|error
func (v *VPNCore) GetBytesUp() int64
func (v *VPNCore) GetBytesDown() int64
func (v *VPNCore) GetLatencyMS() int64
func (v *VPNCore) ImportURI(uri string) (string, error)     // возвращает profileJSON
func (v *VPNCore) ListProfiles() string                     // JSON массив
func (v *VPNCore) DeleteProfile(id string) error
func (v *VPNCore) SetStatusCallback(cb StatusCallback)
func (v *VPNCore) SetKeychainStorage(storage KeychainStorage) // только iOS
func (v *VPNCore) SetPacketFlow(flow PacketTunnelFlow)        // только iOS
func (v *VPNCore) GetSupportedTransports(carrier string) string

// gomobile поддерживает только примитивы и интерфейсы.
// Никаких map, slice структур — только string, int64, bool, []byte, интерфейсы.

type StatusCallback interface {
    OnStatusChanged(status, message string)
}

type PacketTunnelFlow interface {
    ReadPacket() []byte
    WritePacket(data []byte) bool
}

type KeychainStorage interface {
    Set(key, value string) error
    Get(key string) (string, error)
    Delete(key string) error
}
```

---

## 14. Desktop GUI (Fyne v2)

Тема — тёмная, по умолчанию:
- Background: `#0D0D0D`, Surface: `#1A1A1A`
- Primary/accent: `#00E5FF`, Success: `#00C853`, Error: `#FF1744`
- Text: `#FFFFFF`, TextMuted: `#888888`

Правила:
- UI-обновления только в главной горутине.
- Статус из `statusCh` читается в отдельной горутине, обновляется через `canvas.Refresh(widget)`.
- Логи: `binding.NewStringList()` + `widget.NewListWithData()` для автоскролла.

---

## 15. Подписки — 4 формата

```go
// internal/subscription/sub.go

// Формат 1: Base64 (V2Ray стандарт) — декодируй → split "\n" → каждая строка URI
// Формат 2: Clash YAML — парси proxies[] секцию
// Формат 3: sing-box JSON — прямой парсинг
// Формат 4: olcrtc sub — https://github.com/openlibrecommunity/olcrtc/blob/master/docs/sub.md

func DetectFormat(data []byte) SubFormat
func ParseSubscription(data []byte) ([]*core.Profile, error)
func FetchAndParse(url string) ([]*core.Profile, error)

// При загрузке: минимум TLS 1.2, верифицировать сертификат, timeout 15s
// InsecureSkipVerify = true запрещено
```

---

## 16. Парсеры URI

VLESS:
```
vless://<uuid>@<host>:<port>?security=reality&sni=...&pbk=...&sid=...&fp=chrome&flow=xtls-rprx-vision#<name>
```
Обязательно для Reality: параметры `pbk` (public key) и `sid` (short id).

olcrtc:
```
olcrtc://<Carrier>?<Transport>@<RoomID>#<Key>%<ClientID>$<MIMO>
```
Документация: https://github.com/openlibrecommunity/olcrtc/blob/master/docs/uri.md

---

## 17. Kill Switch

```go
type KillSwitch struct { enabled bool }

// При отключении или ошибке — СНАЧАЛА блокируй трафик, ПОТОМ останавливай движок.
// sing-box route: "final": "block" пока не connected
// Android: VpnService.Builder() без addAllowedApplication
// iOS: NEPacketTunnelNetworkSettings без лишних allowedRoutes

func (ks *KillSwitch) Enable()
func (ks *KillSwitch) Disable()
```

---

## 18. Метрики

```go
// internal/core/metrics.go
// TCP connect до сервера каждые 30 секунд

func (m *Manager) PingProfile(p *Profile) (time.Duration, error) {
    start := time.Now()
    conn, err := net.DialTimeout("tcp",
        fmt.Sprintf("%s:%d", addr, port), 5*time.Second)
    if err != nil { return 0, err }
    conn.Close()
    return time.Since(start), nil
}
// Трафик: из sing-box stats API
```

---

## 19. Roadmap

### Этап 1 — Ядро (делай сейчас)
- [ ] `profile.go` — типы данных
- [ ] `storage.go` — JSON + шифрование через keyring
- [ ] `singbox/parse.go` — vless://, ss://, trojan://, vmess://
- [ ] `olcrtc/parse.go` — olcrtc://
- [ ] `singbox/config.go` — генерация JSON конфига
- [ ] `singbox/engine.go` — Start/Stop
- [ ] `olcrtc/engine.go` — Start/Stop
- [ ] `manager.go` — Connect/Disconnect/Status/Import
- [ ] `proxy/socks.go` — SOCKS5 с auth (§8)
- [ ] Тесты на парсеры и генерацию конфига

### Этап 2 — Desktop GUI
- [ ] Fyne app + тема
- [ ] home, profiles, add_profile, logs, settings экраны
- [ ] TUN режим Linux/Windows

### Этап 3 — Android
- [ ] `mobile/mobile.go` — gomobile API
- [ ] Android Studio + VpnService.kt
- [ ] Тест на устройстве

### Этап 4 — iOS
- [ ] `mage ios` → xcframework
- [ ] Xcode: App target + Extension target
- [ ] `ConnectIOS()` через sing-box platform.Interface
- [ ] `PacketTunnelProvider.swift`
- [ ] `KeychainBridge.swift`
- [ ] `GetSupportedTransports()` — только datachannel на iOS (§10)
- [ ] Запрос Network Extension entitlement у Apple
- [ ] TestFlight распространение

### Этап 5 — Безопасность и полировка
- [ ] Kill Switch
- [ ] Автозапуск
- [ ] Метрики трафика
- [ ] Все форматы подписок
- [ ] Пинг всех профилей

---

## 20. Тестирование

Обязательные тесты:
- `singbox/parse_test.go` — вариации vless://, ss://, trojan://, vmess://
- `olcrtc/parse_test.go` — olcrtc:// URI
- `singbox/config_test.go` — валидация JSON через sing-box парсер
- `storage_test.go` — шифрование/дешифрование, сохранение/загрузка
- `socks_test.go` — тест должен падать при анонимном подключении к SOCKS5

Не тестируем: Fyne UI (нет headless), реальные сетевые соединения (только mock).

---

## 21. Чего НЕ делать

1. Не запускай SOCKS5/HTTP прокси без аутентификации — критическая уязвимость (§8)
2. Не храни UUID/Password/Key в открытом тексте в JSON на диске
3. Не используй `InsecureSkipVerify = true` нигде
4. Не вызывай blocking операции в UI-горутине
5. Не включай xray HandlerService или ClashAPI без secret токена
6. Не хардкодь порты — только через Settings с дефолтами 2080/2081
7. Не используй `os.Exit()` вне `main.go`
8. Не поддерживай vp8channel/seichannel на iOS — риск OOM kill Extension (§10)
9. Не добавляй зависимости без обсуждения
10. Не игнорируй `context.Done()` в горутинах движков
11. Не использй `_ = err` — все ошибки обрабатывать

---

## 22. Справочные ссылки

| Ресурс | URL |
|--------|-----|
| sing-box конфиг | https://sing-box.sagernet.org/configuration/ |
| sing-box experimental API | https://github.com/sagernet/sing-box/tree/dev/experimental |
| olcrtc URI | https://github.com/openlibrecommunity/olcrtc/blob/master/docs/uri.md |
| olcrtc sub | https://github.com/openlibrecommunity/olcrtc/blob/master/docs/sub.md |
| olcrtc настройки | https://github.com/openlibrecommunity/olcrtc/blob/master/docs/settings.md |
| Fyne docs | https://docs.fyne.io |
| gomobile | https://pkg.go.dev/golang.org/x/mobile/cmd/gomobile |
| pion/webrtc | https://github.com/pion/webrtc |
| SOCKS5 уязвимость PoC | https://github.com/runetfreedom/per-app-split-bypass-poc |
| Happ (iOS архитектура) | https://github.com/hiddify/hiddify-app |
| olcbox (Kotlin UX паттерны) | https://github.com/alananisimov/olcbox |

---

*Версия: 2.0 | Май 2026*