import Foundation

// CategoryService 分类服务，提供分类相关的公开接口
class CategoryService {
    static let shared = CategoryService()

    private init() {}

    // 获取所有分类（公开接口）
    func all() async throws -> CateMenuResponse {
        return try await APIClient.shared.request(path: "/blog/cate/all", body: EmptyRequest())
    }
}
