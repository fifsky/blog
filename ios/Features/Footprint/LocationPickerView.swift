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
/// 支持点击地图选择位置和关键字搜索地点，确认后返回坐标
struct LocationPickerView: View {

    // MARK: - 绑定属性

    /// 选中的纬度（Binding）
    @Binding var latitude: String

    /// 选中的经度（Binding）
    @Binding var longitude: String

    // MARK: - 状态

    /// 当前选中的位置坐标
    @State private var selectedCoordinate: CLLocationCoordinate2D?

    /// 地图显示区域
    @State private var mapRegion: MKCoordinateRegion

    /// 当前地图位置（用于 Map initialPosition）
    @State private var cameraPosition: MapCameraPosition

    /// 搜索框文本
    @State private var searchText = ""

    /// 是否展开搜索结果面板
    @State private var isSearching = false

    /// 搜索补全视图模型
    @State private var searchCompleter = LocationSearchCompleter()

    /// 是否正在执行地点解析（从搜索结果跳转）
    @State private var isResolving = false

    /// 环境变量：dismiss
    @Environment(\.dismiss) private var dismiss

    // MARK: - 初始化

    /// 初始化位置选择器
    /// - Parameters:
    ///   - latitude: 纬度 Binding
    ///   - longitude: 经度 Binding
    init(latitude: Binding<String>, longitude: Binding<String>) {
        _latitude = latitude
        _longitude = longitude

        // 如果已有坐标，以该坐标为中心；否则使用默认位置
        let lat = Double(latitude.wrappedValue) ?? 35.86
        let lon = Double(longitude.wrappedValue) ?? 104.19

        let region = MKCoordinateRegion(
            center: CLLocationCoordinate2D(latitude: lat, longitude: lon),
            span: MKCoordinateSpan(latitudeDelta: 0.05, longitudeDelta: 0.05)
        )

        _mapRegion = State(initialValue: region)
        _cameraPosition = State(initialValue: .region(region))

        // 如果已有坐标，预选
        if Double(latitude.wrappedValue) != nil && Double(longitude.wrappedValue) != nil {
            _selectedCoordinate = State(
                initialValue: CLLocationCoordinate2D(latitude: lat, longitude: lon)
            )
        }
    }

    var body: some View {
        ZStack(alignment: .bottom) {
            // 地图视图
            mapContent

            // 底部搜索栏 + 搜索结果面板
            VStack(spacing: 8) {
                // 搜索结果下拉面板（在搜索栏上方）
                if isSearching && !searchText.isEmpty {
                    searchResultsPanel
                }

                // 底部玻璃搜索栏
                bottomSearchBar
            }
        }
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

    /// 地图主内容，支持点击选择位置
    /// 通过 MapReader 将屏幕坐标转换为地理坐标，兼容 iOS 17+
    private var mapContent: some View {
        MapReader { proxy in
            Map(position: $cameraPosition) {
                // 如果已选中位置，显示标注
                if let coordinate = selectedCoordinate {
                    Annotation("已选位置", coordinate: coordinate) {
                        VStack(spacing: 0) {
                            Image(systemName: "mappin.circle.fill")
                                .font(.title)
                                .foregroundStyle(.red)

                            Image(systemName: "triangle.fill")
                                .font(.caption)
                                .foregroundStyle(.red)
                                .rotationEffect(.degrees(180))
                                .offset(y: -4)
                        }
                    }
                    .annotationTitles(.hidden)
                }
            }
            .mapStyle(.standard(elevation: .realistic))
            .mapControls {
                MapUserLocationButton()
                MapPitchToggle()
            }
            .gesture(
                // 点击地图时收起搜索面板
                SpatialTapGesture()
                    .onEnded { value in
                        if isSearching {
                            isSearching = false
                        }
                        guard let coordinate = proxy.convert(value.location, from: .local) else { return }
                        selectCoordinate(coordinate)
                    }
            )
        }
        .ignoresSafeArea(edges: .bottom)
    }

    /// 设置选中的坐标并更新地图区域
    private func selectCoordinate(_ coordinate: CLLocationCoordinate2D) {
        selectedCoordinate = coordinate
        cameraPosition = .region(
            MKCoordinateRegion(
                center: coordinate,
                span: MKCoordinateSpan(latitudeDelta: 0.01, longitudeDelta: 0.01)
            )
        )
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
                    .onChange(of: searchText) { _, newValue in
                        searchCompleter.search(newValue)
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
            .background(.ultraThinMaterial, in: Capsule())

            // 确认按钮（未选位为灰色、已选为蓝色）
            Button {
                confirmLocation()
            } label: {
                Image(systemName: "checkmark")
                    .font(.system(size: 18, weight: .bold))
                    .foregroundStyle(.white)
                    .frame(width: 44, height: 44)
                    .background(
                        selectedCoordinate == nil
                            ? Color(.systemGray4)
                            : Color.blue
                        , in: Circle()
                    )
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
            .background(.ultraThinMaterial, in: RoundedRectangle(cornerRadius: 16))
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
            isSearching = false
            selectCoordinate(coordinate)
        }
    }

    // MARK: - 操作方法

    /// 确认选择的位置，回传坐标（保留 6 位小数）并关闭
    private func confirmLocation() {
        guard let coordinate = selectedCoordinate else { return }
        latitude = String(format: "%.6f", coordinate.latitude)
        longitude = String(format: "%.6f", coordinate.longitude)
        dismiss()
    }
}
