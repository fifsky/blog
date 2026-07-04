import Foundation

/// API 错误类型
enum APIError: LocalizedError {

    /// 请求被客户端主动取消
    case cancelled

    /// 未授权（需要重新登录）
    case unauthorized

    /// 无效的响应（无法解析）
    case invalidResponse

    /// 服务器错误（包含错误码、错误信息和 HTTP 状态码）
    case serverError(code: String?, message: String, statusCode: Int?)

    /// 客户端错误（网络、编码、未知本地错误等）
    case clientError(code: String, message: String, underlying: Error?)

    /// 解码错误（JSON 解析失败）
    case decodingError(Error)

    /// 上传失败（如 OSS 上传失败）
    case uploadFailed

    /// 是否为可静默忽略的取消错误
    var isCancelled: Bool {
        if case .cancelled = self {
            return true
        }
        return false
    }

    /// 统一错误码，便于外部只判断 APIError
    var code: String {
        switch self {
        case .cancelled:
            "CANCELLED"
        case .unauthorized:
            "UNAUTHORIZED"
        case .invalidResponse:
            "INVALID_RESPONSE"
        case .serverError(let code, _, let statusCode):
            code ?? statusCode.map(String.init) ?? "SERVER_ERROR"
        case .clientError(let code, _, _):
            code
        case .decodingError:
            "DECODING_ERROR"
        case .uploadFailed:
            "UPLOAD_FAILED"
        }
    }

    // MARK: - LocalizedError

    var errorDescription: String? {
        switch self {
        case .cancelled:
            nil
        case .unauthorized:
            "登录已过期，请重新登录"
        case .invalidResponse:
            "服务器返回了无效的响应"
        case .serverError(_, let message, _):
            message
        case .clientError(_, let message, _):
            message
        case .decodingError(let error):
            "数据解析失败：\(error.localizedDescription)"
        case .uploadFailed:
            "文件上传失败"
        }
    }

    /// 将任意 Error 归一化为 APIError
    static func normalized(_ error: Error) -> APIError {
        if let apiError = error as? APIError {
            return apiError
        }

        if error is CancellationError {
            return .cancelled
        }

        if let urlError = error as? URLError {
            if urlError.code == .cancelled {
                return .cancelled
            }

            return .clientError(
                code: String(urlError.errorCode),
                message: "网络连接失败：\(urlError.localizedDescription)",
                underlying: urlError
            )
        }

        return .clientError(
            code: "UNKNOWN",
            message: error.localizedDescription,
            underlying: error
        )
    }

    /// 带业务前缀的展示文案
    func displayMessage(prefix: String? = nil) -> String {
        let message = errorDescription ?? ""
        guard let prefix, !prefix.isEmpty else {
            return message
        }
        return "\(prefix)：\(message)"
    }
}
