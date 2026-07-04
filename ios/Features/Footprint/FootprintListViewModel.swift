import Foundation

/// 足迹列表视图模型
@Observable
class FootprintListViewModel: APIErrorPresentable {

    // MARK: - 状态

    /// 足迹列表
    var footprints: [Footprint] = []

    /// 是否正在加载
    var isLoading = false

    /// 是否正在刷新
    var isRefreshing = false

    /// 错误信息
    var errorMessage: String?

    /// 是否显示错误
    var showError = false

    /// 当前页码
    private var currentPage = 1

    /// 是否还有更多数据
    var hasMore = true

    // MARK: - 私有属性

    private let service = FootprintService.shared

    /// 每页数量
    private let pageSize = 20

    // MARK: - 数据加载

    /// 加载所有足迹数据（使用公开接口获取全部数据）
    func loadFootprints() async {
        guard !isLoading else { return }
        isLoading = true
        errorMessage = nil
        showError = false

        do {
            let response = try await service.all()
            footprints = response.footprints
        } catch {
            handleAPIError(error)
        }

        isLoading = false
    }

    /// 下拉刷新
    func refresh() async {
        guard !isRefreshing else { return }
        isRefreshing = true

        do {
            let response = try await service.all()
            footprints = response.footprints
        } catch {
            handleAPIError(error)
        }

        isRefreshing = false
    }

    /// 加载更多数据（分页）
    func loadMore() async {
        guard !isLoading, hasMore else { return }
        isLoading = true

        do {
            let response = try await service.list(page: currentPage + 1)
            let newList = response.list

            if newList.isEmpty {
                hasMore = false
            } else {
                currentPage += 1
                footprints.append(contentsOf: newList)
            }
        } catch {
            handleAPIError(error)
        }

        isLoading = false
    }
}
