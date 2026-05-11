# 🚀 Быстрый старт

## Для пользователей

### Windows
```bash
# Скачать olcvpn-windows.zip
# Распаковать
# Запустить olcvpn.exe
```

### macOS
```bash
# Скачать olcvpn-macos.zip
# Распаковать
# Запустить olcvpn
```

### Linux
```bash
wget https://github.com/openlibrecommunity/olcvpn/releases/latest/download/olcvpn-linux.tar.gz
tar -xzf olcvpn-linux.tar.gz
./olcvpn
```

### Android
```bash
# Скачать olcvpn.apk
# Установить
# Запустить
```

## Для разработчиков

### Быстрый запуск
```bash
go run ./cmd/olcvpn/main_v2.go
```

### Сборка релизов

**Windows:**
```bash
build-release.bat
```

**Linux/macOS:**
```bash
./build-release.sh
```

**Или через mage:**
```bash
mage release
```

Все файлы будут в папке `release/`:
- `olcvpn-windows.zip` - Windows
- `olcvpn-macos.zip` - macOS (Universal)
- `olcvpn-linux.tar.gz` - Linux
- `olcvpn.apk` - Android
- `olcvpn-ios-framework.zip` - iOS framework

### Автоматический релиз

При создании тега `v*` GitHub Actions автоматически соберёт все платформы:

```bash
git tag v0.1.0
git push origin v0.1.0
```

Релиз появится на https://github.com/openlibrecommunity/olcvpn/releases

## Тесты

```bash
go test ./...
```

## Документация

См. [CLAUDE.md](CLAUDE.md)
