import SwiftUI

/// 文章列表视图模型
@MainActor
@Observable
class ArticleListViewModel: APIErrorPresentable {

    // MARK: - 状态

    /// 文章列表
    var articles: [Article] = []

    /// 是否正在加载（首帧默认 true，避免冷启动时先闪一帧空态）
    var isLoading = true

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

    /// 当前用于列表请求的关键词（空串表示未处于搜索态）
    private(set) var currentKeyword: String = ""

    /// 是否处于搜索态（当前关键词非空）
    var isSearching: Bool {
        !currentKeyword.isEmpty
    }

    // MARK: - 私有属性

    private var currentPage = 0
    /// 是否已执行过首次加载（配合 isLoading 初始 true，避免 guard !isLoading 提前返回）
    private var hasLoaded = false
    private let articleService = ArticleService.shared

    // MARK: - 数据加载

    /// 加载文章列表（首次加载）
    func load() async {
        // 用 hasLoaded 守卫，避免与 isLoading 初始值 true 冲突（guard !isLoading 会提前返回）
        guard !hasLoaded else { return }
        hasLoaded = true
        isLoading = true
        currentPage = 0
        hasMore = true

        do {
            let response = try await articleService.list(page: 1, keyword: requestKeyword)
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
            let response = try await articleService.list(page: 1, keyword: requestKeyword)
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
            let response = try await articleService.list(page: nextPage, keyword: requestKeyword)
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

    // MARK: - 搜索

    /// 应用搜索关键词并重新加载列表
    /// - Parameter keyword: 搜索关键词（自动 trim，空串等价于清除搜索）
    func applySearch(_ keyword: String) {
        let trimmed = keyword.trimmingCharacters(in: .whitespaces)
        currentKeyword = trimmed
        Task { await reload() }
    }

    /// 清除搜索态，恢复全部列表
    func clearSearch() {
        currentKeyword = ""
        Task { await reload() }
    }

    /// 重新加载第 1 页（搜索/清除搜索复用，不受首次加载守卫限制）
    private func reload() async {
        isLoading = true
        currentPage = 0
        hasMore = true

        do {
            let response = try await articleService.list(page: 1, keyword: requestKeyword)
            articles = response.list
            currentPage = 1
            hasMore = articles.count < response.total
        } catch {
            handleAPIError(error)
        }

        isLoading = false
    }

    /// 转换为 Service 层接受的 keyword 参数（空串转 nil）
    private var requestKeyword: String? {
        currentKeyword.isEmpty ? nil : currentKeyword
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
