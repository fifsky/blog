import Foundation

// ArticleService 文章服务，包含公开和管理端的文章接口
class ArticleService {
    static let shared = ArticleService()

    private init() {}

    // 获取文章列表（管理端）
    func list(page: Int, type: Int? = nil, status: String? = nil, keyword: String? = nil) async throws -> ArticleListResponse {
        let request = ArticleListRequest(
            page: page,
            type: type,
            status: status,
            keyword: keyword
        )
        return try await APIClient.shared.request(path: Config.adminArticleListPath, body: request, auth: true)
    }

    // 获取文章详情（管理端）
    func detail(id: Int) async throws -> Article {
        let request = ArticleDetailRequest(id: id)
        return try await APIClient.shared.request(path: Config.adminArticleDetailPath, body: request, auth: true)
    }

    // 创建文章
    func create(cateId: Int, type: Int, title: String, content: String, tags: [String]? = nil) async throws -> IDResponse {
        let request = ArticleCreateRequest(
            cate_id: cateId,
            type: type,
            title: title,
            url: nil,
            content: content,
            status: nil,
            tags: tags
        )
        return try await APIClient.shared.request(path: Config.adminArticleCreatePath, body: request, auth: true)
    }

    // 更新文章
    func update(id: Int, cateId: Int? = nil, type: Int? = nil, title: String? = nil, content: String? = nil, tags: [String]? = nil) async throws -> IDResponse {
        let request = ArticleUpdateRequest(
            id: id,
            cate_id: cateId,
            type: type,
            title: title,
            url: nil,
            content: content,
            status: nil,
            tags: tags
        )
        return try await APIClient.shared.request(path: Config.adminArticleUpdatePath, body: request, auth: true)
    }

    // 删除文章（支持批量）
    func delete(ids: [Int]) async throws {
        let request = ArticleDeleteRequest(ids: ids)
        // 返回 google.protobuf.Empty，只需要检查不抛错即可
        let _: EmptyResponse = try await APIClient.shared.request(path: Config.adminArticleDeletePath, body: request, auth: true)
    }
}
