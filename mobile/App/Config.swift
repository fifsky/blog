import Foundation

/// 应用全局配置
enum Config {

    // MARK: - 服务器地址

    /// API 基础地址
    static let baseURL = "https://api.fifsky.com"

    // MARK: - 认证接口

    /// 用户登录
    static let loginPath = "/blog/login"

    /// 用户登出
    static let logoutPath = "/blog/logout"

    // MARK: - AI 接口

    /// AI 智能创建提醒（自动生成 cron 和内容）
    static let aiRemindCreatePath = "/blog/admin/ai/remind/create"

    // MARK: - 文章接口

    /// 获取文章列表
    static let articleListPath = "/blog/article/list"

    /// 获取文章详情
    static let articleGetPath = "/blog/article/get"

    /// 创建文章
    static let articleCreatePath = "/blog/article/create"

    /// 更新文章
    static let articleUpdatePath = "/blog/article/update"

    /// 删除文章
    static let articleDeletePath = "/blog/article/delete"

    // MARK: - 心情接口

    /// 获取心情列表
    static let moodListPath = "/blog/mood/list"

    /// 创建心情
    static let moodCreatePath = "/blog/mood/create"

    /// 删除心情
    static let moodDeletePath = "/blog/mood/delete"

    // MARK: - 提醒接口

    /// 获取提醒列表
    static let reminderListPath = "/blog/reminder/list"

    /// 创建提醒
    static let reminderCreatePath = "/blog/reminder/create"

    /// 更新提醒
    static let reminderUpdatePath = "/blog/reminder/update"

    /// 删除提醒
    static let reminderDeletePath = "/blog/reminder/delete"

    // MARK: - 足迹接口

    /// 获取足迹列表
    static let footprintListPath = "/blog/footprint/list"

    /// 创建足迹
    static let footprintCreatePath = "/blog/footprint/create"

    /// 删除足迹
    static let footprintDeletePath = "/blog/footprint/delete"

    // MARK: - 文件上传

    /// 文件上传（multipart/form-data）
    static let uploadPath = "/blog/admin/upload"

    /// OSS 预签名上传
    static let ossPresignPath = "/blog/admin/oss/presign"
}
