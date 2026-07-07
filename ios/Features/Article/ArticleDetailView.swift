import SwiftUI
import Textual

/// 修复表格内图片显示不全的单元格样式
///
/// Textual 在 `StructuredText` 顶层定义 `.textContainer` 坐标空间，
/// 图片附件按整个文本容器的宽度缩放（`ImageAttachment.sizeThatFits` 使用 `min(proposedWidth, imageWidth)`）。
/// 表格单元格比文本容器窄（多列表格），图片按容器宽度渲染后溢出单元格被裁剪。
/// 此样式在单元格级别重新定义 `.textContainer` 坐标空间，使 `TextFragment` 读到的容器宽度为单元格宽度，
/// 图片按单元格宽度缩放，不再溢出裁剪。
private struct TableImageFixCellStyle: StructuredText.TableCellStyle {
    func makeBody(configuration: Configuration) -> some View {
        configuration.label
            .fontWeight(configuration.row == 0 ? .semibold : .regular)
            .textual.lineSpacing(.fontScaled(0.471))
            .coordinateSpace(.named("textContainer"))
            .textual.padding(.fontScaled(0.588))
    }
}

/// 文章详情视图
struct ArticleDetailView: View {

    let articleId: Int

    @State private var viewModel: ArticleDetailViewModel?

    /// 当前评论草稿目标
    @State private var draftTarget: CommentDraftTarget = .new

    /// 是否展示评论输入弹层
    @State private var showCommentInput = false

    /// 需要刷新评论列表
    @State private var refreshCommentTrigger = false

    private let commentsAnchor = "article-comments-section"

    var body: some View {
        ZStack(alignment: .bottom) {
            Color(.systemGroupedBackground)
                .ignoresSafeArea()

            if let viewModel, let article = viewModel.article {
                ScrollViewReader { proxy in
                    ScrollView {
                        VStack(spacing: 0) {
                            articleSection(article: article, viewModel: viewModel)

                            // 文章区与评论区之间的灰色间隔条（透出底层 systemGroupedBackground）
                            Color(.systemGroupedBackground)
                                .frame(height: 12)

                            commentsSection(article: article)
                                .id(commentsAnchor)

                            Spacer(minLength: 56)
                        }
                    }
                    // 整个 ScrollView 底色用白色：顶部安全区与文章区块均为纯白，
                    // 两区块之间靠上面的 12pt systemGroupedBackground 间隔条区分
                    .background(Color(.systemBackground))
                    .scrollDismissesKeyboard(.interactively)
                    .safeAreaInset(edge: .bottom) {
                        CommentBottomBar(
                            onCompose: {
                                openCommentInput(.new)
                            },
                            onScrollToComments: {
                                withAnimation(.easeInOut(duration: 0.25)) {
                                    proxy.scrollTo(commentsAnchor, anchor: .top)
                                }
                            }
                        )
                    }
                }
            } else if let viewModel, viewModel.isLoading {
                ProgressView("加载中...")
                    .frame(maxWidth: .infinity, maxHeight: .infinity)
            } else {
                ProgressView("加载中...")
                    .frame(maxWidth: .infinity, maxHeight: .infinity)
            }

            if let viewModel, let article = viewModel.article, showCommentInput {
                commentInputOverlay(article: article)
                    .zIndex(2)
            }
        }
        .animation(.easeInOut(duration: 0.2), value: showCommentInput)
        .navigationTitle("文章详情")
        .navigationBarTitleDisplayMode(.inline)
        .toolbar(.hidden, for: .tabBar)
        .alert("错误", isPresented: Binding(
            get: { viewModel?.showError ?? false },
            set: { if !$0 { viewModel?.showError = false } }
        )) {
            Button("确定", role: .cancel) {}
        } message: {
            Text(viewModel?.errorMessage ?? "未知错误")
        }
        .onAppear {
            if viewModel == nil {
                viewModel = ArticleDetailViewModel(articleId: articleId)
            }
            Task {
                await viewModel?.load()
            }
        }
    }

    /// 文章内容区块
    /// - Parameters:
    ///   - article: 文章详情
    ///   - viewModel: 文章详情视图模型
    /// - Returns: 文章内容视图
    private func articleSection(article: Article, viewModel: ArticleDetailViewModel) -> some View {
        VStack(alignment: .leading, spacing: 18) {
            Text(article.title)
                .font(.system(size: 30, weight: .semibold))
                .multilineTextAlignment(.leading)
                .fixedSize(horizontal: false, vertical: true)

            articleMetaRow(article: article, viewModel: viewModel)

            StructuredText(markdown: article.content)
                .multilineTextAlignment(.leading)
                .textual.textSelection(.enabled)
                .textual.overflowMode(.scroll)
                .textual.tableStyle(.overflow)
                .textual.tableCellStyle(TableImageFixCellStyle())
                .frame(maxWidth: .infinity, alignment: .leading)
                .padding(.top, 4)

            if !article.tags.isEmpty {
                tagSection(tags: article.tags)
                    .padding(.top, 8)
            }

            updateTimeSection(text: viewModel.updateTime(for: article.updated_at))
                .padding(.top, 28)
        }
        .frame(maxWidth: .infinity, alignment: .leading)
        .padding(.horizontal, 16)
        .padding(.top, 18)
        .padding(.bottom, 26)
        .background(Color(.systemBackground))
    }

    /// 文章分类、发布时间和阅读量元信息
    /// - Parameters:
    ///   - article: 文章详情
    ///   - viewModel: 文章详情视图模型
    /// - Returns: 元信息视图
    private func articleMetaRow(article: Article, viewModel: ArticleDetailViewModel) -> some View {
        HStack(spacing: 10) {
            if let cate = article.cate {
                metaItem(systemImage: "folder", text: cate.name)
                    .lineLimit(1)
            }

            metaItem(systemImage: "clock", text: viewModel.relativeTime(for: article.created_at))
                .layoutPriority(1)

            metaItem(systemImage: "eye", text: "\(article.view_num) 次浏览")
                .layoutPriority(1)
        }
        .font(.system(size: 13))
        .foregroundStyle(.secondary)
        .lineLimit(1)
        .minimumScaleFactor(0.85)
    }

    /// 单个文章元信息项
    /// - Parameters:
    ///   - systemImage: 系统图标名称
    ///   - text: 展示文本
    /// - Returns: 元信息项视图
    private func metaItem(systemImage: String, text: String) -> some View {
        HStack(spacing: 4) {
            Image(systemName: systemImage)
                .font(.system(size: 12, weight: .regular))
            Text(text)
        }
    }

    /// 标签区块
    /// - Parameter tags: 标签列表
    /// - Returns: 标签视图
    private func tagSection(tags: [String]) -> some View {
        HStack(alignment: .top, spacing: 8) {
            Image(systemName: "tag.fill")
                .font(.system(size: 16, weight: .regular))
                .foregroundStyle(Color(.systemGray3))
                .frame(width: 22, height: 22)
                .padding(.top, 2)

            FlowLayout(spacing: 8) {
                ForEach(tags, id: \.self) { tag in
                    Text(tag)
                        .font(.system(size: 13))
                        .foregroundStyle(.secondary)
                        .padding(.horizontal, 10)
                        .padding(.vertical, 5)
                        .background(Color(.secondarySystemBackground), in: Capsule())
                }
            }
        }
    }

    /// 更新时间区块
    /// - Parameter text: 更新时间文本
    /// - Returns: 更新时间视图
    private func updateTimeSection(text: String) -> some View {
        HStack(spacing: 10) {
            Rectangle()
                .fill(Color(.systemGray4))
                .frame(maxWidth: 74)
                .frame(height: 1)

            Text("更新于 \(text)")
                .font(.system(size: 13))
                .foregroundStyle(.secondary)
                .lineLimit(1)
                .fixedSize(horizontal: true, vertical: false)
                .layoutPriority(1)

            Rectangle()
                .fill(Color(.systemGray4))
                .frame(maxWidth: 74)
                .frame(height: 1)
        }
        .frame(maxWidth: .infinity)
        .padding(.horizontal, 12)
    }

    /// 评论区块
    /// - Parameter article: 文章详情
    /// - Returns: 评论视图
    private func commentsSection(article: Article) -> some View {
        VStack(alignment: .leading, spacing: 0) {
            CommentListView(
                postId: article.id,
                refreshTrigger: refreshCommentTrigger,
                onReply: { target in
                    openCommentInput(target)
                }
            )
        }
        .frame(maxWidth: .infinity, alignment: .leading)
        .padding(.horizontal, 16)
        .padding(.top, 22)
        .padding(.bottom, 28)
        .background(Color(.systemBackground))
    }

    /// 评论输入遮罩
    /// - Parameter article: 文章详情
    /// - Returns: 输入弹层视图
    private func commentInputOverlay(article: Article) -> some View {
        ZStack(alignment: .bottom) {
            Color.black.opacity(0.38)
                .ignoresSafeArea()
                .onTapGesture {
                    closeCommentInput()
                }

            CommentInputView(
                postId: article.id,
                target: draftTarget,
                onSuccess: {
                    closeCommentInput()
                    refreshCommentTrigger.toggle()
                },
                onCancel: {
                    closeCommentInput()
                }
            )
            .transition(.move(edge: .bottom).combined(with: .opacity))
        }
        .ignoresSafeArea(.container, edges: .bottom)
    }

    /// 打开评论输入弹层
    /// - Parameter target: 评论目标
    private func openCommentInput(_ target: CommentDraftTarget) {
        draftTarget = target
        showCommentInput = true
    }

    /// 关闭评论输入弹层
    private func closeCommentInput() {
        showCommentInput = false
        draftTarget = .new
    }
}
