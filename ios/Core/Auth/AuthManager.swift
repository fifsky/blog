import Foundation

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
            } else {
                KeychainService.delete(key: Self.tokenKey)
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
                KeychainService.set(key: Self.expiresAtKey, value: "\(date.timeIntervalSince1970)")
            } else {
                KeychainService.delete(key: Self.expiresAtKey)
            }
        }
    }

    /// 当前登录用户
    var currentUser: User?

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
        self.accessToken = KeychainService.get(key: Self.tokenKey)
    }

    // MARK: - 登录

    /// 用户登录
    /// - Parameters:
    ///   - userName: 用户名
    ///   - password: 密码
    ///   - totpCode: TOTP 验证码（可选）
    /// - Returns: 登录响应
    func login(userName: String, password: String, totpCode: String? = nil) async throws -> LoginResponse {
        let request = LoginRequest(
            user_name: userName,
            password: password,
            totp_code: totpCode
        )

        let response: LoginResponse = try await APIClient.shared.request(
            path: Config.loginPath,
            body: request
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

    /// 清除登录状态
    func logout() {
        accessToken = nil
        expiresAt = nil
        currentUser = nil
    }
}
