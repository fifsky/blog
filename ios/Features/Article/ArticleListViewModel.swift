import SwiftUI

/// 文章列表视图模型
@Observable
class ArticleListViewModel: APIErrorPresentable {

    // MARK: - 状态

    /// 文章列表
    var articles: [Article] = []

    /// 是否正在加载（首次）
    var isLoading = false

    /// 是否正在下拉刷新
    var isRefreshing = false

    /// 是否正在加载更多
    var isLoadingMore = false

    /// 是否还有更多数据
    var hasMore = true

    /// 错误信息
    var errorMessage: String?

    /// 是否显示错误弹窗
    var showError = false

    // MARK: - 私有属性

    private var currentPage = 0
    private let articleService = ArticleService.shared

    // MARK: - 数据加载

    /// 加载文章列表（首次加载）
    func load() async {
        guard !isLoading else { return }
        isLoading = true
        currentPage = 0
        hasMore = true

        do {
            let response = try await articleService.list(page: 1)
            articles = response.list
            currentPage = 1
            hasMore = articles.count < response.total
        } catch {
            handleAPIError(error)
        }

        isLoading = false
    }

    /// 下拉刷新
    func refresh() async {
        guard !isRefreshing else { return }
        isRefreshing = true
        currentPage = 0
        hasMore = true

        do {
            let response = try await articleService.list(page: 1)
            articles = response.list
            currentPage = 1
            hasMore = articles.count < response.total
        } catch {
            handleAPIError(error)
        }

        isRefreshing = false
    }

    /// 加载更多（分页）
    func loadMore() async {
        guard !isLoadingMore, hasMore else { return }
        isLoadingMore = true
        let nextPage = currentPage + 1

        do {
            let response = try await articleService.list(page: nextPage)
            // 基于已加载 id 去重，避免 Tab 切换 onAppear 重复触发导致数据翻倍
            let existingIds = Set(articles.map { $0.id })
            let newItems = response.list.filter { !existingIds.contains($0.id) }
            articles.append(contentsOf: newItems)
            currentPage = nextPage
            hasMore = !response.list.isEmpty && articles.count < response.total
        } catch {
            handleAPIError(error)
        }

        isLoadingMore = false
    }

    // MARK: - 辅助方法

    /// 解析文章创建时间字符串为相对时间
    /// - Parameter dateString: RFC3339 格式的日期字符串
    /// - Returns: 相对时间描述（如 "3分钟前"）
    func relativeTime(for dateString: String) -> String {
        let formatter = ISO8601DateFormatter()
        formatter.formatOptions = [.withInternetDateTime, .withFractionalSeconds]
        if let date = formatter.date(from: dateString) {
            return date.relativeString()
        }
        // 尝试不带小数秒的格式
        formatter.formatOptions = [.withInternetDateTime]
        if let date = formatter.date(from: dateString) {
            return date.relativeString()
        }
        return dateString
    }
}
