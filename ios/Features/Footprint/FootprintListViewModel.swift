import Foundation

/// 足迹列表视图模型
@MainActor
@Observable
class FootprintListViewModel: APIErrorPresentable {

    // MARK: - 状态

    /// 足迹列表
    var footprints: [Footprint] = []

    /// 是否正在加载
    var isLoading = false

    /// 是否正在刷新
    var isRefreshing = false

    /// 是否正在加载更多
    var isLoadingMore = false

    /// 错误信息
    var errorMessage: String?

    /// 是否显示错误
    var showError = false

    /// 当前页码
    private var currentPage = 0

    /// 是否还有更多数据
    var hasMore = true

    // MARK: - 私有属性

    private let service = FootprintService.shared

    // MARK: - 数据加载

    /// 加载足迹列表（第一页）
    func loadFootprints() async {
        guard !isLoading else { return }
        isLoading = true
        currentPage = 0
        hasMore = true
        errorMessage = nil
        showError = false

        do {
            let response = try await service.list(page: 1)
            footprints = response.list
            currentPage = 1
            hasMore = footprints.count < response.total
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
            let response = try await service.list(page: 1)
            footprints = response.list
            currentPage = 1
            hasMore = footprints.count < response.total
        } catch {
            handleAPIError(error)
        }

        isRefreshing = false
    }

    /// 加载更多数据（分页）
    func loadMore() async {
        guard !isLoading, !isLoadingMore, hasMore else { return }
        isLoadingMore = true
        let nextPage = currentPage + 1

        do {
            let response = try await service.list(page: nextPage)
            let existingIds = Set(footprints.map { $0.id })
            let newItems = response.list.filter { !existingIds.contains($0.id) }
            footprints.append(contentsOf: newItems)
            currentPage = nextPage
            hasMore = !response.list.isEmpty && footprints.count < response.total
        } catch {
            handleAPIError(error)
        }

        isLoadingMore = false
    }
}
