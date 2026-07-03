import SwiftUI

/// 列表页自定义大标题 Header
///
/// 独立于 List/ScrollView 之外，放在 VStack 顶部，占固定高度。
/// 不参与系统 NavigationTitle 大标题布局，也不使用 overlay/safeAreaPadding 等 hack。
///
/// 典型用法：
/// ```
/// VStack(spacing: 0) {
///     ListPageHeader(title: "提醒", subtitle: "进行中")
///     List { ... }
/// }
/// ```
struct ListPageHeader: View {

    /// 主标题（如"博文"/"心情"/"提醒"）
    let title: String

    /// 副标题（可选，如分组数、简介）
    var subtitle: String? = nil

    /// 底部间距（默认 12；含子标题的页面可设为 0，由子标题自身间距控制）
    var bottomPadding: CGFloat = 12

    var body: some View {
        VStack(alignment: .leading, spacing: 4) {
            Text(title)
                .font(.system(size: 34, weight: .bold))
                .foregroundStyle(Color.themePrimary)
                .accessibilityAddTraits(.isHeader)

            if let subtitle, !subtitle.isEmpty {
                Text(subtitle)
                    .font(.subheadline)
                    .foregroundStyle(.secondary)
            }
        }
        .frame(maxWidth: .infinity, alignment: .leading)
        // 顶部留白：透明导航栏已提供 Safe Area 避让，这里仅做美观间距
        .padding(.top, 8)
        .padding(.bottom, bottomPadding)
    }
}
