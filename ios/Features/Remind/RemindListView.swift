import SwiftUI

/// 提醒列表视图
///
/// 布局结构：透明导航栏 + List insetGrouped 分组 + safeAreaInset 固定 Header
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
                    // 分组列表（保持原生 insetGrouped 风格）
                    List {
                        // 分组提醒
                        ForEach(viewModel.groupedReminds, id: \.status) { group in
                            Section {
                                ForEach(group.items) { remind in
                                    remindRow(remind)
                                        // 左滑操作：待确认事项显示「标记完成」，其他显示「删除」
                                        .swipeActions(edge: .trailing, allowsFullSwipe: false) {
                                            if remind.status == RemindStatus.pending.rawValue {
                                                Button {
                                                    Task { await viewModel.markAsDoneDirect(remind: remind) }
                                                } label: {
                                                    Label("完成", systemImage: "checkmark.circle")
                                                }
                                                .tint(.green)

                                                // 不使用 role: .destructive：系统会预判删除并触发 cell 收缩动画，
                                                // 导致左滑删除确认时行闪烁。改用 .tint(.red) 保持红色外观但避开系统动画。
                                                Button {
                                                    viewModel.confirmDelete(remind: remind)
                                                } label: {
                                                    Label("删除", systemImage: "trash")
                                                }
                                                .tint(.red)
                                            } else {
                                                Button {
                                                    viewModel.confirmDelete(remind: remind)
                                                } label: {
                                                    Label("删除", systemImage: "trash")
                                                }
                                                .tint(.red)
                                            }
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
                            .listRowBackground(Color.clear)
                        } else if viewModel.hasMore {
                            Color.clear
                                .frame(height: 1)
                                .onAppear {
                                    Task { await viewModel.loadMore() }
                                }
                                .listRowSeparator(.hidden)
                                .listRowBackground(Color.clear)
                        }
                    }
                    .listStyle(.insetGrouped)
                    // 隐藏 List 默认背景，露出底层装饰背景图
                    .scrollContentBackground(.hidden)
                    // 消除 List 顶部默认 contentInset，与其他 ScrollView 页面顶部对齐
                    .contentMargins(.top, 0, for: .scrollContent)
                    // 这样 Header 内部的 padding(.horizontal, 16) 直接生效，与博文/心情页逐像素对齐，
                    // 注意：Header 不随列表滚动，始终固定在顶部（提醒类页面体验更佳）。
                    .safeAreaInset(edge: .top, spacing: 0) {
                        ListPageHeader(title: "提醒", bottomPadding: 0)
                    }
                    .refreshable {
                        await viewModel.refresh()
                    }
                }
            }
            // 背景图放在 .background 中，铺满屏幕
            .background(PageBackground(imageName: "remind_bg").ignoresSafeArea())
            // 导航栏透明：让背景图自然透出，但保留系统 Toolbar 按钮的原生玻璃质感
            .toolbarBackground(.hidden, for: .navigationBar)
            .toolbarBackground(.visible, for: .navigationBar)
            .toolbar {
                // 右上角 + 按钮：使用原生 ToolbarItem，获得系统玻璃质感/高亮/动画
                ToolbarItem(placement: .topBarTrailing) {
                    Button {
                        showEditor = true
                    } label: {
                        Image(systemName: "plus")
                    }
                }
            }
            .navigationTitle("")
            .navigationBarTitleDisplayMode(.inline)
            .sheet(isPresented: $showEditor, onDismiss: {
                Task { await viewModel.refresh() }
            }) {
                NavigationStack {
                    RemindEditorView()
                }
            }
            // 使用 .alert(item:)：由 remindToDelete 是否为 nil 自动驱动弹窗显隐，
            // 相比 .alert(isPresented:) 减少一次 View Diff，配合 L1（去 destructive role）缓解左滑闪烁
            .alert("删除提醒", isPresented: Binding(
                get: { viewModel.remindToDelete != nil },
                set: { if !$0 { viewModel.remindToDelete = nil } }
            )) {
                Button("取消", role: .cancel) {
                    viewModel.remindToDelete = nil
                }
                Button("删除", role: .destructive) {
                    Task { await viewModel.deleteRemind() }
                }
            } message: {
                Text("确定要删除这条提醒吗？此操作不可撤销。")
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

    /// 单条提醒行（保持原来的 List 行样式）
    private func remindRow(_ remind: Remind) -> some View {
        HStack(spacing: 12) {
            VStack(alignment: .leading, spacing: 4) {
                // 提醒内容
                Text(remind.content)
                    .font(.body)
                    .foregroundStyle(Color.themePrimary)
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
            // 仅待确认（PENDING）状态点击可标记完成，直接执行无需弹窗
            if remind.status == RemindStatus.pending.rawValue {
                Task { await viewModel.markAsDoneDirect(remind: remind) }
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
