# ⚠️ Требования для сборки

## Windows

Для сборки на Windows нужен **MinGW-w64** (для CGO):

### Установка MinGW-w64

**Вариант 1: Через MSYS2 (рекомендуется)**
```bash
# Скачать MSYS2: https://www.msys2.org/
# Установить и запустить MSYS2 MINGW64
pacman -S mingw-w64-x86_64-gcc
```

**Вариант 2: Через Chocolatey**
```bash
choco install mingw
```

**Вариант 3: Прямая загрузка**
- Скачать: https://github.com/niXman/mingw-builds-binaries/releases
- Распаковать в `C:\mingw64`
- Добавить `C:\mingw64\bin` в PATH

### Проверка установки
```bash
gcc --version
```

## Linux

```bash
# Ubuntu/Debian
sudo apt-get install gcc libgl1-mesa-dev xorg-dev

# Fedora
sudo dnf install gcc mesa-libGL-devel libXcursor-devel libXrandr-devel libXinerama-devel libXi-devel

# Arch
sudo pacman -S gcc mesa libxcursor libxrandr libxinerama libxi
```

## macOS

```bash
xcode-select --install
```

## Альтернатива: Сборка без GUI

Если не нужен GUI, можно собрать CLI версию без CGO:

```bash
CGO_ENABLED=0 go build -o olcvpn-cli ./cmd/olcvpn-cli/
```

## Готовые бинарники

Скачать готовые бинарники можно с GitHub Releases:
https://github.com/openlibrecommunity/olcvpn/releases
