import Foundation

// MARK: - 通知名

/// 用户退出登录通知，根视图监听后切回登录页
extension Notification.Name {
    static let didLogout = Notification.Name("AuthManagerDidLogout")
}

// MARK: - 模型定义

/// 登录请求
struct LoginRequest: Encodable {
    let user_name: String
    let password: String
    let totp_code: String?
}

/// 登录响应
/// protojson 编码：int64 字段(expires_at)输出为字符串，message 字段(user)未赋值时为 null
struct LoginResponse: Decodable {
    let access_token: String?
    let user: User?
    let require_totp: Bool?
    /// protojson 将 int64 编码为字符串
    let expires_at: String?
}

/// 用户信息（登录返回的完整用户对象，对应后端 UserItem）
struct User: Decodable, Identifiable {
    let id: Int
    let name: String
    let nick_name: String
    let email: String
    let status: String
    let type: Int
    let created_at: String
    let updated_at: String
}

// MARK: - AuthManager

/// 认证管理器
/// 负责登录状态管理、Token 存储
@MainActor
@Observable
class AuthManager {

    static let shared = AuthManager()

    // MARK: - Keychain 存储键名

    private static let tokenKey = "access_token"
    private static let expiresAtKey = "expires_at"

    // MARK: - 属性

    /// 当前访问令牌（从 Keychain 读取）
    var accessToken: String? {
        didSet {
            if let token = accessToken {
                KeychainService.set(key: Self.tokenKey, value: token)
                #if DEBUG
                DebugAuthCache.set(key: Self.tokenKey, value: token)
                #endif
            } else {
                KeychainService.delete(key: Self.tokenKey)
                #if DEBUG
                DebugAuthCache.delete(key: Self.tokenKey)
                #endif
            }
        }
    }

    /// Token 过期时间（Unix 时间戳）
    var expiresAt: Date? {
        get {
            guard let str = KeychainService.get(key: Self.expiresAtKey),
                  let timestamp = Double(str)
            else { return nil }
            return Date(timeIntervalSince1970: timestamp)
        }
        set {
            if let date = newValue {
                let value = "\(date.timeIntervalSince1970)"
                KeychainService.set(key: Self.expiresAtKey, value: value)
                #if DEBUG
                DebugAuthCache.set(key: Self.expiresAtKey, value: value)
                #endif
            } else {
                KeychainService.delete(key: Self.expiresAtKey)
                #if DEBUG
                DebugAuthCache.delete(key: Self.expiresAtKey)
                #endif
            }
        }
    }

    /// 当前登录用户
    var currentUser: User?

    /// 认证接口服务
    private let authService = AuthService.shared

    /// 是否已登录（Token 存在且未过期）
    var isLoggedIn: Bool {
        guard let token = accessToken, !token.isEmpty else { return false }
        if let expiresAt {
            return Date() < expiresAt
        }
        // 没有过期时间信息，仅检查 Token 是否存在
        return true
    }

    private init() {
        // 初始化时从 Keychain 恢复 Token
        #if DEBUG
        let token = KeychainService.get(key: Self.tokenKey) ?? DebugAuthCache.get(key: Self.tokenKey)
        self.accessToken = token
        if KeychainService.get(key: Self.tokenKey) == nil, let token {
            KeychainService.set(key: Self.tokenKey, value: token)
        }
        if KeychainService.get(key: Self.expiresAtKey) == nil,
           let expiresAt = DebugAuthCache.get(key: Self.expiresAtKey) {
            KeychainService.set(key: Self.expiresAtKey, value: expiresAt)
        }
        #else
        self.accessToken = KeychainService.get(key: Self.tokenKey)
        #endif
    }

    // MARK: - 登录

    /// 用户登录
    /// - Parameters:
    ///   - userName: 用户名
    ///   - password: 密码
    ///   - totpCode: TOTP 验证码（可选）
    /// - Returns: 登录响应
    func login(userName: String, password: String, totpCode: String? = nil) async throws -> LoginResponse {
        let response = try await authService.login(
            userName: userName,
            password: password,
            totpCode: totpCode
        )

        // 需要二次验证时（无 token），不存储登录信息
        if response.require_totp == true {
            return response
        }

        // 存储 Token 和过期时间
        if let token = response.access_token, !token.isEmpty {
            accessToken = token
        }
        if let expiresStr = response.expires_at,
           let timestamp = TimeInterval(expiresStr) {
            expiresAt = Date(timeIntervalSince1970: timestamp)
        }
        currentUser = response.user

        return response
    }

    // MARK: - 登出

    /// 清除登录状态，并广播退出登录通知
    func logout() {
        accessToken = nil
        expiresAt = nil
        currentUser = nil
        NotificationCenter.default.post(name: .didLogout, object: nil)
    }
}

#if DEBUG
/// 调试登录缓存
///
/// 仅用于模拟器/Debug 构建兜底：当重新编译安装导致 Keychain 读不到旧 token 时，
/// 用应用容器里的镜像恢复一次，避免每次调试都重新登录。
private enum DebugAuthCache {

    private static let namespace = "debug_auth_cache"

    /// 读取调试缓存
    static func get(key: String) -> String? {
        UserDefaults.standard.string(forKey: cacheKey(key))
    }

    /// 写入调试缓存
    static func set(key: String, value: String) {
        UserDefaults.standard.set(value, forKey: cacheKey(key))
    }

    /// 删除调试缓存
    static func delete(key: String) {
        UserDefaults.standard.removeObject(forKey: cacheKey(key))
    }

    private static func cacheKey(_ key: String) -> String {
        "\(namespace).\(key)"
    }
}
#endif
