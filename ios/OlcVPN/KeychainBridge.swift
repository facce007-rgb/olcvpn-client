import Foundation
import Security
import VPNCore

class KeychainBridge: NSObject, MobileKeychainStorage {
    private let service = "com.olc.vpn"

    func set(_ key: String?, value: String?) throws {
        guard let key = key, let value = value else {
            throw NSError(domain: "com.olc.vpn.keychain", code: 1, userInfo: [NSLocalizedDescriptionKey: "Invalid key or value"])
        }

        let data = value.data(using: .utf8)!

        let query: [String: Any] = [
            kSecClass as String: kSecClassGenericPassword,
            kSecAttrService as String: service,
            kSecAttrAccount as String: key,
            kSecAttrAccessGroup as String: "group.com.olc.vpn"
        ]

        SecItemDelete(query as CFDictionary)

        var addQuery = query
        addQuery[kSecValueData as String] = data
        addQuery[kSecAttrAccessible as String] = kSecAttrAccessibleAfterFirstUnlock

        let status = SecItemAdd(addQuery as CFDictionary, nil)

        if status != errSecSuccess {
            throw NSError(domain: "com.olc.vpn.keychain", code: Int(status), userInfo: [NSLocalizedDescriptionKey: "Failed to save to keychain"])
        }
    }

    func get(_ key: String?) throws -> String {
        guard let key = key else {
            throw NSError(domain: "com.olc.vpn.keychain", code: 1, userInfo: [NSLocalizedDescriptionKey: "Invalid key"])
        }

        let query: [String: Any] = [
            kSecClass as String: kSecClassGenericPassword,
            kSecAttrService as String: service,
            kSecAttrAccount as String: key,
            kSecAttrAccessGroup as String: "group.com.olc.vpn",
            kSecReturnData as String: true,
            kSecMatchLimit as String: kSecMatchLimitOne
        ]

        var result: AnyObject?
        let status = SecItemCopyMatching(query as CFDictionary, &result)

        if status == errSecSuccess, let data = result as? Data, let value = String(data: data, encoding: .utf8) {
            return value
        } else if status == errSecItemNotFound {
            throw NSError(domain: "com.olc.vpn.keychain", code: Int(status), userInfo: [NSLocalizedDescriptionKey: "Key not found"])
        } else {
            throw NSError(domain: "com.olc.vpn.keychain", code: Int(status), userInfo: [NSLocalizedDescriptionKey: "Failed to read from keychain"])
        }
    }

    func delete(_ key: String?) throws {
        guard let key = key else {
            throw NSError(domain: "com.olc.vpn.keychain", code: 1, userInfo: [NSLocalizedDescriptionKey: "Invalid key"])
        }

        let query: [String: Any] = [
            kSecClass as String: kSecClassGenericPassword,
            kSecAttrService as String: service,
            kSecAttrAccount as String: key,
            kSecAttrAccessGroup as String: "group.com.olc.vpn"
        ]

        let status = SecItemDelete(query as CFDictionary)

        if status != errSecSuccess && status != errSecItemNotFound {
            throw NSError(domain: "com.olc.vpn.keychain", code: Int(status), userInfo: [NSLocalizedDescriptionKey: "Failed to delete from keychain"])
        }
    }
}
