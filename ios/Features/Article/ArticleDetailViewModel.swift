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
        if let date = Date.parseAPIString(dateString) {
            return date.relativeString()
        }
        return dateString
    }

    /// 格式化文章更新时间
    /// - Parameter dateString: 后端时间字符串
    /// - Returns: yyyy-MM-dd HH:mm 格式的更新时间
    func updateTime(for dateString: String) -> String {
        Date.formattedAPIString(dateString, format: "yyyy-MM-dd HH:mm")
    }
}
