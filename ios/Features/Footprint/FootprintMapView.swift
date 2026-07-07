import SwiftUI
import MapKit

/// 足迹地图视图
/// 使用 iOS 17+ MapKit API 展示所有足迹在地图上的位置
/// 注意：需要在 Info.plist 中添加 NSLocationWhenInUseUsageDescription
struct FootprintMapView: View {

    @State private var viewModel = FootprintMapViewModel()

    /// 新建足迹编辑器弹窗
    @State private var showEditor = false

    /// 照片浏览器相关状态
    @State private var photoBrowserURLs: [String] = []
    @State private var photoBrowserIndex = 0

    /// 标记用户点击了照片，等待 sheet 收起后再 push 浏览器
    @State private var pendingPhotoBrowse = false

    /// UIKit 导航壳入口
    @Environment(\.appNavigator) private var navigator

    var body: some View {
        ZStack {
            // 地图主视图
            mapContent

            // 加载指示器
            if viewModel.isLoading && viewModel.footprints.isEmpty {
                ProgressView("加载中...")
                    .padding()
                    .background(.ultraThinMaterial, in: RoundedRectangle(cornerRadius: 12))
            }
        }
        .navigationTitle("足迹地图")
        .navigationBarTitleDisplayMode(.inline)
        .toolbar {
            ToolbarItem(placement: .topBarTrailing) {
                Button {
                    showEditor = true
                } label: {
                    Image(systemName: "plus")
                }
            }
        }
        .task {
            await viewModel.loadFootprints()
        }
        .sheet(isPresented: $viewModel.showDetail) {
            footprintDetailSheet
        }
        .sheet(isPresented: $showEditor) {
            NavigationStack {
                FootprintEditorView(
                    footprint: nil,
                    onSave: {
                        showEditor = false
                        // 新建后刷新足迹地图
                        Task { await viewModel.loadFootprints() }
                    },
                    onCancel: {
                        showEditor = false
                    }
                )
            }
        }
        .alert("错误", isPresented: $viewModel.showError) {
            Button("确定", role: .cancel) {}
        } message: {
            Text(viewModel.errorMessage ?? "未知错误")
        }
        // 点击照片后先关闭 sheet，sheet 收起后再 push 浏览器
        .onChange(of: viewModel.showDetail) { _, isShown in
            if !isShown && pendingPhotoBrowse {
                pendingPhotoBrowse = false
                let urls = photoBrowserURLs
                let index = photoBrowserIndex
                let placeName = viewModel.selectedFootprint?.name ?? "照片"
                // 等待 sheet 收起动画完成后再 push，避免动画重叠
                DispatchQueue.main.asyncAfter(deadline: .now() + 0.25) {
                    navigator.push(
                        PhotoBrowserView(
                            photoURLs: urls,
                            initialIndex: index,
                            placeName: placeName
                        ),
                        onPop: {
                            if viewModel.selectedFootprint != nil {
                                viewModel.showDetail = true
                            }
                        }
                    )
                }
            }
        }
    }

    // MARK: - 地图内容

    /// 地图主内容
    private var mapContent: some View {
        Map(initialPosition: .region(viewModel.mapRegion)) {
            // 遍历所有足迹，在地图上添加标注
            ForEach(viewModel.footprints) { footprint in
                if let coordinate = footprint.coordinate {
                    Annotation(
                        footprint.name ?? "未命名",
                        coordinate: coordinate,
                        anchor: UnitPoint(x: 0.5, y: 0.91)
                    ) {
                        FootprintPhotoMarkerView(footprint: footprint)
                            .onTapGesture {
                                viewModel.selectFootprint(footprint)
                            }
                    }
                    .annotationTitles(.hidden)
                }
            }
        }
        .mapStyle(.standard(elevation: .realistic))
        .mapControls {
            MapUserLocationButton()
        }
    }

    // MARK: - 视图切换按钮（已移除：原切换列表为 TODO 空实现，刷新无可见反馈，均无实际作用）

    // MARK: - 足迹详情 Sheet

    /// 足迹详情弹窗
    private var footprintDetailSheet: some View {
        NavigationStack {
            Group {
                if let footprint = viewModel.selectedFootprint {
                    footprintDetailView(footprint)
                }
            }
            .toolbar {
                ToolbarItem(placement: .topBarTrailing) {
                    Button("完成") {
                        viewModel.showDetail = false
                    }
                }
            }
        }
    }

    /// 单个足迹的详情视图
    private func footprintDetailView(_ footprint: Footprint) -> some View {
        ScrollView {
            VStack(alignment: .leading, spacing: 16) {
                // 足迹名称
                Text(footprint.name ?? "未命名")
                    .font(.title2)
                    .bold()

                // 日期和坐标信息
                HStack(spacing: 16) {
                    if let date = footprint.date, !date.isEmpty {
                        Label(date, systemImage: "calendar")
                            .font(.subheadline)
                            .foregroundStyle(.secondary)
                    }
                    if let lat = footprint.latitude, let lon = footprint.longitude,
                       !(lat.isEmpty && lon.isEmpty) {
                        Label("\(lat), \(lon)", systemImage: "location")
                            .font(.subheadline)
                            .foregroundStyle(.secondary)
                    }
                }

                // 描述
                if let desc = footprint.description, !desc.isEmpty {
                    Text(desc)
                        .font(.body)
                        .foregroundStyle(.secondary)
                }

                // 照片缩略图
                if let photos = footprint.photos, !photos.isEmpty {
                    photoGrid(photos)
                }

                // 关联链接
                if let url = footprint.url, !url.isEmpty {
                    Link(destination: URL(string: url) ?? URL(string: "about:blank")!) {
                        HStack {
                            Image(systemName: "safari")
                            let label = footprint.url_label ?? ""
                            Text(label.isEmpty ? url : label)
                                .lineLimit(1)
                            Spacer()
                            Image(systemName: "arrow.up.right.square")
                        }
                        .font(.subheadline)
                    }
                }
            }
            .padding()
        }
        .presentationDetents([.medium, .large])
        .presentationDragIndicator(.visible)
    }

    /// 照片网格视图
    private func photoGrid(_ photos: [FootprintPhoto]) -> some View {
        let urls = photos.compactMap { $0.src }
        return ScrollView(.horizontal, showsIndicators: false) {
            HStack(spacing: 8) {
                ForEach(Array(photos.enumerated()), id: \.element.src) { index, photo in
                    let thumb = photo.thumbnail ?? ""
                    let src = photo.src ?? ""
                    AsyncImage(url: URL(string: thumb.isEmpty ? src : thumb)) { phase in
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
                    .frame(width: 120, height: 120)
                    .clipShape(RoundedRectangle(cornerRadius: 8))
                    .onTapGesture {
                        // 记录要浏览的照片，先收起 sheet，收起后再 push 全屏浏览器
                        photoBrowserURLs = urls
                        photoBrowserIndex = index
                        pendingPhotoBrowse = true
                        viewModel.showDetail = false
                    }
                }
            }
        }
    }
}

// MARK: - 足迹标注视图

/// 自定义地图标注视图：照片气泡 + 精确坐标点
struct FootprintPhotoMarkerView: View {

    /// 足迹数据
    let footprint: Footprint

    var body: some View {
        VStack(spacing: 4) {
            photoBubble
            coordinateDot
        }
        .frame(width: 58, height: 76)
        .accessibilityLabel(footprint.name ?? "足迹")
    }

    /// 照片气泡
    private var photoBubble: some View {
        Group {
            if let url = footprint.markerPhotoURL {
                AsyncImage(url: url) { phase in
                    switch phase {
                    case .success(let image):
                        image
                            .resizable()
                            .aspectRatio(contentMode: .fill)
                    case .failure:
                        fallbackBubble
                    default:
                        ProgressView()
                            .controlSize(.small)
                            .frame(maxWidth: .infinity, maxHeight: .infinity)
                            .background(Color(.secondarySystemBackground))
                    }
                }
            } else {
                fallbackBubble
            }
        }
        .frame(width: 58, height: 58)
        .clipShape(Circle())
        .overlay {
            Circle()
                .stroke(.white, lineWidth: 3)
        }
        .shadow(color: .black.opacity(0.22), radius: 5, x: 0, y: 2)
    }

    /// 无照片或加载失败时的兜底气泡
    private var fallbackBubble: some View {
        ZStack {
            markerColor
            Text(footprint.name?.first.map(String.init) ?? "足")
                .font(.title3)
                .fontWeight(.semibold)
                .foregroundStyle(.white)
        }
    }

    /// 坐标点：Annotation 的 anchor 对齐到这个圆点中心
    private var coordinateDot: some View {
        Circle()
            .fill(Color.blue)
            .frame(width: 14, height: 14)
            .overlay {
                Circle()
                    .stroke(.white, lineWidth: 3)
            }
            .shadow(color: .black.opacity(0.24), radius: 3, x: 0, y: 1)
    }

    /// 标记颜色
    private var markerColor: Color {
        Color(hex: footprint.marker_color ?? "#0A84FF") ?? .blue
    }
}

// MARK: - Footprint 扩展

extension Footprint {

    /// 地图标注使用的第一张照片 URL
    var markerPhotoURL: URL? {
        guard let firstPhoto = photos?.first else { return nil }
        let thumb = firstPhoto.thumbnail?.trimmingCharacters(in: .whitespacesAndNewlines) ?? ""
        let src = firstPhoto.src?.trimmingCharacters(in: .whitespacesAndNewlines) ?? ""
        let urlString = thumb.isEmpty ? src : thumb
        guard !urlString.isEmpty else { return nil }
        return URL(string: urlString)
    }

    /// 将字符串经纬度转换为 CLLocationCoordinate2D（GCJ-02 转 WGS-84，适配 MapKit）
    var coordinate: CLLocationCoordinate2D? {
        guard let latStr = latitude, let lonStr = longitude,
              let lat = Double(latStr), let lon = Double(lonStr) else { return nil }
        // 库存的是高德 GCJ-02 坐标，MapKit 用 WGS-84，需转换
        return CoordinateTransform.gcj02ToWgs84(
            CLLocationCoordinate2D(latitude: lat, longitude: lon)
        )
    }
}

// MARK: - Color + 十六进制扩展

extension Color {

    /// 从十六进制字符串初始化颜色
    /// 支持格式：#RRGGBB、RRGGBB、#RGB、RGB
    init?(hex: String) {
        var hexSanitized = hex.trimmingCharacters(in: .whitespacesAndNewlines)
        hexSanitized = hexSanitized.replacingOccurrences(of: "#", with: "")

        // 处理 3 位简写：#RGB -> #RRGGBB（每个字符重复一次）
        if hexSanitized.count == 3 {
            hexSanitized = hexSanitized.map { "\($0)\($0)" }.joined()
        }

        guard hexSanitized.count == 6 else { return nil }

        var rgb: UInt64 = 0
        Scanner(string: hexSanitized).scanHexInt64(&rgb)

        self.init(
            red: Double((rgb & 0xFF0000) >> 16) / 255.0,
            green: Double((rgb & 0x00FF00) >> 8) / 255.0,
            blue: Double(rgb & 0x0000FF) / 255.0
        )
    }

    /// 系统蓝色（默认标记颜色）
    static let systemBlue = Color.blue
}
