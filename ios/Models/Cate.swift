// 分类菜单相关模型

import Foundation

/// 分类菜单项（公开接口 /blog/cate/all，URL 中仅含 domain，无数字 ID）
struct CateMenuItem: Decodable {
    let url: String
    let content: String
}

/// 分类菜单响应
struct CateMenuResponse: Decodable {
    let list: [CateMenuItem]
}

/// 分类列表项（管理端接口 /blog/admin/cate/list，包含数字 ID）
struct CateItem: Decodable, Identifiable {
    let id: Int
    let name: String
    let desc: String?
    let domain: String?
    let num: Int?
}

/// 分类列表响应（管理端）
struct CateListResponse: Decodable {
    let list: [CateItem]
    let total: Int
}
