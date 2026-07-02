import Foundation

/// 心情编辑视图模型
@Observable
class MoodEditorViewModel {

    // MARK: - 状态

    /// 心情内容
    var content = ""

    /// 是否正在保存
    var isSaving = false

    /// 是否保存成功（用于关闭视图）
    var didSave = false

    /// 错误信息
    var errorMessage: String?

    /// 是否显示错误弹窗
    var showError = false

    /// 是否显示表情选择器
    var showEmojiPicker = false

    // MARK: - 私有属性

    /// 编辑模式下的已有心情（nil 表示新建）
    private let existingMood: Mood?
    private let moodService = MoodService.shared

    // MARK: - 计算属性

    /// 是否为编辑模式
    var isEditing: Bool {
        existingMood != nil
    }

    /// 标题
    var title: String {
        isEditing ? "编辑心情" : "记录心情"
    }

    /// 保存按钮是否可用
    var isSaveButtonDisabled: Bool {
        content.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty || isSaving
    }

    // MARK: - 初始化

    /// 初始化编辑器
    /// - Parameter mood: 已有的心情对象，传入时为编辑模式，nil 为新建模式
    init(mood: Mood? = nil) {
        self.existingMood = mood
        self.content = mood?.content ?? ""
    }

    // MARK: - 操作

    /// 在光标位置插入表情
    /// - Parameter emoji: 要插入的表情字符
    func insertEmoji(_ emoji: String) {
        content.append(emoji)
    }

    /// 保存心情
    func save() async {
        guard !isSaveButtonDisabled else { return }
        isSaving = true

        do {
            if let mood = existingMood {
                // 编辑模式：更新
                _ = try await moodService.update(id: mood.id, content: content)
            } else {
                // 新建模式：创建
                _ = try await moodService.create(content: content)
            }
            didSave = true
        } catch {
            errorMessage = error.localizedDescription
            showError = true
        }

        isSaving = false
    }
}
