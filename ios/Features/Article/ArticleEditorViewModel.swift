import SwiftUI
import PhotosUI

/// 文章编辑器视图模型
@Observable
class ArticleEditorViewModel {

    // MARK: - 状态

    /// 文章标题
    var title = ""

    /// 文章分类 ID
    var cateId: Int = 0

    /// 文章分类名称（用于显示）
    var cateName: String = "请选择分类"

    /// 文章标签（逗号分隔的字符串）
    var tagsText = ""

    /// 文章内容
    var content = ""

    /// 是否正在保存
    var isSaving = false

    /// 是否正在加载分类列表
    var isLoadingCategories = false

    /// 错误信息
    var errorMessage: String?

    /// 是否显示错误弹窗
    var showError = false

    /// 保存成功标记
    var isSaved = false

    /// 是否正在上传图片
    var isUploadingImage = false

    // MARK: - 编辑模式

    /// 是否为编辑模式（有文章则为编辑）
    let isEditing: Bool

    /// 正在编辑的文章（编辑模式）
    private let article: Article?

    // MARK: - 私有属性

    private let articleService = ArticleService.shared
    private let categoryService = CategoryService.shared
    private let uploadService = UploadService.shared

    /// 分类列表
    var categories: [CateMenuItem] = []

    // MARK: - 初始化

    init(article: Article? = nil) {
        self.article = article
        self.isEditing = article != nil

        // 编辑模式：填充现有数据
        if let article {
            self.title = article.title
            self.cateId = article.cate_id
            self.cateName = article.cate?.name ?? "未分类"
            self.tagsText = article.tags.joined(separator: ", ")
            self.content = article.content
        }
    }

    // MARK: - 数据加载

    /// 加载分类列表
    func loadCategories() async {
        guard !isLoadingCategories else { return }
        isLoadingCategories = true

        do {
            let response = try await categoryService.all()
            categories = response.list
        } catch {
            errorMessage = "加载分类失败：\(error.localizedDescription)"
            showError = true
        }

        isLoadingCategories = false
    }

    // MARK: - 保存

    /// 保存文章（创建或更新）
    func save() async {
        guard !isSaving else { return }

        // 验证
        guard !title.trimmingCharacters(in: .whitespaces).isEmpty else {
            errorMessage = "请输入文章标题"
            showError = true
            return
        }
        guard cateId > 0 else {
            errorMessage = "请选择文章分类"
            showError = true
            return
        }
        guard !content.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty else {
            errorMessage = "请输入文章内容"
            showError = true
            return
        }

        isSaving = true

        // 解析标签
        let tags = tagsText
            .split(separator: ",")
            .map { $0.trimmingCharacters(in: .whitespaces) }
            .filter { !$0.isEmpty }

        do {
            if isEditing, let article {
                // 更新文章
                _ = try await articleService.update(
                    id: article.id,
                    cateId: cateId,
                    title: title,
                    content: content,
                    tags: tags.isEmpty ? nil : tags
                )
            } else {
                // 创建文章
                _ = try await articleService.create(
                    cateId: cateId,
                    type: 0,
                    title: title,
                    content: content,
                    tags: tags.isEmpty ? nil : tags
                )
            }
            isSaved = true
        } catch {
            errorMessage = error.localizedDescription
            showError = true
        }

        isSaving = false
    }

    // MARK: - 图片上传

    /// 上传图片并插入 Markdown 格式
    /// - Parameter item: 从 PhotosPicker 选中的图片
    func uploadAndInsertImage(item: PhotosPickerItem) async {
        guard let data = try? await item.loadTransferable(type: Data.self) else {
            errorMessage = "无法读取选中的图片"
            showError = true
            return
        }

        isUploadingImage = true

        do {
            let filename = "image_\(Int(Date().timeIntervalSince1970)).png"
            let imageUrl = try await uploadService.uploadImage(imageData: data, filename: filename)
            // 在光标位置插入 Markdown 图片语法
            let imageMarkdown = "\n![image](\(imageUrl))\n"
            content += imageMarkdown
        } catch {
            errorMessage = "图片上传失败：\(error.localizedDescription)"
            showError = true
        }

        isUploadingImage = false
    }

    // MARK: - Markdown 工具栏操作

    /// 在内容中插入 Markdown 语法
    /// - Parameter syntax: 要插入的语法包裹对，如 ("**", "**") 表示粗体
    func insertMarkdown(syntax: (prefix: String, suffix: String)) {
        content += "\(syntax.prefix)文本\(syntax.suffix)"
    }

    /// 插入链接语法
    func insertLink() {
        content += "[链接文字](https://)"
    }

    /// 插入行内代码
    func insertCode() {
        content += "`代码`"
    }
}
