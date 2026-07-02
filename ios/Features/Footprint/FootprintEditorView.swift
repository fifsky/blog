import SwiftUI
import PhotosUI
import MapKit

/// 足迹编辑器视图
/// 支持新建和编辑足迹，包含名称、描述、日期、标记颜色、位置选择和照片上传
struct FootprintEditorView: View {

    @State private var viewModel: FootprintEditorViewModel

    /// 名称输入框焦点（新建模式自动聚焦）
    @FocusState private var nameFocused: Bool

    /// 保存成功的回调
    var onSave: () -> Void

    /// 取消的回调
    var onCancel: () -> Void

    /// 初始化视图
    /// - Parameters:
    ///   - footprint: 可选，传入时进入编辑模式
    ///   - onSave: 保存成功回调
    ///   - onCancel: 取消回调
    init(footprint: Footprint? = nil, onSave: @escaping () -> Void, onCancel: @escaping () -> Void) {
        _viewModel = State(initialValue: FootprintEditorViewModel(footprint: footprint))
        self.onSave = onSave
        self.onCancel = onCancel
    }

    var body: some View {
        Form {
            // MARK: - 基本信息
            Section {
                // 地点名称（必填）
                TextField("地点名称", text: $viewModel.name)
                    .font(.body)
                    .focused($nameFocused)

                // 描述（可选）
                TextEditor(text: $viewModel.descriptionText)
                    .font(.body)
                    .frame(minHeight: 80, maxHeight: 150)
                    .scrollContentBackground(.hidden)
                    .overlay(alignment: .topLeading) {
                        // 占位提示文字
                        if viewModel.descriptionText.isEmpty {
                            Text("描述（可选）")
                                .font(.body)
                                .foregroundStyle(Color(.placeholderText))
                                .allowsHitTesting(false)
                                .padding(.top, 8)
                        }
                    }
            } header: {
                Text("基本信息")
            }

            // MARK: - 日期
            Section {
                DatePicker(
                    "日期",
                    selection: Binding(
                        get: {
                            // 从日期字符串解析为 Date
                            let formatter = DateFormatter()
                            formatter.dateFormat = "yyyy-MM-dd"
                            return formatter.date(from: viewModel.dateString) ?? Date()
                        },
                        set: { newDate in
                            // 将 Date 转换为字符串
                            let formatter = DateFormatter()
                            formatter.dateFormat = "yyyy-MM-dd"
                            viewModel.dateString = formatter.string(from: newDate)
                        }
                    ),
                    displayedComponents: .date
                )
            } header: {
                Text("日期")
            }

            // MARK: - 标记颜色
            Section {
                colorPickerGrid
            } header: {
                Text("标记颜色")
            }

            // MARK: - 位置
            Section {
                // 位置显示
                HStack {
                    Image(systemName: "location")
                        .foregroundStyle(.secondary)
                    Text(viewModel.locationDescription)
                        .foregroundStyle(
                            viewModel.latitude.isEmpty ? .tertiary : .primary
                        )
                    Spacer()
                }

                // 选择位置按钮
                Button {
                    viewModel.showLocationPicker = true
                } label: {
                    Label(
                        viewModel.latitude.isEmpty ? "选择位置" : "更改位置",
                        systemImage: "map"
                    )
                }
            } header: {
                Text("位置")
            }

            // MARK: - 照片
            Section {
                // 照片选择器
                PhotosPicker(
                    selection: $viewModel.selectedPhotoItems,
                    maxSelectionCount: 9,
                    matching: .images,
                    photoLibrary: .shared()
                ) {
                    HStack {
                        Image(systemName: "photo.on.rectangle.angled")
                        Text("选择照片（最多9张）")
                        Spacer()
                        if viewModel.selectedPhotoItems.count > 0 {
                            Text("\(viewModel.selectedPhotoItems.count)")
                                .foregroundStyle(.secondary)
                        }
                    }
                }

                // 已选照片预览
                if !viewModel.selectedPhotoItems.isEmpty {
                    selectedPhotosPreview
                }

                // 已有照片（编辑模式）
                if !viewModel.uploadedPhotoURLs.isEmpty {
                    existingPhotosPreview
                }
            } header: {
                Text("照片")
            } footer: {
                Text("新建足迹时，照片将在保存时上传。")
            }
        }
        .navigationTitle(viewModel.isEditing ? "编辑足迹" : "新建足迹")
        .navigationBarTitleDisplayMode(.inline)
        .toolbar {
            ToolbarItem(placement: .topBarLeading) {
                Button("取消") {
                    onCancel()
                }
            }

            ToolbarItem(placement: .topBarTrailing) {
                Button("保存") {
                    Task { await viewModel.save() }
                }
                .disabled(!viewModel.canSave || viewModel.isSaving)
                .fontWeight(.semibold)
            }
        }
        .onChange(of: viewModel.selectedPhotoItems) {
            Task { await viewModel.loadSelectedPhotos() }
        }
        .onChange(of: viewModel.isSaveSuccess) {
            if viewModel.isSaveSuccess {
                onSave()
            }
        }
        .alert("错误", isPresented: $viewModel.showError) {
            Button("确定", role: .cancel) {}
        } message: {
            Text(viewModel.errorMessage ?? "未知错误")
        }
        .sheet(isPresented: $viewModel.showLocationPicker) {
            NavigationStack {
                LocationPickerView(
                    latitude: $viewModel.latitude,
                    longitude: $viewModel.longitude
                )
            }
        }
        .overlay {
            // 保存/上传中的遮罩
            if viewModel.isSaving || viewModel.isUploading {
                savingOverlay
            }
        }
        .onAppear {
            // 新建模式自动聚焦名称输入框
            if !viewModel.isEditing {
                nameFocused = true
            }
        }
    }

    // MARK: - 颜色选择器

    /// 标记颜色网格选择器
    private var colorPickerGrid: some View {
        LazyVGrid(
            columns: Array(repeating: GridItem(.flexible(), spacing: 12), count: 4),
            spacing: 12
        ) {
            ForEach(FootprintEditorViewModel.markerColors, id: \.hex) { colorInfo in
                colorButton(name: colorInfo.name, hex: colorInfo.hex)
            }
        }
        .padding(.vertical, 4)
    }

    /// 单个颜色选择按钮
    private func colorButton(name: String, hex: String) -> some View {
        Button {
            viewModel.markerColor = hex
        } label: {
            VStack(spacing: 4) {
                Circle()
                    .fill(Color(hex: hex) ?? .blue)
                    .frame(width: 36, height: 36)
                    .overlay(
                        viewModel.markerColor == hex
                            ? Circle()
                                .stroke(.white, lineWidth: 2)
                            : nil
                    )
                    .overlay(
                        viewModel.markerColor == hex
                            ? Image(systemName: "checkmark")
                                .font(.caption)
                                .foregroundStyle(.white)
                            : nil
                    )

                Text(name)
                    .font(.caption2)
                    .foregroundStyle(.secondary)
            }
        }
        .buttonStyle(.plain)
    }

    // MARK: - 照片预览

    /// 已选新照片的缩略图预览
    private var selectedPhotosPreview: some View {
        ScrollView(.horizontal, showsIndicators: false) {
            HStack(spacing: 8) {
                ForEach(Array(viewModel.photoPreviews.enumerated()), id: \.offset) { index, data in
                    if let image = UIImage(data: data) {
                        Image(uiImage: image)
                            .resizable()
                            .aspectRatio(contentMode: .fill)
                            .frame(width: 80, height: 80)
                            .clipShape(RoundedRectangle(cornerRadius: 8))
                            .overlay(alignment: .topTrailing) {
                                // 删除按钮
                                Button {
                                    viewModel.removeSelectedPhoto(at: index)
                                } label: {
                                    Image(systemName: "xmark.circle.fill")
                                        .font(.title3)
                                        .foregroundStyle(.white, .black.opacity(0.6))
                                }
                            }
                    }
                }
            }
        }
    }

    /// 已有照片的缩略图（编辑模式，从服务器加载）
    private var existingPhotosPreview: some View {
        ScrollView(.horizontal, showsIndicators: false) {
            HStack(spacing: 8) {
                ForEach(viewModel.uploadedPhotoURLs, id: \.self) { urlString in
                    AsyncImage(url: URL(string: urlString)) { phase in
                        switch phase {
                        case .success(let image):
                            image
                                .resizable()
                                .aspectRatio(contentMode: .fill)
                        case .failure:
                            Image(systemName: "photo")
                                .foregroundStyle(.gray)
                        default:
                            ProgressView()
                        }
                    }
                    .frame(width: 80, height: 80)
                    .clipShape(RoundedRectangle(cornerRadius: 8))
                }
            }
        }
    }

    // MARK: - 保存遮罩

    /// 保存/上传中的半透明遮罩
    private var savingOverlay: some View {
        ZStack {
            Color.black.opacity(0.3)
                .ignoresSafeArea()

            VStack(spacing: 16) {
                ProgressView()
                    .controlSize(.large)

                if viewModel.isUploading {
                    Text(viewModel.uploadProgress)
                        .font(.subheadline)
                        .foregroundStyle(.white)
                } else {
                    Text("保存中...")
                        .font(.subheadline)
                        .foregroundStyle(.white)
                }
            }
            .padding(24)
            .background(.ultraThinMaterial, in: RoundedRectangle(cornerRadius: 16))
        }
    }
}
