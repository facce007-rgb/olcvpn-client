@echo off
REM GitHub Release Creation Script
REM Создаёт релиз v1.0.0 с бинарниками и source code

echo 🚀 Creating GitHub Release v1.0.0
echo.

REM Проверяем наличие release директории
if not exist release\ (
    echo ❌ Release directory not found! Run: mage release
    exit /b 1
)

echo 📦 Release artifacts:
dir /b release\

echo.
echo 📝 Creating release on GitHub...
echo.
echo MANUAL STEPS (GitHub CLI not available):
echo.
echo 1. Go to: https://github.com/yanisplugg/olcvpn-client/releases/new
echo.
echo 2. Fill in:
echo    - Tag: v1.0.0 (select existing)
echo    - Title: OLC VPN Client v1.0.0
echo    - Description:
echo.
echo      ## 🎉 First Release
echo.
echo      Cross-platform VPN client with sing-box and olcrtc engines.
echo.
echo      ### ✨ Features
echo      - sing-box engine: VLESS+Reality, Shadowsocks, Trojan, VMess
echo      - olcrtc engine: WebRTC tunnel via Russian white services
echo      - Windows desktop GUI (Fyne)
echo      - Android support (Kotlin/Compose)
echo      - iOS support (SwiftUI)
echo      - Subscription support (4 formats)
echo      - SOCKS5 with authentication (security fix)
echo.
echo      ### 📦 Downloads
echo      - **Windows**: olcvpn-windows.zip
echo      - **macOS**: olcvpn-macos.zip (Universal Binary)
echo      - **Linux**: olcvpn-linux.tar.gz
echo      - **Android**: olcvpn-android.apk
echo      - **iOS**: Build from source (see BUILD.md)
echo.
echo      ### 📚 Documentation
echo      - [Build Instructions](BUILD.md)
echo      - [Architecture](CLAUDE.md)
echo.
echo 3. Upload files from release\ directory:
echo    - olcvpn-windows.zip
echo    - olcvpn-macos.zip
echo    - olcvpn-linux.tar.gz
echo    - olcvpn-android.apk (if exists)
echo.
echo 4. Check "Set as the latest release"
echo.
echo 5. Click "Publish release"
echo.
echo 📌 Source code (zip/tar.gz) will be added automatically by GitHub!
echo.
pause
