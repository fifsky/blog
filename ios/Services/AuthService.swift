import Foundation

// AuthService 认证服务，处理用户登录相关接口
class AuthService {
    static let shared = AuthService()

    private init() {}

    // 登录接口
    func login(userName: String, password: String, totpCode: String? = nil) async throws -> LoginResponse {
        let request = LoginRequest(user_name: userName, password: password, totp_code: totpCode)
        // 标识 iOS 客户端，后端据此发放 30 天有效期的 token
        return try await APIClient.shared.request(
            path: Config.loginPath,
            body: request,
            headers: ["X-Client-Type": "ios"]
        )
    }
}
