import Foundation

// CommentService 评论服务，包含公开和管理端的评论接口
class CommentService {
    static let shared = CommentService()

    private init() {}

    // 获取某篇文章的全部评论（公开接口）
    func list(postId: Int) async throws -> CommentListResponse {
        let request = CommentListRequest(post_id: postId)
        return try await APIClient.shared.request(path: Config.commentListPath, body: request)
    }

    // 创建评论（公开接口）
    func create(postId: Int, name: String, content: String, email: String? = nil, website: String? = nil, pid: Int = 0, replyName: String? = nil) async throws -> CommentCreateResponse {
        let request = CommentCreateRequest(
            post_id: postId,
            name: name,
            email: email,
            website: website,
            pid: pid,
            reply_name: replyName,
            content: content
        )
        return try await APIClient.shared.request(path: Config.commentCreatePath, body: request)
    }

    // 获取评论列表（管理端）
    func adminList(page: Int, keyword: String? = nil) async throws -> AdminCommentListResponse {
        let request = AdminCommentListRequest(page: page, keyword: keyword)
        return try await APIClient.shared.request(path: Config.adminCommentListPath, body: request, auth: true)
    }

    // 删除评论（支持批量）
    func delete(ids: [Int]) async throws {
        let request = CommentDeleteRequest(ids: ids)
        // 返回 google.protobuf.Empty，只需要检查不抛错即可
        let _: EmptyResponse = try await APIClient.shared.request(path: Config.adminCommentDeletePath, body: request, auth: true)
    }
}
