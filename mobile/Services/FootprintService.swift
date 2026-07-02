import Foundation

// FootprintService 足迹服务，包含公开和管理端的足迹接口
class FootprintService {
    static let shared = FootprintService()

    private init() {}

    // 获取足迹列表（管理端）
    func list(page: Int) async throws -> FootprintListResponse {
        let request = FootprintListRequest(page: page)
        return try await APIClient.shared.request(path: "/blog/admin/footprint/list", body: request, auth: true)
    }

    // 创建足迹
    func create(params: FootprintCreateRequest) async throws -> IDResponse {
        return try await APIClient.shared.request(path: "/blog/admin/footprint/create", body: params, auth: true)
    }

    // 更新足迹
    func update(params: FootprintUpdateRequest) async throws -> IDResponse {
        return try await APIClient.shared.request(path: "/blog/admin/footprint/update", body: params, auth: true)
    }

    // 删除足迹
    func delete(id: Int) async throws {
        let request = FootprintDeleteRequest(id: id)
        // 返回 google.protobuf.Empty，只需要检查不抛错即可
        let _: EmptyResponse = try await APIClient.shared.request(path: "/blog/admin/footprint/delete", body: request, auth: true)
    }

    // 获取所有足迹数据（公开接口，用于地图展示）
    func all() async throws -> GetFootprintsResponse {
        return try await APIClient.shared.request(path: "/blog/travel/footprints", body: EmptyRequest())
    }
}
