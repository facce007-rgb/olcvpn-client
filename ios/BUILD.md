# iOS Build Instructions

## Requirements

- macOS 13+ (Ventura or later)
- Xcode 15+
- Go 1.23+
- Apple Developer Account (for Network Extension entitlement)

## Setup

1. **Install Go and gomobile:**
   ```bash
   brew install go
   go install golang.org/x/mobile/cmd/gomobile@latest
   gomobile init
   ```

2. **Build Go framework:**
   ```bash
   cd olcvpn-client
   mage ios
   ```
   This creates `ios/Frameworks/VPNCore.xcframework`

3. **Request Network Extension entitlement:**
   - Go to https://developer.apple.com/account
   - Navigate to Certificates, Identifiers & Profiles
   - Request `com.apple.developer.networking.networkextension` capability
   - Add `packet-tunnel-provider` to your App ID

4. **Open Xcode project:**
   ```bash
   open ios/OlcVPN.xcodeproj
   ```

5. **Configure signing:**
   - Select OlcVPN target
   - Go to Signing & Capabilities
   - Select your Team
   - Ensure Network Extension capability is enabled

6. **Configure Extension target:**
   - Select OlcVPNExtension target
   - Configure signing with same Team
   - Verify Network Extension entitlement is present

7. **Build and run:**
   - Select a physical device (Network Extension doesn't work in Simulator)
   - Press Cmd+R to build and run

## Project Structure

```
ios/
├── OlcVPN/                    # Main app target
│   ├── ContentView.swift      # Main UI
│   ├── ProfilesView.swift     # Profile management
│   ├── VPNManager.swift       # NETunnelProviderManager wrapper
│   └── KeychainBridge.swift   # iOS Keychain implementation
├── OlcVPNExtension/           # Network Extension target (separate process!)
│   └── PacketTunnelProvider.swift
└── Frameworks/
    └── VPNCore.xcframework    # Go code compiled via gomobile
```

## Important Notes

### Two Separate Targets

iOS VPN apps require **two targets**:
1. **Main app** (OlcVPN) — UI, manages VPN connection
2. **Network Extension** (OlcVPNExtension) — runs VPN, separate process

The Go core runs in the Extension, not the main app.

### Memory Constraints

Network Extensions have strict memory limits:
- Old devices: ~50MB
- New devices: ~200MB

**OOM kills are silent** — VPN just stops working.

### olcrtc Transport Support

On iOS, only `datachannel` transport is supported:
- ✅ `wbstream + datachannel`
- ✅ `jazz + datachannel`
- ❌ `telemost + vp8channel` (software VP8 encoder too heavy)
- ❌ `telemost + seichannel` (H.264 from Go problematic)

See CLAUDE.md §10 for detailed analysis.

### Routing Loop Prevention

iOS automatically excludes Network Extension sockets from the tunnel.
No need for manual `protect()` calls like on Android.

### Keychain Storage

iOS uses native Keychain instead of `go-keyring`:
- Secrets stored in iOS Keychain
- Accessed via `KeychainBridge.swift`
- Survives app reinstalls (if configured)

## Distribution

### TestFlight (Recommended for Russia)

App Store removes VPN apps in Russia. Use TestFlight:
- Up to 10,000 users
- No RKN review
- Free with Apple Developer account

### App Store (Outside Russia)

Standard App Store submission works outside Russia.

### Sideloading

Users can install via:
- AltStore
- Sideloadly
- Xcode (for developers)

## Troubleshooting

### "Network Extension entitlement not found"

Request the entitlement from Apple Developer Portal first.

### VPN connects but no traffic

Check:
1. Extension target has correct bundle ID
2. Extension entitlements file exists
3. `PacketTunnelProvider` is correctly configured

### Memory crashes

Monitor memory in Xcode Instruments. If crashes occur:
1. Reduce sing-box features
2. Disable olcrtc video transports
3. Optimize buffer sizes

### Build errors

```bash
# Clean and rebuild
mage clean
rm -rf ios/Frameworks/VPNCore.xcframework
mage ios
```

## References

- [Apple Network Extension Guide](https://developer.apple.com/documentation/networkextension)
- [Hiddify iOS Implementation](https://github.com/hiddify/hiddify-app) — reference architecture
- [sing-box iOS Platform Interface](https://github.com/sagernet/sing-box/tree/dev/experimental/libbox)
