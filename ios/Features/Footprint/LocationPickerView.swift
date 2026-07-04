import SwiftUI
import MapKit

/// 位置搜索补全视图模型
/// 使用 MKLocalSearchCompleter 实现地址关键字搜索
@Observable
private class LocationSearchCompleter: NSObject, MKLocalSearchCompleterDelegate {

    /// 搜索结果
    var results: [MKLocalSearchCompletion] = []

    /// 是否正在搜索
    var isLoading = false

    private let completer = MKLocalSearchCompleter()

    override init() {
        super.init()
        completer.delegate = self
        // 仅搜索地点，过滤掉商业类别等噪音
        completer.resultTypes = .address
    }

    /// 更新搜索关键字
    /// - Parameter query: 搜索词
    func search(_ query: String) {
        if query.isEmpty {
            results = []
            isLoading = false
            return
        }
        isLoading = true
        completer.queryFragment = query
    }

    // MARK: - MKLocalSearchCompleterDelegate

    func completerDidUpdateResults(_ completer: MKLocalSearchCompleter) {
        results = completer.results
        isLoading = false
    }

    func completer(_ completer: MKLocalSearchCompleter, didFailWithError error: Error) {
        isLoading = false
    }
}

/// 位置选择器视图
/// 顶部返回按钮 + 全屏地图 + 底部玻璃质感搜索栏（定位按钮 + 确认按钮）
/// 采用「中心固定图钉 + 地图拖动微调」交互（类似微信发送位置）：
/// - 图钉始终位于地图屏幕中心
/// - 拖动地图时图钉抬起（缩小），停止时落下（弹回）动画
/// - 搜索后定位到结果坐标，再通过拖动微调
/// - 点击右侧确认按钮以地图中心点作为最终选取坐标
struct LocationPickerView: View {

    // MARK: - 绑定属性

    /// 选中的纬度（Binding）
    @Binding var latitude: String

    /// 选中的经度（Binding）
    @Binding var longitude: String

    // MARK: - 状态

    /// 当前选中的位置坐标（始终等于地图中心）
    @State private var selectedCoordinate: CLLocationCoordinate2D?

    /// 地图显示区域
    @State private var mapRegion: MKCoordinateRegion

    /// 当前地图位置（用于 Map position）
    @State private var cameraPosition: MapCameraPosition

    /// 搜索框文本
    @State private var searchText = ""

    /// 是否正在拖动地图（用于图钉抬起动画）
    @State private var isDraggingMap = false

    /// 搜索补全视图模型
    @State private var searchCompleter = LocationSearchCompleter()

    /// 是否正在执行地点解析（从搜索结果跳转）
    @State private var isResolving = false

    /// 是否聚焦搜索框（用于显示/隐藏搜索结果面板）
    @FocusState private var isSearchFieldFocused: Bool

    /// 环境变量：dismiss
    @Environment(\.dismiss) private var dismiss

    // MARK: - 初始化

    /// 初始化位置选择器
    /// - Parameters:
    ///   - latitude: WGS-84 纬度 Binding（由调用方缓存，全程不做坐标转换）
    ///   - longitude: WGS-84 经度 Binding
    init(latitude: Binding<String>, longitude: Binding<String>) {
        _latitude = latitude
        _longitude = longitude

        // 如果已有坐标，以该坐标为中心；否则使用默认位置
        let lat = Double(latitude.wrappedValue) ?? 35.86
        let lon = Double(longitude.wrappedValue) ?? 104.19

        // 全程 WGS-84，无需坐标转换（调用方 FootprintEditorViewModel 已缓存 WGS-84）
        let coordinate = CLLocationCoordinate2D(latitude: lat, longitude: lon)
        let region = MKCoordinateRegion(
            center: coordinate,
            span: MKCoordinateSpan(latitudeDelta: 0.05, longitudeDelta: 0.05)
        )

        _mapRegion = State(initialValue: region)
        _cameraPosition = State(initialValue: .region(region))

        // 如果已有坐标，预选
        if Double(latitude.wrappedValue) != nil && Double(longitude.wrappedValue) != nil {
            _selectedCoordinate = State(initialValue: coordinate)
        }
    }

    var body: some View {
        ZStack(alignment: .bottom) {
            // 地图视图
            mapContent

            // 中心固定图钉覆盖层（始终位于地图屏幕中心）
            centerPinOverlay

            // 底部搜索栏 + 搜索结果面板
            VStack(spacing: 8) {
                // 搜索结果下拉面板（在搜索栏上方）
                if isSearchFieldFocused && !searchText.isEmpty {
                    searchResultsPanel
                }

                // 底部搜索栏
                bottomSearchBar
            }
            // 浮于图钉之上：搜索结果列表展开时不应被中心图钉遮挡
            .zIndex(2)
            .ignoresSafeArea(.keyboard, edges: .bottom)
        }
        // 点击空白区域收起键盘
        .hideKeyboardOnTap()
        .navigationTitle("选择位置")
        .navigationBarTitleDisplayMode(.inline)
        .toolbar {
            // 原生返回按钮，与其他页面风格一致
            ToolbarItem(placement: .topBarLeading) {
                Button("取消") {
                    dismiss()
                }
            }
        }
    }

    // MARK: - 地图内容

    /// 地图主内容
    /// 图钉固定在屏幕中心，拖动地图时记录中心点作为选中坐标
    /// 通过 onMapCameraChange 监听地图拖动状态与中心点
    private var mapContent: some View {
        Map(position: $cameraPosition) {
            // 中心固定图钉模型下，无需在地图上添加移动标注
        }
        .mapStyle(.standard(elevation: .realistic))
        .mapControls {
            MapUserLocationButton()
            MapPitchToggle()
        }
        .onMapCameraChange { context in
            // 实时跟踪地图中心点作为选中坐标
            selectedCoordinate = context.region.center
            // 连续变化时标记为拖动中（图钉抬起）
            isDraggingMap = true
        }
        .onMapCameraChange(frequency: .onEnd) { _ in
            // 拖动结束，图钉落下
            isDraggingMap = false
        }
        .ignoresSafeArea(edges: .bottom)
        .zIndex(0)
    }

    /// 设置地图中心并同步选中坐标
    /// - Parameter coordinate: 目标坐标
    private func centerMap(on coordinate: CLLocationCoordinate2D) {
        selectedCoordinate = coordinate
        cameraPosition = .region(
            MKCoordinateRegion(
                center: coordinate,
                span: MKCoordinateSpan(latitudeDelta: 0.01, longitudeDelta: 0.01)
            )
        )
    }

    // MARK: - 中心固定图钉

    /// 地图屏幕中心固定图钉覆盖层
    /// 使用 📍 emoji，其自带朝下尖角，尖端天然对准地图中心坐标
    /// 拖动地图时图钉抬起（缩小并轻微上移），停止时落下（弹回原大小）
    private var centerPinOverlay: some View {
        Text("📍")
            .font(.system(size: 40))
            // 拖动时抬起：缩小并上移，模拟「被吸起」效果
            .scaleEffect(isDraggingMap ? 0.85 : 1.0)
            .offset(y: isDraggingMap ? -6 : 0)
            .animation(.spring(response: 0.3, dampingFraction: 0.6), value: isDraggingMap)
            // 📍 字形尖端在底部，对准地图屏幕中心需上移约半个字形高度
            .offset(y: -20)
            .shadow(color: .black.opacity(0.2), radius: 3, x: 0, y: 1)
            .frame(maxWidth: .infinity, maxHeight: .infinity)
            .allowsHitTesting(false)
            .zIndex(1)
    }

    // MARK: - 底部搜索栏

    /// 底部玻璃质感搜索栏：搜索输入 + 确认按钮
    private var bottomSearchBar: some View {
        HStack(spacing: 8) {
            // 搜索输入框
            HStack(spacing: 8) {
                Image(systemName: "magnifyingglass")
                    .foregroundStyle(.secondary)

                TextField("搜索地点", text: $searchText)
                    .textInputAutocapitalization(.never)
                    .autocorrectionDisabled()
                    .submitLabel(.search)
                    .focused($isSearchFieldFocused)
                    .onChange(of: searchText) { _, newValue in
                        searchCompleter.search(newValue)
                    }
                    .onSubmit {
                        // 点击键盘「搜索」按钮：解析首个匹配结果并定位
                        resolveFirstResult()
                    }

                if !searchText.isEmpty {
                    Button {
                        searchText = ""
                        searchCompleter.search("")
                    } label: {
                        Image(systemName: "xmark.circle.fill")
                            .foregroundStyle(.secondary)
                    }
                }
            }
            .padding(.horizontal, 14)
            .frame(height: 44)
            .background(Color.white.opacity(0.9), in: Capsule())

            // 确认按钮（始终可点击，以当前地图中心作为选取坐标）
            Button {
                confirmLocation()
            } label: {
                Image(systemName: "checkmark")
                    .font(.system(size: 18, weight: .bold))
                    .foregroundStyle(.white)
                    .frame(width: 44, height: 44)
                    .background(Color.blue, in: Circle())
            }
            .disabled(selectedCoordinate == nil)
        }
        .padding(.horizontal, 16)
        .padding(.bottom, 16)
    }

    /// 搜索结果下拉面板
    private var searchResultsPanel: some View {
        ScrollView {
            LazyVStack(spacing: 0) {
                if searchCompleter.isLoading {
                    ProgressView()
                        .padding()
                } else if searchCompleter.results.isEmpty {
                    Text("无匹配地点")
                        .font(.subheadline)
                        .foregroundStyle(.secondary)
                        .padding()
                } else {
                    ForEach(searchCompleter.results, id: \.self) { completion in
                        Button {
                            resolveCompletion(completion)
                        } label: {
                            VStack(alignment: .leading, spacing: 2) {
                                Text(completion.title)
                                    .font(.body)
                                    .foregroundStyle(.primary)
                                if !completion.subtitle.isEmpty {
                                    Text(completion.subtitle)
                                        .font(.caption)
                                        .foregroundStyle(.secondary)
                                        .lineLimit(1)
                                }
                            }
                            .frame(maxWidth: .infinity, alignment: .leading)
                            .padding(.horizontal, 16)
                            .padding(.vertical, 10)
                        }
                        Divider()
                    }
                }
            }
            .background(Color.white.opacity(0.8), in: RoundedRectangle(cornerRadius: 16))
        }
        .frame(maxHeight: 280)
        .padding(.horizontal, 16)
    }

    /// 解析搜索补全项为坐标并定位
    /// - Parameter completion: 搜索补全结果
    private func resolveCompletion(_ completion: MKLocalSearchCompletion) {
        isResolving = true
        let request = MKLocalSearch.Request(completion: completion)
        let search = MKLocalSearch(request: request)
        search.start { response, _ in
            isResolving = false
            guard let coordinate = response?.mapItems.first?.placemark.coordinate else { return }
            searchText = ""
            isSearchFieldFocused = false
            centerMap(on: coordinate)
        }
    }

    /// 解析当前搜索词的首个结果并定位（键盘「搜索」按钮触发）
    private func resolveFirstResult() {
        // 优先使用已补全的首个结果，避免重复请求
        if let first = searchCompleter.results.first {
            resolveCompletion(first)
            return
        }
        // 补全列表为空时按关键字发起一次完整搜索
        guard !searchText.trimmingCharacters(in: .whitespaces).isEmpty else { return }
        isResolving = true
        let request = MKLocalSearch.Request()
        request.naturalLanguageQuery = searchText
        let search = MKLocalSearch(request: request)
        search.start { response, _ in
            isResolving = false
            guard let coordinate = response?.mapItems.first?.placemark.coordinate else { return }
            searchText = ""
            isSearchFieldFocused = false
            centerMap(on: coordinate)
        }
    }

    // MARK: - 操作方法

    /// 确认选择的位置，回传 WGS-84 坐标（保留 6 位小数）并关闭
    /// 全程 WGS-84 无坐标转换，由调用方 FootprintEditorViewModel 正向算 GCJ-02 供后端
    private func confirmLocation() {
        guard let coordinate = selectedCoordinate else { return }
        latitude = String(format: "%.6f", coordinate.latitude)
        longitude = String(format: "%.6f", coordinate.longitude)
        dismiss()
    }
}
