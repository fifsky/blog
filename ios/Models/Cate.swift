// 分类菜单相关模型

import Foundation

/// 分类菜单项
struct CateMenuItem: Decodable {
    let url: String
    let content: String
}

/// 分类菜单响应
struct CateMenuResponse: Decodable {
    let list: [CateMenuItem]
}
