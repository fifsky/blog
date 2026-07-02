import SwiftUI

/// 心情列表视图
struct MoodListView: View {

    @State private var viewModel = MoodListViewModel()

    /// 是否显示心情编辑器（新建/编辑共用）
    @State private var showEditor = false

    /// 当前编辑的心情（nil 表示新建）
    @State private var editingMood: Mood?

    /// 长按操作的目标心情
    @State private var actionMood: Mood?

    var body: some View {
        NavigationStack {
            Group {
                if viewModel.isLoading && viewModel.moods.isEmpty {
                    // 首次加载中
                    ProgressView("加载中...")
                } else if viewModel.moods.isEmpty && !viewModel.isLoading {
                    // 空状态
                    ContentUnavailableView {
                        Label("还没有心情", systemImage: "heart.text.square")
                    } description: {
                        Text("点击右上角 + 记录你的第一条心情")
                    }
                } else {
                    // 心情卡片列表
                    ScrollView {
                        LazyVStack(spacing: 12) {
                            ForEach(viewModel.moods) { mood in
                                moodCard(mood)
                            }

                            // 加载更多指示器
                            if viewModel.isLoadingMore {
                                HStack {
                                    Spacer()
                                    ProgressView()
                                        .padding()
                                    Spacer()
                                }
                            } else if viewModel.hasMore {
                                // 触底加载更多
                                Color.clear
                                    .frame(height: 1)
                                    .onAppear {
                                        Task { await viewModel.loadMore() }
                                    }
                            }
                        }
                        .padding(.horizontal, 16)
                        .padding(.top, 12)
                    }
                    .refreshable {
                        await viewModel.refresh()
                    }
                }
            }
            .navigationTitle("心情")
            .toolbar {
                ToolbarItem(placement: .primaryAction) {
                    Button {
                        editingMood = nil
                        showEditor = true
                    } label: {
                        Image(systemName: "plus")
                    }
                }
            }
            // 新建/编辑共用 sheet
            .sheet(isPresented: $showEditor, onDismiss: {
                editingMood = nil
                // 编辑器关闭后刷新列表（涵盖新建、编辑场景）
                Task { await viewModel.refresh() }
            }) {
                NavigationStack {
                    MoodEditorView(mood: editingMood)
                }
            }
            // 长按弹出的操作菜单
            .confirmationDialog("心情操作",
                                isPresented: Binding(
                                    get: { actionMood != nil },
                                    set: { if !$0 { actionMood = nil } }
                                ),
                                titleVisibility: .visible) {
                Button("编辑") {
                    if let mood = actionMood {
                        editingMood = mood
                        actionMood = nil
                        showEditor = true
                    }
                }
                Button("删除", role: .destructive) {
                    if let mood = actionMood {
                        viewModel.confirmDelete(mood: mood)
                        actionMood = nil
                    }
                }
                Button("取消", role: .cancel) {
                    actionMood = nil
                }
            }
            .alert("删除心情", isPresented: $viewModel.showDeleteConfirmation) {
                Button("取消", role: .cancel) {
                    viewModel.moodToDelete = nil
                }
                Button("删除", role: .destructive) {
                    Task { await viewModel.deleteMood() }
                }
            } message: {
                Text("确定要删除这条心情吗？此操作不可撤销。")
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

    /// 单条心情卡片
    private func moodCard(_ mood: Mood) -> some View {
        VStack(alignment: .leading, spacing: 10) {
            // 心情内容
            Text(mood.content)
                .font(.body)
                .frame(maxWidth: .infinity, alignment: .leading)
                .fixedSize(horizontal: false, vertical: true)

            // 相对时间
            HStack(spacing: 4) {
                Image(systemName: "clock")
                    .font(.caption2)
                Text(parseRelativeTime(mood.created_at))
            }
            .font(.caption)
            .foregroundStyle(.secondary)
        }
        .padding(16)
        .frame(maxWidth: .infinity, alignment: .leading)
        .background(Color(.secondarySystemBackground))
        .clipShape(RoundedRectangle(cornerRadius: 14))
        .contentShape(RoundedRectangle(cornerRadius: 14))
        // 长按弹出操作菜单
        .onLongPressGesture {
            actionMood = mood
        }
    }

    /// 解析 RFC3339 时间字符串并返回相对时间描述
    private func parseRelativeTime(_ dateString: String) -> String {
        let formatter = ISO8601DateFormatter()
        formatter.formatOptions = [.withInternetDateTime, .withFractionalSeconds]
        if let date = formatter.date(from: dateString) {
            return date.relativeString()
        }
        // 降级：尝试不带毫秒的格式
        formatter.formatOptions = [.withInternetDateTime]
        if let date = formatter.date(from: dateString) {
            return date.relativeString()
        }
        return dateString
    }
}
