#!/bin/bash
# Manual GitHub Release Creation Script

echo "🚀 Creating GitHub Release v1.0.0 manually"
echo ""

# Проверяем наличие release файлов
if [ ! -d "release" ]; then
    echo "❌ No release directory found!"
    echo "Run: mage release"
    exit 1
fi

echo "📦 Available files:"
ls -lh release/

echo ""
echo "Creating release via GitHub API..."

# Создаём релиз
RELEASE_DATA=$(cat <<EOF
{
  "tag_name": "v1.0.0",
  "name": "OLC VPN Client v1.0.0",
  "body": "## 🎉 First Stable Release\n\nCross-platform VPN client with sing-box and olcrtc engines.\n\n### ✨ Features\n- **sing-box engine**: VLESS+Reality, Shadowsocks, Trojan, VMess, Hysteria2\n- **olcrtc engine**: WebRTC tunnel via Russian white services (WB Stream, Jazz, Telemost)\n- **Windows**: Desktop GUI with Fyne\n- **macOS**: Build from source (requires macOS host)\n- **Linux**: Build from source (requires Linux host)\n- **Android**: Native app with Kotlin/Compose\n- **iOS**: Build from source (see BUILD.md)\n- **Subscription support**: 4 formats (Base64, Clash, sing-box, olcrtc)\n- **Security**: SOCKS5 with authentication, encrypted profile storage\n\n### 📦 Downloads\n- **Windows**: \`olcvpn-windows-amd64.zip\` (21MB)\n- **Other platforms**: Build from source (see BUILD.md)\n\n### 🍎 iOS Build Instructions\niOS requires building from source on macOS:\n\n1. Install Xcode 15+\n2. Install Go 1.23+\n3. Install gomobile: \`go install golang.org/x/mobile/cmd/gomobile@latest && gomobile init\`\n4. Build framework: \`mage ios\`\n5. Open \`ios/OlcVPN.xcodeproj\` in Xcode\n6. Request Network Extension entitlement from Apple Developer Portal\n7. Build and run on device\n\nSee [ios/BUILD.md](ios/BUILD.md) for detailed instructions.\n\n### 📚 Documentation\n- [Build Instructions](BUILD.md)\n- [Architecture & Guidelines](CLAUDE.md)\n- [GitHub Repository](https://github.com/yanisplugg/olcvpn-client)\n\n### 🔒 Security Notes\n- SOCKS5 proxy requires authentication (fixes per-app split bypass vulnerability)\n- Profile secrets encrypted with system keyring\n- No plaintext credentials in storage\n\n---\n\n**Source code** is automatically included by GitHub (zip/tar.gz).",
  "draft": false,
  "prerelease": false
}
EOF
)

# Нужен GitHub token
if [ -z "$GITHUB_TOKEN" ]; then
    echo "❌ GITHUB_TOKEN not set!"
    echo ""
    echo "Manual steps:"
    echo "1. Go to: https://github.com/yanisplugg/olcvpn-client/releases/new"
    echo "2. Select tag: v1.0.0"
    echo "3. Title: OLC VPN Client v1.0.0"
    echo "4. Upload: release/olcvpn-windows-amd64.zip"
    echo "5. Copy description from above"
    echo "6. Publish release"
    exit 1
fi

# Создаём релиз через API
RELEASE_RESPONSE=$(curl -s -X POST \
  -H "Authorization: token $GITHUB_TOKEN" \
  -H "Accept: application/vnd.github.v3+json" \
  https://api.github.com/repos/yanisplugg/olcvpn-client/releases \
  -d "$RELEASE_DATA")

UPLOAD_URL=$(echo "$RELEASE_RESPONSE" | grep -o '"upload_url": "[^"]*' | cut -d'"' -f4 | sed 's/{?name,label}//')

if [ -z "$UPLOAD_URL" ]; then
    echo "❌ Failed to create release!"
    echo "$RELEASE_RESPONSE"
    exit 1
fi

echo "✅ Release created!"
echo "Uploading assets..."

# Загружаем Windows билд
for file in release/*.zip; do
    if [ -f "$file" ]; then
        filename=$(basename "$file")
        echo "Uploading $filename..."
        curl -s -X POST \
          -H "Authorization: token $GITHUB_TOKEN" \
          -H "Content-Type: application/zip" \
          --data-binary @"$file" \
          "${UPLOAD_URL}?name=${filename}"
        echo "✅ $filename uploaded"
    fi
done

echo ""
echo "✅ Release published!"
echo "View at: https://github.com/yanisplugg/olcvpn-client/releases/tag/v1.0.0"
