import SwiftUI

/// 提醒编辑视图（支持新建和编辑）
struct RemindEditorView: View {

    /// 传入已有的提醒对象以进入编辑模式，nil 为新建模式
    var existingRemind: Remind?

    @State private var viewModel: RemindEditorViewModel
    @Environment(\.dismiss) private var dismiss

    /// 提醒内容输入框焦点（新建模式自动聚焦）
    @FocusState private var contentFocused: Bool

    init(remind: Remind? = nil) {
        self.existingRemind = remind
        _viewModel = State(initialValue: RemindEditorViewModel(remind: remind))
    }

    var body: some View {
        Form {
            // MARK: - 创建模式（仅新建模式显示）
            if !viewModel.isEditing {
                Section {
                    Picker("创建方式", selection: $viewModel.createMode) {
                        Label("AI 智能识别", systemImage: "sparkles")
                            .tag(RemindCreateMode.smart)
                        Label("手动设置", systemImage: "slider.horizontal.3")
                            .tag(RemindCreateMode.manual)
                    }
                    .pickerStyle(.segmented)
                    .labelsHidden()
                } footer: {
                    Text(viewModel.createMode == .smart
                        ? "输入提醒描述，AI 自动识别时间并生成重复规则"
                        : "手动设置 Cron 表达式或选择预设频率")
                }
            }

            // MARK: - 提醒内容
            Section {
                TextField("提醒内容", text: $viewModel.content, axis: .vertical)
                    .lineLimit(3...6)
                    .focused($contentFocused)
            } header: {
                Text("内容")
            } footer: {
                if !viewModel.isEditing && viewModel.createMode == .smart {
                    Text("如：每天早上9点提醒我喝水、每周一例会前10分钟提醒")
                }
            }

            // MARK: - Cron 表达式设置（手动模式或编辑模式）
            if viewModel.showsCronSettings {
                Section {
                    // 快捷预设按钮
                    VStack(alignment: .leading, spacing: 10) {
                        Text("快捷设置")
                            .font(.subheadline)
                            .foregroundStyle(.secondary)

                        ScrollView(.horizontal, showsIndicators: false) {
                            HStack(spacing: 8) {
                                ForEach(RemindEditorViewModel.cronPresets, id: \.expression) { preset in
                                    Button {
                                        viewModel.selectPreset(preset.expression)
                                    } label: {
                                        HStack(spacing: 4) {
                                            Text(preset.label)
                                                .font(.subheadline)
                                        }
                                        .padding(.horizontal, 14)
                                        .padding(.vertical, 8)
                                        .background(
                                            viewModel.cronExpression == preset.expression && !viewModel.isCustomCron
                                                ? Color.accentColor
                                                : Color(.systemGray5)
                                        )
                                        .foregroundStyle(
                                            viewModel.cronExpression == preset.expression && !viewModel.isCustomCron
                                                ? .white
                                                : .primary
                                        )
                                        .clipShape(Capsule())
                                    }
                                }

                                // 自定义按钮
                                Button {
                                    viewModel.switchToCustomCron()
                                } label: {
                                    HStack(spacing: 4) {
                                        Image(systemName: "slider.horizontal.3")
                                        Text("自定义")
                                            .font(.subheadline)
                                    }
                                    .padding(.horizontal, 14)
                                    .padding(.vertical, 8)
                                    .background(
                                        viewModel.isCustomCron
                                            ? Color.accentColor
                                            : Color(.systemGray5)
                                    )
                                    .foregroundStyle(
                                        viewModel.isCustomCron
                                            ? .white
                                            : .primary
                                    )
                                    .clipShape(Capsule())
                                }
                            }
                        }
                    }

                    // 自定义 cron 输入框
                    if viewModel.isCustomCron {
                        TextField("Cron 表达式", text: $viewModel.cronExpression)
                            .textInputAutocapitalization(.never)
                            .autocorrectionDisabled()
                            .font(.system(.body, design: .monospaced))
                    }

                    // cron 预览信息
                    if !viewModel.cronExpression.isEmpty {
                        HStack(spacing: 6) {
                            Image(systemName: "clock")
                                .font(.caption)
                                .foregroundStyle(.secondary)
                            Text("Cron: \(viewModel.cronExpression)")
                                .font(.caption)
                                .foregroundStyle(.secondary)
                        }
                    }
                } header: {
                    Text("重复规则")
                } footer: {
                    if viewModel.isCustomCron {
                        Text("Cron 格式：分 时 日 月 星期（如 0 9 * * * 表示每天9点）")
                    }
                }
            }
        }
        // 点击空白收起 + 拖拽下滑交互式收起键盘
        .hideKeyboardOnTap()
        .scrollDismissesKeyboard(.interactively)
        .navigationTitle(viewModel.title)
        .navigationBarTitleDisplayMode(.inline)
        .toolbar {
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
            if existingRemind == nil {
                contentFocused = true
            }
        }
    }
}
