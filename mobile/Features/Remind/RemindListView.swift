import SwiftUI

/// 提醒列表视图
struct RemindListView: View {

    @State private var viewModel = RemindListViewModel()

    /// 是否显示新建提醒编辑器
    @State private var showEditor = false

    var body: some View {
        NavigationStack {
            Group {
                if viewModel.isLoading && viewModel.reminds.isEmpty {
                    // 首次加载中
                    ProgressView("加载中...")
                } else if viewModel.reminds.isEmpty && !viewModel.isLoading {
                    // 空状态
                    ContentUnavailableView {
                        Label("还没有提醒", systemImage: "bell.badge")
                    } description: {
                        Text("点击右上角 + 创建你的第一条提醒")
                    }
                } else {
                    // 分组列表
                    List {
                        ForEach(viewModel.groupedReminds, id: \.status) { group in
                            Section {
                                ForEach(group.items) { remind in
                                    remindRow(remind)
                                }
                                .onDelete { indexSet in
                                    if let index = indexSet.first {
                                        let remind = group.items[index]
                                        viewModel.confirmDelete(remind: remind)
                                    }
                                }
                            } header: {
                                HStack {
                                    Text(group.title)
                                    Spacer()
                                    // 状态计数标签
                                    Text("\(group.items.count)")
                                        .font(.caption)
                                        .foregroundStyle(.secondary)
                                        .padding(.horizontal, 8)
                                        .padding(.vertical, 2)
                                        .background(Color(.systemGray5))
                                        .clipShape(Capsule())
                                }
                            }
                        }

                        // 加载更多
                        if viewModel.isLoadingMore {
                            HStack {
                                Spacer()
                                ProgressView()
                                    .padding()
                                Spacer()
                            }
                            .listRowSeparator(.hidden)
                        } else if viewModel.hasMore {
                            Color.clear
                                .frame(height: 1)
                                .onAppear {
                                    Task { await viewModel.loadMore() }
                                }
                                .listRowSeparator(.hidden)
                        }
                    }
                    .listStyle(.insetGrouped)
                    .refreshable {
                        await viewModel.refresh()
                    }
                }
            }
            .navigationTitle("提醒")
            .toolbar {
                ToolbarItem(placement: .primaryAction) {
                    Button {
                        showEditor = true
                    } label: {
                        Image(systemName: "plus")
                    }
                }
            }
            .sheet(isPresented: $showEditor, onDismiss: {
                // 编辑器关闭后刷新列表（涵盖新建、编辑、删除场景）
                Task { await viewModel.refresh() }
            }) {
                NavigationStack {
                    RemindEditorView()
                }
            }
            .alert("删除提醒", isPresented: $viewModel.showDeleteConfirmation) {
                Button("取消", role: .cancel) {
                    viewModel.remindToDelete = nil
                }
                Button("删除", role: .destructive) {
                    Task { await viewModel.deleteRemind() }
                }
            } message: {
                Text("确定要删除这条提醒吗？此操作不可撤销。")
            }
            .confirmationDialog("标记完成", isPresented: $viewModel.showCompleteAction, titleVisibility: .visible) {
                Button("标记为已完成") {
                    Task { await viewModel.markAsDone() }
                }
                Button("取消", role: .cancel) {
                    viewModel.remindToComplete = nil
                }
            } message: {
                if let remind = viewModel.remindToComplete {
                    Text("确定将「\(remind.content)」标记为已完成吗？")
                }
            }
            .alert("错误", isPresented: $viewModel.showError) {
                Button("确定", role: .cancel) {}
            } message: {
                Text(viewModel.errorMessage ?? "未知错误")
            }
            .task {
                await viewModel.load()
            }
        }
    }

    // MARK: - 子视图

    /// 单条提醒行
    private func remindRow(_ remind: Remind) -> some View {
        HStack(spacing: 12) {
            VStack(alignment: .leading, spacing: 4) {
                // 提醒内容
                Text(remind.content)
                    .font(.body)
                    .lineLimit(2)

                // 下次提醒时间（可读格式）
                if !remind.next_time.isEmpty {
                    HStack(spacing: 4) {
                        Image(systemName: "clock")
                            .font(.caption2)
                        Text("下次提醒：" + formatNextTime(remind.next_time))
                    }
                    .font(.caption)
                    .foregroundStyle(.secondary)
                }
            }

            Spacer()

            // 状态标签
            statusBadge(remind.status)
        }
        .contentShape(Rectangle())
        .onTapGesture {
            let status = remind.status
            // 仅 ACTIVE 和 PENDING 状态可以标记完成
            if status == RemindStatus.active.rawValue || status == RemindStatus.pending.rawValue {
                viewModel.requestComplete(remind: remind)
            }
        }
    }

    /// 状态标签视图
    private func statusBadge(_ status: String) -> some View {
        Text(statusLabel(status))
            .font(.caption2)
            .fontWeight(.medium)
            .padding(.horizontal, 8)
            .padding(.vertical, 4)
            .background(statusColor(status).opacity(0.15))
            .foregroundStyle(statusColor(status))
            .clipShape(Capsule())
    }

    /// 状态显示文字
    private func statusLabel(_ status: String) -> String {
        switch status {
        case RemindStatus.active.rawValue: return "进行中"
        case RemindStatus.pending.rawValue: return "待确认"
        case RemindStatus.done.rawValue: return "已完成"
        default: return status
        }
    }

    /// 状态对应颜色
    private func statusColor(_ status: String) -> Color {
        switch status {
        case RemindStatus.active.rawValue: return .blue
        case RemindStatus.pending.rawValue: return .orange
        case RemindStatus.done.rawValue: return .green
        default: return .gray
        }
    }

    /// 将 RFC3339 时间字符串格式化为可读形式（如 "7月3日 09:00"）
    private func formatNextTime(_ dateString: String) -> String {
        let formatter = ISO8601DateFormatter()
        formatter.formatOptions = [.withInternetDateTime, .withFractionalSeconds]
        var date = formatter.date(from: dateString)
        if date == nil {
            // 降级：尝试不带毫秒的格式
            formatter.formatOptions = [.withInternetDateTime]
            date = formatter.date(from: dateString)
        }
        guard let date else { return dateString }
        return date.toString(format: "M月d日 HH:mm")
    }
}
