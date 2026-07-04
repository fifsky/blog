import Foundation

/// 提醒创建模式
enum RemindCreateMode: Hashable {
    /// AI 智能模式：用户输入自然语言描述，AI 自动生成 cron 和内容
    case smart
    /// 手动模式：用户手动设置 cron 表达式或预设
    case manual
}

/// 提醒编辑视图模型
@Observable
class RemindEditorViewModel {

    // MARK: - 状态

    /// 创建模式（仅新建模式有效，默认 AI 智能模式）
    var createMode: RemindCreateMode = .smart

    /// 提醒内容
    var content = ""

    /// cron 表达式
    var cronExpression = ""

    /// 是否使用自定义 cron（而非预设）
    var isCustomCron = false

    /// 是否正在保存
    var isSaving = false

    /// 是否保存成功（用于关闭视图）
    var didSave = false

    /// 错误信息
    var errorMessage: String?

    /// 是否显示错误弹窗
    var showError = false

    // MARK: - 私有属性

    /// 编辑模式下的已有提醒（nil 表示新建）
    private let existingRemind: Remind?
    private let remindService = RemindService.shared

    // MARK: - 计算属性

    /// 是否为编辑模式
    var isEditing: Bool {
        existingRemind != nil
    }

    /// 标题
    var title: String {
        isEditing ? "编辑提醒" : "新建提醒"
    }

    /// 保存按钮是否可用
    var isSaveButtonDisabled: Bool {
        content.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty || isSaving
    }

    /// 是否显示 cron 设置区域（手动模式或编辑模式）
    var showsCronSettings: Bool {
        isEditing || createMode == .manual
    }

    // MARK: - 预设 cron 选项

    /// cron 预设列表
    static let cronPresets: [(label: String, expression: String)] = [
        ("每天", "0 9 * * *"),
        ("工作日", "0 9 * * 1-5"),
        ("每周一", "0 9 * * 1"),
        ("每月1日", "0 9 1 * *"),
    ]

    // MARK: - 初始化

    /// 初始化编辑器
    /// - Parameter remind: 已有的提醒对象，传入时为编辑模式，nil 为新建模式
    init(remind: Remind? = nil) {
        self.existingRemind = remind
        self.content = remind?.content ?? ""
        self.cronExpression = remind?.cron ?? ""
        // 判断现有 cron 是否为预设
        self.isCustomCron = !Self.cronPresets.contains { $0.expression == self.cronExpression }
    }

    // MARK: - 操作

    /// 选择预设 cron
    func selectPreset(_ expression: String) {
        cronExpression = expression
        isCustomCron = false
    }

    /// 切换到自定义 cron 输入
    func switchToCustomCron() {
        isCustomCron = true
    }

    /// 保存提醒
    func save() async {
        guard !isSaveButtonDisabled else { return }
        isSaving = true

        do {
            if let remind = existingRemind {
                // 编辑模式：仅支持手动更新
                let finalCron = isCustomCron ? cronExpression.trimmingCharacters(in: .whitespaces) : cronExpression
                _ = try await remindService.update(
                    id: remind.id,
                    cron: finalCron.isEmpty ? nil : finalCron,
                    content: content,
                    status: nil
                )
            } else if createMode == .smart {
                // 新建模式 - AI 智能：AI 自动生成 cron 和内容
                _ = try await remindService.smartCreate(content: content)
            } else {
                // 新建模式 - 手动：按用户设置的 cron 创建
                let finalCron = isCustomCron ? cronExpression.trimmingCharacters(in: .whitespaces) : cronExpression
                _ = try await remindService.create(
                    cron: finalCron.isEmpty ? nil : finalCron,
                    content: content
                )
            }
            didSave = true
        } catch {
            errorMessage = error.localizedDescription
            showError = true
        }

        isSaving = false
    }
}

