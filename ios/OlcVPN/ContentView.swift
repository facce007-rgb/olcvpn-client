import SwiftUI

struct ContentView: View {
    @StateObject private var vpnManager = VPNManager.shared
    @State private var selectedTab = 0

    var body: some View {
        TabView(selection: $selectedTab) {
            HomeView()
                .tabItem {
                    Label("Home", systemImage: "house.fill")
                }
                .tag(0)

            ProfilesView()
                .tabItem {
                    Label("Profiles", systemImage: "list.bullet")
                }
                .tag(1)

            LogsView()
                .tabItem {
                    Label("Logs", systemImage: "doc.text")
                }
                .tag(2)

            SettingsView()
                .tabItem {
                    Label("Settings", systemImage: "gear")
                }
                .tag(3)
        }
        .accentColor(Color(red: 0.161, green: 0.235, blue: 0.627)) // #293CA0 - Hiddify primary
    }
}

// HomeView в стиле Hiddify/v2RayTun
struct HomeView: View {
    @StateObject private var vpnManager = VPNManager.shared

    var body: some View {
        NavigationView {
            ZStack {
                // Фон
                Color(red: 0.07, green: 0.07, blue: 0.07)
                    .ignoresSafeArea()

                VStack(spacing: 0) {
                    // Карточка профиля сверху (как в Hiddify)
                    ProfileCard()
                        .padding(.horizontal, 16)
                        .padding(.top, 16)

                    Spacer()

                    // Большая круглая кнопка подключения (как в Hiddify)
                    ConnectionButton()

                    // Индикатор задержки под кругом
                    Divider()
                        .background(Color.gray.opacity(0.3))
                        .padding(.horizontal, 32)
                        .padding(.top, 24)

                    Text(vpnManager.latencyMS > 0 && vpnManager.latencyMS < 65000
                         ? "\(vpnManager.latencyMS) ms"
                         : "- ms")
                        .font(.title3)
                        .fontWeight(.bold)
                        .foregroundColor(vpnManager.latencyMS > 0 && vpnManager.latencyMS < 200
                                       ? Color(red: 0.18, green: 0.49, blue: 0.196)
                                       : .white)
                        .padding(.top, 8)

                    Spacer()

                    // Футер с метриками (как в Hiddify)
                    Divider()
                        .background(Color.gray.opacity(0.3))

                    HStack(spacing: 8) {
                        Text("↑ \(formatBytes(vpnManager.bytesUp))")
                        Text("•")
                        Text("↓ \(formatBytes(vpnManager.bytesDown))")
                    }
                    .font(.body)
                    .foregroundColor(.white)
                    .padding(.vertical, 16)
                }
            }
            .navigationTitle("OLC VPN")
            .navigationBarTitleDisplayMode(.inline)
        }
    }
}

// Карточка профиля (как в Hiddify)
struct ProfileCard: View {
    @StateObject private var vpnManager = VPNManager.shared

    var body: some View {
        VStack(alignment: .leading, spacing: 8) {
            Text("No Profile Selected")
                .font(.headline)
                .fontWeight(.bold)
                .foregroundColor(.white)

            Text("Tap + to add profile")
                .font(.subheadline)
                .foregroundColor(.gray)
        }
        .frame(maxWidth: .infinity, alignment: .leading)
        .padding()
        .background(Color(red: 0.12, green: 0.12, blue: 0.12))
        .cornerRadius(12)
    }
}

// Большая круглая кнопка подключения (как в Hiddify)
struct ConnectionButton: View {
    @StateObject private var vpnManager = VPNManager.shared

    var body: some View {
        ZStack {
            // Круг с обводкой
            Circle()
                .stroke(statusColor, lineWidth: 16)
                .frame(width: 200, height: 200)

            // Заливка для connected/connecting
            if vpnManager.status == .connected || vpnManager.status == .connecting {
                Circle()
                    .fill(statusColor.opacity(0.3))
                    .frame(width: 200, height: 200)
            }

            // Текст статуса в центре
            VStack(spacing: 4) {
                Text(statusText)
                    .font(.title2)
                    .fontWeight(.bold)
                    .foregroundColor(.white)
                    .multilineTextAlignment(.center)
            }
        }
        .onTapGesture {
            if vpnManager.isConnected {
                vpnManager.disconnect()
            } else if vpnManager.status != .connecting {
                // TODO: Connect to first profile
            }
        }
    }

    private var statusColor: Color {
        switch vpnManager.status {
        case .connected:
            return Color(red: 0.18, green: 0.49, blue: 0.196) // #2E7D32 - Green 800
        case .connecting, .reasserting:
            return Color(red: 1, green: 0.757, blue: 0.027) // #FFC107 - Amber
        case .disconnecting:
            return Color.gray
        default:
            return Color(red: 0.247, green: 0.318, blue: 0.71) // #3F51B5 - Indigo
        }
    }

    private var statusText: String {
        switch vpnManager.status {
        case .connected:
            return "Connected"
        case .connecting, .reasserting:
            return "Connecting..."
        case .disconnecting:
            return "Disconnecting..."
        default:
            return "Tap to Connect"
        }
    }
}

private func formatBytes(_ bytes: Int64) -> String {
    let unit: Int64 = 1024
    if bytes < unit {
        return "\(bytes) B"
    }
    var div: Int64 = unit
    var exp = 0
    var n = bytes / unit
    while n >= unit {
        div *= unit
        exp += 1
        n /= unit
    }
    let units = ["K", "M", "G", "T", "P", "E"]
    return String(format: "%.1f %@B", Double(bytes) / Double(div), units[exp])
}
