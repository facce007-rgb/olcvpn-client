import SwiftUI

struct LogsView: View {
    @State private var logs: [LogEntry] = []
    @State private var autoScroll = true

    var body: some View {
        NavigationView {
            VStack {
                ScrollViewReader { proxy in
                    ScrollView {
                        LazyVStack(alignment: .leading, spacing: 4) {
                            ForEach(logs) { log in
                                LogRow(log: log)
                                    .id(log.id)
                            }
                        }
                        .padding()
                    }
                    .onChange(of: logs.count) { _ in
                        if autoScroll, let lastLog = logs.last {
                            withAnimation {
                                proxy.scrollTo(lastLog.id, anchor: .bottom)
                            }
                        }
                    }
                }

                HStack {
                    Toggle("Auto-scroll", isOn: $autoScroll)
                        .toggleStyle(SwitchToggleStyle(tint: Color(red: 0, green: 0.898, blue: 1)))

                    Spacer()

                    Button("Clear") {
                        logs.removeAll()
                    }
                    .foregroundColor(Color(red: 1, green: 0.09, blue: 0.267))
                }
                .padding()
                .background(Color(red: 0.1, green: 0.1, blue: 0.1))
            }
            .navigationTitle("Logs")
            .background(Color(red: 0.05, green: 0.05, blue: 0.05).ignoresSafeArea())
            .onAppear {
                startLogging()
            }
        }
    }

    private func startLogging() {
        // TODO: Subscribe to VPNCore logs
        addLog("VPN client started")
    }

    private func addLog(_ message: String) {
        let timestamp = DateFormatter.localizedString(from: Date(), dateStyle: .none, timeStyle: .medium)
        logs.append(LogEntry(timestamp: timestamp, message: message))
    }
}

struct LogRow: View {
    let log: LogEntry

    var body: some View {
        HStack(alignment: .top, spacing: 8) {
            Text(log.timestamp)
                .font(.system(.caption, design: .monospaced))
                .foregroundColor(.gray)
                .frame(width: 80, alignment: .leading)

            Text(log.message)
                .font(.system(.caption, design: .monospaced))
                .foregroundColor(.white)
        }
    }
}

struct LogEntry: Identifiable {
    let id = UUID()
    let timestamp: String
    let message: String
}
