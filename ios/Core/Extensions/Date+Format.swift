import Foundation

extension Date {

    // MARK: - 格式化输出

    /// 按指定格式输出日期字符串
    /// - Parameter format: 日期格式（如 "yyyy-MM-dd HH:mm:ss"）
    /// - Returns: 格式化后的字符串
    func toString(format: String) -> String {
        let formatter = DateFormatter()
        formatter.dateFormat = format
        formatter.locale = Locale(identifier: "zh_CN")
        return formatter.string(from: self)
    }

    // MARK: - 相对时间

    /// 返回相对时间描述，如 "刚刚"、"3分钟前"、"2小时前"、"昨天"、"3天前"
    func relativeString() -> String {
        let now = Date()
        let interval = now.timeIntervalSince(self)

        // 未来时间
        if interval < 0 {
            return "刚刚"
        }

        let seconds = Int(interval)

        // 一分钟内
        if seconds < 60 {
            return "刚刚"
        }

        // 一小时内
        if seconds < 3600 {
            let minutes = seconds / 60
            return "\(minutes)分钟前"
        }

        // 一天内
        if seconds < 86_400 {
            let hours = seconds / 3600
            return "\(hours)小时前"
        }

        // 两天内
        if seconds < 172_800 {
            return "昨天"
        }

        // 七天内
        if seconds < 604_800 {
            let days = seconds / 86_400
            return "\(days)天前"
        }

        // 超过七天，显示日期
        let calendar = Calendar.current
        if calendar.isDate(self, equalTo: now, toGranularity: .year) {
            return toString(format: "M月d日")
        } else {
            return toString(format: "yyyy年M月d日")
        }
    }

    // MARK: - Unix 时间戳

    /// 从 Unix 时间戳创建 Date
    /// - Parameter timestamp: Unix 时间戳（秒）
    /// - Returns: 对应的 Date，无效值返回 nil
    static func fromUnixTimestamp(_ timestamp: Int64) -> Date? {
        // 处理毫秒级时间戳
        if timestamp > 9_999_999_999 {
            return Date(timeIntervalSince1970: TimeInterval(timestamp) / 1000.0)
        }
        return Date(timeIntervalSince1970: TimeInterval(timestamp))
    }
}
