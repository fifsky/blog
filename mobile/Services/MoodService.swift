import Foundation

// MoodService 心情服务，包含公开和管理端的心情接口
class MoodService {
    static let shared = MoodService()

    private init() {}

    // 获取心情列表（公开接口）
    func list(page: Int) async throws -> MoodListResponse {
        let request = MoodListRequest(page: page)
        return try await APIClient.shared.request(path: "/blog/mood/list", body: request)
    }

    // 创建心情
    func create(content: String) async throws -> IDResponse {
        let request = MoodCreateRequest(content: content)
        return try await APIClient.shared.request(path: "/blog/admin/mood/create", body: request, auth: true)
    }

    // 更新心情
    func update(id: Int, content: String?) async throws -> IDResponse {
        let request = MoodUpdateRequest(id: id, content: content)
        return try await APIClient.shared.request(path: "/blog/admin/mood/update", body: request, auth: true)
    }

    // 删除心情（支持批量）
    func delete(ids: [Int]) async throws {
        let request = MoodDeleteRequest(ids: ids)
        // 返回 google.protobuf.Empty，只需要检查不抛错即可
        let _: EmptyResponse = try await APIClient.shared.request(path: "/blog/admin/mood/delete", body: request, auth: true)
    }
}
