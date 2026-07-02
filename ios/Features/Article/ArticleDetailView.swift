import SwiftUI

/// 文章详情视图
struct ArticleDetailView: View {

    let articleId: Int

    @State private var viewModel: ArticleDetailViewModel?

    /// 评论回复信息
    @State private var replyName: String?
    @State private var replyPid: Int = 0

    /// 需要刷新评论列表
    @State private var refreshCommentTrigger = false

    var body: some View {
        ZStack {
            if let viewModel, let article = viewModel.article {
                // 文章内容
                ScrollView {
                    VStack(alignment: .leading, spacing: 16) {
                        // 标题
                        Text(article.title)
                            .font(.largeTitle)
                            .bold()
                            .multilineTextAlignment(.leading)

                        // 日期和浏览量
                        HStack(spacing: 16) {
                            Label(
                                viewModel.relativeTime(for: article.created_at),
                                systemImage: "clock"
                            )
                            .font(.subheadline)
                            .foregroundStyle(.secondary)

                            Label(
                                "\(article.view_num) 次浏览",
                                systemImage: "eye"
                            )
                            .font(.subheadline)
                            .foregroundStyle(.secondary)
                        }

                        // 分类
                        if let cate = article.cate {
                            HStack(spacing: 6) {
                                Image(systemName: "folder")
                                    .font(.caption)
                                    .foregroundStyle(.secondary)
                                Text(cate.name)
                                    .font(.caption)
                                    .foregroundStyle(.secondary)
                            }
                        }

                        // 标签
                        if !article.tags.isEmpty {
                            ScrollView(.horizontal, showsIndicators: false) {
                                HStack(spacing: 8) {
                                    ForEach(article.tags, id: \.self) { tag in
                                        Text(tag)
                                            .font(.caption)
                                            .padding(.horizontal, 10)
                                            .padding(.vertical, 4)
                                            .background(Color(.tertiarySystemBackground), in: Capsule())
                                            .foregroundStyle(.secondary)
                                    }
                                }
                            }
                        }

                        Divider()

                        // 正文内容
                        // TODO: 集成 swift-markdown-ui 后替换为 Markdown(content)
                        Text(article.content)
                            .font(.body)
                            .lineSpacing(6)
                            .multilineTextAlignment(.leading)
                            .textSelection(.enabled)

                        Divider()

                        // 评论区
                        VStack(alignment: .leading, spacing: 12) {
                            Text("评论")
                                .font(.headline)

                            CommentListView(
                                postId: article.id,
                                refreshTrigger: refreshCommentTrigger,
                                onReply: { name, pid in
                                    replyName = name
                                    replyPid = pid
                                }
                            )
                        }

                        // 底部留白，避免被输入框遮挡
                        Spacer(minLength: 80)
                    }
                    .padding(.horizontal, 16)
                }
                .overlay(alignment: .bottom) {
                    // 评论输入框
                    CommentInputView(
                        postId: article.id,
                        replyName: replyName,
                        pid: replyPid,
                        onSuccess: {
                            // 清除回复状态，触发评论列表刷新
                            replyName = nil
                            replyPid = 0
                            refreshCommentTrigger.toggle()
                        }
                    )
                }
            } else if let viewModel, viewModel.isLoading {
                // 加载中
                ProgressView("加载中...")
                    .frame(maxWidth: .infinity, maxHeight: .infinity)
            }
        }
        .navigationTitle("文章详情")
        .navigationBarTitleDisplayMode(.inline)
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
}
