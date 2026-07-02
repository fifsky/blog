import SwiftUI
import MapKit

/// 位置选择器视图
/// 全屏地图，用户通过点击地图选择位置，点击确认后返回坐标
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

    /// 提示文本
    @State private var hintText = "点击地图选择位置"

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
            _hintText = State(initialValue: "已选择位置，点击地图可更改")
        }
    }

    var body: some View {
        ZStack {
            // 地图视图
            mapContent

            // 顶部信息栏
            VStack {
                topInfoBar
                Spacer()
                bottomButtons
            }
        }
        .navigationTitle("选择位置")
        .navigationBarTitleDisplayMode(.inline)
        .toolbar {
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
                // 使用空间点击手势获取屏幕坐标，再转换为地图地理坐标
                SpatialTapGesture()
                    .onEnded { value in
                        guard let coordinate = proxy.convert(value.location, from: .local) else { return }
                        selectCoordinate(coordinate)
                    }
            )
        }
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
        hintText = "\(coordinate.latitude), \(coordinate.longitude)"
    }

    // MARK: - 顶部信息栏

    /// 顶部坐标信息显示
    private var topInfoBar: some View {
        VStack(spacing: 4) {
            Text(hintText)
                .font(.subheadline)
                .foregroundStyle(.secondary)
                .padding(.horizontal, 16)
                .padding(.vertical, 8)
                .background(.ultraThinMaterial, in: RoundedRectangle(cornerRadius: 10))
        }
        .padding(.top, 8)
    }

    // MARK: - 底部按钮

    /// 底部确认按钮区域
    private var bottomButtons: some View {
        VStack(spacing: 12) {
            // 确认位置按钮
            Button {
                confirmLocation()
            } label: {
                HStack {
                    Image(systemName: "checkmark.circle.fill")
                    Text("确认位置")
                        .fontWeight(.medium)
                }
                .frame(maxWidth: .infinity)
                .padding(.vertical, 14)
                .background(
                    selectedCoordinate == nil
                        ? Color(.systemGray4)
                        : .blue
                    , in: RoundedRectangle(cornerRadius: 12))
                .foregroundStyle(.white)
            }
            .disabled(selectedCoordinate == nil)
            .padding(.horizontal, 20)

            // 重置到当前位置按钮
            Button {
                resetToCurrentLocation()
            } label: {
                HStack(spacing: 6) {
                    Image(systemName: "location.circle")
                    Text("使用当前位置")
                        .font(.subheadline)
                }
                .foregroundStyle(.blue)
            }
            .padding(.bottom, 16)
        }
    }

    // MARK: - 操作方法

    /// 确认选择的位置，回传坐标并关闭
    private func confirmLocation() {
        guard let coordinate = selectedCoordinate else { return }
        latitude = "\(coordinate.latitude)"
        longitude = "\(coordinate.longitude)"
        dismiss()
    }

    /// 重置位置到用户当前所在位置
    private func resetToCurrentLocation() {
        let manager = CLLocationManager()
        manager.requestWhenInUseAuthorization()

        if manager.authorizationStatus == .authorizedWhenInUse
            || manager.authorizationStatus == .authorizedAlways {
            if let location = manager.location {
                selectCoordinate(location.coordinate)
            }
        }
    }
}
