import Foundation

/// Multipart/form-data 请求构建辅助工具
enum MultipartFormData {

    /// 构建 multipart/form-data 格式的 URLRequest
    /// - Parameters:
    ///   - url: 上传目标地址
    ///   - boundary: 分隔符
    ///   - fileData: 文件二进制数据
    ///   - fieldName: 表单字段名（如 "uploadFile"）
    ///   - filename: 文件名
    ///   - mimeType: 文件 MIME 类型（如 "image/jpeg"）
    ///   - token: 可选的认证 Token
    /// - Returns: 构建好的 URLRequest
    static func buildRequest(
        url: URL,
        boundary: String,
        fileData: Data,
        fieldName: String,
        filename: String,
        mimeType: String,
        token: String? = nil
    ) -> URLRequest {
        var request = URLRequest(url: url)
        request.httpMethod = "POST"

        // 设置 Content-Type
        let contentType = "multipart/form-data; boundary=\(boundary)"
        request.setValue(contentType, forHTTPHeaderField: "Content-Type")

        // 携带认证 Token
        if let token {
            request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        }

        // 构建请求体
        var body = Data()

        // 添加文件字段
        body.append("--\(boundary)\r\n".data(using: .utf8)!)
        body.append(
            "Content-Disposition: form-data; name=\"\(fieldName)\"; filename=\"\(filename)\"\r\n"
                .data(using: .utf8)!
        )
        body.append("Content-Type: \(mimeType)\r\n\r\n".data(using: .utf8)!)
        body.append(fileData)
        body.append("\r\n".data(using: .utf8)!)

        // 结束标记
        body.append("--\(boundary)--\r\n".data(using: .utf8)!)

        request.httpBody = body
        return request
    }

    /// 生成随机 boundary 字符串
    static func generateBoundary() -> String {
        let uuid = UUID().uuidString
        return "Boundary-\(uuid)"
    }
}
