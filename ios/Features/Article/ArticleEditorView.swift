import SwiftUI
import PhotosUI

/// 文章编辑器视图
/// 支持创建新文章和编辑已有文章
struct ArticleEditorView: View {

    /// 要编辑的文章（nil 表示创建新文章）
    var article: Article?

    @State private var viewModel: ArticleEditorViewModel?

    /// 分类选择器是否展开
    @State private var showCategoryPicker = false

    /// 图片选择器
    @State private var selectedItem: PhotosPickerItem?

    /// 导航返回标记
    @Environment(\.dismiss) private var dismiss

    var body: some View {
        Form {
            // MARK: - 标题
            Section {
                TextField("文章标题", text: Binding(
                    get: { viewModel?.title ?? "" },
                    set: { viewModel?.title = $0 }
                ))
            }

            // MARK: - 分类
            Section {
                Button {
                    showCategoryPicker = true
                } label: {
                    HStack {
                        Text("分类")
                            .foregroundStyle(.primary)
                        Spacer()
                        Text(viewModel?.cateName ?? "请选择分类")
                            .foregroundStyle(
                                (viewModel?.cateId ?? 0) > 0 ? .primary : .secondary
                            )
                        Image(systemName: "chevron.right")
                            .font(.caption)
                            .foregroundStyle(.secondary)
                    }
                }
            }

            // MARK: - 标签
            Section {
                TextField("标签（逗号分隔）", text: Binding(
                    get: { viewModel?.tagsText ?? "" },
                    set: { viewModel?.tagsText = $0 }
                ))
            } header: {
                Text("标签")
            } footer: {
                Text("多个标签用英文逗号分隔，例如：SwiftUI, iOS, 开发")
            }

            // MARK: - 正文
            Section {
                TextEditor(text: Binding(
                    get: { viewModel?.content ?? "" },
                    set: { viewModel?.content = $0 }
                ))
                .frame(minHeight: 200)

                // Markdown 工具栏
                markdownToolbar
            } header: {
                Text("正文内容")
            } footer: {
                Text("支持 Markdown 格式")
            }
        }
        .navigationTitle(viewModel?.isEditing == true ? "编辑文章" : "新建文章")
        .navigationBarTitleDisplayMode(.inline)
        .toolbar {
            // 保存按钮
            ToolbarItem(placement: .topBarTrailing) {
                Button {
                    Task {
                        await viewModel?.save()
                        // 保存成功后返回上一页
                        if viewModel?.isSaved == true {
                            dismiss()
                        }
                    }
                } label: {
                    if viewModel?.isSaving == true {
                        ProgressView()
                    } else {
                        Text("保存")
                            .fontWeight(.semibold)
                    }
                }
                .disabled(viewModel?.isSaving == true)
            }

            // 取消按钮（编辑模式）
            if viewModel?.isEditing == true {
                ToolbarItem(placement: .cancellationAction) {
                    Button("取消") {
                        dismiss()
                    }
                }
            }
        }
        .alert("错误", isPresented: Binding(
            get: { viewModel?.showError ?? false },
            set: { if !$0 { viewModel?.showError = false } }
        )) {
            Button("确定", role: .cancel) {}
        } message: {
            Text(viewModel?.errorMessage ?? "未知错误")
        }
        .sheet(isPresented: $showCategoryPicker) {
            if let categories = viewModel?.categories {
                CategoryPickerView(
                    categories: categories,
                    selectedId: viewModel?.cateId ?? 0,
                    onSelect: { cate in
                        viewModel?.cateId = extractCateId(from: cate.url)
                        viewModel?.cateName = cate.content
                    }
                )
                .presentationDetents([.medium, .large])
            }
        }
        .onChange(of: selectedItem) { _, newItem in
            if let newItem {
                Task {
                    await viewModel?.uploadAndInsertImage(item: newItem)
                    selectedItem = nil
                }
            }
        }
        .onAppear {
            if viewModel == nil {
                viewModel = ArticleEditorViewModel(article: article)
            }
            Task {
                await viewModel?.loadCategories()
            }
        }
    }

    // MARK: - 子视图

    /// Markdown 编辑工具栏
    private var markdownToolbar: some View {
        ScrollView(.horizontal, showsIndicators: false) {
            HStack(spacing: 12) {
                // 粗体
                Button {
                    viewModel?.insertMarkdown(syntax: ("**", "**"))
                } label: {
                    Text("B")
                        .font(.system(.body, design: .rounded).bold())
                        .frame(width: 32, height: 32)
                }
                .buttonStyle(.bordered)

                // 斜体
                Button {
                    viewModel?.insertMarkdown(syntax: ("*", "*"))
                } label: {
                    Text("I")
                        .font(.system(.body, design: .rounded).italic())
                        .frame(width: 32, height: 32)
                }
                .buttonStyle(.bordered)

                // 行内代码
                Button {
                    viewModel?.insertCode()
                } label: {
                    Image(systemName: "chevron.left.forwardslash.chevron.right")
                        .font(.caption)
                        .frame(width: 32, height: 32)
                }
                .buttonStyle(.bordered)

                // 链接
                Button {
                    viewModel?.insertLink()
                } label: {
                    Image(systemName: "link")
                        .font(.caption)
                        .frame(width: 32, height: 32)
                }
                .buttonStyle(.bordered)

                // 图片上传
                PhotosPicker(selection: $selectedItem, matching: .images) {
                    if viewModel?.isUploadingImage == true {
                        ProgressView()
                            .frame(width: 32, height: 32)
                    } else {
                        Image(systemName: "photo")
                            .font(.caption)
                            .frame(width: 32, height: 32)
                    }
                }
                .buttonStyle(.bordered)
            }
            .padding(.vertical, 4)
        }
    }

    // MARK: - 辅助方法

    /// 从分类 URL 中提取分类 ID
    /// URL 格式通常为 /cate/{id} 或包含分类 ID 的路径
    private func extractCateId(from url: String) -> Int {
        // 从 URL 路径中提取最后一段作为 ID
        let components = url.split(separator: "/")
        if let lastComponent = components.last, let id = Int(lastComponent) {
            return id
        }
        // 尝试从查询参数中提取
        if let urlComponents = URLComponents(string: url),
           let idStr = urlComponents.queryItems?.first(where: { $0.name == "id" })?.value,
           let id = Int(idStr) {
            return id
        }
        return 0
    }
}

// MARK: - 分类选择器

/// 底部弹出的分类选择列表
/// 采用标准 iOS 选择列表样式，选中项带勾选标记
struct CategoryPickerView: View {

    let categories: [CateMenuItem]
    let selectedId: Int
    let onSelect: (CateMenuItem) -> Void

    @Environment(\.dismiss) private var dismiss

    /// 从分类 URL 中提取分类 ID
    private func cateId(from url: String) -> Int {
        let components = url.split(separator: "/")
        if let lastComponent = components.last, let id = Int(lastComponent) {
            return id
        }
        return 0
    }

    var body: some View {
        NavigationStack {
            List {
                ForEach(Array(categories.enumerated()), id: \.offset) { _, cate in
                    let id = cateId(from: cate.url)
                    Button {
                        onSelect(cate)
                        dismiss()
                    } label: {
                        HStack {
                            Text(cate.content)
                                .foregroundStyle(.primary)
                            Spacer()
                            if id == selectedId {
                                Image(systemName: "checkmark")
                                    .foregroundStyle(Color.accentColor)
                            }
                        }
                    }
                }
            }
            .navigationTitle("选择分类")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .cancellationAction) {
                    Button("取消") {
                        dismiss()
                    }
                }
            }
        }
    }
}
