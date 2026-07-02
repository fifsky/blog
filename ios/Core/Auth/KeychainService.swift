import Foundation
import Security

/// Keychain 存储服务
/// 用于安全存储敏感数据（如访问令牌）
enum KeychainService {

    // MARK: - 常量

    /// Keychain 服务标识
    private static let service = "com.fifsky.blogapp"

    // MARK: - 公开方法

    /// 从 Keychain 读取字符串值
    /// - Parameter key: 存储键名
    /// - Returns: 存储的字符串值，不存在则返回 nil
    static func get(key: String) -> String? {
        let query: [String: Any] = [
            kSecClass as String: kSecClassGenericPassword,
            kSecAttrService as String: service,
            kSecAttrAccount as String: key,
            kSecReturnData as String: true,
            kSecMatchLimit as String: kSecMatchLimitOne,
        ]

        var item: CFTypeRef?
        let status = SecItemCopyMatching(query as CFDictionary, &item)

        guard status == errSecSuccess else { return nil }

        guard let data = item as? Data,
              let string = String(data: data, encoding: .utf8)
        else {
            return nil
        }

        return string
    }

    /// 向 Keychain 存储字符串值
    /// - Parameters:
    ///   - key: 存储键名
    ///   - value: 要存储的字符串值
    static func set(key: String, value: String) {
        // 先尝试更新，不存在则新增
        let status = SecItemCopyMatching(
            [
                kSecClass as String: kSecClassGenericPassword,
                kSecAttrService as String: service,
                kSecAttrAccount as String: key,
                kSecReturnData as String: false,
            ] as CFDictionary,
            nil
        )

        let data = value.data(using: .utf8) ?? Data()

        if status == errSecSuccess {
            // 已存在，执行更新
            let updateQuery: [String: Any] = [
                kSecClass as String: kSecClassGenericPassword,
                kSecAttrService as String: service,
                kSecAttrAccount as String: key,
            ]
            let attributes: [String: Any] = [
                kSecValueData as String: data,
            ]
            SecItemUpdate(updateQuery as CFDictionary, attributes as CFDictionary)
        } else {
            // 不存在，执行新增
            let addQuery: [String: Any] = [
                kSecClass as String: kSecClassGenericPassword,
                kSecAttrService as String: service,
                kSecAttrAccount as String: key,
                kSecValueData as String: data,
                kSecAttrAccessible as String: kSecAttrAccessibleAfterFirstUnlock,
            ]
            SecItemAdd(addQuery as CFDictionary, nil)
        }
    }

    /// 从 Keychain 删除指定键
    /// - Parameter key: 存储键名
    static func delete(key: String) {
        let query: [String: Any] = [
            kSecClass as String: kSecClassGenericPassword,
            kSecAttrService as String: service,
            kSecAttrAccount as String: key,
        ]
        SecItemDelete(query as CFDictionary)
    }
}
