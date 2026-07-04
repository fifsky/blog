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

    /// 管理端文章列表
    static let adminArticleListPath = "/blog/admin/article/list"

    /// 管理端文章详情
    static let adminArticleDetailPath = "/blog/admin/article/detail"

    /// 管理端创建文章
    static let adminArticleCreatePath = "/blog/admin/article/create"

    /// 管理端更新文章
    static let adminArticleUpdatePath = "/blog/admin/article/update"

    /// 管理端删除文章
    static let adminArticleDeletePath = "/blog/admin/article/delete"

    // MARK: - 分类接口

    /// 公开分类列表
    static let cateAllPath = "/blog/cate/all"

    /// 管理端分类列表
    static let adminCateListPath = "/blog/admin/cate/list"

    // MARK: - 评论接口

    /// 公开评论列表
    static let commentListPath = "/blog/comment/list"

    /// 公开创建评论
    static let commentCreatePath = "/blog/comment/create"

    /// 管理端评论列表
    static let adminCommentListPath = "/blog/admin/comment/list"

    /// 管理端删除评论
    static let adminCommentDeletePath = "/blog/admin/comment/delete"

    // MARK: - 心情接口

    /// 公开心情列表
    static let moodListPath = "/blog/mood/list"

    /// 管理端创建心情
    static let adminMoodCreatePath = "/blog/admin/mood/create"

    /// 管理端更新心情
    static let adminMoodUpdatePath = "/blog/admin/mood/update"

    /// 管理端删除心情
    static let adminMoodDeletePath = "/blog/admin/mood/delete"

    // MARK: - 提醒接口

    /// 管理端提醒列表
    static let adminRemindListPath = "/blog/admin/remind/list"

    /// 管理端创建提醒
    static let adminRemindCreatePath = "/blog/admin/remind/create"

    /// 管理端更新提醒
    static let adminRemindUpdatePath = "/blog/admin/remind/update"

    /// 管理端删除提醒
    static let adminRemindDeletePath = "/blog/admin/remind/delete"

    // MARK: - 足迹接口

    /// 管理端足迹列表
    static let adminFootprintListPath = "/blog/admin/footprint/list"

    /// 公开全部足迹
    static let footprintAllPath = "/blog/travel/footprints"

    /// 管理端创建足迹
    static let adminFootprintCreatePath = "/blog/admin/footprint/create"

    /// 管理端更新足迹
    static let adminFootprintUpdatePath = "/blog/admin/footprint/update"

    /// 管理端删除足迹
    static let adminFootprintDeletePath = "/blog/admin/footprint/delete"

    // MARK: - 文件上传

    /// 文件上传（multipart/form-data）
    static let uploadPath = "/blog/admin/upload"

    /// OSS 预签名上传
    static let ossPresignPath = "/blog/admin/oss/presign"
}
