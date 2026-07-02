// 评论相关模型

import Foundation

/// 评论（api v1）
struct Comment: Identifiable, Decodable {
    let id: Int
    let pid: Int
    let name: String
    let avatar: String
    let website: String
    let content: String
    let reply_name: String
    let created_at: String      // proto: string（RFC3339 格式）
}

/// 评论列表请求
struct CommentListRequest: Encodable {
    let post_id: Int
}

/// 评论列表响应
struct CommentListResponse: Decodable {
    let list: [Comment]
}

/// 创建评论请求
struct CommentCreateRequest: Encodable {
    let post_id: Int
    let name: String
    let email: String?
    let website: String?
    let pid: Int
    let reply_name: String?
    let content: String
}

/// 创建评论响应
struct CommentCreateResponse: Decodable {
    let id: Int
}

/// 管理端评论
struct AdminComment: Identifiable, Decodable {
    let id: Int
    let post_id: Int
    let pid: Int
    let name: String
    let email: String
    let website: String
    let content: String
    let reply_name: String
    let ip: String
    let created_at: String      // proto: string（RFC3339 格式）
    let post_title: String
    let post_url: String
}

/// 管理端评论列表响应
struct AdminCommentListResponse: Decodable {
    let list: [AdminComment]
    let total: Int
}
