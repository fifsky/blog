import Foundation

// UploadService 上传服务，支持直接上传和 OSS 预签名上传两种方式
class UploadService {
    static let shared = UploadService()

    private init() {}

    // 通过 multipart/form-data 直接上传图片到服务器
    // 返回上传后的图片 URL
    func uploadImage(imageData: Data, filename: String, mimeType: String = "image/jpeg") async throws -> String {
        let boundary = UUID().uuidString
        let url = try await APIClient.shared.buildURL(path: "/blog/admin/upload")

        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("multipart/form-data; boundary=\(boundary)", forHTTPHeaderField: "Content-Type")
        // 添加认证 Token
        if let token = AuthManager.shared.accessToken {
            request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        }

        // 构建 multipart/form-data 请求体
        var body = Data()
        // 添加文件字段
        body.append("--\(boundary)\r\n".data(using: .utf8)!)
        body.append("Content-Disposition: form-data; name=\"uploadFile\"; filename=\"\(filename)\"\r\n".data(using: .utf8)!)
        body.append("Content-Type: \(mimeType)\r\n\r\n".data(using: .utf8)!)
        body.append(imageData)
        body.append("\r\n".data(using: .utf8)!)
        body.append("--\(boundary)--\r\n".data(using: .utf8)!)

        request.httpBody = body

        let (data, response) = try await URLSession.shared.data(for: request)

        guard let httpResponse = response as? HTTPURLResponse else {
            throw APIError.invalidResponse
        }

        guard httpResponse.statusCode == 200 else {
            throw APIError.httpError(statusCode: httpResponse.statusCode)
        }

        // 解析响应中的 url 字段（snake_case 直接匹配 protojson 输出）
        let result = try JSONDecoder().decode(UploadResponse.self, from: data)
        return result.url
    }

    // 通过 OSS 预签名 URL 上传图片
    // Step 1: 获取预签名上传地址（传入 content_type 让后端正确签名）
    // Step 2: PUT 数据到预签名地址
    // 返回 CDN 访问地址
    func uploadViaOSS(imageData: Data, filename: String, contentType: String = "image/jpeg") async throws -> String {
        // 第一步：获取预签名 URL，传入 content_type 让后端匹配签名
        let presignRequest = PresignRequest(filename: filename, content_type: contentType)
        let presignResponse: PresignResponse = try await APIClient.shared.request(
            path: "/blog/admin/oss/presign",
            body: presignRequest,
            auth: true
        )

        // 第二步：PUT 数据到预签名 URL，Content-Type 必须与 presign 时一致
        guard let presignURL = URL(string: presignResponse.url) else {
            throw APIError.invalidResponse
        }
        var putRequest = URLRequest(url: presignURL)
        putRequest.httpMethod = "PUT"
        putRequest.setValue(contentType, forHTTPHeaderField: "Content-Type")
        putRequest.httpBody = imageData

        let (_, putResponse) = try await URLSession.shared.data(for: putRequest)

        guard let httpResponse = putResponse as? HTTPURLResponse, httpResponse.statusCode == 200 else {
            throw APIError.uploadFailed
        }

        return presignResponse.cdn_url
    }
}
