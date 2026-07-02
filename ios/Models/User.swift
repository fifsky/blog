// 用户相关模型

import Foundation

/// 登录请求（API v1 公开评论接口用户字段）
struct UserItem: Decodable {
    let id: Int
    let name: String
    let nick_name: String
    let email: String
    let status: Int
    let type: Int
    let created_at: String     // proto: string
    let updated_at: String      // proto: string
}
