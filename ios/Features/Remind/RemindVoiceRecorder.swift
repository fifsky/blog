import AVFoundation
import Foundation

/// 提醒语音录音器，负责采集音频并导出 base64 数据
@MainActor
@Observable
final class RemindVoiceRecorder {

    /// 是否正在录音
    var isRecording = false

    /// 当前录音音量，用于驱动声纹波形
    var audioLevel = 0.05

    /// 录音错误信息
    var errorMessage: String?

    private var recorder: AVAudioRecorder?
    private var recordingURL: URL?
    private var meterTimer: Timer?

    /// 开始录音
    @discardableResult
    func startRecording() async -> Bool {
        guard !isRecording else { return true }

        errorMessage = nil
        let granted = await AVCaptureDevice.requestAccess(for: .audio)
        guard granted else {
            errorMessage = "需要允许麦克风权限才能语音创建提醒"
            return false
        }

        do {
            reset(deleteRecording: true)

            let audioSession = AVAudioSession.sharedInstance()
            try audioSession.setCategory(
                .playAndRecord,
                mode: .spokenAudio,
                options: [.defaultToSpeaker, .allowBluetoothHFP]
            )
            try audioSession.setActive(true)

            let url = FileManager.default.temporaryDirectory
                .appendingPathComponent("remind-\(UUID().uuidString).m4a")
            let settings: [String: Any] = [
                AVFormatIDKey: Int(kAudioFormatMPEG4AAC),
                AVSampleRateKey: 44_100,
                AVNumberOfChannelsKey: 1,
                AVEncoderAudioQualityKey: AVAudioQuality.high.rawValue,
            ]

            let recorder = try AVAudioRecorder(url: url, settings: settings)
            recorder.isMeteringEnabled = true
            guard recorder.record() else {
                throw RecordingError.startFailed
            }

            self.recorder = recorder
            self.recordingURL = url
            self.isRecording = true
            startMetering()
            return true
        } catch {
            reset(deleteRecording: true)
            errorMessage = "无法开始录音：\(error.localizedDescription)"
            return false
        }
    }

    /// 停止录音并返回 base64 音频
    func stopRecordingBase64() -> String? {
        guard isRecording else { return nil }

        recorder?.stop()
        stopMetering()
        isRecording = false
        audioLevel = 0.05
        try? AVAudioSession.sharedInstance().setActive(false, options: .notifyOthersOnDeactivation)

        guard let recordingURL else {
            errorMessage = "录音文件不存在"
            return nil
        }

        do {
            let data = try Data(contentsOf: recordingURL)
            guard !data.isEmpty else {
                errorMessage = "录音为空，请重新录音"
                return nil
            }
            return data.base64EncodedString()
        } catch {
            errorMessage = "读取录音失败：\(error.localizedDescription)"
            return nil
        }
    }

    /// 清理录音状态
    func reset(deleteRecording: Bool = true) {
        if recorder?.isRecording == true {
            recorder?.stop()
        }
        stopMetering()
        isRecording = false
        audioLevel = 0.05
        recorder = nil
        try? AVAudioSession.sharedInstance().setActive(false, options: .notifyOthersOnDeactivation)

        if deleteRecording, let recordingURL {
            try? FileManager.default.removeItem(at: recordingURL)
            self.recordingURL = nil
        }
        errorMessage = nil
    }

    /// 删除最近一次录音文件
    func deleteRecordingFile() {
        if let recordingURL {
            try? FileManager.default.removeItem(at: recordingURL)
        }
        recordingURL = nil
    }

    private func startMetering() {
        stopMetering()
        meterTimer = Timer.scheduledTimer(withTimeInterval: 0.06, repeats: true) { [weak self] _ in
            Task { @MainActor [weak self] in
                guard let self, let recorder = self.recorder, recorder.isRecording else { return }
                recorder.updateMeters()
                let power = recorder.averagePower(forChannel: 0)
                let normalized = max(0.05, min(1, (Double(power) + 55) / 55))
                self.audioLevel = normalized
            }
        }
    }

    private func stopMetering() {
        meterTimer?.invalidate()
        meterTimer = nil
    }
}

private enum RecordingError: LocalizedError {
    case startFailed

    var errorDescription: String? {
        switch self {
        case .startFailed:
            return "录音设备没有成功启动"
        }
    }
}
