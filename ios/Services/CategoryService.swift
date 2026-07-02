import Foundation

// CategoryService 分类服务，提供分类相关的公开与管理端接口
class CategoryService {
    static let shared = CategoryService()

    private init() {}

    // 获取所有分类（公开接口，仅含 domain URL，无数字 ID）
    func all() async throws -> CateMenuResponse {
        return try await APIClient.shared.request(path: "/blog/cate/all", body: EmptyRequest())
    }

    // 获取分类列表（管理端接口，包含数字 ID，用于文章编辑器选择分类）
    func list() async throws -> CateListResponse {
        return try await APIClient.shared.request(path: "/blog/admin/cate/list", body: EmptyRequest(), auth: true)
    }
}
