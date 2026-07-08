import SwiftUI

/// 提醒编辑视图（支持新建和编辑）
struct RemindEditorView: View {

    /// 传入已有的提醒对象以进入编辑模式，nil 为新建模式
    var existingRemind: Remind?

    @State private var viewModel: RemindEditorViewModel
    @State private var voiceRecorder = RemindVoiceRecorder()
    @State private var isVoicePressActive = false
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
                }
            }

            // MARK: - 提醒内容
            if viewModel.showsSmartVoiceInput {
                Section {
                    smartVoiceInput
                }
            } else {
                Section {
                    TextField("提醒内容", text: $viewModel.content, axis: .vertical)
                        .lineLimit(3...6)
                        .focused($contentFocused)
                } header: {
                    Text("内容")
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
                        Text(viewModel.showsSmartVoiceInput ? "提交" : "保存")
                            .fontWeight(.medium)
                    }
                }
                .disabled(viewModel.isSaveButtonDisabled || voiceRecorder.isRecording)
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
        .onChange(of: viewModel.createMode) { _, newValue in
            if newValue == .smart {
                contentFocused = false
            } else {
                voiceRecorder.reset()
                contentFocused = true
            }
        }
        .onAppear {
            // 仅文本输入场景自动聚焦，语音模式避免直接拉起键盘
            if existingRemind == nil && !viewModel.showsSmartVoiceInput {
                contentFocused = true
            }
        }
        .onDisappear {
            voiceRecorder.reset()
        }
    }

    /// AI 智能模式下的语音录入面板
    private var smartVoiceInput: some View {
        VStack(spacing: 18) {
            if !viewModel.content.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty {
                recognizedTextCard
            } else {
                voiceEmptyState
            }

            voiceRecordingStage

            if let errorMessage = voiceRecorder.errorMessage, !errorMessage.isEmpty {
                Label(errorMessage, systemImage: "exclamationmark.triangle")
                    .font(.caption)
                    .foregroundStyle(.red)
                    .frame(maxWidth: .infinity, alignment: .leading)
            }
        }
        .padding(.vertical, 12)
    }

    /// 固定高度的录音舞台，避免波形出现时挤动按钮位置
    private var voiceRecordingStage: some View {
        ZStack {
            VoiceWaveView(level: voiceRecorder.audioLevel)
                .frame(height: 38)
                .opacity(voiceRecorder.isRecording ? 1 : 0)
                .scaleEffect(voiceRecorder.isRecording ? 1 : 0.92)
                .offset(y: -78)

            voiceRecordButton
        }
        .frame(maxWidth: .infinity)
        .frame(height: 162)
        .animation(.spring(response: 0.28, dampingFraction: 0.86), value: voiceRecorder.isRecording)
    }

    /// 未识别文字时的空状态
    private var voiceEmptyState: some View {
        VStack(spacing: 8) {
            Text(voiceRecorder.isRecording ? "正在聆听..." : "按住说出提醒")
                .font(.headline)
            Text(voiceRecorder.isRecording ? "松开后自动转成文字" : "例如：明天上午九点提醒我喝水")
                .font(.footnote)
                .foregroundStyle(.secondary)
        }
        .frame(maxWidth: .infinity)
        .padding(.top, 4)
    }

    /// 语音识别后的文字确认卡片
    private var recognizedTextCard: some View {
        HStack(alignment: .top, spacing: 12) {
            Text(viewModel.content)
                .font(.body)
                .foregroundStyle(.primary)
                .frame(maxWidth: .infinity, alignment: .leading)
                .textSelection(.enabled)

            Button {
                voiceRecorder.reset()
                viewModel.resetSmartDraft()
            } label: {
                Image(systemName: "xmark.circle.fill")
                    .font(.title3)
                    .foregroundStyle(.secondary)
            }
            .buttonStyle(.plain)
            .accessibilityLabel("清空并重来")
        }
        .padding(14)
        .background(Color(.secondarySystemGroupedBackground))
        .clipShape(RoundedRectangle(cornerRadius: 14, style: .continuous))
    }

    /// 按住录音按钮
    private var voiceRecordButton: some View {
        VStack(spacing: 10) {
            ZStack {
                Circle()
                    .fill(voiceButtonBackground)
                    .frame(width: 84, height: 84)
                    .scaleEffect(voiceRecorder.isRecording ? 1.04 : 1)
                    .animation(.spring(response: 0.25, dampingFraction: 0.72), value: voiceRecorder.isRecording)

                if viewModel.isTranscribing {
                    ProgressView()
                        .controlSize(.large)
                } else {
                    Image(systemName: voiceRecorder.isRecording ? "waveform" : "mic.fill")
                        .font(.system(size: 32, weight: .semibold))
                        .foregroundStyle(voiceRecorder.isRecording ? .red : Color.accentColor)
                }
            }
            .contentShape(Circle())
            .gesture(
                DragGesture(minimumDistance: 0)
                    .onChanged { _ in beginVoicePress() }
                    .onEnded { _ in endVoicePress() }
            )
            .opacity(viewModel.isTranscribing ? 0.72 : 1)
            .allowsHitTesting(!viewModel.isTranscribing)

            Text(voiceButtonText)
                .font(.subheadline)
                .foregroundStyle(.secondary)
        }
        .frame(maxWidth: .infinity)
    }

    private var voiceButtonBackground: Color {
        if viewModel.isTranscribing {
            return Color(.systemGray5)
        }
        return voiceRecorder.isRecording ? Color.red.opacity(0.12) : Color.accentColor.opacity(0.12)
    }

    private var voiceButtonText: String {
        if viewModel.isTranscribing {
            return "正在识别..."
        }
        return voiceRecorder.isRecording ? "松开识别" : "按住说话"
    }

    private func beginVoicePress() {
        guard !isVoicePressActive && !viewModel.isTranscribing else { return }
        isVoicePressActive = true

        Task {
            let started = await voiceRecorder.startRecording()
            if started && !isVoicePressActive {
                await finishVoiceRecording()
            }
        }
    }

    private func endVoicePress() {
        guard isVoicePressActive else { return }
        isVoicePressActive = false

        Task {
            await finishVoiceRecording()
        }
    }

    private func finishVoiceRecording() async {
        guard let audioBase64 = voiceRecorder.stopRecordingBase64() else { return }
        await viewModel.transcribeSpeech(audioBase64: audioBase64)
        voiceRecorder.deleteRecordingFile()
    }
}

/// 录音时的声纹波形
private struct VoiceWaveView: View {
    let level: Double

    private let barCount = 18

    var body: some View {
        TimelineView(.animation(minimumInterval: 0.06)) { context in
            let phase = context.date.timeIntervalSinceReferenceDate * 7
            HStack(alignment: .center, spacing: 4) {
                ForEach(0..<barCount, id: \.self) { index in
                    Capsule()
                        .fill(Color.accentColor.opacity(0.82))
                        .frame(width: 4, height: barHeight(index: index, phase: phase))
                }
            }
            .frame(maxWidth: .infinity)
        }
    }

    private func barHeight(index: Int, phase: Double) -> CGFloat {
        let wave = (sin(phase + Double(index) * 0.62) + 1) / 2
        let base = 8 + wave * 28
        return CGFloat(base * max(0.35, min(level + 0.2, 1.1)))
    }
}
