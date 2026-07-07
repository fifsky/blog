import SwiftUI
import UIKit

/// 主内容视图
/// 使用 UIKit TabBar + NavigationController 承载四个 SwiftUI 标签页
struct ContentView: UIViewControllerRepresentable {

    init() {
        // 仅配置未选中色（主题黑）。
        // 选中色由 AppTabBarController 的实例 tintColor 设置，不写入 UITabBarItemAppearance.selected——
        // MapKit 可能重置 UITabBar.appearance()，实例 tintColor 更稳定。
        let appearance = UITabBarAppearance()
        appearance.configureWithDefaultBackground()
        let item = UITabBarItemAppearance()
        item.normal.iconColor = UIColor(Color.themePrimary)
        item.normal.titleTextAttributes = [.foregroundColor: UIColor(Color.themePrimary)]
        appearance.stackedLayoutAppearance = item
        UITabBar.appearance().standardAppearance = appearance
        UITabBar.appearance().scrollEdgeAppearance = appearance
    }

    func makeUIViewController(context: Context) -> AppTabBarController {
        AppTabBarController()
    }

    func updateUIViewController(_ uiViewController: AppTabBarController, context: Context) {
        // 根 Tab 结构固定，无需随 SwiftUI 状态更新。
    }
}
