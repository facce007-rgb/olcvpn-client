# Android Build Instructions

## Prerequisites

1. **Android Studio** (latest version)
2. **JDK 17** or higher
3. **Android SDK** (API 34)
4. **Go 1.22+** and **gomobile**

## Build Steps

### 1. Build gomobile AAR

```bash
# Install gomobile
go install golang.org/x/mobile/cmd/gomobile@latest
gomobile init

# Build AAR
cd C:\freeclaude
gomobile bind -target android -androidapi 21 -o android/app/libs/vpncore.aar ./mobile/
```

### 2. Open in Android Studio

```bash
# Open the android folder in Android Studio
cd android
# File -> Open -> select android folder
```

### 3. Sync Gradle

Android Studio will automatically sync Gradle dependencies.

### 4. Build APK

```bash
# Debug build
./gradlew assembleDebug

# Release build (requires signing)
./gradlew assembleRelease
```

### 5. Install on Device

```bash
# Via Android Studio: Run -> Run 'app'
# Or via adb:
adb install app/build/outputs/apk/debug/app-debug.apk
```

## Project Structure

```
android/
├── app/
│   ├── src/main/
│   │   ├── kotlin/com/olc/vpn/
│   │   │   ├── MainActivity.kt
│   │   │   ├── OlcVpnApplication.kt
│   │   │   ├── service/
│   │   │   │   └── OlcVpnService.kt
│   │   │   ├── ui/
│   │   │   │   ├── screens/
│   │   │   │   │   ├── HomeScreen.kt
│   │   │   │   │   ├── ProfilesScreen.kt
│   │   │   │   │   ├── LogsScreen.kt
│   │   │   │   │   └── SettingsScreen.kt
│   │   │   │   └── theme/
│   │   │   │       ├── Color.kt
│   │   │   │       ├── Theme.kt
│   │   │   │       └── Type.kt
│   │   │   └── viewmodel/
│   │   │       └── VpnViewModel.kt
│   │   ├── res/
│   │   │   ├── values/
│   │   │   │   ├── strings.xml
│   │   │   │   ├── colors.xml
│   │   │   │   └── themes.xml
│   │   │   └── drawable/
│   │   └── AndroidManifest.xml
│   ├── libs/
│   │   └── vpncore.aar (generated)
│   └── build.gradle
├── build.gradle
└── settings.gradle
```

## Features

✅ Material 3 Design (Compose)
✅ Dark theme by default
✅ VPN Service with foreground notification
✅ Traffic statistics
✅ Profile management
✅ Real-time logs
✅ Settings

## Permissions

- `INTERNET` - Network access
- `ACCESS_NETWORK_STATE` - Check connectivity
- `FOREGROUND_SERVICE` - VPN service
- `BIND_VPN_SERVICE` - VPN permission
- `POST_NOTIFICATIONS` - Show notifications

## Testing

1. Enable Developer Options on Android device
2. Enable USB Debugging
3. Connect device via USB
4. Run from Android Studio

## Troubleshooting

### AAR not found
```bash
# Rebuild AAR
gomobile bind -target android -androidapi 21 -o android/app/libs/vpncore.aar ./mobile/
```

### Gradle sync failed
```bash
# Clean and rebuild
./gradlew clean
./gradlew build
```

### VPN permission denied
- Check AndroidManifest.xml has `BIND_VPN_SERVICE` permission
- User must grant VPN permission in app

## Release Build

1. Generate signing key:
```bash
keytool -genkey -v -keystore olcvpn.keystore -alias olcvpn -keyalg RSA -keysize 2048 -validity 10000
```

2. Add to `app/build.gradle`:
```gradle
android {
    signingConfigs {
        release {
            storeFile file("olcvpn.keystore")
            storePassword "your_password"
            keyAlias "olcvpn"
            keyPassword "your_password"
        }
    }
    buildTypes {
        release {
            signingConfig signingConfigs.release
        }
    }
}
```

3. Build:
```bash
./gradlew assembleRelease
```

Output: `app/build/outputs/apk/release/app-release.apk`
