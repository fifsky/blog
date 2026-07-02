// 公共类型

import Foundation

/// 通用 ID 响应
struct IDResponse: Decodable {
    let id: Int
}

/// 空响应（对应 google.protobuf.Empty，返回 {}）
struct EmptyResponse: Decodable {}

/// 空请求（对应 google.protobuf.Empty，发送 {}）
struct EmptyRequest: Encodable {}

/// 文章删除请求（批量）
struct ArticleDeleteRequest: Encodable {
    let ids: [Int]
}

/// 管理端评论列表请求
struct AdminCommentListRequest: Encodable {
    let page: Int
    let keyword: String?
}

/// 评论删除请求（批量）
struct CommentDeleteRequest: Encodable {
    let ids: [Int]
}

/// 心情删除请求（批量）
struct MoodDeleteRequest: Encodable {
    let ids: [Int]
}

/// 获取足迹全部响应（公开接口）
struct GetFootprintsResponse: Decodable {
    let footprints: [Footprint]
}

/// 足迹删除请求
struct FootprintDeleteRequest: Encodable {
    let id: Int
}

/// 用户摘要信息
struct UserSummary: Decodable, Hashable {
    let id: Int
    let name: String
    let nick_name: String
}

/// 分类摘要信息
struct CateSummary: Decodable, Hashable {
    let id: Int
    let name: String
    let domain: String
}

/// 错误响应
struct ErrorResponse: Decodable {
    let code: String
    let message: String
}
