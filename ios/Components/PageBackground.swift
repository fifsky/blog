import SwiftUI

/// 列表页装饰性背景
///
/// 将水彩背景图铺满整个屏幕（含安全区域），轻微压暗以提升前景内容可读性。
/// 背景固定，不随滚动变化；Header 与 List 均位于背景之上。
struct PageBackground: View {

    /// 背景图资源名（Assets.xcassets 中的 imageset 名称）
    let imageName: String

    var body: some View {
        GeometryReader { _ in
            // 背景图，填充整个屏幕（含安全区域）
            Image(imageName)
                .resizable()
                .scaledToFill()
                // 轻微压暗，提升前景内容可读性
                .overlay(Color.black.opacity(0.05))
                .accessibilityHidden(true)
        }
        // 忽略安全区域，确保铺满整个屏幕（包括状态栏区域）
        .ignoresSafeArea()
    }
}
