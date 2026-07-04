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

    /// 标题输入框焦点（新建模式自动聚焦）
    @FocusState private var titleFocused: Bool

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
                .focused($titleFocused)
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
                // Form 内 Button 必须显式指定 buttonStyle，
                // 否则复杂 label 会导致整行命中区域被 List/Form 吞掉，表现为点击无反应
                .buttonStyle(.borderless)
                // .borderless 会用 accentColor 渲染 label，覆盖 .primary；
                // 显式设主题黑，与表单其它文字一致
                .tint(Color.themePrimary)
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
        // 点击空白收起 + 拖拽下滑交互式收起键盘
        .hideKeyboardOnTap()
        .scrollDismissesKeyboard(.interactively)
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
                        viewModel?.cateId = cate.id
                        viewModel?.cateName = cate.name
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
            // 新建模式自动聚焦标题输入框
            if article == nil {
                titleFocused = true
            }
        }
    }

    // MARK: - 子视图

    /// Markdown 编辑工具栏
    private var markdownToolbar: some View {
        let isUploadingImage = viewModel?.isUploadingImage == true

        return ScrollView(.horizontal, showsIndicators: false) {
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
                    if isUploadingImage {
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
            .foregroundStyle(.primary)
            .padding(.vertical, 4)
        }
    }

    // MARK: - 子视图（无需 URL 解析，分类项自带数字 ID）
}

// MARK: - 分类选择器

/// 底部弹出的分类选择列表
/// 采用标准 iOS 选择列表样式，选中项带勾选标记
struct CategoryPickerView: View {

    let categories: [CateItem]
    let selectedId: Int
    let onSelect: (CateItem) -> Void

    @Environment(\.dismiss) private var dismiss

    var body: some View {
        NavigationStack {
            List {
                ForEach(categories) { cate in
                    Button {
                        onSelect(cate)
                        dismiss()
                    } label: {
                        HStack {
                            Text(cate.name)
                                .foregroundStyle(.primary)
                            Spacer()
                            if cate.id == selectedId {
                                Image(systemName: "checkmark")
                                    .foregroundStyle(Color.themePrimary)
                            }
                        }
                    }
                    // List 内 Button 需显式 buttonStyle，否则点击无反应；
                    // .borderless + .tint(themePrimary) 避免文字被 accentColor 染蓝
                    .buttonStyle(.borderless)
                    .tint(Color.themePrimary)
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
            // sheet 弹出，不继承主 NavigationStack 的 tint，需显式设主题黑
            .tint(Color.themePrimary)
        }
    }
}
