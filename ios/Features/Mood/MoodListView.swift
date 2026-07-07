import SwiftUI

/// 心情列表视图
///
/// 布局结构：原生大标题导航栏 + ScrollView(卡片) + 背景图
struct MoodListView: View {

    @State private var viewModel = MoodListViewModel()

    /// 是否显示心情编辑器（新建/编辑共用）
    @State private var showEditor = false

    /// 当前编辑的心情（nil 表示新建）
    @State private var editingMood: Mood?

    var body: some View {
        // 主滚动内容：卡片，随原生大标题一起滚动
        ScrollView {
            VStack(spacing: 16) {
                contentList
            }
            .padding(.bottom, 16)
        }
        .refreshable {
            await viewModel.refresh()
        }
        // 背景图放在 .background 中，铺满屏幕
        .background(PageBackground(imageName: "moon_bg").ignoresSafeArea())
        .toolbar {
            // 右上角 + 按钮：使用原生 ToolbarItem，获得系统玻璃质感/高亮/动画
            ToolbarItem(placement: .topBarTrailing) {
                Button {
                    editingMood = nil
                    showEditor = true
                } label: {
                    Image(systemName: "plus")
                }
            }
        }
        .navigationTitle("心情")
        .navigationBarTitleDisplayMode(.large)
        // 新建/编辑共用 sheet
        .sheet(isPresented: $showEditor, onDismiss: {
            editingMood = nil
            Task { await viewModel.refresh() }
        }) {
            NavigationStack {
                MoodEditorView(mood: editingMood)
            }
        }
        .alert("删除心情", isPresented: $viewModel.showDeleteConfirmation) {
            Button("取消", role: .cancel) {
                viewModel.moodToDelete = nil
            }
            Button("删除", role: .destructive) {
                guard let mood = viewModel.moodToDelete else { return }
                Task { await viewModel.deleteMood(mood: mood) }
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

    // MARK: - 子视图

    /// 列表主体（区分加载/空/有数据三种状态）
    @ViewBuilder
    private var contentList: some View {
        if viewModel.isLoading && viewModel.moods.isEmpty {
            LoadingView()
        } else if viewModel.moods.isEmpty && !viewModel.isLoading {
            ContentUnavailableView {
                Label("还没有心情", systemImage: "heart.text.square")
            } description: {
                Text("点击右上角 + 记录你的第一条心情")
            }
        } else {
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
            // 卡片横向 padding：与 Header（自管 16）逐像素对齐
            .padding(.horizontal, 16)
        }
    }

    /// 单条心情卡片
    private func moodCard(_ mood: Mood) -> some View {
        VStack(alignment: .leading, spacing: 10) {
            // 心情内容
            Text(mood.content)
                .font(.body)
                .foregroundStyle(Color.themePrimary)
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
        .background(Color(.systemBackground).opacity(0.9))
        .clipShape(RoundedRectangle(cornerRadius: 22))
        .contentShape(RoundedRectangle(cornerRadius: 22))
        // 长按弹出系统 Context Menu（原生毛玻璃/Haptic/动画，跟随卡片）
        .contextMenu {
            Button {
                editingMood = mood
                showEditor = true
            } label: {
                Label("编辑", systemImage: "square.and.pencil")
            }
            Button(role: .destructive) {
                viewModel.confirmDelete(mood: mood)
            } label: {
                Label("删除", systemImage: "trash")
            }
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
