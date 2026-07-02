// 文章相关模型

import Foundation

/// 文章详情
struct Article: Identifiable, Decodable, Hashable {
    let id: Int
    let cate_id: Int
    let type: Int
    let user_id: Int
    let title: String
    let url: String
    let content: String
    let status: String          // proto: string
    let created_at: String      // proto: string（RFC3339 格式）
    let updated_at: String      // proto: string（RFC3339 格式）
    let user: UserSummary?
    let cate: CateSummary?
    let view_num: Int
    let tags: [String]
}

/// 文章列表请求（管理端）
struct ArticleListRequest: Encodable {
    let page: Int?
    let type: Int?
    let status: String?         // proto: string
    let keyword: String?
}

/// 文章列表响应
struct ArticleListResponse: Decodable {
    let list: [Article]
    let total: Int
}

/// 文章详情请求
struct ArticleDetailRequest: Encodable {
    let id: Int
}

/// 创建文章请求
struct ArticleCreateRequest: Encodable {
    let cate_id: Int
    let type: Int
    let title: String
    let url: String?
    let content: String
    let status: String?        // proto: string
    let tags: [String]?
}

/// 更新文章请求
struct ArticleUpdateRequest: Encodable {
    let id: Int
    let cate_id: Int?
    let type: Int?
    let title: String?
    let url: String?
    let content: String?
    let status: String?        // proto: string
    let tags: [String]?
}
