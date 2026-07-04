import Foundation

/// API 客户端（单例）
/// 所有后端接口均使用 POST + JSON Body，响应为 protojson 编码（snake_case 键）
actor APIClient {

    static let shared = APIClient()

    private let session: URLSession
    private let decoder: JSONDecoder
    private let encoder: JSONEncoder

    private init() {
        let config = URLSessionConfiguration.default
        config.timeoutIntervalForRequest = 30
        config.timeoutIntervalForResource = 60
        self.session = URLSession(configuration: config)

        // protojson 使用 snake_case，模型属性直接使用 snake_case 命名
        self.decoder = JSONDecoder()

        self.encoder = JSONEncoder()
    }

    // MARK: - 通用请求方法

    /// 发起 POST JSON 请求
    /// - Parameters:
    ///   - path: API 路径（如 "/blog/login"）
    ///   - body: 请求体（Codable 对象，可为 nil）
    ///   - auth: 是否携带 Authorization 头
    ///   - headers: 附加的自定义请求头（如 X-Client-Type）
    /// - Returns: 解码后的响应对象
    func request<T: Decodable>(
        path: String,
        body: Encodable? = nil,
        auth: Bool = false,
        headers: [String: String]? = nil
    ) async throws -> T {
        guard let url = URL(string: Config.baseURL + path) else {
            throw APIError.invalidResponse
        }

        var urlRequest = URLRequest(url: url)
        urlRequest.httpMethod = "POST"
        urlRequest.setValue("application/json", forHTTPHeaderField: "Content-Type")

        // 设置附加的自定义请求头
        if let headers {
            for (key, value) in headers {
                urlRequest.setValue(value, forHTTPHeaderField: key)
            }
        }

        // 携带认证 Token
        if auth {
            // AuthManager 是 @Observable class（非 actor），accessToken 为同步属性，无需 await
            let token = AuthManager.shared.accessToken
            if let token {
                urlRequest.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
            }
        }

        // 编码请求体
        if let body {
            do {
                urlRequest.httpBody = try encoder.encode(body)
            } catch {
                throw APIError.clientError(
                    code: "ENCODING_ERROR",
                    message: "请求数据编码失败：\(error.localizedDescription)",
                    underlying: error
                )
            }
        }

        // 调试日志：打印实际发送的请求（NSLog 确保进入 unified log）
        #if DEBUG
        let bodyString = urlRequest.httpBody.flatMap { String(data: $0, encoding: .utf8) } ?? "nil"
        NSLog("[API] REQ %@ %@ body: %@", urlRequest.httpMethod ?? "?", urlRequest.url?.absoluteString ?? "?", bodyString)
        #endif

        // 发起请求
        let data: Data
        let httpResponse: HTTPURLResponse
        do {
            (data, httpResponse) = try await send(urlRequest)
        } catch {
            throw APIError.normalized(error)
        }

        // 调试日志：打印响应状态码和内容
        #if DEBUG
        let respString = String(data: data, encoding: .utf8) ?? "<binary>"
        NSLog("[API] RESP %d %@ body: %@", httpResponse.statusCode, urlRequest.url?.lastPathComponent ?? "?", respString)
        #endif

        // 处理非 2xx 状态码
        guard (200..<300).contains(httpResponse.statusCode) else {
            throw responseError(from: data, statusCode: httpResponse.statusCode)
        }

        // 解码成功响应
        do {
            return try decoder.decode(T.self, from: data)
        } catch {
            throw APIError.decodingError(error)
        }
    }

    /// 执行 URLRequest，并统一归一化网络层错误
    func send(_ request: URLRequest) async throws -> (Data, HTTPURLResponse) {
        do {
            let (data, response) = try await session.data(for: request)
            guard let httpResponse = response as? HTTPURLResponse else {
                throw APIError.invalidResponse
            }
            return (data, httpResponse)
        } catch {
            throw APIError.normalized(error)
        }
    }

    /// 解析非 2xx HTTP 响应
    func responseError(from data: Data, statusCode: Int) -> APIError {
        if statusCode == 401 {
            return .unauthorized
        }

        // 尝试解析错误响应 {"code": "...", "message": "..."}
        if let errorInfo = try? decoder.decode(ServerErrorResponse.self, from: data) {
            return .serverError(
                code: errorInfo.code,
                message: errorInfo.message,
                statusCode: statusCode
            )
        }

        return .serverError(
            code: "\(statusCode)",
            message: "请求失败（HTTP \(statusCode)）",
            statusCode: statusCode
        )
    }

    /// 构建 API URL（用于非 JSON 请求，如 multipart/form-data 上传）
    func buildURL(path: String) throws -> URL {
        guard let url = URL(string: Config.baseURL + path) else {
            throw APIError.invalidResponse
        }
        return url
    }
}

// MARK: - 错误响应模型

/// 服务器错误响应结构
/// code 可选：后端业务错误（如 fmt.Errorf）经 FromError 包装后 code 字段为空
private struct ServerErrorResponse: Decodable {
    let code: String?
    let message: String
}
