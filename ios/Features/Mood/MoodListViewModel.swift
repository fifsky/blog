import Foundation

/// 心情列表视图模型
@Observable
class MoodListViewModel: APIErrorPresentable {

    // MARK: - 状态

    /// 心情列表数据
    var moods: [Mood] = []

    /// 是否正在加载
    var isLoading = false

    /// 是否正在刷新
    var isRefreshing = false

    /// 是否正在加载更多
    var isLoadingMore = false

    /// 是否还有更多数据
    var hasMore = true

    /// 错误信息
    var errorMessage: String?

    /// 是否显示错误弹窗
    var showError = false

    /// 删除确认弹窗
    var showDeleteConfirmation = false

    /// 待删除的心情
    var moodToDelete: Mood?

    // MARK: - 私有属性

    private var currentPage = 0
    private let moodService = MoodService.shared

    // MARK: - 数据加载

    /// 加载心情列表（首次加载）
    func load() async {
        guard !isLoading else { return }
        isLoading = true
        currentPage = 0
        hasMore = true

        do {
            let response = try await moodService.list(page: 1)
            moods = response.list
            currentPage = 1
            hasMore = moods.count < response.total
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
            let response = try await moodService.list(page: 1)
            moods = response.list
            currentPage = 1
            hasMore = moods.count < response.total
        } catch {
            handleAPIError(error)
        }

        isRefreshing = false
    }

    /// 加载下一页
    func loadMore() async {
        guard !isLoadingMore, hasMore else { return }
        isLoadingMore = true

        let nextPage = currentPage + 1

        do {
            let response = try await moodService.list(page: nextPage)
            moods.append(contentsOf: response.list)
            currentPage = nextPage
            hasMore = moods.count < response.total
        } catch {
            handleAPIError(error)
        }

        isLoadingMore = false
    }

    // MARK: - 删除

    /// 确认删除心情
    func confirmDelete(mood: Mood) {
        moodToDelete = mood
        showDeleteConfirmation = true
    }

    /// 执行删除心情
    func deleteMood() async {
        guard let mood = moodToDelete else { return }

        do {
            try await moodService.delete(ids: [mood.id])
            moods.removeAll { $0.id == mood.id }
        } catch {
            handleAPIError(error)
        }

        moodToDelete = nil
    }
}
