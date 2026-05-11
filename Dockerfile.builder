# Multi-stage builder для кросс-компиляции OLC VPN Client
FROM golang:1.23-bookworm AS base

# Устанавливаем зависимости для Fyne
RUN apt-get update && apt-get install -y \
    gcc g++ \
    libgl1-mesa-dev \
    xorg-dev \
    gcc-mingw-w64-x86-64 \
    gcc-aarch64-linux-gnu \
    libc6-dev-arm64-cross \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /build

# Копируем исходники
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Сборка Windows
FROM base AS windows-builder
ENV CGO_ENABLED=1
ENV GOOS=windows
ENV GOARCH=amd64
ENV CC=x86_64-w64-mingw32-gcc
ENV CXX=x86_64-w64-mingw32-g++
RUN go build -ldflags="-s -w -H windowsgui" -o /out/olcvpn.exe ./cmd/olcvpn/main.go

# Сборка Linux
FROM base AS linux-builder
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64
RUN go build -ldflags="-s -w" -o /out/olcvpn ./cmd/olcvpn/main.go

# Сборка macOS (требует osxcross, упрощённая версия)
FROM base AS macos-builder
RUN apt-get update && apt-get install -y clang lld && rm -rf /var/lib/apt/lists/*
ENV CGO_ENABLED=1
ENV GOOS=darwin
ENV GOARCH=amd64
ENV CC=clang
# Примечание: полноценная сборка macOS требует osxcross SDK
RUN go build -ldflags="-s -w" -o /out/olcvpn-macos ./cmd/olcvpn/main.go || echo "macOS build requires osxcross"

# Финальный образ с артефактами
FROM scratch AS artifacts
COPY --from=windows-builder /out/olcvpn.exe /windows/
COPY --from=linux-builder /out/olcvpn /linux/
COPY --from=macos-builder /out/olcvpn-macos /macos/
