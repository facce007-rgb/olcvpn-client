import NetworkExtension
import VPNCore

class PacketTunnelProvider: NEPacketTunnelProvider {
    private var core: VPNCore?
    private var packetFlow: PacketFlowBridge?

    override func startTunnel(options: [String : NSObject]?, completionHandler: @escaping (Error?) -> Void) {
        guard let config = protocolConfiguration as? NETunnelProviderProtocol,
              let profileJSON = config.providerConfiguration?["profile"] as? String else {
            completionHandler(NSError(domain: "com.olc.vpn", code: 1, userInfo: [NSLocalizedDescriptionKey: "No profile configuration"]))
            return
        }

        let settings = NEPacketTunnelNetworkSettings(tunnelRemoteAddress: "127.0.0.1")

        let ipv4Settings = NEIPv4Settings(addresses: ["10.0.0.1"], subnetMasks: ["255.255.255.0"])
        ipv4Settings.includedRoutes = [NEIPv4Route.default()]
        settings.ipv4Settings = ipv4Settings

        let dnsSettings = NEDNSSettings(servers: ["1.1.1.1", "8.8.8.8"])
        settings.dnsSettings = dnsSettings

        settings.mtu = 1500

        setTunnelNetworkSettings(settings) { [weak self] error in
            guard let self = self else { return }

            if let error = error {
                completionHandler(error)
                return
            }

            self.core = VPNCore()

            let keychainBridge = KeychainBridge()
            self.core?.setKeychainStorage(keychainBridge)

            self.packetFlow = PacketFlowBridge(packetFlow: self.packetTunnelFlow)
            self.core?.setPacketFlow(self.packetFlow)

            self.core?.setStatusCallback(StatusCallbackImpl { [weak self] status, message in
                if status == "connected" {
                    self?.reasserting = false
                } else if status == "error" {
                    self?.cancelTunnelWithError(NSError(domain: "com.olc.vpn", code: 2, userInfo: [NSLocalizedDescriptionKey: message]))
                }
            })

            do {
                try self.core?.connectIOS(profileJSON)
                completionHandler(nil)
            } catch {
                completionHandler(error)
            }
        }
    }

    override func stopTunnel(with reason: NEProviderStopReason, completionHandler: @escaping () -> Void) {
        do {
            try core?.disconnect()
        } catch {
            NSLog("Error disconnecting: \(error)")
        }
        core = nil
        packetFlow = nil
        completionHandler()
    }

    override func handleAppMessage(_ messageData: Data, completionHandler: ((Data?) -> Void)?) {
        if let message = String(data: messageData, encoding: .utf8) {
            switch message {
            case "status":
                let status = core?.getStatus() ?? "disconnected"
                completionHandler?(status.data(using: .utf8))
            case "metrics":
                let bytesUp = core?.getBytesUp() ?? 0
                let bytesDown = core?.getBytesDown() ?? 0
                let latency = core?.getLatencyMS() ?? 0
                let metrics = "{\\"up\\":\(bytesUp),\\"down\\":\(bytesDown),\\"latency\\":\(latency)}"
                completionHandler?(metrics.data(using: .utf8))
            default:
                completionHandler?(nil)
            }
        } else {
            completionHandler?(nil)
        }
    }
}

class PacketFlowBridge: NSObject, MobilePacketTunnelFlow {
    private let flow: NEPacketTunnelFlow

    init(packetFlow: NEPacketTunnelFlow) {
        self.flow = packetFlow
        super.init()
    }

    func readPacket() -> Data? {
        var packets: [Data] = []
        var protocols: [NSNumber] = []
        flow.readPackets { packets, protocols in
            // Callback
        }
        return packets.first
    }

    func writePacket(_ data: Data?) -> Bool {
        guard let data = data else { return false }
        return flow.writePackets([data], withProtocols: [AF_INET as NSNumber])
    }
}

class StatusCallbackImpl: NSObject, MobileStatusCallback {
    private let callback: (String, String) -> Void

    init(callback: @escaping (String, String) -> Void) {
        self.callback = callback
        super.init()
    }

    func onStatusChanged(_ status: String?, message: String?) {
        callback(status ?? "", message ?? "")
    }
}
