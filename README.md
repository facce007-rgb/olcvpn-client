# 🎉 OLC VPN Client

[![Release](https://img.shields.io/github/v/release/yanisplugg/olcvpn-client?style=flat-square)](https://github.com/yanisplugg/olcvpn-client/releases)
[![CI](https://img.shields.io/github/actions/workflow/status/yanisplugg/olcvpn-client/ci.yml?branch=main&style=flat-square&label=CI)](https://github.com/yanisplugg/olcvpn-client/actions)
[![License](https://img.shields.io/github/license/yanisplugg/olcvpn-client?style=flat-square)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/yanisplugg/olcvpn-client?style=flat-square)](go.mod)

Кросс-платформенный VPN-клиент с поддержкой sing-box и olcrtc.

## 📥 Скачать

**[➡️ Перейти к релизам](../../releases/latest)**

Выберите версию для вашей системы:

| Платформа | Файл | Инструкция |
|-----------|------|------------|
| 🪟 **Windows** | `olcvpn-windows.zip` | Распаковать и запустить `olcvpn.exe` |
| 🐧 **Linux** | `olcvpn-linux.tar.gz` | `tar -xzf olcvpn-linux.tar.gz && ./olcvpn` |
| 🍎 **macOS** | `olcvpn-macos.zip` | Распаковать и запустить |
| 🤖 **Android** | `olcvpn.apk` | Скачать и установить APK |
| 📱 **iOS** | TestFlight | Скоро |

## ✨ Возможности

- ✅ **sing-box** — VLESS, Reality, Shadowsocks, Trojan, VMess, Hysteria2, TUIC
- ✅ **olcrtc** — WebRTC туннель через белые сервисы (WB Stream, Jazz, Telemost)
- ✅ **Kill Switch** — защита от утечек при разрыве соединения
- ✅ **Подписки** — Base64, Clash, sing-box, olcrtc форматы
- ✅ **QR коды** — сканирование и генерация
- ✅ **Material Dark** — современный тёмный интерфейс
- ✅ **Метрики** — скорость, трафик, пинг в реальном времени

## 🚀 Быстрый старт

### Windows
1. Скачайте `olcvpn-windows.zip` из [релизов](../../releases/latest)
2. Распакуйте архив
3. Запустите `olcvpn.exe`

### Linux
```bash
wget https://github.com/yanisplugg/olcvpn-client/releases/latest/download/olcvpn-linux.tar.gz
tar -xzf olcvpn-linux.tar.gz
./olcvpn
```

### Android
1. Скачайте `olcvpn.apk` из [релизов](../../releases/latest)
2. Разрешите установку из неизвестных источников
3. Установите APK
4. Запустите приложение

## 🔧 Сборка из исходников

### Требования
- Go 1.23+
- MinGW-w64 (Windows)
- GCC (Linux/macOS)

### Windows
```bash
git clone https://github.com/yanisplugg/olcvpn-client.git
cd olcvpn-client
install-mingw-fixed.bat  # Установить MinGW
build-release.bat        # Собрать
```

### Linux/macOS
```bash
git clone https://github.com/yanisplugg/olcvpn-client.git
cd olcvpn-client
go build -o olcvpn ./cmd/olcvpn
```

## 📚 Документация

- [CLAUDE.md](CLAUDE.md) — полная техническая документация
- [DEPLOY.md](DEPLOY.md) — инструкция по публикации релизов

## 🛡️ Безопасность

- ✅ SOCKS5 с обязательной аутентификацией
- ✅ AES-256-GCM шифрование секретов
- ✅ Системный keyring (DPAPI/Keychain/Secret Service)
- ✅ TLS 1.3 для всех соединений

## 🤝 Вклад в проект

Pull requests приветствуются! Перед большими изменениями откройте issue для обсуждения.

## 📄 Лицензия

MIT License

## 🔗 Ссылки

- [sing-box](https://github.com/sagernet/sing-box)
- [olcrtc](https://github.com/openlibrecommunity/olcrtc)
- [Fyne](https://fyne.io)

---

**Проект готов к использованию!** 🎯

