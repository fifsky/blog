import SwiftUI

/// 评论列表视图
/// 显示文章下的所有评论，支持主评论和嵌套回复
struct CommentListView: View {

    let postId: Int

    /// 刷新触发器（外部改变时重新加载）
    var refreshTrigger: Bool = false

    /// 回复回调
    var onReply: ((String, Int) -> Void)?

    @State private var viewModel: CommentListViewModel?

    var body: some View {
        Group {
            if let viewModel {
                if viewModel.isLoading && viewModel.comments.isEmpty {
                    // 加载中
                    VStack {
                        Spacer()
                        ProgressView()
                            .frame(height: 40)
                        Spacer()
                    }
                } else if viewModel.comments.isEmpty {
                    // 无评论
                    VStack(spacing: 8) {
                        Image(systemName: "bubble.left")
                            .font(.title2)
                            .foregroundStyle(.secondary)
                        Text("暂无评论，来说点什么吧")
                            .font(.subheadline)
                            .foregroundStyle(.secondary)
                    }
                    .frame(maxWidth: .infinity)
                    .padding(.vertical, 24)
                } else {
                    // 评论列表
                    LazyVStack(alignment: .leading, spacing: 16) {
                        ForEach(viewModel.mainComments) { comment in
                            CommentItemView(
                                comment: comment,
                                relativeTime: viewModel.relativeTime(for: comment.created_at),
                                onReply: {
                                    onReply?(comment.name, comment.id)
                                }
                            )

                            // 嵌套回复
                            let replies = viewModel.replies(for: comment.id)
                            if !replies.isEmpty {
                                VStack(alignment: .leading, spacing: 12) {
                                    ForEach(replies) { reply in
                                        CommentItemView(
                                            comment: reply,
                                            isReply: true,
                                            relativeTime: viewModel.relativeTime(for: reply.created_at),
                                            onReply: {
                                                onReply?(reply.name, reply.id)
                                            }
                                        )
                                    }
                                }
                                .padding(.leading, 20)
                            }
                        }
                    }
                }
            }
        }
        .onAppear {
            if viewModel == nil {
                viewModel = CommentListViewModel(postId: postId)
            }
            Task {
                await viewModel?.load()
            }
        }
        .onChange(of: refreshTrigger) {
            // 外部触发刷新
            Task {
                await viewModel?.load()
            }
        }
        .alert("错误", isPresented: Binding(
            get: { viewModel?.showError ?? false },
            set: { if !$0 { viewModel?.showError = false } }
        )) {
            Button("确定", role: .cancel) {}
        } message: {
            Text(viewModel?.errorMessage ?? "未知错误")
        }
    }
}

// MARK: - 评论项视图

/// 单条评论视图
private struct CommentItemView: View {

    let comment: Comment
    var isReply: Bool = false
    let relativeTime: String
    var onReply: (() -> Void)?

    var body: some View {
        HStack(alignment: .top, spacing: 10) {
            // 头像占位
            Circle()
                .fill(Color(.systemGray4))
                .frame(width: 36, height: 36)
                .overlay {
                    Text(String(comment.name.prefix(1)))
                        .font(.caption)
                        .bold()
                        .foregroundStyle(.white)
                }

            VStack(alignment: .leading, spacing: 4) {
                // 名称和时间
                HStack(alignment: .center, spacing: 8) {
                    Text(comment.name)
                        .font(.subheadline)
                        .bold()

                    Text(relativeTime)
                        .font(.caption2)
                        .foregroundStyle(.secondary)
                }

                // 回复标记
                if !comment.reply_name.isEmpty {
                    Text("回复 \(comment.reply_name)")
                        .font(.caption)
                        .foregroundStyle(.secondary)
                }

                // 评论内容
                Text(comment.content)
                    .font(.subheadline)
                    .foregroundStyle(.primary)
                    .textSelection(.enabled)

                // 回复按钮
                Button {
                    onReply?()
                } label: {
                    Text("回复")
                        .font(.caption2)
                        .foregroundStyle(.secondary)
                }
            }
        }
    }
}
