import Foundation
import MapKit

/// 足迹地图视图模型
@Observable
class FootprintMapViewModel {

    // MARK: - 状态

    /// 足迹列表
    var footprints: [Footprint] = []

    /// 是否正在加载
    var isLoading = false

    /// 错误信息
    var errorMessage: String?

    /// 是否显示错误
    var showError = false

    /// 选中的足迹（用于展示详情 sheet）
    var selectedFootprint: Footprint?

    /// 是否显示足迹详情
    var showDetail = false

    /// 是否显示列表视图
    var showListView = false

    // MARK: - 私有属性

    private let service = FootprintService.shared

    // MARK: - 地图位置

    /// 默认地图中心（中国中心大致位置）
    static let defaultRegion = MKCoordinateRegion(
        center: CLLocationCoordinate2D(latitude: 35.86, longitude: 104.19),
        span: MKCoordinateSpan(latitudeDelta: 40, longitudeDelta: 40)
    )

    /// 地图区域（根据足迹数据自动调整）
    var mapRegion = defaultRegion

    // MARK: - 数据加载

    /// 加载所有足迹数据（公开接口，无需认证）
    func loadFootprints() async {
        guard !isLoading else { return }
        isLoading = true
        errorMessage = nil
        showError = false

        do {
            let response = try await service.all()
            footprints = response.footprints

            // 根据足迹数据调整地图区域
            adjustMapRegion()
        } catch {
            if !error.isCancellation {
                errorMessage = error.localizedDescription
                showError = true
            }
        }

        isLoading = false
    }

    /// 根据已有足迹坐标调整地图显示区域
    private func adjustMapRegion() {
        guard !footprints.isEmpty else {
            mapRegion = Self.defaultRegion
            return
        }

        // 筛选出有效坐标的足迹
        let validItems = footprints.filter {
            guard let latStr = $0.latitude, let lonStr = $0.longitude,
                  let lat = Double(latStr), let lon = Double(lonStr) else { return false }
            return lat != 0 && lon != 0
        }

        guard !validItems.isEmpty else {
            mapRegion = Self.defaultRegion
            return
        }

        // 计算所有足迹的坐标范围（GCJ-02 转 WGS-84，适配 MapKit）
        let wgsCoords = validItems.compactMap { item -> CLLocationCoordinate2D? in
            item.coordinate
        }
        let lats = wgsCoords.map { $0.latitude }
        let lons = wgsCoords.map { $0.longitude }

        let minLat = lats.min() ?? 0
        let maxLat = lats.max() ?? 0
        let minLon = lons.min() ?? 0
        let maxLon = lons.max() ?? 0

        let center = CLLocationCoordinate2D(
            latitude: (minLat + maxLat) / 2,
            longitude: (minLon + maxLon) / 2
        )

        // 计算跨度，留一些边距
        let latDelta = max(maxLat - minLat, 1) * 1.2
        let lonDelta = max(maxLon - minLon, 1) * 1.2

        mapRegion = MKCoordinateRegion(
            center: center,
            span: MKCoordinateSpan(latitudeDelta: min(latDelta, 180), longitudeDelta: min(lonDelta, 180))
        )
    }

    // MARK: - 交互

    /// 选中某个足迹
    func selectFootprint(_ footprint: Footprint) {
        selectedFootprint = footprint
        showDetail = true
    }

    /// 刷新数据
    func refresh() async {
        await loadFootprints()
    }
}
