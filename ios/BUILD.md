# iOS Build Instructions

## Prerequisites

1. **macOS** with Xcode 15+
2. **Go 1.22+** and **gomobile**
3. **Apple Developer Account** (for Network Extension entitlement)

## Build Steps

### 1. Build gomobile xcframework

```bash
# Install gomobile
go install golang.org/x/mobile/cmd/gomobile@latest
gomobile init

# Build xcframework
cd C:\freeclaude
gomobile bind -target ios -o ios/Frameworks/VPNCore.xcframework ./mobile/
```

### 2. Create Xcode Project

1. Open Xcode
2. Create new iOS App project
3. Product Name: `OlcVPN`
4. Bundle Identifier: `com.olc.vpn`
5. Interface: SwiftUI
6. Language: Swift

### 3. Add Network Extension Target

1. File → New → Target
2. Select "Network Extension"
3. Product Name: `OlcVPNExtension`
4. Bundle Identifier: `com.olc.vpn.extension`

### 4. Add Files to Project

Copy these files to Xcode project:

**Main App Target (OlcVPN):**
- `ios/OlcVPN/OlcVPNApp.swift`
- `ios/OlcVPN/ContentView.swift`
- `ios/OlcVPN/ProfilesView.swift`
- `ios/OlcVPN/LogsView.swift`
- `ios/OlcVPN/SettingsView.swift`
- `ios/OlcVPN/VPNManager.swift`
- `ios/OlcVPN/KeychainBridge.swift`
- `ios/OlcVPN/Info.plist`
- `ios/OlcVPN/OlcVPN.entitlements`

**Extension Target (OlcVPNExtension):**
- `ios/OlcVPNExtension/PacketTunnelProvider.swift`
- `ios/OlcVPNExtension/Info.plist`
- `ios/OlcVPNExtension/OlcVPNExtension.entitlements`

### 5. Add VPNCore.xcframework

1. Drag `ios/Frameworks/VPNCore.xcframework` into Xcode project
2. Select both targets (OlcVPN and OlcVPNExtension)
3. Embed & Sign

### 6. Configure Capabilities

**Main App Target:**
1. Signing & Capabilities → + Capability
2. Add "Network Extensions"
3. Add "App Groups" → `group.com.olc.vpn`
4. Add "Keychain Sharing" → `com.olc.vpn`

**Extension Target:**
1. Signing & Capabilities → + Capability
2. Add "Network Extensions" → Packet Tunnel
3. Add "App Groups" → `group.com.olc.vpn`
4. Add "Keychain Sharing" → `com.olc.vpn`

### 7. Request Network Extension Entitlement

1. Go to https://developer.apple.com/contact/request/network-extension/
2. Fill out the form explaining VPN use case
3. Wait for approval (usually 1-2 weeks)
4. Download provisioning profiles with entitlement

### 8. Build and Run

```bash
# Build from Xcode: Product → Build
# Or via command line:
xcodebuild -project ios/OlcVPN.xcodeproj -scheme OlcVPN -configuration Debug
```

### 9. Install on Device

1. Connect iPhone/iPad via USB
2. Select device in Xcode
3. Product → Run
4. Trust developer certificate on device (Settings → General → VPN & Device Management)

## Project Structure

```
ios/
├── OlcVPN/                          # Main app target
│   ├── OlcVPNApp.swift              # App entry point
│   ├── ContentView.swift            # Home screen
│   ├── ProfilesView.swift           # Profiles management
│   ├── LogsView.swift               # Real-time logs
│   ├── SettingsView.swift           # Settings
│   ├── VPNManager.swift             # VPN state management
│   ├── KeychainBridge.swift         # iOS Keychain integration
│   ├── Info.plist
│   └── OlcVPN.entitlements
├── OlcVPNExtension/                 # Network Extension target
│   ├── PacketTunnelProvider.swift   # VPN tunnel implementation
│   ├── Info.plist
│   └── OlcVPNExtension.entitlements
└── Frameworks/
    └── VPNCore.xcframework          # Go core (generated)
```

## Features

✅ SwiftUI with dark theme
✅ Material Design-inspired UI
✅ Network Extension with PacketTunnelProvider
✅ iOS Keychain for secrets
✅ Real-time traffic statistics
✅ Profile management
✅ Logs viewer
✅ Settings

## Permissions

- `com.apple.developer.networking.networkextension` - VPN functionality
- `com.apple.security.application-groups` - Share data between app and extension
- `keychain-access-groups` - Secure storage

## Testing

1. Enable Developer Mode on iOS 16+ (Settings → Privacy & Security → Developer Mode)
2. Connect device via USB
3. Run from Xcode
4. Grant VPN permission when prompted

## Troubleshooting

### xcframework not found
```bash
# Rebuild xcframework
gomobile bind -target ios -o ios/Frameworks/VPNCore.xcframework ./mobile/
```

### Network Extension entitlement missing
- Request entitlement from Apple Developer portal
- Download new provisioning profiles
- Refresh profiles in Xcode (Preferences → Accounts → Download Manual Profiles)

### VPN permission denied
- Check Info.plist has correct NSExtension configuration
- Verify entitlements are properly set
- User must grant VPN permission in app

### Extension crashes with OOM
- Only use `datachannel` transport on iOS (see CLAUDE.md §10)
- Avoid `vp8channel` and `seichannel` — they cause memory issues

## Distribution

### TestFlight (Recommended for Russia)

1. Archive app: Product → Archive
2. Upload to App Store Connect
3. Add to TestFlight
4. Share link with users (up to 10,000 testers)

**Advantages:**
- No App Store review for VPN restrictions
- Free distribution
- Easy updates

### App Store

1. Submit for review
2. Note: VPN apps are removed from Russian App Store by RKN
3. Use non-Russian developer account to avoid removal

### Sideloading

1. Export IPA: Product → Archive → Distribute App → Ad Hoc
2. Users install via AltStore or similar tools
3. Requires re-signing every 7 days (free account) or 1 year (paid account)

## Known Limitations

- **olcrtc on iOS**: Only `datachannel` transport supported (§10 in CLAUDE.md)
- **Memory limit**: Network Extension has 50-200 MB limit depending on device
- **Background**: Extension can be killed by iOS if inactive for 15+ minutes
- **TUN interface**: Must use `platform.Interface` instead of TUN fd (unlike Android)

## References

- [Apple Network Extension Guide](https://developer.apple.com/documentation/networkextension)
- [gomobile documentation](https://pkg.go.dev/golang.org/x/mobile/cmd/gomobile)
- [Hiddify iOS implementation](https://github.com/hiddify/hiddify-app) (reference architecture)
- [sing-box iOS integration](https://github.com/SagerNet/sing-box)
