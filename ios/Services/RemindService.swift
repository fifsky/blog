import Foundation

// RemindService 提醒服务，处理定时提醒的管理接口
class RemindService {
    static let shared = RemindService()

    private init() {}

    // 获取提醒列表
    func list(page: Int) async throws -> RemindListResponse {
        let request = RemindListRequest(page: page)
        return try await APIClient.shared.request(path: Config.adminRemindListPath, body: request, auth: true)
    }

    // 创建提醒
    func create(cron: String?, content: String) async throws -> IDResponse {
        let request = RemindCreateRequest(cron: cron, content: content)
        return try await APIClient.shared.request(path: Config.adminRemindCreatePath, body: request, auth: true)
    }

    // AI 智能创建提醒（AI 根据自然语言描述自动生成 cron 表达式和提醒内容）
    func smartCreate(content: String) async throws -> IDResponse {
        let request = RemindSmartCreateRequest(content: content)
        return try await APIClient.shared.request(path: Config.aiRemindCreatePath, body: request, auth: true)
    }

    // 提醒语音转文字
    func transcribeSpeech(audioBase64: String) async throws -> RemindSpeechTranscribeResponse {
        let request = RemindSpeechTranscribeRequest(audio_base64: audioBase64)
        return try await APIClient.shared.request(path: Config.aiRemindTranscribePath, body: request, auth: true)
    }

    // 更新提醒
    func update(id: Int, cron: String?, content: String?, status: String?) async throws -> IDResponse {
        // 将字符串状态映射为枚举，便于以 proto 枚举值编码
        let statusEnum = status.flatMap { RemindStatus(rawValue: $0.uppercased()) }
        let request = RemindUpdateRequest(id: id, cron: cron, content: content, status: statusEnum)
        return try await APIClient.shared.request(path: Config.adminRemindUpdatePath, body: request, auth: true)
    }

    // 删除提醒
    func delete(id: Int) async throws {
        let request = RemindDeleteRequest(id: id)
        // 返回 google.protobuf.Empty，只需要检查不抛错即可
        let _: EmptyResponse = try await APIClient.shared.request(path: Config.adminRemindDeletePath, body: request, auth: true)
    }
}
