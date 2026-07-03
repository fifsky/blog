import SwiftUI

/// 主内容视图
/// 包含四个标签页：博文、心情、提醒、足迹
struct ContentView: View {

    @State private var selectedTab = 0

    init() {
        // 配置 TabBar：未选中图标/文字颜色为深灰黑
        let appearance = UITabBarAppearance()
        appearance.configureWithDefaultBackground()
        let unselectedItem = UITabBarItemAppearance()
        unselectedItem.normal.iconColor = UIColor(Color.themePrimary)
        unselectedItem.normal.titleTextAttributes = [.foregroundColor: UIColor(Color.themePrimary)]
        unselectedItem.selected.iconColor = UIColor(Color(red: 0x00/255.0, green: 0x7d/255.0, blue: 0xfe/255.0))
        unselectedItem.selected.titleTextAttributes = [.foregroundColor: UIColor(Color(red: 0x00/255.0, green: 0x7d/255.0, blue: 0xfe/255.0))]
        appearance.stackedLayoutAppearance = unselectedItem
        UITabBar.appearance().standardAppearance = appearance
        UITabBar.appearance().scrollEdgeAppearance = appearance
    }

    var body: some View {
        TabView(selection: $selectedTab) {
            // 博文
            NavigationStack {
                ArticleListView()
            }
            .tabItem {
                Label("博文", systemImage: "doc.text")
            }
            .tag(0)

            // 心情
            NavigationStack {
                MoodListView()
            }
            .tabItem {
                Label("心情", systemImage: "face.smiling")
            }
            .tag(1)

            // 提醒
            NavigationStack {
                RemindListView()
            }
            .tabItem {
                Label("提醒", systemImage: "bell")
            }
            .tag(2)

            // 足迹
            NavigationStack {
                FootprintMapView()
            }
            .tabItem {
                Label("足迹", systemImage: "map")
            }
            .tag(3)
        }
        // 选中 tab 的图标/文字颜色
        .tint(Color(red: 0x00/255.0, green: 0x7d/255.0, blue: 0xfe/255.0))
    }
}
