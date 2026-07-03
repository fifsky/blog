import SwiftUI

/// 文章列表视图
///
/// 布局结构：`ZStack(背景图 + ScrollView(Header + Cards) + 浮动三点菜单)`
/// Header 是 ScrollView 的第一个元素，随页面一起滚动（同 Apple Notes/Journal）。
struct ArticleListView: View {

    @State private var viewModel = ArticleListViewModel()

    /// 导航到文章详情
    @State private var selectedArticle: Article?

    /// 导航到文章编辑器
    @State private var showEditor = false

    /// 退出登录确认弹窗
    @State private var showLogoutConfirmation = false

    var body: some View {
        NavigationStack {
            // 主滚动内容：Header + 卡片，作为一个连续页面滚动
            ScrollView {
                VStack(spacing: 16) {
                    // Header 是 ScrollView 第一项，会随页面一起滚动
                    ListPageHeader(title: "博文")

                    contentList
                }
                .padding(.horizontal, 16)
                .padding(.bottom, 16)
            }
            .refreshable {
                await viewModel.refresh()
            }
            // 背景图放在 .background 中，铺满屏幕
            .background(PageBackground(imageName: "article_bg").ignoresSafeArea())
            // 导航栏透明：让背景图自然透出，但保留系统 Toolbar 按钮的原生玻璃质感
            .toolbarBackground(.hidden, for: .navigationBar)
            .toolbarBackground(.visible, for: .navigationBar)
            .toolbar {
                // 右上角三点菜单：使用原生 ToolbarItem，获得系统玻璃质感/高亮/动画
                ToolbarItem(placement: .topBarTrailing) {
                    Menu {
                        Button {
                            showEditor = true
                        } label: {
                            Label("新增博文", systemImage: "square.and.pencil")
                        }

                        Divider()

                        Button(role: .destructive) {
                            showLogoutConfirmation = true
                        } label: {
                            Label("退出登录", systemImage: "rectangle.portrait.and.arrow.right")
                        }
                    } label: {
                        Image(systemName: "ellipsis")
                    }
                }
            }
            .navigationTitle("")
            .navigationBarTitleDisplayMode(.inline)
            .navigationDestination(item: $selectedArticle) { article in
                ArticleDetailView(articleId: article.id)
            }
            .navigationDestination(isPresented: $showEditor) {
                ArticleEditorView()
            }
            .alert("错误", isPresented: $viewModel.showError) {
                Button("确定", role: .cancel) {}
            } message: {
                Text(viewModel.errorMessage ?? "未知错误")
            }
            .alert("退出登录", isPresented: $showLogoutConfirmation) {
                Button("退出登录", role: .destructive) {
                    AuthManager.shared.logout()
                }
                Button("取消", role: .cancel) {}
            } message: {
                Text("确定要退出当前账号吗？")
            }
            .task {
                await viewModel.load()
            }
        }
    }

    // MARK: - 子视图

    /// 列表主体（区分加载/空/有数据三种状态）
    @ViewBuilder
    private var contentList: some View {
        if viewModel.isLoading && viewModel.articles.isEmpty {
            Text(" ")
        } else if viewModel.articles.isEmpty {
            ContentUnavailableView {
                Label("暂无文章", systemImage: "doc.text")
            } description: {
                Text("点击右上角 ⋯ 创建第一篇文章")
            }
        } else {
            LazyVStack(spacing: 12) {
                ForEach(viewModel.articles) { article in
                    Button {
                        selectedArticle = article
                    } label: {
                        ArticleRowView(
                            article: article,
                            relativeTime: viewModel.relativeTime(for: article.created_at),
                            categoryName: article.cate?.name
                        )
                    }
                    .buttonStyle(.plain)
                    .onAppear {
                        if article == viewModel.articles.last {
                            Task { await viewModel.loadMore() }
                        }
                    }
                }

                // 加载更多指示器
                if viewModel.isLoadingMore {
                    HStack(spacing: 6) {
                        ProgressView()
                        Text("加载更多...")
                    }
                    .foregroundStyle(.secondary)
                    .padding(.vertical, 8)
                }
            }
        }
    }
}

// MARK: - 文章行视图

/// 单个文章行（半透明卡片）
private struct ArticleRowView: View {

    let article: Article
    let relativeTime: String
    let categoryName: String?

    var body: some View {
        VStack(alignment: .leading, spacing: 8) {
            // 标题和分类
            HStack(alignment: .top) {
                Text(article.title)
                    .font(.headline)
                    .lineLimit(2)
                    .multilineTextAlignment(.leading)

                Spacer()

                // 分类标签
                if let categoryName, !categoryName.isEmpty {
                    Text(categoryName)
                        .font(.caption2)
                        .padding(.horizontal, 6)
                        .padding(.vertical, 2)
                        .background(Color.accentColor.opacity(0.12), in: Capsule())
                        .foregroundStyle(Color.accentColor)
                }
            }

            // 时间和浏览量（紧凑布局，缩小图标与文字间距）
            HStack(spacing: 12) {
                HStack(spacing: 3) {
                    Image(systemName: "clock")
                    Text(relativeTime)
                }
                .font(.caption)
                .foregroundStyle(.secondary)

                if article.view_num > 0 {
                    HStack(spacing: 3) {
                        Image(systemName: "eye")
                        Text("\(article.view_num)")
                    }
                    .font(.caption)
                    .foregroundStyle(.secondary)
                }
            }

            // 标签
            if !article.tags.isEmpty {
                FlowLayout(spacing: 6) {
                    ForEach(article.tags, id: \.self) { tag in
                        Text(tag)
                            .font(.caption2)
                            .padding(.horizontal, 8)
                            .padding(.vertical, 3)
                            .background(Color(.secondarySystemBackground), in: Capsule())
                            .foregroundStyle(.secondary)
                    }
                }
            }
        }
        .padding(14)
        .frame(maxWidth: .infinity, alignment: .leading)
        // 半透明卡片背景，让底层装饰背景图透出
        .background(Color(.systemBackground).opacity(0.92))
        .clipShape(RoundedRectangle(cornerRadius: 14))
        .contentShape(RoundedRectangle(cornerRadius: 14))
    }
}

// MARK: - 流式布局（用于标签自动换行）

/// 简单的流式布局，用于标签 chip 自动换行显示
private struct FlowLayout: Layout {

    var spacing: CGFloat = 8

    func sizeThatFits(proposal: ProposedViewSize, subviews: Subviews, cache: inout ()) -> CGSize {
        let result = computeLayout(proposal: proposal, subviews: subviews)
        return result.size
    }

    func placeSubviews(in bounds: CGRect, proposal: ProposedViewSize, subviews: Subviews, cache: inout ()) {
        let result = computeLayout(proposal: proposal, subviews: subviews)

        for (index, position) in result.positions.enumerated() {
            subviews[index].place(
                at: CGPoint(x: bounds.minX + position.x, y: bounds.minY + position.y),
                anchor: .topLeading,
                proposal: ProposedViewSize(result.sizes[index])
            )
        }
    }

    // MARK: - 布局计算

    private struct LayoutResult {
        var size: CGSize
        var positions: [CGPoint]
        var sizes: [CGSize]
    }

    private func computeLayout(proposal: ProposedViewSize, subviews: Subviews) -> LayoutResult {
        var positions: [CGPoint] = []
        var sizes: [CGSize] = []
        var currentX: CGFloat = 0
        var currentY: CGFloat = 0
        var rowHeight: CGFloat = 0
        var maxWidth: CGFloat = 0
        let maxWidthConstraint = proposal.width ?? .infinity

        for subview in subviews {
            let size = subview.sizeThatFits(.unspecified)
            sizes.append(size)

            if currentX + size.width > maxWidthConstraint, currentX > 0 {
                // 换行
                currentX = 0
                currentY += rowHeight + spacing
                rowHeight = 0
            }

            positions.append(CGPoint(x: currentX, y: currentY))
            rowHeight = max(rowHeight, size.height)
            currentX += size.width + spacing
            maxWidth = max(maxWidth, currentX)
        }

        return LayoutResult(
            size: CGSize(width: maxWidth, height: currentY + rowHeight),
            positions: positions,
            sizes: sizes
        )
    }
}
