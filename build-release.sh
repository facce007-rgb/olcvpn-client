#!/bin/bash
# Скрипт для быстрой сборки всех релизов

set -e

echo "🚀 OLC VPN Client - Release Builder"
echo "===================================="

# Проверяем зависимости
command -v go >/dev/null 2>&1 || { echo "❌ Go не установлен"; exit 1; }
command -v gomobile >/dev/null 2>&1 || { echo "⚠️  gomobile не установлен, устанавливаю..."; go install golang.org/x/mobile/cmd/gomobile@latest; gomobile init; }

# Создаём директорию release
mkdir -p release
rm -rf release/*

echo ""
echo "📦 Сборка Windows..."
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -o release/tmp/olcvpn.exe ./cmd/olcvpn/main_v2.go
cd release/tmp && zip ../olcvpn-windows.zip olcvpn.exe && cd ../..
rm -rf release/tmp

echo ""
echo "📦 Сборка macOS..."
mkdir -p release/tmp
GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -o release/tmp/olcvpn-amd64 ./cmd/olcvpn/main_v2.go
GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 go build -o release/tmp/olcvpn-arm64 ./cmd/olcvpn/main_v2.go
if command -v lipo >/dev/null 2>&1; then
    lipo -create -output release/tmp/olcvpn release/tmp/olcvpn-amd64 release/tmp/olcvpn-arm64
else
    cp release/tmp/olcvpn-arm64 release/tmp/olcvpn
fi
chmod +x release/tmp/olcvpn
cd release/tmp && zip ../olcvpn-macos.zip olcvpn && cd ../..
rm -rf release/tmp

echo ""
echo "📦 Сборка Linux..."
mkdir -p release/tmp
GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o release/tmp/olcvpn ./cmd/olcvpn/main_v2.go
chmod +x release/tmp/olcvpn
tar -czf release/olcvpn-linux.tar.gz -C release/tmp olcvpn
rm -rf release/tmp

echo ""
echo "📦 Сборка Android AAR..."
mkdir -p android/app/libs
gomobile bind -target android -androidapi 21 -o android/app/libs/vpncore.aar ./mobile/

if [ -f "android/gradlew" ]; then
    echo ""
    echo "📦 Сборка Android APK..."
    cd android
    chmod +x gradlew
    ./gradlew assembleRelease
    cp app/build/outputs/apk/release/app-release-unsigned.apk ../release/olcvpn.apk
    cd ..
    echo "⚠️  APK не подписан. Для подписи используйте jarsigner."
else
    echo "⚠️  Android Gradle не найден, пропускаю сборку APK"
fi

echo ""
echo "📦 Сборка iOS xcframework..."
mkdir -p ios/Frameworks
gomobile bind -target ios -o ios/Frameworks/VPNCore.xcframework ./mobile/
cd ios/Frameworks && zip -r ../../release/olcvpn-ios-framework.zip VPNCore.xcframework && cd ../..

echo ""
echo "✅ Сборка завершена!"
echo ""
echo "📂 Файлы в директории release/:"
ls -lh release/
echo ""
echo "🎉 Готово к распространению!"
