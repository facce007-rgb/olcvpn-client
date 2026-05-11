import SwiftUI

struct SettingsView: View {
    @State private var autoConnect = false
    @State private var killSwitch = false
    @State private var socksPort = "2080"
    @State private var httpPort = "2081"

    var body: some View {
        NavigationView {
            Form {
                Section(header: Text("Connection")) {
                    Toggle("Auto-connect on startup", isOn: $autoConnect)
                    Toggle("Kill Switch", isOn: $killSwitch)
                }

                Section(header: Text("Local Proxy")) {
                    HStack {
                        Text("SOCKS5 Port")
                        Spacer()
                        TextField("2080", text: $socksPort)
                            .keyboardType(.numberPad)
                            .multilineTextAlignment(.trailing)
                            .frame(width: 80)
                    }

                    HStack {
                        Text("HTTP Port")
                        Spacer()
                        TextField("2081", text: $httpPort)
                            .keyboardType(.numberPad)
                            .multilineTextAlignment(.trailing)
                            .frame(width: 80)
                    }
                }

                Section(header: Text("About")) {
                    HStack {
                        Text("Version")
                        Spacer()
                        Text("1.0.0")
                            .foregroundColor(.gray)
                    }

                    HStack {
                        Text("sing-box")
                        Spacer()
                        Text("1.11.0")
                            .foregroundColor(.gray)
                    }

                    HStack {
                        Text("olcrtc")
                        Spacer()
                        Text("latest")
                            .foregroundColor(.gray)
                    }
                }

                Section {
                    Link("GitHub Repository", destination: URL(string: "https://github.com/openlibrecommunity/olcvpn")!)
                    Link("Documentation", destination: URL(string: "https://github.com/openlibrecommunity/olcvpn/blob/master/CLAUDE.md")!)
                }
            }
            .navigationTitle("Settings")
            .background(Color(red: 0.05, green: 0.05, blue: 0.05).ignoresSafeArea())
        }
    }
}
