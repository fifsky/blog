import Foundation

/// 文章详情视图模型
@MainActor
@Observable
class ArticleDetailViewModel: APIErrorPresentable {

    // MARK: - 状态

    /// 文章详情
    var article: Article?

    /// 是否正在加载
    var isLoading = false

    /// 错误信息
    var errorMessage: String?

    /// 是否显示错误弹窗
    var showError = false

    // MARK: - 私有属性

    private let articleId: Int
    private let articleService = ArticleService.shared

    // MARK: - 初始化

    init(articleId: Int) {
        self.articleId = articleId
    }

    // MARK: - 数据加载

    /// 加载文章详情
    func load() async {
        guard !isLoading else { return }
        isLoading = true

        do {
            article = try await articleService.detail(id: articleId)
        } catch {
            handleAPIError(error)
        }

        isLoading = false
    }

    // MARK: - 辅助方法

    /// 解析文章创建时间字符串为相对时间
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
