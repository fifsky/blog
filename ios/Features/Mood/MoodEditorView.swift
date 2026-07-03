import SwiftUI

/// 心情编辑视图（支持新建和编辑）
struct MoodEditorView: View {

    /// 传入已有的心情对象以进入编辑模式，nil 为新建模式
    var existingMood: Mood?

    @State private var viewModel: MoodEditorViewModel
    @Environment(\.dismiss) private var dismiss

    /// 内容输入框焦点（新建模式自动聚焦）
    @FocusState private var contentFocused: Bool

    init(mood: Mood? = nil) {
        self.existingMood = mood
        _viewModel = State(initialValue: MoodEditorViewModel(mood: mood))
    }

    var body: some View {
        VStack(spacing: 0) {
            // 文本编辑区域
            TextEditor(text: $viewModel.content)
                .font(.body)
                .textInputAutocapitalization(.sentences)
                .autocorrectionDisabled()
                .scrollContentBackground(.hidden)
                .focused($contentFocused)
                .padding(.horizontal, 16)
                .padding(.top, 12)
                .overlay(alignment: .topLeading) {
                    // 占位提示文字
                    if viewModel.content.isEmpty {
                        Text("此刻的心情...")
                            .font(.body)
                            .foregroundStyle(.quaternary)
                            .padding(.horizontal, 20)
                            .padding(.top, 18)
                            .allowsHitTesting(false)
                    }
                }
        }
        .navigationTitle(viewModel.title)
        .navigationBarTitleDisplayMode(.inline)
        .toolbar {
            // 表情按钮
            ToolbarItem(placement: .keyboard) {
                // 键盘 accessory bar 紧贴输入法边框，加底部留白避免顶死。
                // padding 加在 Button 上，由系统 accessory bar 容器承担为整体下移效果。
                Button {
                    viewModel.showEmojiPicker = true
                } label: {
                    Text(viewModel.content.isEmpty ? "😊" : "😊")
                        .font(.title3)
                }
            }

            // 保存按钮
            ToolbarItem(placement: .confirmationAction) {
                Button {
                    Task { await viewModel.save() }
                } label: {
                    if viewModel.isSaving {
                        ProgressView()
                    } else {
                        Text("保存")
                            .fontWeight(.medium)
                    }
                }
                .disabled(viewModel.isSaveButtonDisabled)
            }

            // 取消按钮
            ToolbarItem(placement: .cancellationAction) {
                Button("取消") {
                    dismiss()
                }
            }
        }
        .sheet(isPresented: $viewModel.showEmojiPicker) {
            EmojiPicker { emoji in
                viewModel.insertEmoji(emoji)
            }
        }
        .alert("保存失败", isPresented: $viewModel.showError) {
            Button("确定", role: .cancel) {}
        } message: {
            Text(viewModel.errorMessage ?? "未知错误")
        }
        .onChange(of: viewModel.didSave) {
            if viewModel.didSave {
                dismiss()
            }
        }
        .onAppear {
            // 新建模式自动聚焦内容输入框
            if existingMood == nil {
                contentFocused = true
            }
        }
    }
}
