import SwiftUI
import PhotosUI

/// 足迹编辑器视图模型
/// 负责创建和编辑足迹的逻辑，包括照片上传
@Observable
class FootprintEditorViewModel {

    // MARK: - 表单数据

    /// 足迹名称（必填）
    var name = ""

    /// 描述
    var descriptionText = ""

    /// 日期字符串（格式：yyyy-MM-dd）
    var dateString = ""

    /// 标记颜色（十六进制）
    var markerColor = "#FF3B30"

    /// 纬度字符串
    var latitude = ""

    /// 经度字符串
    var longitude = ""

    /// 选中的照片（PhotosPicker 的结果）
    var selectedPhotoItems: [PhotosPickerItem] = []

    /// 已选照片的本地缩略图数据（用于预览）
    var photoPreviews: [Data] = []

    /// 已上传的照片 URL 列表
    var uploadedPhotoURLs: [String] = []

    // MARK: - 状态

    /// 是否正在保存
    var isSaving = false

    /// 是否正在上传照片
    var isUploading = false

    /// 上传进度描述
    var uploadProgress = ""

    /// 错误信息
    var errorMessage: String?

    /// 是否显示错误
    var showError = false

    /// 保存是否成功
    var isSaveSuccess = false

    /// 是否显示位置选择器
    var showLocationPicker = false

    // MARK: - 编辑模式

    /// 正在编辑的足迹（nil 表示新建模式）
    private var editingFootprint: Footprint?

    /// 是否为编辑模式
    var isEditing: Bool {
        editingFootprint != nil
    }

    // MARK: - 预定义颜色

    /// 可选的标记颜色列表
    static let markerColors: [(name: String, hex: String)] = [
        ("红色", "#FF3B30"),
        ("橙色", "#FF9500"),
        ("黄色", "#FFCC00"),
        ("绿色", "#34C759"),
        ("蓝色", "#007AFF"),
        ("紫色", "#AF52DE"),
        ("粉色", "#FF2D55"),
        ("棕色", "#A2845E"),
    ]

    // MARK: - 私有属性

    private let footprintService = FootprintService.shared
    private let uploadService = UploadService.shared

    // MARK: - 初始化

    /// 初始化编辑器视图模型
    /// - Parameter footprint: 如果传入则进入编辑模式，否则为新建模式
    init(footprint: Footprint? = nil) {
        if let footprint {
            editingFootprint = footprint
            name = footprint.name ?? ""
            descriptionText = footprint.description ?? ""
            dateString = footprint.date ?? ""
            markerColor = footprint.marker_color ?? "#FF3B30"
            latitude = footprint.latitude ?? ""
            longitude = footprint.longitude ?? ""
            // 加载已有照片 URL
            uploadedPhotoURLs = (footprint.photos ?? []).compactMap { $0.src }
        }
    }

    // MARK: - 照片处理

    /// 处理从照片选择器选中的照片，加载缩略图用于预览
    func loadSelectedPhotos() async {
        photoPreviews.removeAll()

        for item in selectedPhotoItems {
            if let data = try? await item.loadTransferable(type: Data.self) {
                photoPreviews.append(data)
            }
        }
    }

    /// 移除已选照片（按索引）
    func removeSelectedPhoto(at index: Int) {
        guard index < selectedPhotoItems.count else { return }
        selectedPhotoItems.remove(at: index)
        photoPreviews.remove(at: index)
    }

    // MARK: - 位置处理

    /// 设置选定的位置坐标（保留 6 位小数）
    func setLocation(latitude lat: Double, longitude lon: Double) {
        latitude = String(format: "%.6f", lat)
        longitude = String(format: "%.6f", lon)
    }

    /// 位置信息描述（只展示坐标，逗号分隔，避免文字过长换行）
    var locationDescription: String {
        if latitude.isEmpty || longitude.isEmpty {
            return "未选择位置"
        }
        return "\(latitude), \(longitude)"
    }

    // MARK: - 表单验证

    /// 表单是否可以提交
    var canSave: Bool {
        !name.trimmingCharacters(in: .whitespaces).isEmpty
    }

    // MARK: - 保存

    /// 保存足迹（先上传照片，再创建/更新足迹）
    func save() async {
        guard canSave else { return }
        isSaving = true

        do {
            // 第一步：上传所有新选中的照片
            var allPhotoURLs = uploadedPhotoURLs

            if !selectedPhotoItems.isEmpty {
                isUploading = true
                let newURLs = try await uploadSelectedPhotos()
                allPhotoURLs.append(contentsOf: newURLs)
                isUploading = false
            }

            // 第二步：创建或更新足迹
            if let footprint = editingFootprint {
                // 编辑模式
                let request = FootprintUpdateRequest(
                    id: footprint.id,
                    name: name,
                    description: descriptionText.isEmpty ? nil : descriptionText,
                    longitude: longitude,
                    latitude: latitude,
                    date: dateString.isEmpty ? nil : dateString,
                    marker_color: markerColor,
                    categories: nil,
                    url: nil,
                    url_label: nil,
                    photo_urls: allPhotoURLs.isEmpty ? nil : allPhotoURLs
                )
                _ = try await footprintService.update(params: request)
            } else {
                // 新建模式
                let request = FootprintCreateRequest(
                    name: name,
                    description: descriptionText.isEmpty ? nil : descriptionText,
                    longitude: longitude,
                    latitude: latitude,
                    date: dateString.isEmpty ? nil : dateString,
                    marker_color: markerColor,
                    categories: nil,
                    url: nil,
                    url_label: nil,
                    photo_urls: allPhotoURLs.isEmpty ? nil : allPhotoURLs
                )
                _ = try await footprintService.create(params: request)
            }

            isSaveSuccess = true
        } catch {
            errorMessage = error.localizedDescription
            showError = true
        }

        isSaving = false
        isUploading = false
    }

    // MARK: - 照片上传

    /// 上传所有选中的照片
    /// - Returns: 上传成功后的图片 URL 列表
    private func uploadSelectedPhotos() async throws -> [String] {
        var urls: [String] = []
        let total = selectedPhotoItems.count

        for (index, item) in selectedPhotoItems.enumerated() {
            uploadProgress = "上传照片 \(index + 1)/\(total)..."

            // 加载照片数据
            guard let data = try? await item.loadTransferable(type: Data.self) else {
                continue
            }

            // 生成文件名
            let filename = "footprint_\(Int(Date().timeIntervalSince1970))_\(index).jpg"

            // 上传图片
            let url = try await uploadService.uploadImage(imageData: data, filename: filename)
            urls.append(url)
        }

        uploadProgress = ""
        return urls
    }
}
