import Foundation

/// 评论根节点：主评论和其下平铺回复
struct CommentRoot: Identifiable {
    /// 主评论
    let root: Comment

    /// 主评论下的全部回复
    let replies: [Comment]

    /// 根节点 ID
    var id: Int { root.id }
}

/// 评论列表视图模型
@MainActor
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

    /// 评论根节点列表
    var roots: [CommentRoot] {
        Self.buildRoots(from: comments)
    }

    /// 评论总数
    var totalCount: Int {
        comments.count
    }

    /// 构建两层评论结构，回复的回复也平铺在顶层主评论下
    /// - Parameter comments: 后端返回的评论列表
    /// - Returns: 主评论和回复组成的根节点列表
    static func buildRoots(from comments: [Comment]) -> [CommentRoot] {
        let roots = comments.filter { $0.pid == 0 }
        let replies = comments.filter { $0.pid != 0 }

        return roots
            .sorted { compareDateString($0.created_at, $1.created_at, ascending: false) }
            .map { root in
                let children = replies
                    .filter { $0.pid == root.id }
                    .sorted { compareDateString($0.created_at, $1.created_at, ascending: true) }
                return CommentRoot(root: root, replies: children)
            }
    }

    // MARK: - 辅助方法

    /// 解析评论创建时间字符串为相对时间
    func relativeTime(for dateString: String) -> String {
        if let date = Date.parseAPIString(dateString) {
            return date.relativeString()
        }
        return dateString
    }

    /// 比较两个 API 时间字符串
    /// - Parameters:
    ///   - lhs: 左侧时间字符串
    ///   - rhs: 右侧时间字符串
    ///   - ascending: 是否升序
    /// - Returns: lhs 是否应该排在 rhs 前面
    private static func compareDateString(_ lhs: String, _ rhs: String, ascending: Bool) -> Bool {
        let leftDate = Date.parseAPIString(lhs)
        let rightDate = Date.parseAPIString(rhs)

        switch (leftDate, rightDate) {
        case let (left?, right?):
            return ascending ? left < right : left > right
        default:
            return ascending ? lhs < rhs : lhs > rhs
        }
    }
}
