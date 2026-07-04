import Foundation

/// 提醒列表视图模型
@Observable
class RemindListViewModel: APIErrorPresentable {

    // MARK: - 状态

    /// 提醒列表数据
    var reminds: [Remind] = []

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

    /// 待删除的提醒（非 nil 时驱动 .alert(item:) 弹出删除确认弹窗）
    var remindToDelete: Remind?

    /// 待标记完成的提醒
    var remindToComplete: Remind?

    /// 是否显示标记完成操作菜单
    var showCompleteAction = false

    // MARK: - 私有属性

    private var currentPage = 0
    private let remindService = RemindService.shared

    // MARK: - 数据分组

    /// 按状态分组后的提醒列表
    var groupedReminds: [(status: RemindStatus, title: String, items: [Remind])] {
        let active = reminds.filter { $0.status == RemindStatus.active.rawValue }
        let pending = reminds.filter { $0.status == RemindStatus.pending.rawValue }
        let done = reminds.filter { $0.status == RemindStatus.done.rawValue }

        return [
            (status: .active, title: "进行中", items: active),
            (status: .pending, title: "待确认", items: pending),
            (status: .done, title: "已完成", items: done),
        ].filter { !$0.items.isEmpty }
    }

    // MARK: - 数据加载

    /// 加载提醒列表（首次加载）
    func load() async {
        guard !isLoading else { return }
        isLoading = true
        currentPage = 0
        hasMore = true

        do {
            let response = try await remindService.list(page: 1)
            reminds = response.list
            currentPage = 1
            hasMore = reminds.count < response.total
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
            let response = try await remindService.list(page: 1)
            reminds = response.list
            currentPage = 1
            hasMore = reminds.count < response.total
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
            let response = try await remindService.list(page: nextPage)
            // 基于已加载 id 去重，避免 Tab 切换 onAppear 重复触发导致数据翻倍
            let existingIds = Set(reminds.map { $0.id })
            let newItems = response.list.filter { !existingIds.contains($0.id) }
            reminds.append(contentsOf: newItems)
            currentPage = nextPage
            // 当前页返回为空或不足一页时，判定无更多
            hasMore = !response.list.isEmpty && reminds.count < response.total
        } catch {
            handleAPIError(error)
        }

        isLoadingMore = false
    }

    // MARK: - 删除

    /// 确认删除提醒（仅设置待删除项，由 .alert(item:) 自动驱动弹窗显隐）
    func confirmDelete(remind: Remind) {
        remindToDelete = remind
    }

    /// 执行删除提醒
    func deleteRemind() async {
        guard let remind = remindToDelete else { return }

        do {
            try await remindService.delete(id: remind.id)
            reminds.removeAll { $0.id == remind.id }
        } catch {
            handleAPIError(error)
        }

        remindToDelete = nil
    }

    // MARK: - 标记完成

    /// 请求标记提醒为已完成（已废弃，保留兼容）
    func requestComplete(remind: Remind) {
        remindToComplete = remind
        showCompleteAction = true
    }

    /// 执行标记完成（已废弃，保留兼容）
    func markAsDone() async {
        guard let remind = remindToComplete else { return }

        do {
            _ = try await remindService.update(
                id: remind.id,
                cron: nil,
                content: nil,
                status: RemindStatus.done.rawValue
            )
            // 直接重新加载列表以确保数据一致
            await refresh()
        } catch {
            handleAPIError(error)
        }

        remindToComplete = nil
    }

    /// 直接标记提醒为已完成（不弹窗，直接调用接口）
    /// - Parameter remind: 要标记完成的提醒
    func markAsDoneDirect(remind: Remind) async {
        do {
            _ = try await remindService.update(
                id: remind.id,
                cron: nil,
                content: nil,
                status: RemindStatus.done.rawValue
            )
            // 直接重新加载列表以确保数据一致
            await refresh()
        } catch {
            handleAPIError(error)
        }
    }
}
