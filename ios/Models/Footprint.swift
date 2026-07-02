// 足迹相关模型

import Foundation

/// 足迹照片
struct FootprintPhoto: Decodable {
    let src: String?
    let thumbnail: String?
}

/// 足迹
/// protojson EmitUnpopulated 会输出零值空串，字段均为可选以增强容错
struct Footprint: Identifiable, Decodable {
    let id: Int
    let name: String?
    let description: String?
    let longitude: String?       // proto: string
    let latitude: String?        // proto: string
    let date: String?
    let marker_color: String?
    let categories: [String]?
    let url: String?
    let url_label: String?
    let photos: [FootprintPhoto]?
}

/// 足迹列表请求
struct FootprintListRequest: Encodable {
    let page: Int
}

/// 足迹列表响应
struct FootprintListResponse: Decodable {
    let list: [Footprint]
    let total: Int
}

/// 创建足迹请求
struct FootprintCreateRequest: Encodable {
    let name: String
    let description: String?
    let longitude: String       // proto: string
    let latitude: String        // proto: string
    let date: String?
    let marker_color: String?
    let categories: [String]?
    let url: String?
    let url_label: String?
    let photo_urls: [String]?
}

/// 更新足迹请求
struct FootprintUpdateRequest: Encodable {
    let id: Int
    let name: String?
    let description: String?
    let longitude: String?      // proto: string
    let latitude: String?       // proto: string
    let date: String?
    let marker_color: String?
    let categories: [String]?
    let url: String?
    let url_label: String?
    let photo_urls: [String]?
}
