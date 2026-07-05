import SwiftUI

/// 评论列表视图
/// 显示文章下的所有评论，主评论和回复保持两层结构。
struct CommentListView: View {

    let postId: Int

    /// 刷新触发器（外部改变时重新加载）
    var refreshTrigger: Bool = false

    /// 回复回调
    var onReply: ((CommentDraftTarget) -> Void)?

    @State private var viewModel: CommentListViewModel?

    var body: some View {
        VStack(alignment: .leading, spacing: 20) {
            header

            if let viewModel {
                if viewModel.isLoading && viewModel.comments.isEmpty {
                    loadingView
                } else if viewModel.roots.isEmpty {
                    emptyView
                } else {
                    commentList(viewModel)
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

    /// 评论区标题
    private var header: some View {
        HStack(alignment: .center) {
            Text("全部评论 \(viewModel?.totalCount ?? 0)")
                .font(.system(size: 21, weight: .semibold))
                .foregroundStyle(.primary)

            Spacer()
        }
        .padding(.top, 4)
    }

    /// 加载状态
    private var loadingView: some View {
        HStack {
            Spacer()
            ProgressView()
                .frame(height: 52)
            Spacer()
        }
    }

    /// 空评论状态
    private var emptyView: some View {
        VStack(spacing: 8) {
            Image(systemName: "bubble.left")
                .font(.title2)
                .foregroundStyle(.secondary)
            Text("暂无评论，来说点什么吧")
                .font(.subheadline)
                .foregroundStyle(.secondary)
        }
        .frame(maxWidth: .infinity)
        .padding(.vertical, 30)
    }

    /// 评论列表内容
    /// - Parameter viewModel: 评论列表视图模型
    /// - Returns: 评论列表视图
    private func commentList(_ viewModel: CommentListViewModel) -> some View {
        LazyVStack(alignment: .leading, spacing: 22) {
            ForEach(viewModel.roots) { node in
                VStack(alignment: .leading, spacing: 12) {
                    CommentItemView(
                        comment: node.root,
                        relativeTime: viewModel.relativeTime(for: node.root.created_at),
                        onReply: {
                            onReply?(.reply(
                                rootId: node.root.id,
                                replyName: "",
                                placeholderName: node.root.name
                            ))
                        }
                    )

                    if !node.replies.isEmpty {
                        VStack(alignment: .leading, spacing: 14) {
                            ForEach(node.replies) { reply in
                                CommentItemView(
                                    comment: reply,
                                    isReply: true,
                                    relativeTime: viewModel.relativeTime(for: reply.created_at),
                                    onReply: {
                                        onReply?(.reply(
                                            rootId: node.root.id,
                                            replyName: reply.name,
                                            placeholderName: reply.name
                                        ))
                                    }
                                )
                            }
                        }
                        .padding(.leading, 52)
                    }
                }
                .padding(.bottom, 18)
            }
        }
    }
}

// MARK: - 评论项视图

/// 单条评论视图
private struct CommentItemView: View {

    let comment: Comment
    var isReply = false
    let relativeTime: String
    var onReply: (() -> Void)?

    private var avatarSize: CGFloat {
        isReply ? 32 : 42
    }

    var body: some View {
        HStack(alignment: .top, spacing: 12) {
            CommentAvatarView(comment: comment, size: avatarSize)

            VStack(alignment: .leading, spacing: 7) {
                nameRow

                if !comment.reply_name.isEmpty {
                    replyNameRow
                }

                Text(comment.content)
                    .font(.system(size: isReply ? 15 : 16))
                    .foregroundStyle(.primary)
                    .lineSpacing(3)
                    .fixedSize(horizontal: false, vertical: true)
                    .textSelection(.enabled)

                HStack(spacing: 14) {
                    Text(relativeTime)
                        .font(.caption)
                        .foregroundStyle(.secondary)

                    Spacer()

                    Button {
                        onReply?()
                    } label: {
                        Image(systemName: "bubble.left")
                            .font(.system(size: 15, weight: .regular))
                            .foregroundStyle(.secondary)
                            .frame(width: 28, height: 24)
                    }
                    .buttonStyle(.plain)
                    .accessibilityLabel("回复 \(comment.name)")
                }
            }
            .frame(maxWidth: .infinity, alignment: .leading)
        }
    }

    /// 昵称行
    private var nameRow: some View {
        Text(comment.name)
            .font(.system(size: isReply ? 16 : 17, weight: .medium))
            .foregroundStyle(.secondary)
            .lineLimit(1)
    }

    /// 回复对象行
    private var replyNameRow: some View {
        HStack(spacing: 4) {
            Text("回复")
                .foregroundStyle(.secondary)
            Text("@\(comment.reply_name)")
                .foregroundStyle(.blue)
        }
        .font(.caption)
    }
}

// MARK: - 评论头像

/// 评论头像视图
private struct CommentAvatarView: View {

    let comment: Comment
    let size: CGFloat

    var body: some View {
        Group {
            if let url = URL(string: comment.avatar), !comment.avatar.isEmpty {
                AsyncImage(url: url) { phase in
                    switch phase {
                    case let .success(image):
                        image
                            .resizable()
                            .scaledToFill()
                    default:
                        fallback
                    }
                }
            } else {
                fallback
            }
        }
        .frame(width: size, height: size)
        .clipShape(Circle())
        .overlay {
            Circle()
                .stroke(Color(.systemGray5), lineWidth: 1)
        }
    }

    /// 头像加载失败时的占位
    private var fallback: some View {
        Circle()
            .fill(Color(.systemGray4))
            .overlay {
                Text(String(comment.name.prefix(1)))
                    .font(.system(size: size * 0.4, weight: .bold))
                    .foregroundStyle(.white)
            }
    }
}
