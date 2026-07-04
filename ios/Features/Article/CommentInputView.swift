import SwiftUI

/// 评论输入视图
/// 固定在底部的评论输入框，支持发表评论和回复评论
struct CommentInputView: View {

    /// 文章 ID
    let postId: Int

    /// 回复的目标用户名（nil 表示发表新评论）
    var replyName: String?

    /// 回复的评论 ID（0 表示发表新评论）
    var pid: Int = 0

    /// 发送成功回调
    var onSuccess: (() -> Void)?

    @State private var commentText = ""
    @State private var isSending = false
    @State private var errorMessage: String?
    @State private var showError = false

    private let commentService = CommentService.shared

    var body: some View {
        VStack(spacing: 0) {
            // 回复提示
            if let replyName, !replyName.isEmpty {
                HStack {
                    Image(systemName: "arrow.turn.down.left")
                        .font(.caption)
                    Text("回复 \(replyName)")
                        .font(.caption)
                        .foregroundStyle(.secondary)
                    Spacer()
                    Button {
                        // 由父视图清除回复状态
                    } label: {
                        Image(systemName: "xmark.circle.fill")
                            .font(.caption)
                            .foregroundStyle(.secondary)
                    }
                }
                .padding(.horizontal, 16)
                .padding(.top, 8)
            }

            // 输入框和发送按钮
            HStack(spacing: 8) {
                TextField(
                    replyName != nil ? "回复 \(replyName ?? "")..." : "写评论...",
                    text: $commentText,
                    axis: .vertical
                )
                .lineLimit(1 ... 5)
                .textFieldStyle(.roundedBorder)
                .submitLabel(.send)
                .onSubmit {
                    sendComment()
                }

                // 发送按钮
                Button {
                    sendComment()
                } label: {
                    if isSending {
                        ProgressView()
                            .tint(.blue)
                    } else {
                        Image(systemName: "paperplane.fill")
                            .foregroundStyle(commentText.trimmingCharacters(in: .whitespaces).isEmpty ? .gray : .blue)
                    }
                }
                .disabled(commentText.trimmingCharacters(in: .whitespaces).isEmpty || isSending)
            }
            .padding(.horizontal, 16)
            .padding(.vertical, 10)
            .background(.bar, in: Rectangle())
        }
        .alert("错误", isPresented: $showError) {
            Button("确定", role: .cancel) {}
        } message: {
            Text(errorMessage ?? "未知错误")
        }
    }

    // MARK: - 发送评论

    /// 发送评论
    private func sendComment() {
        let trimmed = commentText.trimmingCharacters(in: .whitespacesAndNewlines)
        guard !trimmed.isEmpty else { return }

        isSending = true

        Task {
            do {
                _ = try await commentService.create(
                    postId: postId,
                    name: AuthManager.shared.currentUser?.nick_name
                        ?? AuthManager.shared.currentUser?.name
                        ?? "匿名用户",
                    content: trimmed,
                    pid: pid,
                    replyName: replyName
                )
                // 发送成功
                commentText = ""
                onSuccess?()
            } catch {
                let apiError = APIError.normalized(error)
                if !apiError.isCancelled {
                    errorMessage = apiError.displayMessage()
                    showError = true
                }
            }

            isSending = false
        }
    }
}
