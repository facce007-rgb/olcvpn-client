import Foundation
import NetworkExtension
import Combine

class VPNManager: ObservableObject {
    static let shared = VPNManager()

    @Published var status: NEVPNStatus = .invalid
    @Published var bytesUp: Int64 = 0
    @Published var bytesDown: Int64 = 0
    @Published var latencyMS: Int64 = 0

    private var manager: NETunnelProviderManager?
    private var statusObserver: NSObjectProtocol?
    private var metricsTimer: Timer?

    private init() {
        loadManager()
        observeStatus()
    }

    deinit {
        if let observer = statusObserver {
            NotificationCenter.default.removeObserver(observer)
        }
        metricsTimer?.invalidate()
    }

    private func loadManager() {
        NETunnelProviderManager.loadAllFromPreferences { [weak self] managers, error in
            if let error = error {
                print("Failed to load VPN manager: \(error)")
                return
            }

            if let manager = managers?.first {
                self?.manager = manager
                self?.status = manager.connection.status
            } else {
                self?.createManager()
            }
        }
    }

    private func createManager() {
        let manager = NETunnelProviderManager()
        manager.localizedDescription = "OLC VPN"

        let proto = NETunnelProviderProtocol()
        proto.providerBundleIdentifier = "com.olc.vpn.extension"
        proto.serverAddress = "OLC VPN Server"
        manager.protocolConfiguration = proto

        manager.isEnabled = true

        manager.saveToPreferences { [weak self] error in
            if let error = error {
                print("Failed to save VPN manager: \(error)")
                return
            }

            self?.manager = manager
            self?.status = manager.connection.status
        }
    }

    private func observeStatus() {
        statusObserver = NotificationCenter.default.addObserver(
            forName: .NEVPNStatusDidChange,
            object: nil,
            queue: .main
        ) { [weak self] notification in
            guard let connection = notification.object as? NEVPNConnection else { return }
            self?.status = connection.status

            if connection.status == .connected {
                self?.startMetricsTimer()
            } else {
                self?.stopMetricsTimer()
                self?.bytesUp = 0
                self?.bytesDown = 0
                self?.latencyMS = 0
            }
        }
    }

    func connect(profileJSON: String) {
        guard let manager = manager else {
            print("VPN manager not loaded")
            return
        }

        let proto = manager.protocolConfiguration as? NETunnelProviderProtocol
        proto?.providerConfiguration = ["profile": profileJSON]

        manager.saveToPreferences { [weak self] error in
            if let error = error {
                print("Failed to save profile: \(error)")
                return
            }

            do {
                try manager.connection.startVPNTunnel()
            } catch {
                print("Failed to start VPN: \(error)")
            }
        }
    }

    func disconnect() {
        manager?.connection.stopVPNTunnel()
    }

    private func startMetricsTimer() {
        metricsTimer?.invalidate()
        metricsTimer = Timer.scheduledTimer(withTimeInterval: 1.0, repeats: true) { [weak self] _ in
            self?.updateMetrics()
        }
    }

    private func stopMetricsTimer() {
        metricsTimer?.invalidate()
        metricsTimer = nil
    }

    private func updateMetrics() {
        guard let session = manager?.connection as? NETunnelProviderSession else { return }

        do {
            try session.sendProviderMessage("metrics".data(using: .utf8)!) { [weak self] response in
                guard let data = response,
                      let json = try? JSONSerialization.jsonObject(with: data) as? [String: Int64] else {
                    return
                }

                DispatchQueue.main.async {
                    self?.bytesUp = json["up"] ?? 0
                    self?.bytesDown = json["down"] ?? 0
                    self?.latencyMS = json["latency"] ?? 0
                }
            }
        } catch {
            print("Failed to get metrics: \(error)")
        }
    }

    var statusText: String {
        switch status {
        case .invalid: return "Not Configured"
        case .disconnected: return "Disconnected"
        case .connecting: return "Connecting..."
        case .connected: return "Connected"
        case .reasserting: return "Reconnecting..."
        case .disconnecting: return "Disconnecting..."
        @unknown default: return "Unknown"
        }
    }

    var isConnected: Bool {
        status == .connected
    }
}
