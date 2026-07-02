// 心情/动态相关模型

import Foundation

/// 心情/动态
struct Mood: Identifiable, Decodable {
    let id: Int
    let content: String
    let user: UserSummary?
    let created_at: String     // proto: string（RFC3339 格式）
    let updated_at: String     // proto: string（RFC3339 格式）
}

/// 心情列表请求
struct MoodListRequest: Encodable {
    let page: Int
}

/// 心情列表响应
struct MoodListResponse: Decodable {
    let list: [Mood]
    let total: Int
}

/// 创建心情请求
struct MoodCreateRequest: Encodable {
    let content: String
}

/// 更新心情请求
struct MoodUpdateRequest: Encodable {
    let id: Int
    let content: String?
}
