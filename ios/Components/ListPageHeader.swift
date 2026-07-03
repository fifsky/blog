import SwiftUI

/// 列表页统一大标题 Header
///
/// 自己负责全部布局（安全区域由外层导航栏避让）：
/// - 左右边距：16
/// - 顶部：8
/// - 底部：bottomPadding（默认 8）
///
/// - ScrollView 页面：把 Header 和内容并列放进 VStack，内容由内容自身负责横向 padding。
/// - List 页面：用 `.safeAreaInset(edge: .top)` 把 Header 固定在 List 上方，脱离 List 行的
///   圆角裁剪，让 Header 内部的 padding 直接生效，与 ScrollView 页面逐像素对齐。
struct ListPageHeader: View {

    /// 主标题（如"博文"/"心情"/"提醒"）
    let title: String

    /// 副标题（可选，如分组数、简介）
    var subtitle: String? = nil

    /// 标题前的图标名（可选，如"book"/"moon"/"bell"）
    var icon: String? = nil

    /// 底部间距（默认 8；紧跟 Section header 的页面可设为 0）
    var bottomPadding: CGFloat = 8

    var body: some View {
        VStack(alignment: .leading, spacing: 6) {
            HStack(spacing: 8) {
                if let icon {
                    Image(systemName: icon)
                        .font(.title3)
                        .foregroundStyle(.secondary)
                }
                Text(title)
                    .font(.system(size: 34, weight: .bold))
                    .foregroundStyle(Color.themePrimary)
                    .accessibilityAddTraits(.isHeader)
            }

            if let subtitle, !subtitle.isEmpty {
                Text(subtitle)
                    .font(.subheadline)
                    .foregroundStyle(.secondary)
            }
        }
        .frame(maxWidth: .infinity, alignment: .leading)
        // Header 自己拥有全部边距：任何页面都不要再叠加 padding
        .padding(.horizontal, 16)
        .padding(.top, 8)
        .padding(.bottom, bottomPadding)
    }
}
