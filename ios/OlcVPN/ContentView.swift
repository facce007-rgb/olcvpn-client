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
        .accentColor(Color(red: 0, green: 0.898, blue: 1))
    }
}

struct HomeView: View {
    @StateObject private var vpnManager = VPNManager.shared

    var body: some View {
        NavigationView {
            ScrollView {
                VStack(spacing: 20) {
                    StatusCard()
                    ControlsCard()
                    MetricsCard()
                }
                .padding()
            }
            .navigationTitle("OLC VPN")
            .background(Color(red: 0.05, green: 0.05, blue: 0.05).ignoresSafeArea())
        }
    }
}

struct StatusCard: View {
    @StateObject private var vpnManager = VPNManager.shared

    var body: some View {
        VStack(alignment: .leading, spacing: 12) {
            Text("Connection Status")
                .font(.headline)
                .foregroundColor(.white)

            HStack {
                Circle()
                    .fill(statusColor)
                    .frame(width: 12, height: 12)

                Text(vpnManager.statusText)
                    .font(.title2)
                    .foregroundColor(.white)
            }

            Text("No profile selected")
                .font(.subheadline)
                .foregroundColor(.gray)
        }
        .frame(maxWidth: .infinity, alignment: .leading)
        .padding()
        .background(Color(red: 0.1, green: 0.1, blue: 0.1))
        .cornerRadius(12)
    }

    private var statusColor: Color {
        switch vpnManager.status {
        case .connected: return Color(red: 0, green: 0.784, blue: 0.325)
        case .connecting, .reasserting: return Color(red: 1, green: 0.757, blue: 0)
        case .disconnecting: return Color.gray
        default: return Color(red: 1, green: 0.09, blue: 0.267)
        }
    }
}

struct ControlsCard: View {
    @StateObject private var vpnManager = VPNManager.shared

    var body: some View {
        VStack(alignment: .leading, spacing: 12) {
            Text("Controls")
                .font(.headline)
                .foregroundColor(.white)

            Button(action: {
                // TODO: Connect to first profile
            }) {
                Text("Connect")
                    .frame(maxWidth: .infinity)
                    .padding()
                    .background(vpnManager.isConnected ? Color.gray : Color(red: 0, green: 0.898, blue: 1))
                    .foregroundColor(.black)
                    .cornerRadius(8)
            }
            .disabled(vpnManager.isConnected)

            Button(action: {
                vpnManager.disconnect()
            }) {
                Text("Disconnect")
                    .frame(maxWidth: .infinity)
                    .padding()
                    .background(vpnManager.isConnected ? Color(red: 1, green: 0.09, blue: 0.267) : Color.gray)
                    .foregroundColor(.white)
                    .cornerRadius(8)
            }
            .disabled(!vpnManager.isConnected)
        }
        .frame(maxWidth: .infinity, alignment: .leading)
        .padding()
        .background(Color(red: 0.1, green: 0.1, blue: 0.1))
        .cornerRadius(12)
    }
}

struct MetricsCard: View {
    @StateObject private var vpnManager = VPNManager.shared

    var body: some View {
        VStack(alignment: .leading, spacing: 12) {
            Text("Metrics")
                .font(.headline)
                .foregroundColor(.white)

            HStack {
                VStack(alignment: .leading) {
                    Text("Upload")
                        .font(.caption)
                        .foregroundColor(.gray)
                    Text(formatBytes(vpnManager.bytesUp))
                        .font(.title3)
                        .foregroundColor(.white)
                }

                Spacer()

                VStack(alignment: .leading) {
                    Text("Download")
                        .font(.caption)
                        .foregroundColor(.gray)
                    Text(formatBytes(vpnManager.bytesDown))
                        .font(.title3)
                        .foregroundColor(.white)
                }

                Spacer()

                VStack(alignment: .leading) {
                    Text("Latency")
                        .font(.caption)
                        .foregroundColor(.gray)
                    Text(vpnManager.latencyMS > 0 ? "\(vpnManager.latencyMS) ms" : "-")
                        .font(.title3)
                        .foregroundColor(.white)
                }
            }
        }
        .frame(maxWidth: .infinity, alignment: .leading)
        .padding()
        .background(Color(red: 0.1, green: 0.1, blue: 0.1))
        .cornerRadius(12)
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
}
