import SwiftUI

/// 主内容视图
/// 包含四个标签页：博文、心情、提醒、足迹
struct ContentView: View {

    @State private var selectedTab = 0

    init() {
        // 仅配置未选中色（主题黑）。
        // 选中色交给 Asset Catalog AccentColor（品牌蓝），不用 UITabBarItemAppearance.selected——
        // MapKit 会重置 UITabBar.appearance()，但不会重置 Asset Catalog 的 AccentColor。
        let appearance = UITabBarAppearance()
        appearance.configureWithDefaultBackground()
        let item = UITabBarItemAppearance()
        item.normal.iconColor = UIColor(Color.themePrimary)
        item.normal.titleTextAttributes = [.foregroundColor: UIColor(Color.themePrimary)]
        appearance.stackedLayoutAppearance = item
        UITabBar.appearance().standardAppearance = appearance
        UITabBar.appearance().scrollEdgeAppearance = appearance
    }

    var body: some View {
        TabView(selection: $selectedTab) {
            // 博文
            NavigationStack {
                ArticleListView()
                    // 子页面 accentColor 用主题黑，不跟随 TabBar 选中色（品牌蓝）
                    .tint(Color.themePrimary)
            }
            .tabItem {
                Label("博文", systemImage: "doc.text")
            }
            .tag(0)

            // 心情
            NavigationStack {
                MoodListView()
                    .tint(Color.themePrimary)
            }
            .tabItem {
                Label("心情", systemImage: "face.smiling")
            }
            .tag(1)

            // 提醒
            NavigationStack {
                RemindListView()
                    .tint(Color.themePrimary)
            }
            .tabItem {
                Label("提醒", systemImage: "bell")
            }
            .tag(2)

            // 足迹
            NavigationStack {
                FootprintMapView()
                    .tint(Color.themePrimary)
            }
            .tabItem {
                Label("足迹", systemImage: "map")
            }
            .tag(3)
        }
        // 选中色由 Asset Catalog AccentColor 控制（品牌蓝 #007DFE）
    }
}
