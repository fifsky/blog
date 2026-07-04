import Foundation

/// 评论列表视图模型
@Observable
class CommentListViewModel: APIErrorPresentable {

    // MARK: - 状态

    /// 评论列表
    var comments: [Comment] = []

    /// 是否正在加载
    var isLoading = false

    /// 错误信息
    var errorMessage: String?

    /// 是否显示错误弹窗
    var showError = false

    // MARK: - 私有属性

    private let postId: Int
    private let commentService = CommentService.shared

    // MARK: - 初始化

    init(postId: Int) {
        self.postId = postId
    }

    // MARK: - 数据加载

    /// 加载评论列表
    func load() async {
        guard !isLoading else { return }
        isLoading = true

        do {
            let response = try await commentService.list(postId: postId)
            comments = response.list
        } catch {
            handleAPIError(error)
        }

        isLoading = false
    }

    // MARK: - 数据整理

    /// 主评论列表（pid == 0）
    var mainComments: [Comment] {
        comments.filter { $0.pid == 0 }
    }

    /// 获取某条主评论的回复列表
    /// - Parameter parentId: 主评论 ID
    /// - Returns: 该主评论下的所有回复
    func replies(for parentId: Int) -> [Comment] {
        comments.filter { $0.pid == parentId }
    }

    // MARK: - 辅助方法

    /// 解析评论创建时间字符串为相对时间
    func relativeTime(for dateString: String) -> String {
        let formatter = ISO8601DateFormatter()
        formatter.formatOptions = [.withInternetDateTime, .withFractionalSeconds]
        if let date = formatter.date(from: dateString) {
            return date.relativeString()
        }
        formatter.formatOptions = [.withInternetDateTime]
        if let date = formatter.date(from: dateString) {
            return date.relativeString()
        }
        return dateString
    }
}
