// 上传相关模型

import Foundation

/// 预签名请求
struct PresignRequest: Encodable {
    let filename: String
    let content_type: String?
}

/// 预签名响应
struct PresignResponse: Decodable {
    let url: String
    let cdn_url: String
}

/// 上传响应
struct UploadResponse: Decodable {
    let url: String
}
