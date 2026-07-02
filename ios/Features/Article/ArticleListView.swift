import SwiftUI

/// 文章列表视图
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
            Group {
                if viewModel.isLoading && viewModel.articles.isEmpty {
                    // 首次加载中
                    ProgressView("加载中...")
                        .frame(maxWidth: .infinity, maxHeight: .infinity)
                } else if viewModel.articles.isEmpty {
                    // 空状态
                    ContentUnavailableView {
                        Label("暂无文章", systemImage: "doc.text")
                    } description: {
                        Text("点击右上角 ⋯ 创建第一篇文章")
                    }
                } else {
                    // 文章列表
                    articleList
                }
            }
            .navigationTitle("博文")
            .toolbar {
                ToolbarItem(placement: .topBarTrailing) {
                    // 三点菜单按钮，点击弹出新增博文 / 退出登录选项
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
                            .font(.system(size: 16, weight: .medium))
                    }
                }
            }
            .refreshable {
                await viewModel.refresh()
            }
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

    /// 文章列表主体
    private var articleList: some View {
        List {
            ForEach(viewModel.articles) { article in
                ArticleRowView(
                    article: article,
                    relativeTime: viewModel.relativeTime(for: article.created_at),
                    categoryName: article.cate?.name
                )
                .onTapGesture {
                    selectedArticle = article
                }
                .onAppear {
                    // 滚动到最后几条时加载更多
                    if article == viewModel.articles.last {
                        Task {
                            await viewModel.loadMore()
                        }
                    }
                }
            }

            // 加载更多指示器
            if viewModel.isLoadingMore {
                HStack {
                    Spacer()
                    ProgressView()
                    Text("加载更多...")
                    Spacer()
                }
                .listRowSeparator(.hidden)
            }
        }
        .listStyle(.plain)
    }
}

// MARK: - 文章行视图

/// 单个文章行
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
        .padding(.vertical, 4)
        .contentShape(Rectangle())
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
