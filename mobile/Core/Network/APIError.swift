import Foundation

/// API 错误类型
enum APIError: LocalizedError {

    /// 未授权（需要重新登录）
    case unauthorized

    /// 无效的响应（无法解析）
    case invalidResponse

    /// 服务器错误（包含错误码和错误信息）
    /// code 可选：后端业务错误可能不返回 code，仅返回 message
    case serverError(code: String?, message: String)

    /// 网络错误（底层网络请求失败）
    case networkError(Error)

    /// 解码错误（JSON 解析失败）
    case decodingError(Error)

    /// HTTP 错误（非 200 状态码且无法解析为 ServerError）
    case httpError(statusCode: Int)

    /// 上传失败（如 OSS 上传失败）
    case uploadFailed

    // MARK: - LocalizedError

    var errorDescription: String? {
        switch self {
        case .unauthorized:
            "登录已过期，请重新登录"
        case .invalidResponse:
            "服务器返回了无效的响应"
        case .serverError(_, let message):
            message
        case .networkError(let error):
            "网络连接失败：\(error.localizedDescription)"
        case .decodingError(let error):
            "数据解析失败：\(error.localizedDescription)"
        case .httpError(let statusCode):
            "请求失败（HTTP \(statusCode)）"
        case .uploadFailed:
            "文件上传失败"
        }
    }
}
