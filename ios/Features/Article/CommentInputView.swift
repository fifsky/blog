import SwiftUI
import UIKit

/// 评论草稿目标
enum CommentDraftTarget: Equatable {
    /// 发表全新评论
    case new

    /// 回复评论，rootId 始终为顶层主评论 ID
    case reply(rootId: Int, replyName: String, placeholderName: String)

    /// 请求中的父评论 ID
    var pid: Int {
        switch self {
        case .new:
            return 0
        case let .reply(rootId, _, _):
            return rootId
        }
    }

    /// 请求中的被回复人昵称
    var requestReplyName: String? {
        switch self {
        case .new:
            return nil
        case let .reply(_, replyName, _):
            return replyName
        }
    }

    /// 输入框占位文案
    var placeholder: String {
        switch self {
        case .new:
            return "输入评论"
        case let .reply(_, _, placeholderName):
            return "回复\(placeholderName)"
        }
    }
}

/// 底部评论工具栏
struct CommentBottomBar: View {

    /// 打开评论输入回调
    var onCompose: () -> Void

    /// 定位到评论区回调
    var onScrollToComments: () -> Void

    var body: some View {
        HStack(spacing: 14) {
            Button(action: onCompose) {
                HStack(spacing: 7) {
                    Image(systemName: "pencil")
                        .font(.system(size: 15, weight: .medium))
                    Text("写评论...")
                        .font(.system(size: 15))
                }
                .foregroundStyle(.secondary)
                .frame(maxWidth: .infinity, alignment: .leading)
                .padding(.horizontal, 12)
                .frame(height: 38)
                .background(Color(.secondarySystemBackground), in: RoundedRectangle(cornerRadius: 6))
            }
            .buttonStyle(.plain)

            Button(action: onScrollToComments) {
                Image(systemName: "bubble.left.and.bubble.right")
                    .font(.system(size: 24, weight: .regular))
                    .foregroundStyle(.primary)
                    .frame(width: 46, height: 38)
            }
            .buttonStyle(.plain)
            .accessibilityLabel("查看评论")
        }
        .padding(.horizontal, 16)
        .padding(.top, 8)
        .padding(.bottom, 8)
        .background {
            Color.white
                .ignoresSafeArea(edges: .bottom)
                .shadow(color: .black.opacity(0.06), radius: 10, y: -2)
        }
    }
}

/// 评论输入视图
/// 底部弹层形式展示，支持发表评论、回复评论和插入常用 emoji。
struct CommentInputView: View {

    /// 文章 ID
    let postId: Int

    /// 评论目标
    let target: CommentDraftTarget

    /// 发送成功回调
    var onSuccess: (() -> Void)?

    /// 取消输入回调
    var onCancel: (() -> Void)?

    @State private var commentText = ""
    @State private var isSending = false
    @State private var showEmoji = false
    @State private var errorMessage: String?
    @State private var showError = false
    @FocusState private var isInputFocused: Bool

    private let maxContentLength = 1000
    private let commentService = CommentService.shared

    private let emojis = [
        "😀", "😄", "😁", "😆", "😅", "😂", "🤣",
        "😊", "😇", "🙂", "😉", "😌", "😍", "🥰",
        "😘", "😋", "😜", "🤪", "🤨", "🧐", "🤓",
        "😎", "🥳", "😢", "😭", "😤", "😠", "🤬",
        "🙄", "😏", "😱", "🤯", "👍", "👎", "👏",
        "🙌", "🙏", "💪", "🤝", "✌️", "❤️", "💔",
        "🔥", "✨", "🎉", "💯", "🌹", "☕", "😴"
    ]

    var body: some View {
        VStack(spacing: 0) {
            VStack(spacing: 10) {
                TextField(target.placeholder, text: $commentText, axis: .vertical)
                    .focused($isInputFocused)
                    .lineLimit(1 ... 4)
                    .font(.system(size: 16))
                    .padding(.horizontal, 12)
                    .padding(.vertical, 8)
                    .background(Color(.secondarySystemBackground), in: RoundedRectangle(cornerRadius: 5))
                    .submitLabel(.send)
                    .onSubmit {
                        sendComment()
                    }
                    .onTapGesture {
                        showEmoji = false
                    }

                HStack(spacing: 18) {
                    Button {
                        toggleEmojiPanel()
                    } label: {
                        Image(systemName: showEmoji ? "keyboard" : "face.smiling")
                            .font(.system(size: 25, weight: .regular))
                            .foregroundStyle(.secondary)
                            .frame(width: 34, height: 34)
                    }
                    .buttonStyle(.plain)
                    .accessibilityLabel(showEmoji ? "显示键盘" : "选择表情")

                    Spacer()

                    Text("\(commentText.count)/\(maxContentLength)")
                        .font(.system(size: 16))
                        .foregroundStyle(.secondary)

                    Button {
                        sendComment()
                    } label: {
                        if isSending {
                            ProgressView()
                                .tint(.white)
                                .frame(width: 58, height: 34)
                        } else {
                            Text("发送")
                                .font(.system(size: 16, weight: .medium))
                                .frame(width: 58, height: 34)
                        }
                    }
                    .foregroundStyle(.white)
                    .background(sendDisabled ? Color.blue.opacity(0.35) : Color.blue.opacity(0.85), in: Capsule())
                    .disabled(sendDisabled)
                }
            }
            .padding(.horizontal, 16)
            .padding(.top, 14)
            .padding(.bottom, 12)

            if showEmoji {
                emojiPanel
                    .transition(.move(edge: .bottom).combined(with: .opacity))
            }
        }
        .frame(maxWidth: .infinity)
        .background(alignment: .top) {
            UnevenRoundedRectangle(topLeadingRadius: 18, topTrailingRadius: 18)
                .fill(Color.white)
        }
        .background(alignment: .bottom) {
            Color.white
                .frame(height: 160)
                .offset(y: 96)
                .ignoresSafeArea(edges: .bottom)
        }
        .onAppear {
            DispatchQueue.main.asyncAfter(deadline: .now() + 0.25) {
                isInputFocused = true
            }
        }
        .onChange(of: commentText) {
            if commentText.count > maxContentLength {
                commentText = String(commentText.prefix(maxContentLength))
            }
        }
        .onChange(of: isInputFocused) { _, focused in
            if focused {
                showEmoji = false
            }
        }
        .alert("错误", isPresented: $showError) {
            Button("确定", role: .cancel) {}
        } message: {
            Text(errorMessage ?? "未知错误")
        }
    }

    /// emoji 面板
    private var emojiPanel: some View {
        ScrollView {
            LazyVGrid(
                columns: Array(repeating: GridItem(.flexible(), spacing: 14), count: 7),
                spacing: 18
            ) {
                ForEach(emojis, id: \.self) { emoji in
                    Button {
                        insertEmoji(emoji)
                    } label: {
                        Text(emoji)
                            .font(.system(size: 30))
                            .frame(width: 38, height: 38)
                    }
                    .buttonStyle(.plain)
                }
            }
            .padding(.horizontal, 18)
            .padding(.vertical, 18)
        }
        .frame(height: 320)
        .background(Color(.systemGroupedBackground))
    }

    /// 是否禁用发送按钮
    private var sendDisabled: Bool {
        commentText.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty || isSending
    }

    /// 切换 emoji 面板
    private func toggleEmojiPanel() {
        if showEmoji {
            showEmoji = false
            isInputFocused = true
        } else {
            isInputFocused = false
            hideKeyboard()
            withAnimation(.easeInOut(duration: 0.2)) {
                showEmoji = true
            }
        }
    }

    /// 插入 emoji 表情
    /// - Parameter emoji: 被选中的 emoji
    private func insertEmoji(_ emoji: String) {
        guard commentText.count + emoji.count <= maxContentLength else { return }
        commentText.append(emoji)
    }

    /// 收起系统键盘
    private func hideKeyboard() {
        UIApplication.shared.sendAction(
            #selector(UIResponder.resignFirstResponder),
            to: nil,
            from: nil,
            for: nil
        )
    }

    /// 发送评论
    private func sendComment() {
        let trimmed = commentText.trimmingCharacters(in: .whitespacesAndNewlines)
        guard !trimmed.isEmpty, !isSending else { return }

        isSending = true
        hideKeyboard()

        Task { @MainActor in
            do {
                let author = currentCommentAuthor()
                _ = try await commentService.create(
                    postId: postId,
                    name: author.name,
                    content: trimmed,
                    email: author.email,
                    website: author.website,
                    pid: target.pid,
                    replyName: target.requestReplyName
                )
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

    /// 当前登录用户的评论身份，和 web 管理员登录后的评论身份保持一致。
    /// - Returns: 昵称、邮箱和站点网址
    @MainActor
    private func currentCommentAuthor() -> (name: String, email: String, website: String) {
        let currentUser = AuthManager.shared.currentUser
        let name = [currentUser?.name, currentUser?.nick_name]
            .compactMap { $0?.trimmingCharacters(in: .whitespacesAndNewlines) }
            .first { !$0.isEmpty } ?? "fifsky"
        let email = currentUser?.email.trimmingCharacters(in: .whitespacesAndNewlines) ?? ""

        return (
            name: name,
            email: email,
            website: "https://fifsky.com"
        )
    }
}
