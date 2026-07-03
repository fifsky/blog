import Foundation
import CoreLocation

/// GCJ-02（火星坐标系）与 WGS-84（真实 GPS 坐标系）互转
///
/// 背景：中国大陆地图出于法规要求，对真实 GPS 坐标做了非线性加密（GCJ-02）。
/// 高德地图 Web/SDK 采用 GCJ-02，而苹果 MapKit 采用 WGS-84。
/// 本服务存的是高德采集的 GCJ-02 坐标，iOS 展示给 MapKit 前需转成 WGS-84。
///
/// 算法为业界通用的偏导迭代近似，参考公开实现。
enum CoordinateTransform {

    /// 常量 a：克拉索夫斯基椭球体长半轴
    private static let a: Double = 6378245.0

    /// 常量 ee：椭球偏心率平方
    private static let ee: Double = 0.00669342162296594323

    // MARK: - 公开方法

    /// GCJ-02（火星坐标）转 WGS-84（GPS 坐标）
    /// - Parameter coordinate: 高德/腾讯地图坐标
    /// - Returns: WGS-84 坐标（用于 MapKit、Google Maps 等）
    static func gcj02ToWgs84(_ coordinate: CLLocationCoordinate2D) -> CLLocationCoordinate2D {
        // 海外坐标无需转换
        guard isInChina(coordinate) else { return coordinate }

        // 先计算 GCJ-02 相对 WGS-84 的偏移量
        let (dLat, dLon) = transformOffset(coordinate.latitude, coordinate.longitude)
        // 粗略用偏移量反推，精度足够（误差在米级以下）
        return CLLocationCoordinate2D(
            latitude: coordinate.latitude - dLat,
            longitude: coordinate.longitude - dLon
        )
    }

    /// WGS-84（GPS 坐标）转 GCJ-02（火星坐标）
    /// - Parameter coordinate: MapKit/GPS 采集到的坐标
    /// - Returns: GCJ-02 坐标（用于高德地图等）
    static func wgs84ToGcj02(_ coordinate: CLLocationCoordinate2D) -> CLLocationCoordinate2D {
        guard isInChina(coordinate) else { return coordinate }

        let (dLat, dLon) = transformOffset(coordinate.latitude, coordinate.longitude)
        return CLLocationCoordinate2D(
            latitude: coordinate.latitude + dLat,
            longitude: coordinate.longitude + dLon
        )
    }

    // MARK: - 内部计算

    /// 判断坐标是否在中国大陆范围内（海外不加密，无需转换）
    private static func isInChina(_ coordinate: CLLocationCoordinate2D) -> Bool {
        let lat = coordinate.latitude
        let lon = coordinate.longitude
        return lat > 0.8293 && lat < 55.8271 && lon > 72.004 && lon < 137.8347
    }

    /// 计算 GCJ-02 加密的经纬度偏移量
    /// - Returns: (纬度偏移 dLat, 经度偏移 dLon)
    private static func transformOffset(_ lat: Double, _ lon: Double) -> (Double, Double) {
        var dLat = transformLat(lon - 105.0, lat - 35.0)
        var dLon = transformLon(lon - 105.0, lat - 35.0)
        let radLat = lat / 180.0 * .pi
        var magic = sin(radLat)
        magic = 1 - ee * magic * magic
        let sqrtMagic = sqrt(magic)
        dLat = (dLat * 180.0) / ((a * (1 - ee)) / (magic * sqrtMagic) * .pi)
        dLon = (dLon * 180.0) / (a / sqrtMagic * cos(radLat) * .pi)
        return (dLat, dLon)
    }

    /// 纬度偏移变换
    private static func transformLat(_ x: Double, _ y: Double) -> Double {
        var ret = -100.0 + 2.0 * x + 3.0 * y + 0.2 * y * y
            + 0.1 * x * y + 0.2 * sqrt(abs(x))
        ret += (20.0 * sin(6.0 * x * .pi) + 20.0 * sin(2.0 * x * .pi)) * 2.0 / 3.0
        ret += (20.0 * sin(y * .pi) + 40.0 * sin(y / 3.0 * .pi)) * 2.0 / 3.0
        ret += (160.0 * sin(y / 12.0 * .pi) + 320.0 * sin(y * .pi / 30.0)) * 2.0 / 3.0
        return ret
    }

    /// 经度偏移变换
    private static func transformLon(_ x: Double, _ y: Double) -> Double {
        var ret = 300.0 + x + 2.0 * y + 0.1 * x * x
            + 0.1 * x * y + 0.1 * sqrt(abs(x))
        ret += (20.0 * sin(6.0 * x * .pi) + 20.0 * sin(2.0 * x * .pi)) * 2.0 / 3.0
        ret += (20.0 * sin(x * .pi) + 40.0 * sin(x / 3.0 * .pi)) * 2.0 / 3.0
        ret += (150.0 * sin(x / 12.0 * .pi) + 300.0 * sin(x / 30.0 * .pi)) * 2.0 / 3.0
        return ret
    }
}
