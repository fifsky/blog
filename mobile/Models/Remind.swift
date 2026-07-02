// 提醒相关模型

import Foundation

/// 提醒状态枚举（proto 枚举值为大写）
enum RemindStatus: String, CaseIterable, Codable {
    case active = "ACTIVE"
    case pending = "PENDING"
    case done = "DONE"
}

/// 提醒
struct Remind: Identifiable, Decodable {
    let id: Int
    let cron: String
    let content: String
    let status: String         // proto: string，原始 JSON 值如 "ACTIVE"/"PENDING"/"DONE"
    let next_time: String       // proto: string（RFC3339 格式）
    let created_at: String      // proto: string（RFC3339 格式）
    let updated_at: String      // proto: string（RFC3339 格式）
}

/// 提醒列表请求
struct RemindListRequest: Encodable {
    let page: Int
}

/// 提醒列表响应
struct RemindListResponse: Decodable {
    let list: [Remind]
    let total: Int
}

/// 创建提醒请求
struct RemindCreateRequest: Encodable {
    let cron: String?
    let content: String
}

/// 更新提醒请求
struct RemindUpdateRequest: Encodable {
    let id: Int
    let cron: String?
    let content: String?
    let status: RemindStatus?  // Encodable 请求，使用枚举编码为字符串
}

/// 删除提醒请求
struct RemindDeleteRequest: Encodable {
    let id: Int
}

/// AI 智能创建提醒请求（AI 自动生成 cron 和内容）
struct RemindSmartCreateRequest: Encodable {
    let content: String
}
