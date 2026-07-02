import SwiftUI
import MapKit

/// 足迹地图视图
/// 使用 iOS 17+ MapKit API 展示所有足迹在地图上的位置
/// 注意：需要在 Info.plist 中添加 NSLocationWhenInUseUsageDescription
struct FootprintMapView: View {

    @State private var viewModel = FootprintMapViewModel()

    /// 新建足迹编辑器弹窗
    @State private var showEditor = false

    /// 切换到列表视图的回调
    var onShowListView: () -> Void

    /// 新建足迹的回调
    var onAddFootprint: () -> Void

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
                    onAddFootprint()
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
    }

    // MARK: - 地图内容

    /// 地图主内容
    private var mapContent: some View {
        Map(initialPosition: .region(viewModel.mapRegion)) {
            // 遍历所有足迹，在地图上添加标注
            ForEach(viewModel.footprints) { footprint in
                if let coordinate = footprint.coordinate {
                    Annotation(footprint.name ?? "未命名", coordinate: coordinate) {
                        FootprintMarkerView(color: footprint.marker_color ?? "#FF3B30")
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
        .overlay(alignment: .bottomTrailing) {
            // 切换视图按钮
            toggleViewButton
        }
    }

    // MARK: - 视图切换按钮

    /// 底部右侧的视图切换按钮
    private var toggleViewButton: some View {
        VStack(spacing: 8) {
            // 切换到列表视图
            Button {
                onShowListView()
            } label: {
                Image(systemName: "list.bullet")
                    .font(.title3)
                    .frame(width: 44, height: 44)
                    .background(.ultraThinMaterial, in: Circle())
            }
            .padding(.trailing, 16)

            // 刷新按钮
            Button {
                Task { await viewModel.refresh() }
            } label: {
                Image(systemName: "arrow.clockwise")
                    .font(.title3)
                    .frame(width: 44, height: 44)
                    .background(.ultraThinMaterial, in: Circle())
            }
        }
        .padding(.bottom, 16)
    }

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
        ScrollView(.horizontal, showsIndicators: false) {
            HStack(spacing: 8) {
                ForEach(photos, id: \.src) { photo in
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
                }
            }
        }
    }
}

// MARK: - 足迹标注视图

/// 自定义地图标注视图：带颜色的圆形标记
struct FootprintMarkerView: View {

    /// 标记颜色（十六进制字符串）
    let color: String

    var body: some View {
        ZStack {
            Circle()
                .fill(markerColor)
                .frame(width: 28, height: 28)
                .shadow(color: .black.opacity(0.3), radius: 2, x: 0, y: 1)

            // 内部白色圆点
            Circle()
                .fill(.white)
                .frame(width: 8, height: 8)
        }
        .overlay(
            // 底部三角箭头
            Triangle()
                .fill(markerColor)
                .frame(width: 12, height: 8)
                .offset(y: 16),
            alignment: .bottom
        )
    }

    /// 将十六进制颜色字符串转换为 SwiftUI Color
    private var markerColor: Color {
        Color(hex: color) ?? .systemBlue
    }
}

// MARK: - 三角形形状

/// 向下的三角形形状（用于标注底部箭头）
private struct Triangle: Shape {
    func path(in rect: CGRect) -> Path {
        var path = Path()
        path.move(to: CGPoint(x: rect.midX, y: rect.maxY))
        path.addLine(to: CGPoint(x: rect.minX, y: rect.minY))
        path.addLine(to: CGPoint(x: rect.maxX, y: rect.minY))
        path.closeSubpath()
        return path
    }
}

// MARK: - Footprint 扩展

extension Footprint {

    /// 将字符串经纬度转换为 CLLocationCoordinate2D
    var coordinate: CLLocationCoordinate2D? {
        guard let latStr = latitude, let lonStr = longitude,
              let lat = Double(latStr), let lon = Double(lonStr) else { return nil }
        return CLLocationCoordinate2D(latitude: lat, longitude: lon)
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
